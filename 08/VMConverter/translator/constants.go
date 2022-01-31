package translator

import "fmt"

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
