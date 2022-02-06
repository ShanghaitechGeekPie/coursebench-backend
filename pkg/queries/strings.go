package queries

// 返回一个字符串的渲染长度
// 非常粗糙的方案，一个ascii字符的长度为1，一个中文字符的长度为2
func GetActualLength(s string) int {
	sum := 0
	for _, c := range s {
		if c > 0x80 {
			sum += 2
		} else {
			sum += 1
		}
	}
	return sum
}
