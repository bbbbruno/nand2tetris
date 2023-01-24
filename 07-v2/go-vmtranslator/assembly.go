package vmtranslate

import (
	"fmt"
)

const END = `(END)
@END
0;JMP
`

var arithmeticAssembly = map[string]string{
	"add": pop("M") + pop("D+M") + push(),
	"sub": pop("M") + pop("M-D") + push(),
	"neg": pop("-M") + push(),
	"and": pop("M") + pop("D&M") + push(),
	"or":  pop("M") + pop("D|M") + push(),
	"not": pop("!M") + push(),
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

func push() string {
	return `@SP
A=M
M=D
@SP
M=M+1
`
}

func pop(comp string) string {
	return fmt.Sprintf(`@SP
M=M-1
A=M
D=%s
`, comp)
}
