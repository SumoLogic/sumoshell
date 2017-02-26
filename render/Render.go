package main

import "fmt"
import "github.com/SumoLogic/sumoshell/util"
import "github.com/SumoLogic/sumoshell/render-util"
import "strings"
import "os/exec"
import "os"
import (
	"os/signal"
	"strconv"
	"syscall"
)

type Renderer struct {
	showRaw     bool
	colWidths   *map[string]int
	cols        *[]string
	height      int64
	width       int64
	rowsPrinted *int64
	inRelation  *bool
	limit       int64
}

func main() {
	m := make(map[string]int)
	cols := []string{}
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdout
	out, _ := cmd.Output()
	wh := strings.Split(string(out), " ")
	var hstr, wstr string
	renderAll := false
	if len(wh) == 1 {
		hstr = "0"
		wstr = "0"
		renderAll = true
	} else {
		wstr = strings.Trim(wh[1], "\n")
		hstr = strings.Trim(wh[0], "\n")
	}
	height, _ := strconv.ParseInt(hstr, 10, 64)
	width, _ := strconv.ParseInt(wstr, 10, 64)

	rows := int64(0)
	inRelation := false

	if len(os.Args) == 2 && os.Args[1] == "noraw" {
		util.ConnectToStdIn(Renderer{false, &m, &cols, height, width, &rows, &inRelation, 20})
	} else if (len(os.Args) == 2 && os.Args[1] == "all") || renderAll {
		c := make(chan os.Signal, 2)

		// Ignore SIGTERM signals. We will stop running anyway when our input is over.
		// This allows us to dump the current state
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
		}()
		rels := []map[string]interface{}{}
		holder := relationHolder{&rels}
		util.ConnectToStdIn(holder)
		r := Renderer{false, &m, &cols, height, width, &rows, &inRelation, -1}
		for _, item := range *holder.lastRelation {
			r.Process(item)
		}
	} else {
		util.ConnectToStdIn(Renderer{true, &m, &cols, height, width, &rows, &inRelation, 20})
	}
}

type relationHolder struct {
	lastRelation *[]map[string]interface{}
}

func (r relationHolder) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		panic("all only supports aggregate data")
	}
	if util.IsStartRelation(inp) {
		slice := []map[string]interface{}{inp}
		*r.lastRelation = slice
	}
	if util.IsRelation(inp) {
		slice := append(*r.lastRelation, inp)
		*r.lastRelation = slice
	}
	if util.IsEndRelation(inp) {
		slice := append(*r.lastRelation, inp)
		*r.lastRelation = slice
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
		if *r.inRelation {
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
		if len(*r.cols) > 0 {
			r.printHeader()
		}
	}
	if util.IsRelation(inp) {
		// If we haven't printed the header yet
		if *r.rowsPrinted >= r.limit && r.limit != -1 {
			return
		}
		if !*r.inRelation {
			panic("Can't get relation before StartRelation")
		}
		colsWidth := render.Columns([]map[string]interface{}{inp})
		colNames := render.ColumnNames(colsWidth)
		*r.cols = colNames
		*r.colWidths = colsWidth
		if *r.rowsPrinted == 0 {
			r.printHeader()
		}
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

func (r Renderer) printHeader() {
	for _, col := range *r.cols {
		width := (*r.colWidths)[col]
		spaces := strings.Repeat(" ", width-len(col))
		fmt.Printf("%v%s", col, spaces)
	}
	*r.rowsPrinted += 1
	fmt.Printf("\n")

}
