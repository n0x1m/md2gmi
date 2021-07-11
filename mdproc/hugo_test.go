package mdproc_test

import (
	"testing"

	"github.com/n0x1m/md2gmi/mdproc"
	"github.com/n0x1m/md2gmi/pipe"
)

func TestMdDocument2Gmi(t *testing.T) {
	t.Parallel()

	document := `---
title: "Dre's log"
---

nox, Latin "night; darkness"

` + "```" + `
	___
	(o,o)  < nox.im
	{` + "`" + `"'}    Fiat lux.
	-"-"-
` + "```" + `

[Gemini](gemini://nox.im) · [RSS](/index.xml) · [About](/about) · [Github](https://github.com/n0x1m)<!-- · [Twitter](https://twitter.com/_noxim) -->

Contact me via ` + "`" + `dre@nox.im` + "`" + `. You may use my [age](/snippets/actually-good-encryption/) public key to send me files securely: ` + "`" + `age1vpyptw64mz2vhtj7tvfh9saj0y8zy8fguety5n3wpmwzpkn0rd6swh02an` + "`" + `.

<!--
First principles bottom up thinker and tinkerer. Contact me via dre@nox.im. Please use my [GnuPG key](/noxim.asc).
Proudly made without PHP, Javascript, Ruby, Python and SQL.
-->

## Posts

`

	gmiout := `# Dre's log

nox, Latin "night; darkness"

` + "```" + `
	___
	(o,o)  < nox.im
	{` + "`" + `"'}    Fiat lux.
	-"-"-
` + "```" + `

Gemini[1] · RSS[2] · About[3] · Github[4]

=> gemini://nox.im 1: Gemini
=> /index.xml 2: RSS
=> /about 3: About
=> https://github.com/n0x1m 4: Github

Contact me via ` + "`" + `dre@nox.im` + "`" + `. You may use my age[1] public key to send me files securely: ` + "`" + `age1vpyptw64mz2vhtj7tvfh9saj0y8zy8fguety5n3wpmwzpkn0rd6swh02an` + "`" + `.

=> /snippets/actually-good-encryption/ 1: age

## Posts

`

	s := pipe.New()
	s.Use(mdproc.Preprocessor())
	s.Use(mdproc.RemoveFrontMatter)
	s.Use(mdproc.FormatHeadings)
	s.Use(mdproc.FormatLinks)
	s.Handle(source(t, document), sink(t, gmiout))
}
