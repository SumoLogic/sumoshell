// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
package main

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/render-util"
	ui "gopkg.in/gizak/termui.v1"
	"strconv"
)

func main() {
	flush := func() error { return nil }
	renderState := render.NewConnectedRenderState(flush)
	createUi(renderState)
}

func createUi(state *render.RenderState) {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	ui.UseTheme("helloworld")

	bc := ui.NewBarChart()
	data := []int{}
	bclabels := []string{}
	bc.Border.Label = "SumoCLI"
	bc.Data = data
	//bc.Width = 26
	bc.Height = ui.TermHeight() - 10
	bc.DataLabels = bclabels
	bc.TextColor = ui.ColorGreen
	bc.BarColor = ui.ColorRed
	bc.NumColor = ui.ColorYellow

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, bc)))

	// calculate layout
	ui.Body.Align()

	done := make(chan bool)
	redraw := make(chan bool)

	update := func() error {
		//fmt.Println("updating")
		bars := render.Columns(*state.Messages)
		columns := render.ColumnNames(bars)
		query, ok := (*state.Meta)["_queryString"]
		if ok && query.(string) != bc.Border.Label {
			bc.Border.Label = query.(string)
		}

		dataCol := render.NumericColumn(columns)
		data := []int{}
		labels := []string{}
		extractor := render.LabelExtractor(columns)
		for _, msg := range *state.Messages {
			labels = append(labels, extractor(msg))
			// go is the worst
			floatVal, _ := strconv.ParseFloat(fmt.Sprint(msg[dataCol]), 64)
			rowData := int(floatVal)
			data = append(data, rowData)
		}
		bc.Data = data
		var numCols int
		if len(data) > 0 {
			numCols = len(data)
		} else {
			numCols = 1
		}

		newBarWidth := (ui.TermWidth() - 10) / numCols
		if int(newBarWidth) != int(bc.BarWidth) {
			bc.BarWidth = newBarWidth
		}

		if len(labels) != len(bc.DataLabels) {
			bc.DataLabels = labels
		}
		redraw <- true
		return nil
	}

	state.Flush = update
	evt := ui.EventCh()

	ui.Render(ui.Body)
	go update()

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == ui.EventResize {
				ui.Body.Width = ui.TermWidth()
				ui.Body.Align()
				go func() { redraw <- true }()
			}
		case <-done:
			return
		case <-redraw:
			ui.Render(ui.Body)
		}
	}
}
