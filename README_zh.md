学习Go语言时实现的集合操作工具库，类似于Java 8 中新增的Stream API。由于Go语言不支持泛型，所以基于反射实现。只用于学习目的，不要用于生产（PS:当然也不会有人用）。

项目地址：https://github.com/tk103331/stream

集合操作包括生成操作、中间操作和终止操作。
生成操作返回值是Steam对象，相当于数据的源头，可以调用Stream的其他方法；中间操作返回值是Stream对象，可以继续调用Stream的方法，即可以链式调用方法；终止操作不能继续调用方法。


下面介绍下这个库的API：

----------

**数据准备**
后面的操作都是基于集合数据的，先准备一些测试数据。

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

### 循环遍历 ForEach ###
循环遍历集合中的每一个元素，需要提供一个包含一个参数的处理函数作为参数，形如  func(o T)，循环遍历时会把每个元素作为处理函数的实参。
ForEach 方法是终止操作。

    func (s *stream) ForEach(actFunc interface{})

例子:

	students := createStudents()
	stream, _ := New(students)

	stream.ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})

输出:

	{id:1, name:Kate, age:16,scores:[67 79 61]}
	{id:2, name:Lee, age:22,scores:[80 76 80]}
	{id:3, name:Lee, age:15,scores:[62 69 68]}
	{id:4, name:Lucy, age:22,scores:[65 97 86]}
	{id:5, name:Mask, age:15,scores:[68 78 67]}
	{id:6, name:Jim, age:20,scores:[68 90 75]}
	{id:7, name:King, age:22,scores:[87 91 89]}
	{id:8, name:Jack, age:16,scores:[91 65 86]}
	{id:9, name:King, age:21,scores:[94 63 93]}
	{id:10, name:Jim, age:20,scores:[64 99 93]}

### 迭代器 Iterate ###
It 方法可以从一个迭代器中创建一个Stream对象，迭代器就是一个迭代产生数据的迭代函数，迭代函数形如 func(prev T) (next T,more bool)，迭代函数的参数为上一个元素的值，返回值是下一个元素的值，和是否还有更多元素。
It 方法是生成操作。

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

### 生成器 Generate ###
Gen 方法可以从一个生成器中创建一个Stream对象，生成器就是一个不断产生数据的生成函数，生成函数形如 func() (next T,more bool)，生成函数没有参数，返回值是下一个元素的值，和是否还有更多元素。
Gen 方法是生成操作。
Gen 方法和It 方法的区别就是，它可以不依赖上一个元素的值。

    func Gen(genFunc interface{}) (*stream, error)

例子:

	stream, _ := Gen(func() (int, bool) {
		x := rand.Intn(10)
		return x, x < 8
	})
	stream.ForEach(func(x int) {
		fmt.Printf("\t%d\n", x)
	})

输出:

	1
	7
	7
	9

### 过滤 Filter ###
Filter 方法对集合中的元素进行过滤，筛选出符合条件的元素，需要提供一个过滤函数，过滤函数形如func(o T) bool，参数为集合中的元素，返回值是表示该元素是否符合条件。
Filter 方法是中间操作。

    func (s *stream) Filter(filterFunc interface{}) *stream

例子:

	students := createStudents()
	stream, _ := New(students)

	stream.Filter(func(s student) bool {
		return s.age > 20
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})

输出:

	{id:2, name:Lee, age:22,scores:[80 76 80]}
	{id:4, name:Lucy, age:22,scores:[65 97 86]}
	{id:7, name:King, age:22,scores:[87 91 89]}
	{id:9, name:King, age:21,scores:[94 63 93]}

### 映射 Map ###
Map 方法可以将集合中的每个元素映射为新的值，从而得到一个新的集合，需要提供一个映射函数，形如func(o T1) T2，参数为集合中的元素，返回值是表示该元素映射的新值。
Map 方法是中间操作。

    func (s *stream) Map(mapFunc interface{}) *stream

例子:

	students := createStudents()
	stream, _ := New(students)

	stream.Map(func(s student) string {
		return s.name
	}).ForEach(func(s string) {
		fmt.Printf("\t%s\n", s)
	})

输出:

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

### 打平映射 FlatMap ###
FlatMap 方法可以将集合中每个元素映射为多个元素，返回新的集合包含映射的所有元素。需要提供一个映射函数，形如 func(o T1) []T2，参数为集合中的元素，返回值是表示该元素映射的新值的集合。
FlatMap 方法是中间操作。
FlatMap 方法和Map 方法的区别在于，它可以将集合中每个元素嵌套的集合打平，合并为新的集合。

    func (s *stream) FlatMap(mapFunc interface{}) *stream

