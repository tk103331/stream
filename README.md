# stream
A Go language implementation of the Java Stream API.

----------

**Preparation**

    type student struct {
    	id int
    	name   string
    	ageint
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
    			id: i + 1,
    			name:   names[rand.Intn(len(names))],
    			age:rnd(15, 26),
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


----------

### ForEach ###



### Iterate ###
It create a stream from a iterator.itFunc: func(prev T) (next T,more bool)

    func It(initValue interface{}, itFunc interface{}) (*stream, error)

Sample:

	stream, _ := It(root, func(n *node) (*node, bool) {
		return n.next, n.next.next != nil
	})
	stream.ForEach(func(n *node) {
		fmt.Printf("\tnode{id:%d}\n", n.id)
	})

Output:

    node{id:1}
    node{id:2}
    node{id:3}
    node{id:4}
    node{id:5}
    node{id:6}
    node{id:7}
    node{id:8}
    node{id:9}
    node{id:10}

### Generate ###
Gen create a stream by invoke genFunc. genFunc: func() (next T,more bool)

    func Gen(genFunc interface{}) (*stream, error)

Sapmle:

	stream, _ := Gen(func() (int, bool) {
		x := rand.Intn(10)
		return x, x < 8
	})
	stream.ForEach(func(x int) {
		fmt.Printf("\t%d\n", x)
	})

Output:


	1
	7
	7
	9




### Filter ###
Filter operation. filterFunc: func(o T) bool

    func (s *stream) Filter(filterFunc interface{}) *stream

Sample:

	students := createStudents()
	stream, _ := New(students)

	stream.Filter(func(s student) bool {
		return s.age > 20
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})

Output:

	{id:2, name:Lee, age:22,scores:[80 76 80]}
	{id:4, name:Lucy, age:22,scores:[65 97 86]}
	{id:7, name:King, age:22,scores:[87 91 89]}
	{id:9, name:King, age:21,scores:[94 63 93]}

### Map ###
Map operation. Map one to one.mapFunc: func(o T1) T2

    func (s *stream) Map(mapFunc interface{}) *stream

Sample:

	students := createStudents()
	stream, _ := New(students)

	stream.Map(func(s student) string {
		return s.name
	}).ForEach(func(s string) {
		fmt.Printf("\t%s\n", s)
	})

Output:

	Kate
	Lee
	Lee
	Lucy
	Mask
	Jim
	King
	Jack
	King
	Jim

### FlatMap ###

FlatMap operation. Map one to many.mapFunc: func(o T1) []T2

    func (s *stream) FlatMap(mapFunc interface{}) *stream

Sample:

	students := createStudents()
	stream, _ := New(students)
	var data []int
	stream.FlatMap(func(s student) []int {
		return s.scores
	}).ToSlice(&data)
	fmt.Printf("\t%v\n", data)

Output:

    [67 79 61 80 76 80 62 69 68 65 97 86 68 78 67 68 90 75 87 91 89 91 65 86 94 63 93 64 99 93]

### Sort ###
Sort operation. lessFunc: func(o1,o2 T) bool

    func (s *stream) Sort(lessFunc interface{}) *stream

Sample:

	students := createStudents()
	stream, _ := New(students)

	stream.Sort(func(s1, s2 student) bool {
		return s1.scores[0]+s1.scores[1]+s1.scores[2] > s2.scores[0]+s2.scores[1]+s2.scores[2]
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})

Output:


	{id:7, name:King, age:22,scores:[87 91 89]}
	{id:10, name:Jim, age:20,scores:[64 99 93]}
	{id:9, name:King, age:21,scores:[94 63 93]}
	{id:4, name:Lucy, age:22,scores:[65 97 86]}
	{id:8, name:Jack, age:16,scores:[91 65 86]}
	{id:2, name:Lee, age:22,scores:[80 76 80]}
	{id:6, name:Jim, age:20,scores:[68 90 75]}
	{id:5, name:Mask, age:15,scores:[68 78 67]}
	{id:1, name:Kate, age:16,scores:[67 79 61]}
	{id:3, name:Lee, age:15,scores:[62 69 68]}

### Distinct ###
Distinct operation. equalFunc: func(o1,o2 T) bool

    func (s *stream) Distinct(equalFunc interface{}) *stream

Sample:

	students := createStudents()
	stream, _ := New(students)

	stream.Map(func(s student) string {
		return s.name
	}).Distinct(func(p1, p2 string) bool {
		return p1 == p2
	}).ForEach(func(s string) {
		fmt.Printf("\t%s\n", s)
	})

### Peek ###
Peek operation. peekFunc: func(o T)

    func (s *stream) Peek(peekFunc interface{}) *stream

Sample:


