## md2gmi

Convert Markdown to Gemini Gemini [gemtext](https://gemini.circumlunar.space/docs/gemtext.gmi)
markup with Go. Working with streams and pipes for UNIX like behavior utilizing Go channels.

### Usage

```
Usage of ./md2gmi:
  -in string
        specify a .md (Markdown) file to read from, otherwise stdin (default)
  -out string
        specify a .gmi (gemtext) file to write to, otherwise stdout (default)
```

### Example

    go get github.com/n0x1m/md2gmi
    cat file.md | md2gmi
    md2gmi -in file.md -out file.gmi