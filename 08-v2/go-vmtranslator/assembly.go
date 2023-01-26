package vmtranslate

import (
	"fmt"
)

const END = `(END)
@END
0;JMP
`

var arithmeticAssembly = map[string]string{
	"add": popCmp("M") + popCmp("D+M") + push(),
	"sub": popCmp("M") + popCmp("M-D") + push(),
	"neg": popCmp("-M") + push(),
	"and": popCmp("M") + popCmp("D&M") + push(),
	"or":  popCmp("M") + popCmp("D|M") + push(),
	"not": popCmp("!M") + push(),
}

const COMPARE_ASSEMBLY = `@SP
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
`

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

func flowAssembly(command, symbol string) string {
	switch command {
	case "label":
		return fmt.Sprintf(`(%s)
`, symbol)
	case "goto":
		return fmt.Sprintf(`@%s
0;JMP
`, symbol)
	case "if-goto":
		return pop() + fmt.Sprintf(`@%s
D;JNE
`, symbol)
	default:
		return ""
	}
}
