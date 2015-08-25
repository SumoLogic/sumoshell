package render

import (
	"fmt"
	"github.com/SumoLogic/sumoshell/util"
	"sort"
	"strings"
)

func uiColumn(column string) bool {
	if column == util.Type {
		return false
	}

    if column == "_raw" {
        return false
    }
	return true
}

func Columns(inp []map[string]interface{}) map[string]int {
	res := make(map[string]int)
	for _, m := range inp {
		for k, v := range m {
			if uiColumn(k) {
				var effectiveLength int
				if length(k) > length(v) {
					effectiveLength = length(k)
				} else {
					effectiveLength = length(v)
				}
				if res[k] < effectiveLength {
					res[k] = effectiveLength
				}
			}
		}
	}
	return res
}

func ColumnNames(columns map[string]int) []string {
	viewNames := []string{}
	for k, _ := range columns {
		viewNames = append(viewNames, k)
	}
	sort.Strings(viewNames)
	return viewNames
}

func NumericColumn(columns []string) string {
	for _, c := range columns {
		if c == "_avg" || c == "_sum" || c == "_count" {
			return c
		}
	}
	return ""
}

func LabelExtractor(columns []string) func(map[string]interface{}) string {
	labelCols := []string{}
	for _, c := range columns {
		if !strings.HasPrefix(c, "_") {
			labelCols = append(labelCols, c)
		}
	}
	return func(inp map[string]interface{}) string {
		res := []string{}
		for _, c := range labelCols {
			res = append(res, fmt.Sprint(inp[c]))
		}
		return strings.Join(res, "-")
	}
}

func length(inp interface{}) int {
	return len(fmt.Sprint(inp)) + 3
}

type RenderState struct {
	Messages *[]map[string]interface{}
	Meta     *map[string]interface{}
	Flush    func() error
}

func NewConnectedRenderState(flush func() error) *RenderState {
	messages := make([]map[string]interface{}, 0)
	meta := make(map[string]interface{})
	state := &RenderState{&messages, &meta, flush}
	go util.ConnectToStdIn(state)
	return state
}

func (state RenderState) Process(inp map[string]interface{}) {
	if util.IsStartRelation(inp) {
		newmap := make([]map[string]interface{}, 0)
		*state.Messages = newmap
	} else if util.IsEndRelation(inp) {
		state.Flush()
	} else if util.IsRelation(inp) {
		*state.Messages = append(*state.Messages, inp)
	} else if util.IsMeta(inp) {
		*state.Meta = inp
	} else if util.IsPlus(inp) {
		*state.Messages = append(*state.Messages, inp)
		state.Flush()
	} else {
		// :-( no type information
	}
}
