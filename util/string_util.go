package util

func TuoFeng2SheXing(src []byte) (out []byte) {
	l := len(src)
	out = []byte{}
	for i := 0; i < l; i = i + 1 {
		// 大写变小写
		if 97 - 32 <= src[i] && src[i] <= 122 - 32 {
			if (i != 0) {
				out = append(out, 95)
			}
			out = append(out, src[i] + 32)
		} else {
			out = append(out, src[i])
		}
	}

	return
}

func SheXing2TuoFeng(src []byte) (out []byte) {
	l := len(src)
	out = make([]byte, l)

	// 首字母小写->大写
	if 97 <= src[0] && src[0] <= 122 {
		out[0] = src[0] - 32
	} else {
		out[0] = src[0]
	}

	del := 0
	for i := 1; i < l; i = i + 1 {
		// 是下划线
		if 95 == src[i] {
			// 下划线的下一个是小写字母
			if 97 <= src[i + 1] && src[i + 1] <= 122 {
				out[i - del] = src[i + 1] - 32
			} else {
				out[i - del] = src[i + 1]
			}
			del++
			i++
		} else {
			out[i - del] = src[i]
		}
	}
	out = out[0 : l - del]
	return
}

func ChunkJoin(str, sub string, l int) string {
	privateKey := ""
	for i, l := 0, len(str) / 64; i < l; i++ {
		privateKey = privateKey + "\n" + str[i * 64:(i + 1) * 64]
	}
	if privateKey == "" {
		return ""
	}
	return privateKey[1:]
}
