package sort

import "testing"

//
func TestSort(t *testing.T) {
	var arrSign = []*Element{
		NewElementStr("x"),
		NewElementStr("15"),
		NewElementStr("1"),
		NewElementStr("276"),
		NewElementNum(7),
		NewElementNum(1),
		NewElementNum(0),
		NewElementStr("deliverycenter"),
	}
	var expected = "xdeliverycenter276157110"

	StringNums(arrSign)

	var strSign = ""
	l := len(arrSign)
	for i := l - 1; i >= 0; i-- {
		strSign += arrSign[i].Str
	}

	t.Log(strSign)

	if strSign != expected {
		t.Errorf("res:%v\nexpected:%v\n", strSign, expected)
	}
}
