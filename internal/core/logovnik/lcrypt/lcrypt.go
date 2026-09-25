package lcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt шифрует строку с использованием секретного ключа
func Encrypt(text string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Генерируем уникальный вектор инициализации (nonce) для каждого сообщения
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Шифруем данные и добавляем nonce в начало среза байт
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(text), nil)

	// Переводим в Base64 для удобного хранения/передачи в виде строки
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt расшифровывает строку обратно
func Decrypt(encryptedText string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// password := "мой_очень_короткий_или_длинный_пароль_123"
// При шифровании: генерируем ключ и соль
// key, salt, _ := DeriveKeyPBKDF2(password, nil)
// fmt.Printf("Ключ (32 байта): %x\n", key)
// fmt.Printf("Соль (сохраните её рядом с шифром): %x\n", salt)

// DeriveKeyPBKDF2 принимает пароль любой длины и возвращает 32-байтовый ключ и соль
// return key, salt, nil
func DeriveKeyPBKDF2(password string, customSalt []byte) (key []byte, salt []byte, err error) {
	if customSalt == nil {
		// Генерируем случайную соль (рекомендуется 8-16 байт)
		salt = make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	} else {
		salt = customSalt
	}

	// 4096 итераций, хэш SHA-256, целевая длина 32 байта
	key, err = pbkdf2.Key(sha256.New, password, salt, 4096, 32)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}
