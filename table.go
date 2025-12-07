package table

import (
	"fmt"
	"strconv"
	"strings"
	"bytes"
)

const (
	EffectNone				= 0
	EffectBold       		= 1 << iota
	EffectDim
	EffectItalic
	EffectUnderline
	EffectBlink
	EffectReverse
	EffectStrikethrough
	EffectOverline
	EffectDoubleUnderline
)

const ansiReset	= "\033[0m"

const sepWidth = 1  // width of vertical separator, e.g. '|' or '┃'
const pad = 1       // spaces on each side of cell content

type Field struct {
	Name			string
	Caption			string
	Flex			int
	Width			int
	Align			string
	Wrap			bool
	Color			int
	BackgroundColor	int
	Effect			int
	Prefix			string
	PrefixColor		int
	PrefixEffect	int
	Suffix			string
	SuffixColor		int
	SuffixEffect	int
	IsNumber		bool
	DecimalPlaces	int
	SortBy			string
	SortAsc			bool
	CountTotal		bool
	maxWidth		int
	drawWidth		int
}

type Group struct {
	Caption			string
	Data			[][]Cell
	SortBy			string
	SortAsc			bool
	HeadColor		int
	HeadBackground	int
	HeadEffect		int
}

type Table struct {
	Width			int
	Fields			[]Field
	Groups			[]Group
	Data			[][]Cell
	HeadColor		int
	HeadBackground	int
	HeadEffect		int
	DrawHeader		bool
	BorderColor		int
	BackgroundColor int
}

type Cell struct {
	Data			string
	Color			int
	BackgroundColor	int
	Effect			int
}

func (t *Table) Push(data ...any) {
	row := make([]Cell, len(t.Fields))
	for i, v := range data {
		switch cell := v.(type) {
			case Cell:
				row[i] = cell
			case *Cell:
				row[i] = *cell
			default:
				row[i] = Cell{
					Data: fmt.Sprint(v), 
					Color: 0, 
					BackgroundColor: 0, 
					Effect: 0, 
				}
		}
	}
	t.Data = append(t.Data, row)
}

func (g *Group) Push(data ...any) {
	row := make([]Cell, len(data))
	for i, v := range data {
		switch cell := v.(type) {
			case Cell:
				row[i] = cell
			case *Cell:
				row[i] = *cell
			default:
				row[i] = Cell{
					Data: fmt.Sprint(v), 
					Color: 0, 
					BackgroundColor: 0, 
					Effect: 0, 
				}
		}
	}
	g.Data = append(g.Data, row)
}

func calculateRowDrawWidths(t *Table) {
	usedSpace := 0
	totalFlex := 0
	flexFields := 0
	for _, f := range t.Fields {
		if f.Width > 0 {
			// have width set
			f.maxWidth = f.Width
			usedSpace = usedSpace + f.Width
			f.drawWidth = f.Width
			continue
		}
		flexFields ++
		if f.Flex == 0 {
			f.Flex = 1
		}
		totalFlex = totalFlex + f.Flex
	}
	if t.Width == 0 {
		t.Width = usedSpace + (flexFields * 3) + 2
	}
	tableWidth := t.Width - 2 - (len(t.Fields) * 2)
	flexWidth := (tableWidth - usedSpace) / totalFlex

	realWidth := 0
	for i := range t.Fields {
		f := &t.Fields[i]
		if f.Width > 0 {
			f.drawWidth = f.Width
		} else {
			f.drawWidth = flexWidth * f.Flex
		}
		realWidth += f.drawWidth + 2
	}

	realWidth += 1 + len(t.Fields)
	diff := t.Width - realWidth

	t.Fields[0].drawWidth += diff
}

