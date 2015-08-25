package main

import "fmt"
import "github.com/SumoLogic/sumoshell/util"
import "github.com/SumoLogic/sumoshell/render-util"
import "strings"
import "os/exec"
import "os"
import "strconv"

type Renderer struct {
	showRaw     bool
	colWidths   *map[string]int
	cols        *[]string
	height      int64
	rowsPrinted *int64
}

func main() {
	m := make(map[string]int)
	cols := []string{}
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdout
	out, _ := cmd.Output()
	wh := strings.Split(string(out), " ")
	hstr := strings.Trim(wh[0], "\n")
	rows := int64(0)
	height, _ := strconv.ParseInt(hstr, 10, 64)

	if len(os.Args) == 2 && os.Args[1] == "noraw" {
		util.ConnectToStdIn(Renderer{false, &m, &cols, height, &rows})
	} else {
		util.ConnectToStdIn(Renderer{true, &m, &cols, height, &rows})
	}
}

func (r Renderer) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		var printed = false
		colsWidth := render.Columns([]map[string]interface{}{inp})
		colNames := render.ColumnNames(colsWidth)
		*r.cols = colNames
		*r.colWidths = colsWidth
		for _, col := range colNames {
			v, _ := inp[col]
			vStr := fmt.Sprint(v)
			spaces := strings.Repeat(" ", (len(col)+3+colsWidth[col])-len(vStr))
			fmt.Printf("[%s=%v]%s", col, vStr, spaces)
		}
		if printed {
			fmt.Printf(";")
		}
		if r.showRaw {
			fmt.Printf(util.ExtractRaw(inp))
		}
		fmt.Print("\n")
	}
	if util.IsStartRelation(inp) {
		fmt.Println("======")
		for _, col := range *r.cols {
			width := (*r.colWidths)[col]
			spaces := strings.Repeat(" ", width-len(col))
			fmt.Printf("%v%s", col, spaces)
		}
		*r.rowsPrinted += 2
		fmt.Printf("\n")
	}
	if util.IsRelation(inp) {
		colsWidth := render.Columns([]map[string]interface{}{inp})
		colNames := render.ColumnNames(colsWidth)
		*r.cols = colNames
		*r.colWidths = colsWidth
		for _, col := range colNames {
			v, _ := inp[col]
			vStr := fmt.Sprint(v)
			spaces := strings.Repeat(" ", colsWidth[col]-len(vStr))
			fmt.Printf("%v%s", vStr, spaces)
		}
		*r.rowsPrinted += 1
		fmt.Printf("\n")
	}

	if util.IsEndRelation(inp) {
		if *r.rowsPrinted < r.height-1 {
			for i := *r.rowsPrinted; i < r.height; i++ {
				fmt.Printf("\n")
			}
		}
		*r.rowsPrinted = 0
	}
}
