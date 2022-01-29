package translator

import (
	"fmt"
	"math/rand"
)

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
M=0
`, comp)
}

func operateDouble(comp string) string {
	return pop("M") + pop(comp) + push()
}

func operateSingle(comp string) string {
	return pop(comp) + push()
}

func operateCompare(jump string) string {
	return pop("M") + pop("M-D") + compare(jump) + push()
}

func compare(jump string) string {
	label := rand.Intn(1000000)
	return fmt.Sprintf(`@TRUE%d
D;%s
D=0
@FINAL%d
0;JMP
(TRUE%d)
D=-1
(FINAL%d)
`, label, jump, label, label, label)
}

var segmentLabelMap = map[string]string{
	"argument": "ARG",
	"local":    "LCL",
	"this":     "THIS",
	"that":     "THAT",
}

func memoryPush(segment string, index int) string {
	label := segmentLabelMap[segment]
	return getSegment(label, index) + push()
}

func memoryPop(segment string, index int) string {
	label := segmentLabelMap[segment]
	return setAddress(label, index) + pop("M") + setSegment()
}

func getSegment(label string, index int) string {
	return fmt.Sprintf(`@%d
D=A
@%s
A=M+D
D=M
`, index, label)
}

func setAddress(label string, index int) string {
	return fmt.Sprintf(`@%d
D=A
@%s
D=M+D
@R13
M=D
`, index, label)
}

func setSegment() string {
	return `@R13
A=M
M=D
@R13
M=0
`
}

func constant(index int) string {
	return fmt.Sprintf(`@%d
D=A
`, index)
}

func end() string {
	return `(END)
@END
0;JMP
`
}
