@0
D=A
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
@LCL
D=M+D
@R13
M=D
@SP
M=M-1
A=M
D=M
M=0
@R13
A=M
M=D
@R13
M=0
(LOOP_START)
@0
D=A
@ARG
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
@LCL
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
M=0
@SP
M=M-1
A=M
D=M+D
M=0
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
@LCL
D=M+D
@R13
M=D
@SP
M=M-1
A=M
D=M
M=0
@R13
A=M
M=D
@R13
M=0
@0
D=A
@ARG
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
M=0
@SP
M=M-1
A=M
D=M-D
M=0
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
@ARG
D=M+D
@R13
M=D
@SP
M=M-1
A=M
D=M
M=0
@R13
A=M
M=D
@R13
M=0
@0
D=A
@ARG
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
M=0
@LOOP_START
D;JNE
@0
D=A
@LCL
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
(END)
@END
0;JMP
