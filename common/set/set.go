package set

// Set 是一个泛型集合类型
type Set[T comparable] map[T]struct{}

// NewSet 创建一个新的集合
func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

// Add 向集合中添加元素
func (s Set[T]) Add(items ...T) {

	for _, item := range items {
		(s)[item] = struct{}{}
	}
}

// Remove 从集合中移除元素
func (s Set[T]) Remove(item T) {
	delete(s, item)
}

// Contains 检查集合中是否包含指定元素
func (s Set[T]) Contains(item T) bool {
	_, exists := s[item]
	return exists
}

// Size 返回集合中元素的数量
func (s Set[T]) Size() int {
	return len(s)
}

// Values 返回集合中所有元素的切片
func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))
	for item := range s {
		values = append(values, item)
	}
	return values
}
func (s Set[T]) ToArr() []T {
	if len(s) == 0 {
		return []T{}
	}
	out := make([]T, 0, len(s))
	for key, _ := range s {
		out = append(out, key)
	}
	return out
}

func Array2Set[T comparable](ins []T) Set[T] {
	set := NewSet[T]()
	for _, in := range ins {
		set.Add(in)
	}
	return set
}
