package stream

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

var StrictMode bool

type Stream struct {
	ops  []op
	data []interface{}
	res  reflect.Type
}

type op struct {
	typ string
	fun reflect.Value
	idx bool
}

type FuncSorter struct {
	data []interface{}
	fun  reflect.Value
}

func (s *FuncSorter) Len() int           { return len(s.data) }
func (s *FuncSorter) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s *FuncSorter) Less(i, j int) bool { return call(s.fun, s.data[i], s.data[j])[0].Bool() }

// New create a stream from a slice
func New(arr interface{}) (*Stream, error) {
	ops := make([]op, 0)
	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() == reflect.Ptr {
		arrValue = arrValue.Elem()
	}
	if arrValue.Kind() == reflect.Slice || arrValue.Kind() == reflect.Array {
		for i := 0; i < arrValue.Len(); i++ {
			dataValue.Set(reflect.Append(dataValue, arrValue.Index(i)))
		}
	} else {
		return nil, errors.New("the type of arr parameter must be Array or Slice")
	}

	return &Stream{ops: ops, data: data, res: arrValue.Type().Elem()}, nil
}

// Of create a stream from some values
func Of(args ...interface{}) (*Stream, error) {
	return New(args)
}

// Ints create a stream from some int64 values.
func Ints(args ...int64) (*Stream, error) {
	return New(args)
}

// Floats create a stream from some float64 values.
func Floats(args ...float64) (*Stream, error) {
	return New(args)
}

// Strings create a stream from some string values.
func Strings(args ...string) (*Stream, error) {
	return New(args)
}

// It create a stream from an iterator. itFunc: func(prev T) (next T,more bool).
func It(initValue interface{}, itFunc interface{}) (*Stream, error) {
	funcValue := reflect.ValueOf(itFunc)

	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()
	prev := reflect.ValueOf(initValue)
	for {
		out := funcValue.Call([]reflect.Value{prev})
		dataValue.Set(reflect.Append(dataValue, out[0]))
		if !out[1].Bool() {
			break
		}
		prev = out[0]
	}
	return New(data)
}

// Gen create a stream by invoke genFunc. genFunc: func() (next T,more bool)
func Gen(genFunc interface{}) (*Stream, error) {
	funcValue := reflect.ValueOf(genFunc)
	if StrictMode {
		err := validateFunc(funcValue, []reflect.Type{}, []reflect.Type{})
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s, must be like func(prev T) (next T,more bool)", err.Error()))
		}
	}
	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()
	for {
		out := call(funcValue)
		dataValue.Set(reflect.Append(dataValue, out[0]))
		if !out[1].Bool() {
			break
		}
	}
	return New(data)
}

// GenN create a stream by invoke genFunc N times. genFunc: func() (ele T)
func GenN(num int, genFunc interface{}) (*Stream, error) {
	if num < 0 {
		return nil, errors.New("num is negative")
	}
	funcValue := reflect.ValueOf(genFunc)
	data := make([]interface{}, num)
	for i := 0; i < num; i++ {
		out := call(funcValue, i)
		data[i] = out[0].Interface()
	}
	return New(data)
}

func (s *Stream) Reset() *Stream {
	s.ops = make([]op, 0)
	return s
}

// Filter operation. filterFunc: func(o T) bool
func (s *Stream) Filter(filterFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(filterFunc)
	s.ops = append(s.ops, op{typ: "filter", fun: funcValue})
	return s
}

// FilterIndex operation with index. filterFunc: func(o T, i int) bool
func (s *Stream) FilterIndex(filterFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(filterFunc)
	s.ops = append(s.ops, op{typ: "filter", fun: funcValue, idx: true})
	return s
}

// Map operation. Map one to one
// mapFunc: func(o T1) T2
func (s *Stream) Map(mapFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(mapFunc)
	s.ops = append(s.ops, op{typ: "map", fun: funcValue})
	return s
}

// MapIndex operation with index. Map one to one
// mapFunc: func(o T1, i int) T2
func (s *Stream) MapIndex(mapFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(mapFunc)
	s.ops = append(s.ops, op{typ: "map", fun: funcValue, idx: true})
	return s
}

// FlatMap operation. Map one to many
// mapFunc: func(o T1) []T2
func (s *Stream) FlatMap(mapFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(mapFunc)
	s.ops = append(s.ops, op{typ: "flatMap", fun: funcValue})
	return s
}

