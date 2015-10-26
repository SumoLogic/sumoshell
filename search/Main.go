package search

import "os"
import "bufio"
import "io"
import "github.com/SumoLogic/sumoshell/util"
import "strings"

func BuildAndConnect(args []string) {
	if len(args[1:]) > 0 {
		read(args[1])
	} else {
		read("")
	}
}

func read(filterString string) {
	r, w := io.Pipe()
	handler := util.NewRawInputHandler(w)
	go util.ConnectToReader(SumoFilter{filterString, util.NewJsonWriter()}, r)
	bio := bufio.NewReader(os.Stdin)
	var line, hasMoreInLine, err = bio.ReadLine()
	for err != io.EOF || hasMoreInLine {
		handler.Process(line)
		line, hasMoreInLine, err = bio.ReadLine()
	}
	handler.Flush()
}

type SumoFilter struct {
	param  string
	output *util.JsonWriter
}

func (filt SumoFilter) Process(inp map[string]interface{}) {
	if strings.Contains(util.ExtractRaw(inp), filt.param) {
		filt.output.Write(inp)
	}
}
