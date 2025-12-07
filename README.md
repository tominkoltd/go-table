# 🧱 **TableWriter — A Lightweight ANSI Table Renderer for Go**

A fast, dependency-free, Unicode-aware, ANSI-styled table renderer for terminal applications.
Designed for simplicity, predictable output, and fully customizable formatting.

Useful for:

- CLI apps
- Process managers
- Debuggers
- Monitoring tools
- Log viewers
- Anything that needs **clean terminal tables** with **colors + Unicode**

No external libraries. No magic. 100% Go.

## ✨ Features

- **Simple declarative column definitions**
- **Automatic flexible column sizing (`Flex` / `Width`)**
- **Clean ANSI colors + background + effects (bold, underline, etc.)**
- **Per-prefix and per-suffix styling**
- **Unicode-safe alignment** (`č, ž, ä`, box-drawing chars, etc.)
- **Number formatting with decimal rounding**
- **Pure Go, zero dependencies**
- **User-friendly API (`Push(...)`) with auto type handling**
- **Nice Unicode borders (┏━┯━ etc.)**
- **Optional header rendering**
- **Works perfectly with full-width terminals**

## 📦 Installation

```
go get github.com/yourusername/table
```

(Replace with your actual repo path.)

## 🚀 Quick Example

```go
package main

import (
    "fmt"
    "os"
    "github.com/yourusername/table"
    term "golang.org/x/term"
)

func main() {
    w, _, _ := term.GetSize(int(os.Stdout.Fd()))

    workers := table.Table{
        Width:      w,
        DrawHeader: true,
        BorderColor: 6,
    }

    workers.Fields = []table.Field{
        {Caption: "PID", Flex: 1, Align: "center", Color: 11},
        {Caption: "Amount", Flex: 1, Align: "right", Prefix: "£",
         PrefixColor: 11, IsNumber: true, DecimalPlaces: 2, Effect: table.EffectBold},
        {Caption: "Description", Flex: 2},
    }

    workers.Push("1", "45.334", "hello world")
    workers.Push("3", table.Cell{Data: "986", Color: 25}, "some text")

    fmt.Println(workers.Draw())
}
```

### Output

```
┏━━━━━━━━┯━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━┓
┃  PID   │   Amount │      Description       ┃
┣━━━━━━━━┿━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━┫
┃   1    │    £45.33 │ hello world            ┃
┃   3    │   £986.00 │ some text              ┃
┗━━━━━━━━┷━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━┛
```

(Colors hidden in GitHub preview.)

## 🛠 Field Options

Each column is defined using a simple struct:

```go
type Field struct {
    Caption         string
    Flex            int
    Width           int
    Align           string
    Color           int
    BackgroundColor int
    Effect          int

    Prefix          string
    PrefixColor     int
    PrefixEffect    int

    Suffix          string
    SuffixColor     int
    SuffixEffect    int

    IsNumber        bool
    DecimalPlaces   int
}
```

### Highlight: Number Formatting

Automatically rounds + formats numbers:

```go
{ Caption: "Amount", IsNumber: true, DecimalPlaces: 2 }
```

- `"45.3339"` → `"45.33"`
- `"986"` → `"986.00"`

## 🖌 ANSI Styling

Built-in bitmask effects:

```
EffectBold
EffectDim
EffectItalic
EffectUnderline
EffectBlink
EffectReverse
EffectStrikethrough
EffectOverline
EffectDoubleUnderline
```

Use them like:

```go
Effect: table.EffectBold | table.EffectUnderline
```

## 🎨 Colors

You can color:

- fields
- cells
- prefixes
- suffixes
- headers
- borders

All using 256-color ANSI:

```go
Color: 11
BackgroundColor: 0
BorderColor: 6
```

## 📐 Unicode-Safe Alignment

`drawField()` uses `[]rune` internally:

- Correctly handles `č, ď, ž, ä`
- Works with box drawing (`┏━┯┓`)
- No external dependencies

## 📤 API

### Push rows

```go
table.Push("1", "45.33", "description")
table.Push(table.Cell{Data: "999", Color: 12})
```

You can mix raw values and styled cells.

### Draw the table

```go
output := table.Draw()
fmt.Println(output)
```

## 🧩 Roadmap

- Group support (`Group.Push`, headers, totals)
- Custom border themes (ASCII, light, heavy)
- Multi-line / wrapping cells
- Per-column min/max widths
- Automatic numeric detection
- Sorting by field/group
- Row separators
- Markdown / CSV export modes
- Emoji width handling (optional)

## 💬 Why Another Tablewriter?

Because existing Go table libraries are:

- heavy
- verbose
- full of `.SetXxx()` chains
- not Unicode-safe
- not ANSI-clean
- too big for small/medium CLI apps

This library is:

- tiny
- dependency-free
- predictable
- beautifully configurable
- perfect for real CLI applications

## 📝 License

MIT License (or whichever you choose)