// FlatMapIndex operation with index. Map one to many
// mapFunc: func(o T1) []T2
func (s *Stream) FlatMapIndex(mapFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(mapFunc)
	s.ops = append(s.ops, op{typ: "flatMap", fun: funcValue, idx: true})
	return s
}

// Sort operation. lessFunc: func(o1,o2 T) bool
func (s *Stream) Sort(lessFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(lessFunc)
	s.ops = append(s.ops, op{typ: "sort", fun: funcValue})
	return s
}

// Distinct operation. equalFunc: func(o1,o2 T) bool
func (s *Stream) Distinct(equalFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(equalFunc)
	s.ops = append(s.ops, op{typ: "distinct", fun: funcValue})
	return s
}

// Peek operation. peekFunc: func(o T)
func (s *Stream) Peek(peekFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(peekFunc)
	s.ops = append(s.ops, op{typ: "peek", fun: funcValue})
	return s
}

// PeekIndex operation with index. peekFunc: func(o T)
func (s *Stream) PeekIndex(peekFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(peekFunc)
	s.ops = append(s.ops, op{typ: "peek", fun: funcValue, idx: true})
	return s
}

// Call operation. Call function with the data.
// callFunc: func()
func (s *Stream) Call(callFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(callFunc)
	s.ops = append(s.ops, op{typ: "call", fun: funcValue})
	return s
}

// Check operation. Check if should be continue process data.
// checkFunc: func(o []T) bool ,checkFunc must return if should be continue process data.
func (s *Stream) Check(checkFunc interface{}) *Stream {
	funcValue := reflect.ValueOf(checkFunc)
	s.ops = append(s.ops, op{typ: "check", fun: funcValue})
	return s
}

// Limit operation.
func (s *Stream) Limit(num int) *Stream {
	if num < 0 {
		num = 0
	}
	funcValue := reflect.ValueOf(func() int { return num })
	s.ops = append(s.ops, op{typ: "limit", fun: funcValue})
	return s
}

// Skip operation.
func (s *Stream) Skip(num int) *Stream {
	if num < 0 {
		num = 0
	}
	funcValue := reflect.ValueOf(func() int { return num })
	s.ops = append(s.ops, op{typ: "skip", fun: funcValue})
	return s
}

// collect operation.
func (s *Stream) collect() []interface{} {
	result := s.data
	for _, op := range s.ops {
		if len(result) == 0 {
			break
		}
		switch op.typ {
		case "filter":
			result = doFilter(result, op)
		case "peek":
			each(result, op.fun, emptyeachfunc, op.idx)
		case "map":
			result = doMap(result, op)
		case "flatMap":
			result = doFlatMap(result, op)
		case "aggMap":

		case "sort":
			sort.Sort(&FuncSorter{data: result, fun: op.fun})
		case "distinct":
			result = doDistinct(result, op)
		case "limit":
			result = doLimit(op, result)
		case "skip":
			result = doSkip(op, result)
		case "call":
			call(op.fun)
		case "check":
			out := call(op.fun, result)
			if !out[0].Bool() {
				break
			}
		}
	}
	return result
}

func doSkip(op op, result []interface{}) []interface{} {
	skip := int(call(op.fun)[0].Int())
	if skip > len(result) {
		skip = len(result)
	}
	temp := result
	return temp[skip:]
}

func doLimit(op op, result []interface{}) []interface{} {
	limit := int(call(op.fun)[0].Int())
	if limit > len(result) {
		limit = len(result)
	}
	temp := result
	return temp[:limit]
}

func doDistinct(result []interface{}, op op) []interface{} {
	temp := make([]interface{}, 0)
	temp = append(temp, result[0])
	for _, it := range result {
		found := false
		for _, it2 := range temp {
			out := call(op.fun, it, it2)
			if out[0].Bool() {
				found = true
			}
		}
		if !found {
			temp = append(temp, it)
		}
	}
	return temp
}

func doFlatMap(result []interface{}, op op) []interface{} {
	temp := make([]interface{}, 0)
	tempVlaue := reflect.ValueOf(&temp).Elem()
	each(result, op.fun, func(i int, it interface{}, out []reflect.Value) bool {
		for i := 0; i < out[0].Len(); i++ {
			tempVlaue.Set(reflect.Append(tempVlaue, out[0].Index(i)))
		}
		return true
	}, op.idx)
	return temp
}

func doMap(result []interface{}, op op) []interface{} {
	temp := make([]interface{}, 0)
	tempValue := reflect.ValueOf(&temp).Elem()
	each(result, op.fun, func(i int, it interface{}, out []reflect.Value) bool {
		tempValue.Set(reflect.Append(tempValue, out[0]))
		return true
	}, op.idx)
	return temp
}

