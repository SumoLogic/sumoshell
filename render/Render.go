package main

import (
	"fmt"
	"log"
	"github.com/jroimartin/gocui"
    "github.com/SumoLogic/sumoshell/render-util"
)

// Need to determine all the columns in the current data, then render based on that



func makeLayout(state *render.RenderState) layoutFunc {
	return func(g *gocui.Gui) error {
		_, maxY := g.Size()
		viewColumns := render.Columns(*state.Messages)
		viewNames := render.ColumnNames(viewColumns)
		var pos = 0
		//log.Println(viewColumns)
		for _, key := range(viewNames) {
			width := viewColumns[key]
			if v, err := g.SetView(key, pos, 0, pos + width, maxY); err != nil {
		        if err != gocui.ErrorUnkView {
		            return err
		        }
		        fmt.Fprintln(v, key)
    		}
    		pos += width + 1
		}
    	update(viewNames, state, g)
    	return nil
	}
}


func update(columns []string, state *render.RenderState, g *gocui.Gui) {
	for _, column := range columns {
		v, _ := g.View(column)
		v.Clear()
		fmt.Fprintln(v, column)
		for _, row := range *state.Messages {
			fmt.Fprintln(v, row[column])
		}
	}
	
}

type layoutFunc func(*gocui.Gui) (error)

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.Quit
}

func main() {
    var err error
    g := gocui.NewGui()
    if err := g.Init(); err != nil {
        log.Panicln(err)
    }
    defer g.Close()

    renderState := render.NewConnectedRenderState(g.Flush)
    g.SetLayout(makeLayout(renderState))
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        log.Panicln(err)
    }



    //go
   
    g.MainLoop()
  
    if err != nil && err != gocui.Quit {
        log.Panicln(err)
    }

}


