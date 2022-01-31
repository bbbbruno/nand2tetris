package translator

import "fmt"

func defineLabel(label string, currentFuncName string) string {
	return fmt.Sprintf(`(%s%s)
`, funcname(currentFuncName), label)
}

func goTo(label string, currentFuncName string) string {
	return fmt.Sprintf(`@%s%s
0;JMP
`, funcname(currentFuncName), label)
}

func ifGoTo(label string, currentFuncName string) string {
	return fmt.Sprintf(`@%s%s
D;JNE
`, funcname(currentFuncName), label)
}

func funcname(name string) (funcname string) {
	if name != "" {
		funcname = name + "$"
	}
	return funcname
}
