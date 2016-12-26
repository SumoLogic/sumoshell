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
	width       int64
	rowsPrinted *int64
	inRelation *bool
}

func main() {
	m := make(map[string]int)
	cols := []string{}
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdout
	out, _ := cmd.Output()
	wh := strings.Split(string(out), " ")
	hstr := strings.Trim(wh[0], "\n")
	wstr := strings.Trim(wh[1], "\n")
	height, _ := strconv.ParseInt(hstr, 10, 64)
	width, _ := strconv.ParseInt(wstr, 10, 64)

	rows := int64(0)
	inRelation := false

	if len(os.Args) == 2 && os.Args[1] == "noraw" {
		util.ConnectToStdIn(Renderer{false, &m, &cols, height, width, &rows, &inRelation})
	} else {
		util.ConnectToStdIn(Renderer{true, &m, &cols, height, width, &rows, &inRelation})
	}
}

func (r Renderer) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		var printed = false
		charsPrinted := int64(0)
		colsWidth := render.Columns([]map[string]interface{}{inp})
		colNames := render.ColumnNames(colsWidth)
		*r.cols = colNames
		*r.colWidths = colsWidth
		for _, col := range colNames {
			v, _ := inp[col]
			vStr := fmt.Sprint(v)
			spaces := strings.Repeat(" ", colsWidth[col]-len(vStr))
			finalStr := fmt.Sprintf("[%s=%v]%s", col, vStr, spaces)
			if charsPrinted+int64(len(finalStr)) > r.width {
				availableChars := r.width - charsPrinted
				if availableChars > 3 {
					fmt.Printf(finalStr[:availableChars-3])
					fmt.Printf("...")
					charsPrinted = r.width
				}
			} else {
				charsPrinted += int64(len(finalStr))
				fmt.Printf(finalStr)
			}
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
		if (*r.inRelation) {
			panic("Already in relation")
		}
		*r.inRelation = true
		for i := int64(0); i < *r.rowsPrinted; i++ {
			// Clear the row
			fmt.Printf("\033[2K")
			// Go up one row
			fmt.Printf("\033[1A")
		}
		*r.rowsPrinted = 0
		for _, col := range *r.cols {
			width := (*r.colWidths)[col]
			spaces := strings.Repeat(" ", width-len(col))
			fmt.Printf("%v%s", col, spaces)
		}
		*r.rowsPrinted += 1
		fmt.Printf("\n")
	}
	if util.IsRelation(inp) {
		if (!*r.inRelation) {
			panic("Can't get relation before StartRelation")
		}
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
		*r.inRelation = false
	}
}
