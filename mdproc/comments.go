package mdproc

import (
	"bytes"
	"regexp"

	"github.com/n0x1m/md2gmi/pipe"
)

func RemoveComments(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		re := regexp.MustCompile(`(?s)<!--(.*?)-->`)

		for b := range in {
			data := b.Payload()
			touched := false

			for _, match := range re.FindAllSubmatch(data, -1) {
				data = bytes.Replace(data, match[0], []byte(""), 1)
				touched = true
			}

			if touched && len(bytes.TrimSpace(data)) == 0 {
				continue
			}
			out <- pipe.NewItem(b.Index(), data)
		}

		close(out)
	}()

	return out
}
