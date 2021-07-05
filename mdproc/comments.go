package mdproc

import (
	"bytes"
	"regexp"

	"github.com/n0x1m/md2gmi/pipe"
)

func RemoveComments(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		re := regexp.MustCompile(`<!--.*-->`)

		for b := range in {
			data := b.Payload()
			for _, match := range re.FindAllSubmatch(data, -1) {
				data = bytes.Replace(data, match[0], []byte(""), 1)
			}
			out <- pipe.NewItem(b.Index(), append(bytes.TrimSpace(data), '\n'))
		}

		close(out)
	}()

	return out
}
