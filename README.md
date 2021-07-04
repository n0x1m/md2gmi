## md2gmi

Convert Markdown to Gemini [gemtext](https://gemini.circumlunar.space/docs/gemtext.gmi) markup with
Go. Working with streams and pipes for UNIX like behavior utilizing Go channels. Processing streams
line by line is deliberately slightly more challenging than it needs to be to play around with go
state machines.

<!-- testing markdown, this should be deleted, below merged -->
See the [gemini
protocol](https://gemini.circumlunar.space/) and the [protocol
spec](https://gemini.circumlunar.space/docs/specification.gmi).

Internally md2gmi does a 1st pass that constructs the core layout for gemtext. This is then streamed
to the 2nd pass line by line. The 2nd pass will convert links and stream line by line to the output.

###Usage

```plain
Usage of ./md2gmi:
  -f string
        specify a .md (Markdown) file to read from, otherwise stdin (default)
  -o string
        specify a .gmi (gemtext) file to write to, otherwise stdout (default)
```

### Example

    go get github.com/n0x1m/md2gmi
    cat file.md | md2gmi
    md2gmi -in file.md -out file.gmi