例子:

	students := createStudents()
	stream, _ := New(students)
	var data []int
	stream.FlatMap(func(s student) []int {
		return s.scores
	}).ToSlice(&data)
	fmt.Printf("\t%v\n", data)

输出:

    [67 79 61 80 76 80 62 69 68 65 97 86 68 78 67 68 90 75 87 91 89 91 65 86 94 63 93 64 99 93]

### 排序 Sort ###
Sort 方法根据一定队则对集合中的元素进行排序，参数为比较函数，形如func(o1,o2 T) bool，参数为集合中的两个元素，返回值为第一参数是否小于第二个参数。排序算法使用sort中的排序算法。
Sort 方法是中间操作。

    func (s *stream) Sort(lessFunc interface{}) *stream

例子:

	students := createStudents()
	stream, _ := New(students)

	stream.Sort(func(s1, s2 student) bool {
		return s1.scores[0]+s1.scores[1]+s1.scores[2] > s2.scores[0]+s2.scores[1]+s2.scores[2]
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})

输出:

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

### 去重 Distinct ###
Distinct 方法会对集合中的元素进行比较，并将重复的元素过滤掉. 参数为比较函数，形如 func(o1,o2 T) bool，参数为集合中的两个元素，返回值为两个元素是否相等。
Distinct 方法是中间操作。

    func (s *stream) Distinct(equalFunc interface{}) *stream

例子:

	students := createStudents()
	stream, _ := New(students)

	stream.Map(func(s student) string {
		return s.name
	}).Distinct(func(p1, p2 string) bool {
		return p1 == p2
	}).ForEach(func(s string) {
		fmt.Printf("\t%s\n", s)
	})

输出:

	Kate
	Lee
	Lucy
	Mask
	Jim
	King
	Jack

### 提取 Peek ###
Peek 方法遍历集合的每个元素，执行一定的处理，处理函数形如func(o T)，参数为集合每一个元素，没有返回值。
Peek 方法和 ForEach 方法的区别，它是一个中间操作，可以继续调用Stream的其他方法。

    func (s *stream) Peek(peekFunc interface{}) *stream

例子:

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

输出:

	filter by age % 2 == 0
	{id:1, name:Kate, age:16,scores:[67 79 61]}
	{id:2, name:Lee, age:22,scores:[80 76 80]}
	{id:4, name:Lucy, age:22,scores:[65 97 86]}
	{id:6, name:Jim, age:20,scores:[68 90 75]}
	{id:7, name:King, age:22,scores:[87 91 89]}
	{id:8, name:Jack, age:16,scores:[91 65 86]}
	{id:10, name:Jim, age:20,scores:[64 99 93]}
	filter by age > 18
	{id:2, name:Lee, age:22,scores:[80 76 80]}
	{id:4, name:Lucy, age:22,scores:[65 97 86]}
	{id:6, name:Jim, age:20,scores:[68 90 75]}
	{id:7, name:King, age:22,scores:[87 91 89]}
	{id:10, name:Jim, age:20,scores:[64 99 93]}

### 调用 Call ###
Call 方法可以在Stream对象执行过程中拿到集合的所有数据，可以对中间结果做一些处理，参数为处理函数，形如func(o []T)，参数为整个集合的数据。
Call 方法为中间操作。

    func (s *stream) Call(callFunc interface{}) *stream

### 检查 Check ###
Check 方法可以在Stream对象执行过程中检查是否需要进行后续操作，参数为判断函数，形如func(o []T) bool，参数为整个集合的数据，返回值为是否继续处理数据。
Check 方法为中间操作。
Check 方法与Call 方法的区分是，它可以终止整个Steam的执行。

    func (s *stream) Check(checkFunc interface{}) *stream

### 限制 Limit ###
Limit 方法可以限制集合中元素的数量，参数为显示的数量。
Limit 方法为中间操作。

    func (s *stream) Limit(num int) *stream

例子:
	
	students := createStudents()
	stream, _ := New(students)

	stream.Limit(5).Call(func() {
		fmt.Println("\tlimit by 5")
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})

输出:

	limit by 5
	{id:1, name:Kate, age:16,scores:[67 79 61]}
	{id:2, name:Lee, age:22,scores:[80 76 80]}
	{id:3, name:Lee, age:15,scores:[62 69 68]}
	{id:4, name:Lucy, age:22,scores:[65 97 86]}
	{id:5, name:Mask, age:15,scores:[68 78 67]}

### 跳过 Skip ###
Skip 方法可以在处理过程中跳过指定数目的元素，参数为跳过的数量.

    func (s *stream) Skip(num int) *stream

例子:

	stream.Skip(5).Call(func() {
		fmt.Println("\tskip by 5")
	}).ForEach(func(s student) {
		fmt.Printf("\t%s\n", s.String())
	})

