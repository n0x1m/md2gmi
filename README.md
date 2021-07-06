## md2gmi

Convert Markdown to Gemini [gemtext](https://gemini.circumlunar.space/docs/gemtext.gmi) markup with
Go. Working with streams and pipes for UNIX like behavior utilizing Go channels. Processing streams
line by line is slightly more complex than it needs to be as I was toying with channels and state
machines.

Internally md2gmi does a 1st pass that constructs the blocks of single lines for gemtext from one or
multiple lines of an input stream. These blocks are then streamed to the 2nd passes. The 2nd pass
will convert hugo front matters, links, fix headings etc. These stages/passes can be composed and
chained with go pipelines. The output sink is either a file or stdout.

### Usage

```plain
Usage of ./md2gmi:
  -i string
        specify a .md (Markdown) file to read from, otherwise stdin (default)
  -o string
        specify a .gmi (gemtext) file to write to, otherwise stdout (default)
```

### Example

    go get github.com/n0x1m/md2gmi
    cat file.md | md2gmi
    md2gmi -i file.md -o file.gmi

The top part of this readme parses from

```markdown
Convert Markdown to Gemini [gemtext](https://gemini.circumlunar.space/docs/gemtext.gmi) markup with
Go. Working with streams and pipes for UNIX like behavior utilizing Go channels. Processing streams
line by line is slightly more complex than it needs to be as I'm playing with channels and state
machines here.

> this is
a quote

<!-- testing markdown, this should be deleted, below merged -->
See the [gemini
protocol](https://gemini.circumlunar.space/) and the [protocol
spec](https://gemini.circumlunar.space/docs/specification.gmi).
```

to

```markdown
Convert Markdown to Gemini gemtext[1] markup with Go. Working with streams and pipes for UNIX like behavior utilizing Go channels. Processing streams line by line is slightly more complex than it needs to be as I'm playing with channels and state machines here.

=> https://gemini.circumlunar.space/docs/gemtext.gmi 1: gemtext

> this is a quote
See the gemini protocol[1] and the protocol spec[2].

=> https://gemini.circumlunar.space/ 1: gemini protocol
=> https://gemini.circumlunar.space/docs/specification.gmi 2: protocol spec
```