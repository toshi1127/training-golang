package ex16

func Join(sep string, a ...string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	case 2:
		return a[0] + sep + a[1]
	case 3:
		return a[0] + sep + a[1] + sep + a[2]
	}
	n := len(sep) * (len(a) - 1) //どれくらいの長さが必要か計算
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}