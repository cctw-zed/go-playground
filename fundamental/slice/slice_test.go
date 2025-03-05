package slice

import (
	"fmt"
	"testing"
)

// 测试切片的切片
func TestSliceOfSlice(t *testing.T) {
	sli1 := make([]int, 10)
	for i := 0; i < 10; i++ {
		sli1[i] = i
	}
	sli2 := sli1[:4]
	sli3 := sli2[6:8]
	sli4 := sli1[6:8]

	t.Log(sli1, sli2, sli3, sli4)

	// 从结果来看
	// 切片的切片也只是基于底层数组来构造
}

// 切片扩容是按照当前切片容量进行扩容的，而不是数组的长度
// 切片扩容后的长度不一定是2的幂，这取决于切片的初始容量
// 当前切片扩容会影响底层数组，但是不会影响另一个切片的 len 和 cap
func TestAppend(t *testing.T) {
	u := []int{11, 12, 13, 14, 15}
	fmt.Println("array:", u) // [11, 12, 13, 14, 15]
	s := u[1:3]
	fmt.Printf("slice(len=%d, cap=%d): %v\n", len(s), cap(s), s) // [12, 13]
	s = append(s, 24)
	fmt.Println("after append 24, array:", u)
	fmt.Printf("after append 24, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)
	s = append(s, 25)
	fmt.Println("after append 25, array:", u)
	fmt.Printf("after append 25, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)
	s = append(s, 26)
	fmt.Println("after append 26, array:", u)
	fmt.Printf("after append 26, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)

	s[0] = 22
	fmt.Println("after reassign 1st elem of slice, array:", u)
	fmt.Printf("after reassign 1st elem of slice, slice(len=%d, cap=%d): %v\n", len(s), cap(s), s)
}
