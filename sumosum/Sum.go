package main

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/group"
	"github.com/SumoLogic/sumoshell/util"
	"os"
	"strconv"
	"time"
)

type sum struct {
	sum *float64
	key string
	// DO NOT MODIFY BASE
	base   map[string]interface{}
	output func(map[string]interface{})
}

func makeSum(key string) sum {
	sumV := 0.0
	return sum{&sumV, key, make(map[string]interface{}), util.NewJsonWriter().Write}

}

func aggregateSum(output grouper.Merger, key string, base map[string]interface{}) util.SumoAggOperator {
	sumV := 0.0
	ticker := time.NewTicker(100 * time.Millisecond)
	sumOp := sum{&sumV, key, base, output.Write}
	go flush(sumOp, ticker)
	return sumOp
}

func main() {
	relevantArgs := os.Args[1:]

	if len(relevantArgs) == 1 {
		ticker := time.NewTicker(1 * time.Second)
		avg := makeSum(relevantArgs[0])
		go flush(avg, ticker)
		util.ConnectToStdIn(avg)
		avg.Flush()
	} else if len(relevantArgs) > 1 {
		key := relevantArgs[0]
		//_ := relevantArgs[1]
		keyFields := relevantArgs[2:]
		agg := grouper.NewAggregate(aggregateSum, keyFields, key)
		util.ConnectToStdIn(agg)
		agg.Flush()
	} else {
		fmt.Println("Need a argument to average (`sum field`)")
	}

}

func flush(sumOp sum, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			sumOp.Flush()
		}
	}
}

func (sumOp sum) Flush() {
	sumOp.output(util.CreateStartRelation())
	sumOp.output(util.CreateRelation(currentState(sumOp)))
	sumOp.output(util.CreateEndRelation())
}

func currentState(s sum) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, val := range s.base {
		ret[key] = val
	}
	ret["_sum"] = *s.sum
	return ret
}

func (s sum) Process(inp map[string]interface{}) {
	v, keyInMap := inp[s.key]
	if keyInMap {
		f, keyIsNumber := strconv.ParseFloat(fmt.Sprint(v), 64)
		if keyIsNumber == nil {
			*s.sum += f
		}
	}
}
