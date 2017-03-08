package util

import (
	"github.com/bysir-zl/bygo/log"
	"net/http"
	"testing"
)

type Handler struct {
}

func (p Handler) ServeHTTP(rsp http.ResponseWriter, req *http.Request) () {
	req.ParseMultipartForm(1024 * 1024)
	fileHeader := req.MultipartForm.File["asset"][0]
	file, err := fileHeader.Open()
	if err != nil {
		log.Error("test", err)
		return
	}

	hash, filePath, err := SaveFormFile("D://", file, fileHeader)
	if err != nil {
		log.Error("test", err)
		return
	}
	log.Info("test", hash, filePath)
}

func TestSaveFormFile(t *testing.T) {
	err := http.ListenAndServe("127.0.0.1:3333", Handler{})
	if err != nil {
		t.Fatal(err)
	}
}
