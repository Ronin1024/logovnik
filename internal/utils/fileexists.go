package utils

import (
	"errors"
	"os"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true // Файл существует
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // Файла точно нет
	}
	// Любая другая ошибка (например, проблемы с правами доступа)
	// В реальном коде здесь лучше возвращать саму ошибку `err`
	return false
}