func doFilter(result []interface{}, op op) []interface{} {
	temp := make([]interface{}, 0)
	each(result, op.fun, func(i int, it interface{}, out []reflect.Value) bool {
		if out[0].Bool() {
			temp = append(temp, it)
		}
		return true
	}, op.idx)
	return temp
}

// Exec operation.
func (s *Stream) Exec() {
	s.collect()
}

// ToSlice operation. targetSlice must be a pointer.
func (s *Stream) ToSlice(targetSlice interface{}) error {
	data := s.collect()
	targetValue := reflect.ValueOf(targetSlice)
	if targetValue.Kind() != reflect.Ptr {
		return errors.New("target slice must be a pointer")
	}
	sliceValue := reflect.Indirect(targetValue)
	for _, it := range data {
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(it)))
	}
	return nil
}

// ForEach executes a provided function once for each array element,and terminate the stream.
// actFunc: func(o T)
func (s *Stream) ForEach(actFunc interface{}) {
	data := s.collect()
	each(data, reflect.ValueOf(actFunc), emptyeachfunc, false)
}

// ForEachIndex executes a provided function once for each array element,and terminate the stream.
// actFunc: func(o T, i int)
func (s *Stream) ForEachIndex(actFunc interface{}) {
	data := s.collect()
	each(data, reflect.ValueOf(actFunc), emptyeachfunc, true)
}

func (s *Stream) all(matchFunc interface{}, idx bool) bool {
	data := s.collect()
	allMatch := true
	each(data, reflect.ValueOf(matchFunc), func(i int, it interface{}, out []reflect.Value) bool {
		if !out[0].Bool() {
			allMatch = false
			return false
		}
		return true
	}, idx)
	return allMatch
}

// AllMatch operation.
// matchFunc: func(o T) bool
func (s *Stream) AllMatch(matchFunc interface{}) bool {
	return s.all(matchFunc, false)
}

// AllMatchIndex operation with index.
// matchFunc: func(o T, i int) bool
func (s *Stream) AllMatchIndex(matchFunc interface{}) bool {
	return s.all(matchFunc, true)
}

func (s *Stream) any(matchFunc interface{}, idx bool) bool {
	data := s.collect()
	anyMatch := false
	each(data, reflect.ValueOf(matchFunc), func(i int, it interface{}, out []reflect.Value) bool {
		if out[0].Bool() {
			anyMatch = true
			return false
		}
		return true
	}, idx)
	return anyMatch
}

// AnyMatch operation. matchFunc: func(o T) bool
func (s *Stream) AnyMatch(matchFunc interface{}) bool {
	return s.any(matchFunc, false)
}

// AnyMatchIndex operation with index. matchFunc: func(o T, i int) bool
func (s *Stream) AnyMatchIndex(matchFunc interface{}) bool {
	return s.any(matchFunc, true)
}

func (s *Stream) none(matchFunc interface{}, idx bool) bool {
	data := s.collect()
	noneMatch := true
	each(data, reflect.ValueOf(matchFunc), func(i int, it interface{}, out []reflect.Value) bool {
		if out[0].Bool() {
			noneMatch = false
			return false
		}
		return true
	}, idx)
	return noneMatch
}

// NoneMatch operation. matchFunc: func(o T) bool
func (s *Stream) NoneMatch(matchFunc interface{}) bool {
	return s.none(matchFunc, false)
}

// NoneMatchIndex operation with index. matchFunc: func(o T) bool
func (s *Stream) NoneMatchIndex(matchFunc interface{}) bool {
	return s.none(matchFunc, true)
}

// Count operation.Return the count of elements in stream.
func (s *Stream) Count() int {
	return len(s.collect())
}

// Group operation. Group values by key.
// Parameter groupFunc: func(o T1) (key T2,value T3). Return map[T2]T3
func (s *Stream) group(groupFunc interface{}, idx bool) map[interface{}][]interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(groupFunc)
	result := make(map[interface{}][]interface{})
	for i, it := range data {
		var out []reflect.Value
		if idx {
			out = call(funcValue, it, i)
		} else {
			out = call(funcValue, it)
		}

		key := out[0].Interface()
		slice, ok := result[key]
		if !ok {
			slice = make([]interface{}, 1)
		}
		slice = append(slice, out[1].Interface())
		result[key] = slice
	}
	return result
}

