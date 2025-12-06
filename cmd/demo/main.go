package main

import (
    "github.com/tominkoltd/go-table"
	"fmt"
	"os"
	"golang.org/x/term"
)

func main() {
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))
	

	workersTable := table.Table{Width: w}

	workersTable.Fields = append(workersTable.Fields, table.Field{
		Caption			: "PID",
		Flex			: 1,
		Align			: "right",
	})
	workersTable.Fields = append(workersTable.Fields, table.Field{
		Caption			: "Amount",
		Flex			: 1,
		Align			: "right",
		Prefix			: "£",
		IsNumber		: true,
		DecimalPlaces	: 2,
		Effects			: []table.Effect{
			{
				Color		: 31,
				Effect		: table.EffectBold | table.EffectBlink,
			},
		},
	})

	workersTable.Push("1", "45")
	workersTable.Push("3", "986")

	tbl := workersTable.Draw()

	fmt.Println(tbl)
}