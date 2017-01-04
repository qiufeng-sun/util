// 字符串和数字排序
// 1.如果能转成数字, 则转成数字比较;2.数字排在字符串前;3.升序排列
package sort

import (
	"fmt"
	"sort"
	"strconv"
)

//
type Element struct {
	Num    int
	Str    string
	IsNum  bool
	CanNum bool
}

//
func NewElementNum(num int) *Element {
	return &Element{
		Num:    num,
		Str:    fmt.Sprintf("%v", num),
		IsNum:  true,
		CanNum: true,
	}
}

func NewElementStr(str string) *Element {
	canNum := true
	num, e := strconv.Atoi(str)
	if e != nil {
		canNum = false
	}

	return &Element{
		Num:    num,
		Str:    str,
		IsNum:  false,
		CanNum: canNum,
	}
}

//
type ElementSlice []*Element

//
func (this ElementSlice) Len() int { return len(this) }

func (this ElementSlice) Less(i, j int) bool {
	ei := this[i]
	ej := this[j]

	if ei.CanNum && ej.CanNum {
		return ei.Num < ej.Num
	}

	if !ei.IsNum && !ej.IsNum {
		return ei.Str < ej.Str
	}

	return ei.CanNum
}

func (this ElementSlice) Swap(i, j int) { this[i], this[j] = this[j], this[i] }

//
func StringNums(a []*Element) { sort.Sort(ElementSlice(a)) }
