package main

import (
    "github.com/tominkoltd/go-table"
	"fmt"
	"os"
	"golang.org/x/term"
	"time"
	"strconv"
)

func main() {
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))
	

	workersTable := table.Table{Width: w, DrawHeader: true, BorderColor: 6}

	workersTable.Fields = append(workersTable.Fields, table.Field{
		Caption			: "PID",
		Flex			: 1,
		Align			: "center",
		Color			: 11,
	})
	workersTable.Fields = append(workersTable.Fields, table.Field{
		Caption			: "Amount",
		Flex			: 1,
		Align			: "right",
		Prefix			: "£",
		PrefixColor		: 11,
		IsNumber		: true,
		DecimalPlaces	: 2,
		Effect			: table.EffectBold,
	})
	workersTable.Fields = append(workersTable.Fields, table.Field{
		Caption			: "Description",
		Flex			: 2,
	})

	cc := &table.Cell{Data: "986", Color: 25}
	
	workersTable.Push("1", "45.334", "hjhgjhJHgjkh k hjcsbvlghj dsf")
	workersTable.Push("3", cc)

	for i:=0; i < 10; i++ {
		cc.Data=strconv.Itoa(i)
		tbl := workersTable.Draw(true)
		fmt.Println(tbl)
		time.Sleep(500 * time.Millisecond)
	}
}