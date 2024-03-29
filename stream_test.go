package stream

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

type student struct {
	id     int
	name   string
	age    int
	scores []int
}

func (s *student) String() string {
	return fmt.Sprintf("{id:%d, name:%s, age:%d,scores:%v}", s.id, s.name, s.age, s.scores)
}

func createStudents() []student {
	names := []string{"Tom", "Kate", "Lucy", "Jim", "Jack", "King", "Lee", "Mask"}
	students := make([]student, 10)
	rnd := func(start, end int) int { return rand.Intn(end-start) + start }
	for i := 0; i < 10; i++ {
		students[i] = student{
			id:     i + 1,
			name:   names[rand.Intn(len(names))],
			age:    rnd(15, 26),
			scores: []int{rnd(60, 100), rnd(60, 100), rnd(60, 100)},
		}
	}
	return students
}

type node struct {
	id   int
	next *node
}

func createNodes() *node {
	i := 10
	n := &node{id: i}
	for i > 0 {
		i--
		n = &node{id: i, next: n}
	}
	return n
}

func TestNew(t *testing.T) {
	students := createStudents()
	stream, _ := New(students)

	fmt.Println(t.Name() + ":")
	stream.ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	fmt.Println()
}

func TestNewP(t *testing.T) {
	students := createStudents()
	_, err := New(&students)
	fmt.Println(err)
}

func TestNewErr(t *testing.T) {

	_, err := New(1)
	fmt.Println(err)
}

func TestInts(t *testing.T) {
	ints, _ := Ints(1, 2, 3)
	ints.ForEach(func(i int64) {
		fmt.Print(i)
	})
}

func TestFloats(t *testing.T) {
	floats, _ := Floats(1.1, 2.2, 3.3)
	floats.ForEach(func(f float64) {
		fmt.Print(f)
	})
}

func TestStrings(t *testing.T) {
	strs, _ := Strings("a", "b", "c")
	strs.ForEach(func(s string) {
		fmt.Print(s)
	})
}

func TestIterate(t *testing.T) {
	root := createNodes()

	fmt.Println(t.Name() + ":")
	stream, _ := It(root, func(n *node) (*node, bool) {
		return n.next, n.next.next != nil
	})
	stream.ForEach(func(n *node) {
		fmt.Printf("\tnode{id:%d}\n", n.id)
	})
	fmt.Println()
}

func TestGenerate(t *testing.T) {
	fmt.Println(t.Name() + ":")
	stream, _ := Gen(func() (int, bool) {
		x := rand.Intn(10)
		return x, x < 8
	})
	stream.ForEach(func(x int) {
		fmt.Printf("\t%d\n", x)
	})
	fmt.Println()
}

func TestGenerateN(t *testing.T) {
	fmt.Println(t.Name() + ":")
	stream, _ := GenN(3, func(idx int) int {
		return idx * 2
	})
	stream.ForEach(func(x int) {
		fmt.Printf("\t%d\n", x)
	})
	fmt.Println()
}