输出:

	skip by 5
	{id:6, name:Jim, age:20,scores:[68 90 75]}
	{id:7, name:King, age:22,scores:[87 91 89]}
	{id:8, name:Jack, age:16,scores:[91 65 86]}
	{id:9, name:King, age:21,scores:[94 63 93]}
	{id:10, name:Jim, age:20,scores:[64 99 93]}

### 全部匹配 AllMatch ###
AllMatch 判断集合中的元素是否都符合条件，需要提供一个判断函数，形如 func(o T) bool ， 参数为集合中的元素，返回值为是否条件。
AllMatch 方法为终止操作，返回值为是否所有都符合条件。

    func (s *stream) AllMatch(matchFunc interface{}) bool

### 任一匹配 AnyMatch ###
AnyMatch 判断集合中的元素是否有任一元素符合条件，需要提供一个判断函数，形如 func(o T) bool ，参数为集合中的元素，返回值为是否条件。
AnyMatch 方法为终止操作，返回值为是否有任一元素符合条件。

    func (s *stream) AnyMatch(matchFunc interface{}) bool

### 全不匹配 NoneMatch ###
NoneMatch 判断集合中的元素是否所有元素都不符合条件，需要提供一个判断函数，形如 func(o T) bool ，参数为集合中的元素，返回值为是否条件。
NoneMatch 方法为终止操作，返回值为是否所有元素都不符合条件。

    func (s *stream) NoneMatch(matchFunc interface{}) bool

例子:

	students := createStudents()
	stream, _ := New(students)

	r1 := stream.AllMatch(func(s student) bool {
		return s.age > 20
	})
	stream.Reset()
	r2 := stream.AnyMatch(func(s student) bool {
		return s.name == "Jim"
	})
	stream.Reset()
	r3 := stream.NoneMatch(func(s student) bool {
		return s.scores[0]+s.scores[1]+s.scores[2] > 270
	})
	fmt.Printf("\tAllMatch: %t, AnyMatch: %t, NoneMatch: %t \n", r1, r2, r3)

输出:

    AllMatch: false, AnyMatch: true, NoneMatch: true

### 计数 Count ###
Count 返回集合中元素的数量。
Count 为终止操作。

    func (s *stream) Count() int

例子:

	students := createStudents()
	stream, _ := New(students)

	r := stream.Count()
	fmt.Printf("\t%d\n", r)

输出:

    10

### 分组 Group ###
Group 方法可以根据规则，将集合中的元素进行分组，需要提供一个分组函数，形如func(o T1) (key T2,value T3)，参数为集合中的元素，返回值为分组的key和value。
Group 方法为终止操作，返回值为分组的map。

    func (s *stream) Group(groupFunc interface{}) interface{}\

### 最大值 Max ###
Max 方法返回集合中最大的元素，需要提供一个比较函数，形如func(o1,o2 T) bool，参数为集合中的两个元素，返回值为第一参数是否小于第二个参数。
Max 方法为终止操作。

    func (s *stream) Max(lessFunc interface{}) interface{}

### 最小值 Min ###
Min 方法返回集合中最大的元素，需要提供一个比较函数，形如func(o1,o2 T) bool，参数为集合中的两个元素，返回值为第一参数是否小于第二个参数。
Min 方法为终止操作。

    func (s *stream) Min(lessFunc interface{}) interface{}

例子:

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

输出:

    Max: {7 King 22 [87 91 89]}, Min: {3 Lee 15 [62 69 68]} 

### 最先匹配 First ###
First 方法返回第一个符合条件的元素，需要提供一个匹配函数，形如 func(o T) bool，参数为集合中的元素，返回值表示该元素是否匹配条件。
First 为终止操作。

    func (s *stream) First(matchFunc interface{}) interface{}

### 最后匹配 Last ###
First 方法返回第一个符合条件的元素，需要提供一个匹配函数，形如 func(o T) bool，参数为集合中的元素，返回值表示该元素是否匹配条件。
First 为终止操作。

    func (s *stream) Last(matchFunc interface{}) interface{}

### 规约 Reduce ###
Reduce 方法可以基于一个初始值，遍历将规约函数应用于集合中的每个元素，得到最终结果，规约函数形如 func(r T2,o T) T2，参数为前面的元素计算结果和当前元素，返回值为新的结果。
Reduce 为终止操作，返回值为规约计算后的结果。

    func (s *stream) Reduce(initValue interface{}, reduceFunc interface{}) interface{}

例子:

	students := createStudents()
	stream, _ := New(students)

	r := 0
	r = stream.Map(func(s student) int {
		return s.scores[0]
	}).Reduce(r, func(sum int, i int) int {
		return sum + i
	}).(int)
	fmt.Printf("\t%d\n", r)

输出:

    746

