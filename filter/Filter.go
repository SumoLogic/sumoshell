package main

import (
	"github.com/SumoLogic/sumoshell/util"
	"log"
	"os"
)

type FilterOperator struct {
	key    string
	value  string
	output *util.JsonWriter
}

const genericError = "where takes arguments like: `where x = y`"

func main() {
	if len(os.Args) < 4 {
		log.Printf("Error! Not enough arguments provided.")
		log.Printf(genericError)
		return
	}

	key := os.Args[1]
	value := os.Args[3]

	util.ConnectToStdIn(FilterOperator{key, value, util.NewJsonWriter()})
}

func (w FilterOperator) Process(inp map[string]interface{}) {
	if util.IsPlus(inp) {
		v, exists := inp[w.key]
		if exists && v == w.value {
			w.output.Write(inp)
		}
	}
}
