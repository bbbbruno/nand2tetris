package translator

import "fmt"

func defineFunc(name string) string {
	return fmt.Sprintf(`(%s)
`, name)
}

func execReturn() string {
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
M=0
@ARG
A=M
M=D
@R13
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
@R13
M=0
@R14
D=M
M=0
A=D
0;JMP
`
}
