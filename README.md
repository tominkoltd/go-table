# 🧱 Table Go Module — A Lightweight ANSI Table Renderer for Go
**Version 0.3**  
https://ko-fi.com/vanthomas ❤️

A fast, dependency-free, Unicode-aware, ANSI-styled table renderer for terminal applications.  
Designed for simplicity, predictable output, and fully customizable formatting.

Useful for:

- CLI apps  
- Debuggers  
- Monitoring dashboards  
- Log/worker supervisors  
- Live-updating tables (NEW in v0.3)  
- Anything needing clean colored terminal tables  

No external libraries. No magic. 100% Go.

---

## ✨ Features

- Simple declarative column definitions  
- Automatic flexible column sizing (`Flex` / `Width`)  
- ANSI colors + effects (bold, underline, reverse, etc.)  
- Per-prefix and per-suffix styling  
- Unicode-safe alignment  
- Number formatting with rounding  
- Pure Go, zero dependencies  
- User-friendly `Push(...)` API  
- Beautiful Unicode borders  
- Optional header drawing  
- Stable terminal width handling  

### 🔥 New in v0.3

- **Live table redraw** using ANSI cursor-up (`Draw(true)`)  
- **Pointer-based cell updates** (mutate cell → table updates automatically)  
- **`Draw()` = normal output, `Draw(true)` = redraw mode**  
- Flicker-free refreshing without clearing lines  
- Perfect for worker dashboards, counters, animations, etc.

Example usage of the new feature:

```go
cc := &table.Cell{Data: "986", Color: 25}

workersTable.Push("1", "45.334", "normal row")
workersTable.Push("3", cc)

for i := 0; i < 10; i++ {
    cc.Data = strconv.Itoa(i)
    fmt.Println(workersTable.Draw(true))
    time.Sleep(500 * time.Millisecond)
}
```

---

## 📦 Installation

```bash
go get github.com/tominkoltd/go-table
```

---

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

Supports:

- captions / headers  
- flexible or fixed widths  
- left / right / center alignment  
- foreground + background colors  
- ANSI text effects  
- numeric mode with decimal rounding  
- per-column prefixes and suffixes (e.g. currency)

---

## 🎨 ANSI Effects

Bitmask flags:

- `EffectBold`  
- `EffectDim`  
- `EffectItalic`  
- `EffectUnderline`  
- `EffectBlink`  
- `EffectReverse`  
- `EffectStrikethrough`  
- `EffectOverline`  
- `EffectDoubleUnderline`  

Combine with OR, for example:

```go
Effect: table.EffectBold | table.EffectUnderline
```

---

## 📤 API Overview

### Push rows

```go
table.Push("1", "45.33", "description")
table.Push(table.Cell{Data: "999", Color: 12})
table.Push(&table.Cell{Data: "986", Color: 25})
```

You can mix raw values and styled cells, including pointers for live updates.

### Draw the table

```go
output := table.Draw()      // normal render
output := table.Draw(true)  // redraw mode (moves cursor up)
fmt.Println(output)
```

- `Draw()`  
  - Renders the table normally (no cursor movement).

- `Draw(true)`  
  - Moves the cursor up by the number of lines drawn previously  
  - Allows in-place redraw of the same table area  

---

## 🧩 Typical Use Cases

- Worker / process overviews  
- Queue / job statistics  
- Live counters and metrics  
- CLI dashboards and status pages  
- Anywhere you’d otherwise spam `fmt.Println` but want it pretty

---

## 💬 Why Another Tablewriter?

Many existing Go table libraries are:

- heavy  
- verbose  
- full of `.SetXxx()` chains  
- not Unicode-safe  
- not ANSI-clean  
- not designed for live redraw

This library is:

- tiny  
- dependency-free  
- predictable  
- ANSI-focused  
- built for real-world CLI apps

---

## 💖 Support Development

If this library saves you time or looks nice in your terminal:

👉 https://ko-fi.com/vanthomas

---

## 📝 License

MIT License
