package mdproc

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/n0x1m/md2gmi/pipe"
)

func FormatLinks(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		fenceOn := false

		for b := range in {
			data := b.Payload()
			if isFence(data) {
				fenceOn = !fenceOn
			}

			if fenceOn {
				out <- pipe.NewItem(b.Index(), b.Payload())

				continue
			}
			out <- pipe.NewItem(b.Index(), formatLinks(b.Payload()))
		}

		close(out)
	}()

	return out
}

func formatLinks(data []byte) []byte {
	// find link name and url
	var buffer []byte

	re := regexp.MustCompile(`!?\[([^\]*]*)\]\(([^ ]*)\)`)

	for i, match := range re.FindAllSubmatch(data, -1) {
		replaceWithIndex := append(match[1], fmt.Sprintf("[%d]", i+1)...)
		data = bytes.Replace(data, match[0], replaceWithIndex, 1)
		// append entry to buffer to be added later
		link := fmt.Sprintf("=> %s %d: %s\n", match[2], i+1, match[1])
		buffer = append(buffer, link...)
	}
	// append links to that paragraph
	if len(buffer) > 0 {
		data = append(data, buffer...)
		data = append(data, []byte("\n")...)
	}

	return data
}
