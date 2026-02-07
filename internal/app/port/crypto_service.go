package port

// CryptoService provides encryption/decryption for sensitive data.
type CryptoService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
