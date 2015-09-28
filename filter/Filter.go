package filter

import (
	"github.com/SumoLogic/sumoshell/util"
)

type FilterOperator struct {
	key    string
	value  string
	output *util.JsonWriter
}

const genericError = "filter takes arguments like: `filter x = y`"

func Build(args []string) (util.SumoOperator, error) {
	if len(args) < 4 {
		return nil, util.ParseError("Error! Not enough arguments provided.\n" + genericError)
	}

	key := args[1]
	eq := args[2]
	if eq != "=" {
		return nil, util.ParseError("Expected `=` found `" + eq + "`\n" + genericError)
	}
	value := args[3]

	return &FilterOperator{key, value, util.NewJsonWriter()}, nil
}

func (w FilterOperator) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		v, exists := inp[w.key]
		if exists && v == w.value {
			w.output.Write(inp)
		}
	}
}