func (t *Table) Draw() string {
	if len(t.Fields) == 0 {
		return ""
	}

	// calculate row draw widths
	calculateRowDrawWidths(t)
	
	var table bytes.Buffer

	// table top border
	table.Write(getAnsi(t.BorderColor, t.BackgroundColor, 0))
	
	table.WriteString("┏━")
	for _, f := range t.Fields {
		table.WriteString(strings.Repeat("━", f.drawWidth))
		table.WriteString("━┯━")
	}
	table.Truncate(table.Len() - len("┯━"))
	table.WriteString("┓\n")

	// table header
	if t.DrawHeader {
		table.WriteString("┃ ")
		for _, f := range t.Fields {
			table.Write(getAnsi(f.Color, t.BackgroundColor, 0))
			table.WriteString(drawField(f.Caption, f.drawWidth, f.Align))
			table.Write(getAnsi(t.BorderColor, t.BackgroundColor, 0))
			table.WriteString(" │ ")
		}
		table.Truncate(table.Len() - len("│ "))
		table.WriteString("┃\n")

		table.WriteString("┣━")
		for _, f := range t.Fields {
			table.WriteString(strings.Repeat("━", f.drawWidth))
			table.WriteString("━┿━")
		}
		table.Truncate(table.Len() - len("┿━"))
		table.WriteString("┫\n")
	}

	// groups TODO

	// data
	for i := range t.Data {
		table.WriteString("┃ ")
		for f := range t.Fields {
			fg := t.Fields[f].Color
			
			if t.Data[i][f].Color > 0 {
				fg = t.Data[i][f].Color
			}
			bg := t.Fields[f].BackgroundColor
			if t.Data[i][f].BackgroundColor > 0 {
				bg = t.Data[i][f].BackgroundColor
			}
			if bg == 0 {
				bg = t.BackgroundColor
			}
			ef := t.Fields[f].Effect
			if t.Data[i][f].Effect > 0 {
				ef = t.Data[i][f].Effect
			}

			// data
			cltxt := t.Data[i][f].Data
			if t.Fields[f].IsNumber {
				nm, err := strconv.ParseFloat(cltxt, 64)
				if err != nil {
					cltxt = "-"
				} else {
					cltxt = strconv.FormatFloat(nm, 'f', t.Fields[f].DecimalPlaces, 64)
				}
			}

			fLen := t.Fields[f].drawWidth

			// prefix
			if t.Fields[f].Prefix != "" {
				prs := []rune(t.Fields[f].Prefix)
				fLen -= len(prs)

				pfg := fg
				pef := ef

				if t.Fields[f].PrefixColor > 0 {
					pfg = t.Fields[f].PrefixColor
				}
				if t.Fields[f].PrefixEffect > 0 {
					pef = t.Fields[f].PrefixEffect
				}
				table.Write(getAnsi(pfg, bg, pef))
				table.WriteString(t.Fields[f].Prefix)
				table.Write(getAnsi(0, 0, 0))
			}

			if t.Fields[f].Suffix != "" {
				srs := []rune(t.Fields[f].Suffix)
				fLen -= len(srs)
			}

			table.Write(getAnsi(fg, bg, ef))
			table.WriteString(drawField(cltxt, fLen, t.Fields[f].Align))

			// suffix
			if t.Fields[f].Suffix != "" {
				if t.Fields[f].SuffixColor > 0 {
					fg = t.Fields[f].SuffixColor
				}
				if t.Fields[f].SuffixEffect > 0 {
					ef = t.Fields[f].SuffixEffect
				}
				table.Write(getAnsi(fg, bg, ef))
				table.WriteString(t.Fields[f].Suffix)
			}

			table.Write(getAnsi(0, 0, 0))
			table.Write(getAnsi(t.BorderColor, t.BackgroundColor, 0))
			table.WriteString(" │ ")
		}
		table.Truncate(table.Len() - len("│ "))
		table.WriteString("┃\n")
	}

	// table bottom border
	table.WriteString("┗━")
	for _, f := range t.Fields {
		table.WriteString(strings.Repeat("━", f.drawWidth))
		table.WriteString("━┷━")
	}
	table.Truncate(table.Len() - len("┷━"))
	table.WriteString("┛\n")
	table.Write(getAnsi(0, 0, 0))

	return table.String()
}

func drawField(s string, w int, a string) string {
	rs := []rune(s)

	if len(rs) > w {
		rs = rs[:w]
	}

    curWidth := len(rs)
    d := w - curWidth
    result := string(rs)

    if d > 0 {
        switch {
        case a == "" || strings.EqualFold(a, "Left"):
            result = result + strings.Repeat(" ", d)
        case strings.EqualFold(a, "Right"):
            result = strings.Repeat(" ", d) + result
        case strings.EqualFold(a, "Center"):
            ld := d / 2
            rd := d - ld
            result = strings.Repeat(" ", ld) + result + strings.Repeat(" ", rd)
        }
    }

    return result
}

func getAnsi(color int, background int, effect int) []byte {
	if color == 0 && background == 0 && effect == 0 {
		return []byte{27, 91, 48, 109}
	}
	var code bytes.Buffer
	code.WriteString("\033[")

	if effect == 0 {
		code.WriteString("0;")
	}
	if effect&EffectBold != 0 {
		code.WriteString("1;")
	}
	if effect&EffectDim != 0 {
		code.WriteString("2;")
	}
	if effect&EffectItalic != 0 {
		code.WriteString("3;")
	}
	if effect&EffectUnderline != 0 {
		code.WriteString("4;")
	}
	if effect&EffectBlink != 0 {
		code.WriteString("5;")
	}
	if effect&EffectReverse != 0 {
		code.WriteString("7;")
	}
	if effect&EffectStrikethrough != 0 {
		code.WriteString("9;")
	}
	if effect&EffectOverline != 0 {
		code.WriteString("53;")
	}
	if effect&EffectDoubleUnderline != 0 {
		code.WriteString("21;")
	}
	if color > 0 {
		code.WriteString("38;5;" + strconv.Itoa(color) + ";")
	} else {
		code.WriteString("39;")
	}
	if background > 0 {
		code.WriteString("48;5;" + strconv.Itoa(background) + ";")
	} else {
		code.WriteString("49;")
	}
	if code.Len() == 0 {
		return nil
	}
	code.Bytes()[code.Len()-1] = 'm'
	return code.Bytes()
}