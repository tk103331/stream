package stream

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type person struct {
	name  string
	age   int
	birth time.Time
}

func (p *person) String() string {
	return fmt.Sprintf("{name:%s, age:%d, birth: %s}", p.name, p.age, strconv.Itoa(p.birth.Year()))
}

func newPersons() []person {
	persons := make([]person, 10)
	for i := 0; i < 10; i++ {
		persons[i] = person{
			name:  "Name-" + strconv.Itoa(i),
			age:   i,
			birth: time.Now().AddDate(-i, 0, 0),
		}
	}
	return persons
}

func TestNew(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	fmt.Printf("persons: %v \n", persons)
	fmt.Printf("stream: %v \n", stream)
}

func TestFilter(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r := stream.Filter(func(p person) bool {
		return p.age > 5
	}).Collect()
	fmt.Printf("Filter result : %#v \n", r)
}

func TestMap(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r := stream.Map(func(p person) string {
		return p.name
	}).Collect()
	fmt.Printf("Map result : %#v \n", r)
}

func TestSort(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r := stream.Sort(func(p1, p2 person) bool {
		return p1.age < p2.age
	}).Collect()
	fmt.Printf("Sort result : %#v \n", r)
}

func TestDistinct(t *testing.T) {
	persons := []person{
		person{name: "Tom"},
		person{name: "King"},
		person{name: "Tom"},
	}
	stream, _ := New(persons)
	r := stream.Distinct(func(p1, p2 person) bool {
		return p1.name == p2.name
	}).Collect()
	fmt.Printf("Distinct result : %#v \n", r)
}

func TestEach(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	fmt.Println("ForEach")
	stream.ForEach(func(p person) {
		fmt.Println(p)
	})
}

func TestMatch(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r1 := stream.AllMatch(func(p person) bool {
		return p.age > 5
	})
	stream.Reset()
	r2 := stream.AnyMatch(func(p person) bool {
		return p.name == "Name-1"
	})
	stream.Reset()
	r3 := stream.NoneMatch(func(p person) bool {
		return p.birth.Year() == 2015
	})
	fmt.Printf("AllMatch: %t, AnyMatch: %t, NoneMatch: %t \n", r1, r2, r3)
}

func TestCount(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r := stream.Filter(func(p person) bool {
		return p.age > 5
	}).Count()
	fmt.Printf("Count: %d \n", r)
}

func TestMaxMin(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r1 := stream.Filter(func(p person) bool {
		return p.age > 5
	}).Max(func(p1, p2 person) bool {
		return p1.age < p2.age
	})
	stream.Reset()
	r2 := stream.Filter(func(p person) bool {
		return p.age < 5
	}).Min(func(p1, p2 person) bool {
		return p1.age < p2.age
	})
	fmt.Printf("Max: %v, Min: %v \n", r1, r2)
}

func TestPeek(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)
	fmt.Println("peek: ")
	stream.Filter(func(p person) bool {
		return p.age%2 == 0
	}).Peek(func(p person) {
		fmt.Println(p)
	}).Filter(func(p person) bool {
		return p.age > 5
	}).Peek(func(p person) {
		fmt.Println(p)
	}).Exec()
}

func TestLimitSkip(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r1 := stream.Limit(5).Collect()
	stream.Reset()
	r2 := stream.Skip(5).Collect()
	fmt.Printf("Limit: %v, skip: %v \n", r1, r2)
}

type sum struct {
	value int
}

func TestReduce(t *testing.T) {
	persons := newPersons()
	stream, _ := New(persons)

	r := 0
	r = stream.Map(func(p person) int {
		return p.age
	}).Reduce(r, func(sum int, i int) int {
		return sum + i
	}).(int)
	fmt.Printf("reduce: %v \n", r)
}

func TestOf(t *testing.T) {
	stream, _ := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)
	r := int(0)
	r = stream.Filter(func(i int) bool {
		return i%2 == 0
	}).Map(func(i int) int {
		return i * 2
	}).Reduce(r, func(sum int, i int) int {
		return sum + i
	}).(int)
	fmt.Printf("test of: %d \n", r)
}

func TestInts(t *testing.T) {
	stream, _ := Ints(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)
	r := int64(0)
	r = stream.Filter(func(i int64) bool {
		return i%2 == 0
	}).Map(func(i int64) int64 {
		return i * 2
	}).Reduce(r, func(sum int64, i int64) int64 {
		return sum + i
	}).(int64)
	fmt.Printf("Ints: %d \n", r)
}

func TestFloats(t *testing.T) {
	stream, _ := Floats(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)
	r := float64(0.0)
	r = stream.Filter(func(i float64) bool {
		return i > 5
	}).Map(func(i float64) float64 {
		return i * 2
	}).Reduce(r, func(sum float64, i float64) float64 {
		return sum + i
	}).(float64)
	fmt.Printf("Floats: %f \n", r)
}
