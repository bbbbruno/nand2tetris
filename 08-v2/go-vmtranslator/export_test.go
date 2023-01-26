package vmtranslate

import "math/rand"

func (t *Translator) ExportSetRandomizer(r *rand.Rand) {
	t.r = r
}

func (t *Translator) ExportSetFunctionName(fn string) {
	t.fn = fn
}
