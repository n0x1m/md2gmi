## md2gmi

Convert Markdown to Gemini Gemini "gemtext" markup with Go.

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
