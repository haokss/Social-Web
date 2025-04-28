package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
)

// 生成文件名
func GenerateFileName(uid uint, file *multipart.FileHeader) (string, error) {
	// 生成哈希值
	hash, err := GenerateMD5Hash(file)
	if err != nil {
		return "", err
	}

	// 获取文件后缀
	ext := filepath.Ext(file.Filename)

	// 生成唯一文件名
	fileName := fmt.Sprintf("%d%s%s", uid, hash, ext)
	return fileName, nil
}
