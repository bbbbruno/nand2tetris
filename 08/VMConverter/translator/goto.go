package translator

import "fmt"

func defineLabel(label string) string {
	return fmt.Sprintf(`(%s)
`, label)
}

func goTo(label string) string {
	return fmt.Sprintf(`@%s
0;JMP
`, label)
}

func ifGoTo(label string) string {
	return fmt.Sprintf(`@%s
D;JNE
`, label)
}
