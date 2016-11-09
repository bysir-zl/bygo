package util

import (
	"path/filepath"
	"os"
	"strings"
)

// 返回所有遍历的文件路径
// dirPth
// suffixs: 不包含.的文件后缀
func WalkDir(dirPth string, suffixs []string) (files []string, err error) {
	files = make([]string, 0, 30)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		//遍历目录
		if err != nil {
			return err
		}
		if fi.IsDir() {
			// 忽略目录
			return nil
		}
		splitL := strings.Split(fi.Name(), ".")
		suf := splitL[len(splitL) - 1]
		suf = strings.ToLower(suf)

		if ItemInArray(suf, suffixs) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}
