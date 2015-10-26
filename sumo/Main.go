package main

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/average"
	"github.com/SumoLogic/sumoshell/count"
	"github.com/SumoLogic/sumoshell/filter"
	"github.com/SumoLogic/sumoshell/parse"
	"github.com/SumoLogic/sumoshell/search"
	"github.com/SumoLogic/sumoshell/sum"
	"github.com/SumoLogic/sumoshell/util"
	"os"
	"time"
)

type Builder func([]string) (util.SumoOperator, error)
type AggBuilder func([]string) (util.SumoAggOperator, error)

var operators = map[string]Builder{
	"parse":  parse.Build,
	"filter": filter.Build,
}

var aggOperators = map[string]AggBuilder{
	"count":   count.Build,
	"average": average.Build,
	"sum":     sum.Build,
}

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Arguments expected")
	} else {
		selectingArg := args[1]
		actualArgs := os.Args[1:]
		nonAggWorked := connectNonAggOperator(selectingArg, actualArgs)
		if nonAggWorked {
			return
		}

		aggWorked := connectAggOperator(selectingArg, actualArgs)
		if aggWorked {
			return
		}

		if selectingArg == "search" {
			search.BuildAndConnect(actualArgs)
			return
		}
		fmt.Println("Operator " + selectingArg + " unrecognized")
	}
}

func connectAggOperator(selector string, args []string) bool {
	aggBuilder, aggOk := aggOperators[selector]
	if !aggOk {
		return false
	}

	aggOperator, err := aggBuilder(args)
	if err != nil {
		fmt.Println(err)
	} else {
		ticker := time.NewTicker(100 * time.Millisecond)
		go flush(aggOperator, ticker)
		util.ConnectToStdIn(aggOperator)
		// Flush when the stream completes to ensure all data is accounted for
		aggOperator.Flush()
	}
	return true
}

func connectNonAggOperator(selector string, args []string) bool {
	builder, ok := operators[selector]
	if ok {
		operator, err := builder(args)
		handleErrorOrWire(operator, err)
	}
	return ok
}

func handleErrorOrWire(operator util.SumoOperator, err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		util.ConnectToStdIn(operator)
	}
}

func flush(aggOp util.SumoAggOperator, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			aggOp.Flush()
		}
	}
}
