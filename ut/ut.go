package ut

func Add(a, b int) int {
	return a + b
}

func Echo(str string) string {
	return str
}

func Append(a, b []int) []int {
	return append(a, b...)
}
