package main

import "fmt"
import "github.com/rcoh/sumo-line/util"
import "strings"

type Renderer struct {}

func main() {
	util.ConnectToStdIn(Renderer{})
}

func (r Renderer) Process(inp map[string]interface{}) {
	var printed = false
	for k,v := range inp {
		if !strings.HasPrefix(k, "_") {
			vStr := fmt.Sprint(v)
			fmt.Printf("[%v=%s]", k, vStr)
			printed = true
		}
	}
	if printed {
		fmt.Printf(";")
	}
	fmt.Printf(util.ExtractRaw(inp))
	fmt.Print("\n")
}