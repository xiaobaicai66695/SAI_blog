package service

import (
	"SAI_blog/repository"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func UpdateIco(uid int64, file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	if !(ext == "" || ext == ".jpg" || ext == ".png" || ext == ".jpeg") {
		return "", fmt.Errorf("格式错误")
	}
	filePath := fmt.Sprintf("./static/user_ico/%d%s", uid, ext)
	//filePath := fmt.Sprintf("/root/project/SAI_blog/static/user_ico/%d%s", uid, ext)
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}
	errChan := make(chan error)
	go func() {
		errChan <- SaveImages(file, filePath)
	}()
	newFilePath := strings.TrimPrefix(filePath, ".")
	repository.UpdateUserIco(uid, newFilePath)
	err := <-errChan
	return newFilePath, err
}

func SaveImages(file *multipart.FileHeader, filePath string) error {
	uploadFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadFile.Close()
	saveFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer saveFile.Close()
	_, err = io.Copy(saveFile, uploadFile)
	return err
}
