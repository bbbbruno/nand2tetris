package vmtranslate

import (
	"fmt"
	"math/rand"
	"time"
)

type Translator struct {
	f string
	r *rand.Rand
}

func NewTranslator(filename string) *Translator {
	return &Translator{f: filename, r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (t *Translator) Translate(cmd *Cmd) string {
	switch cmd.Type {
	case Arithmetic:
		return t.translateArithmetic(cmd)
	case Push:
		return t.translatePush(cmd)
	case Pop:
		return t.translatePop(cmd)
	case Flow:
		return t.translateFlow(cmd)
	default:
		return ""
	}
}

func (t *Translator) translateArithmetic(cmd *Cmd) string {
	switch cmd.Command {
	case "eq", "gt", "lt":
		return fmt.Sprintf(COMPARE_ASSEMBLY, t.r.Intn(1_000_000), t.r.Intn(1_000_000), jmpAssembly[cmd.Command])
	default:
		return arithmeticAssembly[cmd.Command]
	}
}

func (t *Translator) translatePush(cmd *Cmd) string {
	return pushAssembly(cmd.Arg1, cmd.Arg2, t.f)
}

func (t *Translator) translatePop(cmd *Cmd) string {
	return popAssembly(cmd.Arg1, cmd.Arg2, t.f)
}

func (t *Translator) translateFlow(cmd *Cmd) string {
	return flowAssembly(cmd.Command, cmd.Arg1)
}
