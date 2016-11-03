package util

func TuoFeng2SheXing() {

}

func SheXing2TuoFeng(src []byte) (out []byte) {
	l := len(src)
	out = make([]byte, l)

	// 小写->大写
	if 97 <= src[0] && src[0] <= 122 {
		out[0] = src[0] - 32
	}

	del := 0
	for i := 1; i < l; i = i + 1 {
		// 小写
		if 97 <= src[i] && src[i] <= 122 {
			out[i - del] = src[i]
		} else if 65 <= src[i] && src[i] <= 90 {
			// 大写
			out[i - del] = src[i]
		} else {
			if 97 <= src[i + 1] && src[i + 1] <= 122 {
				out[i - del] = src[i + 1] - 32
			} else {
				out[i - del] = src[i + 1]
			}

			del++
			i++
		}
	}
	out = out[0 : l - del]
	return
}
