package sum

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/group"
	"github.com/SumoLogic/sumoshell/util"
	"strconv"
	"sync"
)

type sum struct {
	sum *float64
	key string
	// DO NOT MODIFY BASE
	base   map[string]interface{}
	output func(map[string]interface{})
	mu *sync.Mutex
}

func makeSum(key string) sum {
	sumV := 0.0
	return sum{&sumV, key, make(map[string]interface{}), util.NewJsonWriter().Write, &sync.Mutex{}}
}

func aggregateSum(output grouper.Merger, key string, base map[string]interface{}) util.SumoAggOperator {
	sumV := 0.0
	return sum{&sumV, key, base, output.Write, &sync.Mutex{}}
}

func Build(args []string) (util.SumoAggOperator, error) {
	// Fist arg is "sum"
	args = args[1:]
	if len(args) == 1 {
		return makeSum(args[0]), nil
	} else if len(args) > 1 {
		key := args[0]
		//_ := relevantArgs[1]
		keyFields := args[2:]
		return grouper.NewAggregate(aggregateSum, keyFields, key), nil
	} else {
		return nil, util.ParseError("Need a argument to average (`sum field`)")
	}
}

func (sumOp sum) Flush() {
	sumOp.mu.Lock()
	defer sumOp.mu.Unlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()
	v, keyInMap := inp[s.key]
	if keyInMap {
		f, keyIsNumber := strconv.ParseFloat(fmt.Sprint(v), 64)
		if keyIsNumber == nil {
			*s.sum += f
		}
	}
}
