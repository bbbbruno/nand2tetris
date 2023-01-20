package assemble

import (
	"bufio"
	"fmt"
	"io"
)

func Assemble(in io.Reader, out io.Writer) error {
	lines, st, err := ParseLabels(in)
	if err != nil {
		return err
	}
	cmds, err := ParseLines(lines, st)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(out)
	for _, cmd := range cmds {
		if _, err := fmt.Fprintln(w, cmd.Binary()); err != nil {
			return fmt.Errorf("failed to write to the file: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush the buffer: %w", err)
	}

	return nil
}
