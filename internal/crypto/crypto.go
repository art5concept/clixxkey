package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"runtime"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	KeyLen  = 32 // 256 bits
	SaltLen = 16
)

// Deriva una clave usando Argon2id
func DeriveKey(password, salt []byte) []byte {
	// return argon2.IDKey(password, salt, 1, 64*1024, 4, KeyLen)
	return argon2.IDKey(password, salt, 1, 128*1024, 4, KeyLen)
}

// Genera un salt aleatorio
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLen)
	_, err := rand.Read(salt)
	return salt, err
}

// Detecta si el CPU soporta AES-NI (simplificado)
func HasAESNI() bool {
	return runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64"
}

// Cifra usando AES-256-GCM o XChaCha20-Poly1305 según disponibilidad
func Encrypt(data, key []byte) ([]byte, error) {
	if HasAESNI() {
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		nonce := make([]byte, aesgcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}
		ciphertext := aesgcm.Seal(nonce, nonce, data, nil)
		return append([]byte{0x01}, ciphertext...), nil // 0x01 = AES-GCM
	} else {
		aead, err := chacha20poly1305.NewX(key)
		if err != nil {
			return nil, err
		}
		nonce := make([]byte, chacha20poly1305.NonceSizeX)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}
		ciphertext := aead.Seal(nonce, nonce, data, nil)
		return append([]byte{0x02}, ciphertext...), nil // 0x02 = XChaCha20
	}
}

// Descifra usando el modo correcto según el prefijo
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < 1 {
		return nil, errors.New("ciphertext too short")
	}
	mode := ciphertext[0]
	ciphertext = ciphertext[1:]
	switch mode {
	case 0x01: // AES-GCM
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		nonceSize := aesgcm.NonceSize()
		if len(ciphertext) < nonceSize {
			return nil, errors.New("ciphertext too short for nonce")
		}
		nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
		return aesgcm.Open(nil, nonce, ct, nil)
	case 0x02: // XChaCha20-Poly1305
		aead, err := chacha20poly1305.NewX(key)
		if err != nil {
			return nil, err
		}
		nonceSize := chacha20poly1305.NonceSizeX
		if len(ciphertext) < nonceSize {
			return nil, errors.New("ciphertext too short for nonce")
		}
		nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
		return aead.Open(nil, nonce, ct, nil)
	default:
		return nil, errors.New("unknown encryption mode")
	}
}
