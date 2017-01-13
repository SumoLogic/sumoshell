package average

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/group"
	"github.com/SumoLogic/sumoshell/util"
	"strconv"
	"sync"
)

const Avg = "_avg"

type average struct {
	samples *int
	sum     *float64
	key     string
	// DO NOT MODIFY BASE
	base   map[string]interface{}
	mu     *sync.Mutex
	output func(map[string]interface{})
}

func makeAverage(key string) average {
	samp := 0
	v := 0.0
	mu := sync.Mutex{}
	return average{&samp, &v, key, make(map[string]interface{}), &mu, util.NewJsonWriter().Write}
}

func aggregateAverage(output grouper.Merger, key string, base map[string]interface{}) util.SumoAggOperator {
	samp := 0
	v := 0.0
	mu := sync.Mutex{}
	return average{&samp, &v, key, base, &mu, output.Write}
}

func Build(args []string) (util.SumoAggOperator, error) {
	args = args[1:]

	if len(args) == 1 {
		return makeAverage(args[0]), nil
	} else if len(args) > 1 {
		key := args[0]
		//_ := args[1]
		keyFields := args[2:]
		return grouper.NewAggregate(aggregateAverage, keyFields, key, Avg), nil
	} else {
		return nil, util.ParseError("Need a argument to average (`avg keyname`)")
	}
}

func (avg average) Flush() {
	avg.mu.Lock()
	defer avg.mu.Unlock()
	if *avg.samples > 0 {
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
	ret[Avg] = *a.sum / float64(*a.samples)
	return ret
}

func (a average) Process(inp map[string]interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()
	v, keyInMap := inp[a.key]
	if keyInMap {
		f, keyIsNumber := strconv.ParseFloat(fmt.Sprint(v), 64)
		if keyIsNumber == nil {
			*a.samples += 1
			*a.sum += f
		}
	}
}