// Group operation. Group values by key.
// Parameter groupFunc: func(o T1) (key T2,value T3). Return map[T2][]T3
func (s *Stream) Group(groupFunc interface{}) map[interface{}][]interface{} {
	return s.group(groupFunc, false)
}

// GroupIndex operation with index. Group values by key.
// Parameter groupFunc: func(o T1) (key T2,value T3). Return map[T2][]T3
func (s *Stream) GroupIndex(groupFunc interface{}) map[interface{}][]interface{} {
	return s.group(groupFunc, true)
}

// Max operation.lessFunc: func(o1,o2 T) bool
func (s *Stream) Max(lessFunc interface{}) interface{} {
	funcValue := reflect.ValueOf(lessFunc)
	data := s.collect()
	var max interface{}
	if len(data) > 0 {
		max = data[0]
		for i := 1; i < len(data); i++ {
			out := call(funcValue, max, data[i])
			if out[0].Bool() {
				max = data[i]
			}
		}
	}
	return max
}

// Min operation.lessFunc: func(o1,o2 T) bool
func (s *Stream) Min(lessFunc interface{}) interface{} {
	funcValue := reflect.ValueOf(lessFunc)
	data := s.collect()
	var min interface{}
	if len(data) > 0 {
		min = data[0]
		for i := 1; i < len(data); i++ {
			out := call(funcValue, data[i], min)
			if out[0].Bool() {
				min = data[i]
			}
		}
	}
	return min
}

// First operation. matchFunc: func(o T) bool
func (s *Stream) First(matchFunc interface{}) interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(matchFunc)
	for _, it := range data {
		out := call(funcValue, it)
		if out[0].Bool() {
			return it
		}
	}
	return nil
}

// Last operation. matchFunc: func(o T) bool
func (s *Stream) Last(matchFunc interface{}) interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(matchFunc)
	for i := len(data) - 1; i >= 0; i-- {
		it := data[i]
		out := call(funcValue, it)
		if out[0].Bool() {
			return it
		}
	}
	return nil
}

func (s *Stream) reduce(initValue interface{}, reduceFunc interface{}, idx bool) interface{} {
	data := s.collect()
	funcValue := reflect.ValueOf(reduceFunc)
	result := initValue
	rValue := reflect.ValueOf(&result).Elem()
	for i, it := range data {
		if idx {
			out := call(funcValue, result, it, i)
			rValue.Set(out[0])
		} else {
			out := call(funcValue, result, it)
			rValue.Set(out[0])
		}
	}
	return result
}

// Reduce operation. reduceFunc: func(r T2,o T) T2
func (s *Stream) Reduce(initValue interface{}, reduceFunc interface{}) interface{} {
	return s.reduce(initValue, reduceFunc, false)
}

// ReduceIndex operation with index. reduceFunc: func(r T2,o T,i int) T2
func (s *Stream) ReduceIndex(initValue interface{}, reduceFunc interface{}) interface{} {
	return s.reduce(initValue, reduceFunc, true)
}

// eachfunc is the function for each method,return if should continue loop
type eachfunc func(int, interface{}, []reflect.Value) bool

// emptyeachfunc the empty eachfunc, return true
var emptyeachfunc = func(int, interface{}, []reflect.Value) bool { return true }

func each(data []interface{}, fun reflect.Value, act eachfunc, idx bool) {
	for i, it := range data {
		if idx {
			out := call(fun, it, i)
			if !act(i, it, out) {
				break
			}
		} else {
			out := call(fun, it)
			if !act(i, it, out) {
				break
			}
		}
	}
}

func call(fun reflect.Value, args ...interface{}) []reflect.Value {
	in := make([]reflect.Value, len(args))
	for i, a := range args {
		in[i] = reflect.ValueOf(a).Convert(fun.Type().In(i))
	}
	return fun.Call(in)
}

func validateFunc(fn reflect.Value, in []reflect.Type, out []reflect.Type) error {
	fnType := fn.Type()
	if fn.Kind() != reflect.Func {
		return errors.New("func invalid")
	}
	if fnType.NumIn() != len(in) {
		return errors.New("func in num invalid")
	}
	if fnType.NumOut() != len(out) {
		return errors.New("func out num invalid")
	}
	for i := 0; i < fnType.NumIn(); i++ {
		if fnType.In(i) != in[i] {
			return errors.New("func in type invalid")
		}
	}
	for i := 0; i < fnType.NumOut(); i++ {
		if fnType.Out(i) != out[i] {
			return errors.New("func out type invalid")
		}
	}
	return nil
}
