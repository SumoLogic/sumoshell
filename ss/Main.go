package main

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/average"
	"github.com/SumoLogic/sumoshell/count"
	"github.com/SumoLogic/sumoshell/filter"
	"github.com/SumoLogic/sumoshell/parse"
	"github.com/SumoLogic/sumoshell/sum"
	"github.com/SumoLogic/sumoshell/util"
	"os"
	"time"
)

type Builder func([]string) (util.SumoOperator, error)
type AggBuilder func([]string) (util.SumoAggOperator, error)

func main() {
	operators := map[string]Builder{
		"parse":  parse.Build,
		"filter": filter.Build,
	}

	aggOperators := map[string]AggBuilder{
		"count":   count.Build,
		"average": average.Build,
		"sum":     sum.Build,
	}

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Arguments expected")
	} else {
		selectingArg := args[1]
		builder, ok := operators[selectingArg]
		if !ok {
			aggBuilder, aggOk := aggOperators[selectingArg]
			if !aggOk {
				fmt.Println("Operator " + selectingArg + " not found")
				return
			}
			aggOperator, err := aggBuilder(os.Args[1:])
			if err != nil {
				fmt.Println(err)
			} else {
				ticker := time.NewTicker(100 * time.Millisecond)
				go flush(aggOperator, ticker)
				util.ConnectToStdIn(aggOperator)
				// Flush when the stream completes to ensure all data is accounted for
				aggOperator.Flush()
			}
			return
		}

		operator, err := builder(os.Args[1:])
		if err != nil {
			fmt.Println(err)
		} else {
			util.ConnectToStdIn(operator)
		}
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
