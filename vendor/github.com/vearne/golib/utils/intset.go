package utils

type IntSet struct {
	InternalMap map[int]int
}

func NewIntSet() *IntSet {
	return &IntSet{InternalMap: make(map[int]int)}
}
func (set *IntSet) Add(i int) {
	set.InternalMap[i] = 1
}

func (set *IntSet) AddAll(itemSlice []int) {
	for _, item := range itemSlice {
		set.InternalMap[item] = 1
	}
}

func (set *IntSet) Remove(key int) {
	delete(set.InternalMap, key)
}

func (set *IntSet) Size() int {
	return len(set.InternalMap)
}

func (set *IntSet) Intersection(set2 *IntSet) *IntSet {
	result := NewIntSet()
	for key, _ := range set.InternalMap {
		if _, ok := set2.InternalMap[key]; ok {
			result.Add(key)
		}
	}
	return result
}

func (set *IntSet) ToSlice() []int {
	res := make([]int, 0, 5)
	for key := range set.InternalMap {
		res = append(res, key)
	}
	return res
}
