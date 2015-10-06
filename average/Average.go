package average

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/group"
	"github.com/SumoLogic/sumoshell/util"
	"strconv"
)

type average struct {
	samples *int
	sum     *float64
	key     string
	// DO NOT MODIFY BASE
	base   map[string]interface{}
	ready  *bool
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
	return average{&samp, &v, key, base, &ready, output.Write}
}

func Build(args []string) (util.SumoAggOperator, error) {
	if len(args) == 1 {
		return makeAverage(args[0]), nil
	} else if len(args) > 1 {
		key := args[0]
		//_ := args[1]
		keyFields := args[2:]
		return grouper.NewAggregate(aggregateAverage, keyFields, key), nil
	} else {
		return nil, util.ParseError("Need a argument to average (`avg keyname`)")
	}
}

func (avg average) Flush() {
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
