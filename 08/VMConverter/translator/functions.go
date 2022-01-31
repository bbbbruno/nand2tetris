package translator

import (
	"fmt"
	"math/rand"
)

func setSP() string {
	return `@256
D=A
@SP
M=D
`
}

func defineFunc(name string) string {
	return fmt.Sprintf(`(%s)
`, name)
}

func returnFunc() string {
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

func callFunc(name string, numargs int) string {
	label := rand.Intn(1000000)
	return fmt.Sprintf(`@RETURN%d
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
@%d
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
@%s
0;JMP
(RETURN%d)
`, label, numargs, name, label)
}
