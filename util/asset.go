package util

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// 保存上传文件
func SaveFormFile(root string, mFile multipart.File, fHeader *multipart.FileHeader) (hashCode ,filePath string, err error) {
	defer mFile.Close()
	hash := md5.New()
	if _, err = io.Copy(hash, mFile); err != nil {
		return
	}

	hashByte := hash.Sum(nil)
	// 获取hash
	hashCode = fmt.Sprintf("%x", hashByte)

	fileNameByte := []byte(fHeader.Filename)
	extType := string(fileNameByte[bytes.LastIndexByte(fileNameByte, '.')+1:])

	// 构建文件保存路径与文件名
	filePath = root + buildFilePath(hashCode, extType, 2)
	// 生成目录
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "","", err
		}
	}
	// copy 文件
	nFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer nFile.Close()

	mFile.Seek(0, io.SeekStart)
	if _, err = io.Copy(nFile, mFile); err != nil {
		return
	}

	return
}

// 将hash码分割,作为层级
func buildFilePath(hash, extType string, level int) string {
	path := ChunkJoin(hash, "/", 2)[:3*level] + hash + "." + extType
	return path
}