func TestFilter(t *testing.T) {
	fmt.Println(t.Name() + ": by age > 20")

	students := createStudents()
	stream, _ := New(students)

	stream.Filter(func(s student) bool {
		return s.age > 20
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	fmt.Printf("\n")
}

func TestFilterIndex(t *testing.T) {
	fmt.Println(t.Name() + ": by age > 20")

	students := createStudents()
	stream, _ := New(students)

	stream.FilterIndex(func(s student, i int) bool {
		return s.age > 20 && i%2 == 0
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	fmt.Printf("\n")
}

func TestMap(t *testing.T) {
	fmt.Println(t.Name() + ": by name")
	students := createStudents()
	stream, _ := New(students)

	stream.Map(func(s student) string {
		return s.name
	}).ForEach(func(s string) {
		fmt.Printf("\t%s\n", s)
	})
	fmt.Println()
}

func TestMapIndex(t *testing.T) {
	fmt.Println(t.Name() + ": by name")
	students := createStudents()
	stream, _ := New(students)

	stream.MapIndex(func(s student, i int) string {
		return strconv.Itoa(i) + ":" + s.name
	}).ForEach(func(s string) {
		fmt.Printf("\t%s\n", s)
	})
	fmt.Println()
}

func TestFlatMap(t *testing.T) {
	fmt.Println(t.Name() + ": by scores")
	students := createStudents()
	stream, _ := New(students)
	var data []int
	stream.FlatMap(func(s student) []int {
		return s.scores
	}).ToSlice(&data)
	fmt.Printf("\t%v\n", data)
}

func TestFlatMapIndex(t *testing.T) {
	fmt.Println(t.Name() + ": by scores")
	students := createStudents()
	stream, _ := New(students)
	var data []int
	stream.FlatMapIndex(func(s student, idx int) []int {
		return s.scores
	}).ToSlice(&data)
	fmt.Printf("\t%v\n", data)
}

func TestSort(t *testing.T) {
	fmt.Println(t.Name() + ": by scores desc")
	students := createStudents()
	stream, _ := New(students)

	stream.Sort(func(s1, s2 student) bool {
		return s1.scores[0]+s1.scores[1]+s1.scores[2] > s2.scores[0]+s2.scores[1]+s2.scores[2]
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	fmt.Println()
}

func TestDistinct(t *testing.T) {
	fmt.Println(t.Name() + ": by name")
	students := createStudents()
	stream, _ := New(students)

	stream.Map(func(s student) string {
		return s.name
	}).Distinct(func(p1, p2 string) bool {
		return p1 == p2
	}).ForEach(func(s string) {
		fmt.Printf("\t%s\n", s)
	})
	fmt.Println()
}

func TestForEach(t *testing.T) {
	fmt.Println(t.Name() + ": by name")
	students := createStudents()
	stream, _ := New(students)

	stream.ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	fmt.Println()
}

func TestForEachIndex(t *testing.T) {
	fmt.Println(t.Name() + ": by name")
	students := createStudents()
	stream, _ := New(students)

	stream.ForEachIndex(func(s student, i int) {
		fmt.Printf("\t%d : %s\n", i, s.String())
	})
	fmt.Println()
}

func TestMatch(t *testing.T) {
	fmt.Println(t.Name() + ":")
	students := createStudents()
	stream, _ := New(students)

	r1 := stream.AllMatch(func(s student) bool {
		return s.age > 20
	})
	stream.Reset()
	r11 := stream.AllMatchIndex(func(s student, i int) bool {
		return s.age > 20 && i%2 == 0
	})
	stream.Reset()
	r2 := stream.AnyMatch(func(s student) bool {
		return s.name == "Jim"
	})
	r22 := stream.AnyMatchIndex(func(s student, i int) bool {
		return s.name == "Jim" && i%2 == 0
	})
	stream.Reset()
	r3 := stream.NoneMatch(func(s student) bool {
		return s.scores[0]+s.scores[1]+s.scores[2] > 270
	})
	r33 := stream.NoneMatchIndex(func(s student, i int) bool {
		return s.scores[0]+s.scores[1]+s.scores[2] > 270 && i%2 == 0
	})
	fmt.Printf("\tAllMatch: %t %t, AnyMatch: %t %t, NoneMatch: %t %t \n", r1, r11, r2, r22, r3, r33)
}

func TestCount(t *testing.T) {
	fmt.Println(t.Name() + ":")
	students := createStudents()
	stream, _ := New(students)

	r := stream.Count()
	fmt.Printf("\t%d\n", r)
}

func TestMaxMin(t *testing.T) {
	fmt.Println(t.Name() + ": by scores")
	students := createStudents()
	stream, _ := New(students)

	r1 := stream.Max(func(s1, s2 student) bool {
		return s1.scores[0]+s1.scores[1]+s1.scores[2] < s2.scores[0]+s2.scores[1]+s2.scores[2]
	})
	stream.Reset()
	r2 := stream.Min(func(s1, s2 student) bool {
		return s1.scores[0]+s1.scores[1]+s1.scores[2] < s2.scores[0]+s2.scores[1]+s2.scores[2]
	})
	fmt.Printf("\tMax: %v, Min: %v \n", r1, r2)
}

func TestPeek(t *testing.T) {
	fmt.Println(t.Name() + ":")
	students := createStudents()
	stream, _ := New(students)

	stream.Filter(func(s student) bool {
		return s.age%2 == 0
	}).Call(func() {
		fmt.Println("\tfilter by age % 2 == 0")
	}).Peek(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	}).Filter(func(s student) bool {
		return s.age > 18
	}).Call(func() {
		fmt.Println("\tfilter by age > 18")
	}).Peek(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	}).Exec()
}

func TestPeekIndex(t *testing.T) {
	fmt.Println(t.Name() + ":")
	students := createStudents()
	stream, _ := New(students)

	stream.FilterIndex(func(s student, i int) bool {
		return s.age%2 == 0
	}).Call(func() {
		fmt.Println("\tfilter by age % 2 == 0")
	}).PeekIndex(func(s student, i int) {
		fmt.Printf("\t%s\n", s.String())
	}).FilterIndex(func(s student, i int) bool {
		return s.age > 18
	}).Call(func() {
		fmt.Println("\tfilter by age > 18")
	}).PeekIndex(func(s student, i int) {
		fmt.Printf("\t%s\n", s.String())
	}).Exec()
}

func TestCheck(t *testing.T) {
	fmt.Println(t.Name() + ":")
	students := createStudents()
	stream, _ := New(students)
	stream.Filter(func(s student) bool {
		return s.age%2 == 0
	}).Check(func(sts []interface{}) bool {
		return len(sts) > 2
	}).ForEach(func(s student) {
		fmt.Println(s.String())
	})
}

func TestLimitSkip(t *testing.T) {
	fmt.Println(t.Name() + ":")
	students := createStudents()
	stream, _ := New(students)

	stream.Limit(5).Call(func() {
		fmt.Println("\tlimit by 5")
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	stream.Reset()
	stream.Skip(5).Call(func() {
		fmt.Println("\tskip by 5")
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	fmt.Println()
}

func TestLimitSkipNeg(t *testing.T) {
	fmt.Println(t.Name() + ":")
	students := createStudents()
	stream, _ := New(students)

	stream.Limit(-1).Call(func() {
		fmt.Println("\tlimit by 5")
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	stream.Reset()
	stream.Skip(-1).Call(func() {
		fmt.Println("\tskip by 5")
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})
	fmt.Println()
}

func TestReduce(t *testing.T) {
	fmt.Println(t.Name() + ": sum of scores[0]")
	students := createStudents()
	stream, _ := New(students)

	r := 0
	r = stream.Map(func(s student) int {
		return s.scores[0]
	}).Reduce(r, func(sum int, i int) int {
		return sum + i
	}).(int)
	fmt.Printf("\t%d\n", r)
}

func TestReduceIndex(t *testing.T) {
	fmt.Println(t.Name() + ": sum of scores[0]")
	students := createStudents()
	stream, _ := New(students)

	r := 0
	r = stream.Map(func(s student) int {
		return s.scores[0]
	}).ReduceIndex(r, func(sum int, o int, i int) int {
		return sum + o
	}).(int)
	fmt.Printf("\t%d\n", r)
}

func TestOf(t *testing.T) {
	fmt.Print(t.Name() + ":  ")
	stream, _ := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)

	stream.ForEach(func(i int) {
		fmt.Printf("%d ", i)
	})
	fmt.Println()
}

func TestToSlice(t *testing.T) {
	fmt.Print(t.Name() + ":  ")
	stream, _ := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)

	slice := make([]int, 0)
	stream.ToSlice(&slice)
	fmt.Println(slice)
	fmt.Println()
}

func TestPointer(t *testing.T) {
	fmt.Print(t.Name() + ":  ")
	students := createStudents()
	studentPs := make([]*student, len(students))
	for i, s := range students {
		studentPs[i] = &s
	}
	r := 0
	stream, _ := New(studentPs)
	r = stream.Filter(func(s *student) bool {
		return s.age > 20
	}).FlatMap(func(s *student) []*int {
		intPs := make([]*int, len(s.scores))
		for i, n := range s.scores {
			intPs[i] = &n
		}
		return intPs
	}).Reduce(r, func(sum int, i *int) int {
		return sum + *i
	}).(int)
	fmt.Println(r)
}

func TestGroup(t *testing.T) {
	fmt.Print(t.Name() + ":  ")
	students := createStudents()
	stream, _ := New(students)

	group := stream.Group(func(s student) (int, student) {
		return s.age, s
	})
	fmt.Println(group)
}

func TestGroupIndex(t *testing.T) {
	fmt.Print(t.Name() + ":  ")
	students := createStudents()
	stream, _ := New(students)

	group := stream.GroupIndex(func(s student, i int) (int, student) {
		return i, s
	})
	fmt.Println(group)
}

func TestFirstLast(t *testing.T) {
	fmt.Print(t.Name() + ":  ")
	students := createStudents()
	stream, _ := New(students)

	first := stream.First(func(s student) bool {
		return s.age > 18
	})
	fmt.Println(first)

	stream, _ = New(students)
	last := stream.Last(func(s student) bool {
		return s.age > 18
	})
	fmt.Println(last)
}

func TestValidateFunc(t *testing.T) {
	fmt.Println(t.Name() + ":  ")
	fn1 := func() {}
	err1 := validateFunc(reflect.ValueOf(fn1), []reflect.Type{}, []reflect.Type{})
	fmt.Println(fmt.Sprintf("validate 'func() {}' by in() out(): %t", err1 == nil))

	fn2 := func(i int) {}
	err2 := validateFunc(reflect.ValueOf(fn2), []reflect.Type{reflect.TypeOf(0)}, []reflect.Type{})
	fmt.Println(fmt.Sprintf("validate 'func(int) {}' by in(int) out(): %t", err2 == nil))

	fn3 := func() int { return 0 }
	err3 := validateFunc(reflect.ValueOf(fn3), []reflect.Type{}, []reflect.Type{reflect.TypeOf(0)})
	fmt.Println(fmt.Sprintf("validate 'func() int {}' by in() out(int): %t", err3 == nil))

	fn4 := func(i int) string { return strconv.Itoa(i) }
	err4 := validateFunc(reflect.ValueOf(fn4), []reflect.Type{reflect.TypeOf(0)}, []reflect.Type{reflect.TypeOf("")})
	fmt.Println(fmt.Sprintf("validate 'func(int) string {}' by in(int) out(string): %t", err4 == nil))
}
