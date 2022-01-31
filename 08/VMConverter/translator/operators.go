package translator

import (
	"fmt"
	"math/rand"
)

func push(comp string) string {
	return fmt.Sprintf(`@SP
A=M
M=%s
@SP
M=M+1
`, comp)
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
	return pop("M") + pop(comp) + push("D")
}

func operateSingle(comp string) string {
	return pop(comp) + push("D")
}

func operateCompare(jump string) string {
	return pop("M") + pop("M-D") + compare(jump) + push("D")
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
