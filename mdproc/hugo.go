package mdproc

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/n0x1m/md2gmi/pipe"
)

func RemoveFrontMatter(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		// delete the entire front matter
		re := regexp.MustCompile(`---.*---`)
		// but parse out the title as we want to reinject it
		re2 := regexp.MustCompile(`title:[ "]*([a-zA-Z0-9 :!'@#$%^&*)(]+)["]*`)

		for b := range in {
			data := b.Payload()
			for _, match := range re.FindAllSubmatch(data, -1) {
				data = bytes.Replace(data, match[0], []byte(""), 1)
				for _, title := range re2.FindAllSubmatch(match[0], 1) {
					// add title
					data = []byte(fmt.Sprintf("# %s\n\n", title[1]))
				}
			}
			out <- pipe.NewItem(b.Index(), data)
		}

		close(out)
	}()

	return out
}
