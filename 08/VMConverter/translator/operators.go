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
	"pointer":  "R3",
	"temp":     "R5",
}

func memoryPush(segment string, index int) string {
	label, comp := parseSegment(segment)
	return getSegment(label, index, comp) + push()
}

func memoryPop(segment string, index int) string {
	label, comp := parseSegment(segment)
	return setAddress(label, index, comp) + pop("M") + setSegment()
}

func parseSegment(segment string) (label, comp string) {
	label = segmentLabelMap[segment]
	switch segment {
	case "pointer", "temp":
		comp = "A+D"
	default:
		comp = "M+D"
	}
	return label, comp
}

func getSegment(label string, index int, comp string) string {
	return fmt.Sprintf(`@%d
D=A
@%s
A=%s
D=M
`, index, label, comp)
}

func setAddress(label string, index int, comp string) string {
	return fmt.Sprintf(`@%d
D=A
@%s
D=%s
@R13
M=D
`, index, label, comp)
}

func setSegment() string {
	return `@R13
A=M
M=D
@R13
M=0
`
}

func getConstant(index int) string {
	return fmt.Sprintf(`@%d
D=A
`, index)
}

func getStatic(filename string, index int) string {
	return fmt.Sprintf(`@%s.%d
D=M
`, filename, index)
}

func setStatic(filename string, index int) string {
	return fmt.Sprintf(`@%s.%d
M=D
`, filename, index)
}

func end() string {
	return `(END)
@END
0;JMP
`
}
