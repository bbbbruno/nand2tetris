package vmtranslate

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func initAssembly() string {
	return `@256
D=A
@SP
M=D
` + functionAssembly("call", "Sys.init", 0, rand.New(rand.NewSource(time.Now().UnixNano())))
}

var arithmeticAssembly = map[string]string{
	"add": popCmp("M") + popCmp("D+M") + push(),
	"sub": popCmp("M") + popCmp("M-D") + push(),
	"neg": popCmp("-M") + push(),
	"and": popCmp("M") + popCmp("D&M") + push(),
	"or":  popCmp("M") + popCmp("D|M") + push(),
	"not": popCmp("!M") + push(),
}

var jmpAssembly = map[string]string{
	"eq": "JEQ",
	"gt": "JGT",
	"lt": "JLT",
}

var symbolAssembly = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
	"pointer":  "R3",
	"temp":     "R5",
}

func push() string {
	return `@SP
A=M
M=D
@SP
M=M+1
`
}

func pop() string {
	return popCmp("M")
}

func popCmp(comp string) string {
	return fmt.Sprintf(`@SP
M=M-1
A=M
D=%s
`, comp)
}

func compareAssembly(command string, r *rand.Rand) string {
	return fmt.Sprintf(`@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@TRUE%[1]d
D;%[3]s
D=0
@FINALLY%[2]d
0;JMP
(TRUE%[1]d)
D=-1
(FINALLY%[2]d)
@SP
A=M
M=D
@SP
M=M+1
`, r.Intn(1_000_000), r.Intn(1_000_000), jmpAssembly[command])
}

func pushAssembly(segment string, index int, fname string) string {
	sym := symbolAssembly[segment]
	switch segment {
	case "constant":
		return fmt.Sprintf(`@%d
D=A
`, index) + push()
	case "local", "argument", "this", "that":
		return fmt.Sprintf(`@%d
D=A
@%s
A=M+D
D=M
`, index, sym) + push()
	case "pointer", "temp":
		return fmt.Sprintf(`@%d
D=A
@%s
A=A+D
D=M
`, index, sym) + push()
	case "static":
		return fmt.Sprintf(`@%s.%d
D=M
`, fname, index) + push()
	default:
		return ""
	}
}

func popAssembly(segment string, index int, fname string) string {
	sym := symbolAssembly[segment]
	switch segment {
	case "local", "argument", "this", "that":
		return fmt.Sprintf(`@%d
D=A
@%s
D=M+D
@R13
M=D
`, index, sym) + pop() + `@R13
A=M
M=D
`
	case "pointer", "temp":
		return fmt.Sprintf(`@%d
D=A
@%s
D=A+D
@R13
M=D
`, index, sym) + pop() + `@R13
A=M
M=D
`
	case "static":
		return pop() + fmt.Sprintf(`@%s.%d
M=D
`, fname, index)
	default:
		return ""
	}
}

func flowAssembly(command, symbol, fn string) string {
	var s string
	if fn != "" {
		s = fn + "$"
	}
	switch command {
	case "label":
		return fmt.Sprintf(`(%s%s)
`, s, symbol)
	case "goto":
		return fmt.Sprintf(`@%s%s
0;JMP
`, s, symbol)
	case "if-goto":
		return pop() + fmt.Sprintf(`@%s%s
D;JNE
`, s, symbol)
	default:
		return ""
	}
}

func functionAssembly(command, symbol string, n int, r *rand.Rand) string {
	switch command {
	case "function":
		return fmt.Sprintf(`(%s)
`, symbol) + strings.Repeat(`@SP
A=M
M=0
@SP
M=M+1
`, n)
	case "call":
		ret := r.Intn(1_000_000)
		return fmt.Sprintf(`@RETURN%[1]d
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@%[2]d
D=A
@5
D=A+D
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@%[3]s
0;JMP
(RETURN%[1]d)
`, ret, n, symbol)
	case "return":
		return `@LCL
D=M
@R13
M=D
@5
D=A
@R13
A=M-D
D=M
@R14
M=D
@SP
M=M-1
A=M
D=M
@ARG
A=M
M=D
@ARG
D=M+1
@SP
M=D
@1
D=A
@R13
A=M-D
D=M
@THAT
M=D
@2
D=A
@R13
A=M-D
D=M
@THIS
M=D
@3
D=A
@R13
A=M-D
D=M
@ARG
M=D
@4
D=A
@R13
A=M-D
D=M
@LCL
M=D
@R14
A=M
0;JMP
`
	default:
		return ""
	}
}
