package main
import (
	"os"
	"github.com/SumoLogic/sumoshell/util"
	"github.com/SumoLogic/sumoshell/group"
	"strconv"
	"fmt"
	"time"
)

type average struct {
	samples *int
	sum *float64
	key string
	// DO NOT MODIFY BASE
	base map[string]interface{}
	ready *bool
	output func(map[string]interface{})
}

func makeAverage(key string) average {
	samp := 0
	v := 0.0
	ready := false
	return average{&samp, &v, key, make(map[string]interface{}), &ready, util.NewJsonWriter().Write}

}

func aggregateAverage(output grouper.Merger, key string, base map[string]interface{}) util.SumoAggOperator {
	samp := 0
	v := 0.0
	ready := false
	ticker := time.NewTicker(100 * time.Millisecond)
	avg := average{&samp, &v, key, base, &ready, output.Write}
	go flush(avg, ticker)
	return avg
}

func main() {
	relevantArgs := os.Args[1:]

	if len(relevantArgs) == 1 {
		ticker := time.NewTicker(1 * time.Second)
		avg := makeAverage(relevantArgs[0])
		go flush(avg, ticker)
		util.ConnectToStdIn(avg)
		avg.Flush()
	} else if len(relevantArgs) > 1 {
		key := relevantArgs[0]
		//_ := relevantArgs[1]
		keyFields := relevantArgs[2:]
		agg := grouper.NewAggregate(aggregateAverage, keyFields, key)
		util.ConnectToStdIn(agg)
		agg.Flush()
	} else {
		fmt.Println("Need a argument to average (`avg keyname`)")
	}

}

func flush(avg average, ticker *time.Ticker) {
	for {
		select {
			case <- ticker.C:
				avg.Flush()
		}
	}
}

func (avg average)Flush() {
	if *avg.ready {
			avg.output(util.CreateStartRelation())
			avg.output(util.CreateRelation(currentState(avg)))
			avg.output(util.CreateEndRelation())
	}
}

func currentState(a average) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, val := range a.base {
		ret[key] = val
	}
	ret["_avg"] = *a.sum / float64(*a.samples)
	return ret
}

func (a average) Process(inp map[string]interface{}) {
	v, keyInMap := inp[a.key]
	if keyInMap {
		f, keyIsNumber := strconv.ParseFloat(fmt.Sprint(v), 64)
		if keyIsNumber == nil {
			*a.samples += 1
			*a.sum += f
			*a.ready = true
		}
	}
}
