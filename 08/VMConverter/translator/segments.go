package translator

import "fmt"

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
