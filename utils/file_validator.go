package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
)

func ValidateFile(file *multipart.FileHeader, maxSize int64, allowedTypes []string) error {
	if !StringInSlice(allowedTypes, file.Header.Get("Content-Type")) {
		return errors.New("file type tidak valid")
	}
	if file.Size > maxSize {
		// format max size to B or KB or MB according to the unit
		maxSizeStr := fmt.Sprintf("%.2f", float64(maxSize)) + " B"
		if maxSize > 1024 {
			maxSizeStr = fmt.Sprintf("%.2f", float64(maxSize)/1024) + " KB"
		}
		if maxSize > 1024*1024 {
			maxSizeStr = fmt.Sprintf("%.2f", float64(maxSize)/1024/1024) + " MB"
		}
		return errors.New("ukuran file maksimal " + maxSizeStr)
	}
	return nil
}
