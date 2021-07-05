package mdproc

import (
	"bytes"
	"regexp"

	"github.com/n0x1m/md2gmi/pipe"
)

func FormatHeadings(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		re := regexp.MustCompile(`^[#]{4,}`)
		re2 := regexp.MustCompile(`^(#+)[^# ]`)

		for b := range in {
			// fix up more than 4 levels
			data := re.ReplaceAll(b.Payload(), []byte("###"))
			// ensure we have a space
			sub := re2.FindSubmatch(data)
			if len(sub) > 0 {
				data = bytes.Replace(data, sub[1], append(sub[1], []byte(" ")...), 1)
			}
			// generally if we deal with a heading, add an extra blank line
			if bytes.HasPrefix(data, []byte("#")) {
				data = append(data, '\n')
			}
			// writeback
			out <- pipe.NewItem(b.Index(), data)
		}

		close(out)
	}()

	return out
}
