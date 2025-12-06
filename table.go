package table

import (
	"fmt"
	"strconv"
	"strings"
	"bytes"
	"os"
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

type Effect struct {
	Function		func(value string) bool
	Color			int
	Background		int
	Effect			int
}

type Field struct {
	Name			string
	Caption			string
	Flex			int
	Width			int
	Align			string
	Wrap			bool
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
	Effects			[]Effect
	maxWidth		int
	drawWidth		int
}

type Group struct {
	Caption			string
	Data			[][]string
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
	Data			[][]string
	BoubleBorder	string
	HeadColor		int
	HeadBackground	int
	HeadEffect		int
	maxWidth		int
}

func (t *Table) Push(data ...any) {
	row := make([]string, len(data))
	for i, v := range data {
		row[i] = fmt.Sprint(v)
	}
	t.Data = append(t.Data, row)
}

func (g *Group) Push(data ...string) {
	row := make([]string, len(data))
	for i, v := range data {
		row[i] = fmt.Sprint(v)
	}
	g.Data = append(g.Data, row)
}

func calculateMaxRowWidths(t *Table) {
	t.maxWidth = 0
	usedSpace := 0
	totalFlex := 0
	for i, f := range t.Fields {
		f.maxWidth = 0
		usedSpace = usedSpace + f.Width
		if f.Width == 0 {
			for _, g := range t.Groups {
				for _, d := range g.Data {
					fieldLen := len(d[i])
					if fieldLen > f.maxWidth {
						f.maxWidth = fieldLen
					
					}
				}
			}
			for _, d := range t.Data {
				fieldLen := len(d[i])
				if fieldLen > f.maxWidth {
					f.maxWidth = fieldLen
				
				}
			}
			f.drawWidth = f.maxWidth
			if f.Flex == 0 {
				f.Flex = 1
			}
			totalFlex = totalFlex + f.Flex
		} else {
			f.maxWidth = f.Width
			f.drawWidth = f.Width
		}
		t.maxWidth = t.maxWidth + f.maxWidth
	}
	if t.Width == 0 {
		t.Width = t.maxWidth
	}
	fmt.Println("test")
}

func (t *Table) Draw() string {
	if len(t.Fields) == 0 {
		return ""
	}

	fmt.Fprintln(os.Stderr, "Calculating rows width\n")

	// calculate row widths
	calculateMaxRowWidths(t)

	// table width
	width := t.Width

	// if not set, calculate with fields
	if width == 0 {
		for i, f := range t.Fields {
			if f.Width == 0 {
				// get widest cell
				widest := 0
				// check table data
				for _, c := range t.Data {
					if len(c) <= i {
						break
					}
					colW := len(c[i]) + 2
					if colW > widest {
						widest = colW
					}
				}
				// check group data
				for _, g := range t.Groups {
					for _, c := range g.Data {
						if len(c) <= i {
							break
						}
						colW := len(c[i]) + 2
						if colW > widest {
							widest = colW
						}
					}
				}
				f.Width = widest
				width += widest
			} else {
				width += f.Width
			}
		}
	} 
	if width == 0 {
		// width is set, calculate widths by flex, width (no width/flex, flex = 1)
		N := len(t.Fields)

		// total non-content width: all vertical separators + all left/right paddings
		nonContent := (N+1)*sepWidth + (2*pad)*N

		// usable content width (sum of all column content widths)
		usable := t.Width - nonContent
		if usable < 0 {
			usable = 0 // clamp
		}

		// Sum fixed content widths and collect flex
		fixedTotal := 0
		type flexItem struct {
			index int
			flex  int
		}
		var flexes []flexItem
		totalFlex := 0

		for i, f := range t.Fields {
			if f.Width > 0 {
				fixedTotal += f.Width
			} else {
				fx := f.Flex
				if fx <= 0 {
					fx = 1
				}
				flexes = append(flexes, flexItem{index: i, flex: fx})
				totalFlex += fx
			}
		}

		// Remaining content width to distribute among flex columns
		remaining := usable - fixedTotal
		if remaining < 0 {
			remaining = 0 // (optional) later you can implement proportional shrink
		}

		// Distribute remaining by flex with integer carry to avoid drift
		if totalFlex > 0 && remaining > 0 {
			carry := 0
			for _, it := range flexes {
				numer := remaining*it.flex + carry
				w := numer / totalFlex
				carry = numer % totalFlex
				t.Fields[it.index].Width = w
			}	
		} else {
			for _, it := range flexes {
				t.Fields[it.index].Width = 0
			}
		}
	}
	
	var table bytes.Buffer

	// table header
	table.WriteString("┏━")
	for _, f := range t.Fields {
		table.WriteString(strings.Repeat("━", f.Width))
		table.WriteString("━┯━")
	}
	table.WriteString("━┓\n")



	return table.String()
}

func getAnsi(color int, background int, effect int) string {
	if color == 0 && background == 0 && effect == 0 {
		return ""
	}
	var codes []string
	if effect&EffectBold != 0 {
		codes = append(codes, "1")
	}
	if effect&EffectDim != 0 {
		codes = append(codes, "2")
	}
	if effect&EffectItalic != 0 {
		codes = append(codes, "3")
	}
	if effect&EffectUnderline != 0 {
		codes = append(codes, "4")
	}
	if effect&EffectBlink != 0 {
		codes = append(codes, "5")
	}
	if effect&EffectReverse != 0 {
		codes = append(codes, "7")
	}
	if effect&EffectStrikethrough != 0 {
		codes = append(codes, "9")
	}
	if effect&EffectOverline != 0 {
		codes = append(codes, "53")
	}
	if effect&EffectDoubleUnderline != 0 {
		codes = append(codes, "21")
	}
	if color > 0 {
		codes = append(codes, "38;5;" + strconv.Itoa(color))
	}
	if background > 0 {
		codes = append(codes, "48;5;" + strconv.Itoa(background))
	}
	if len(codes) == 0 {
		return ""
	}
	return "\033[" + strings.Join(codes, ";") + "m"
}