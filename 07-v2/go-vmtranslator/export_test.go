package vmtranslate

import (
	"math/rand"
)

func (c *ArithmeticCmd) ExportSetRandomizer(r *rand.Rand) {
	c.r = r
}
