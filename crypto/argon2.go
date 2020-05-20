package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// Argon2Config required configuration to perform Argon2 hashing algorithm
type Argon2Config struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultArgon2Config returns a basic configuration for Argon2 actions
func DefaultArgon2Config() *Argon2Config {
	return &Argon2Config{
		Memory:  64 * 1024, // 64MB
		Time:    1,
		Threads: 4,
		KeyLen:  32,
	}
}

// Argon2HashString hashes a simple string with Argon2 algorithm, returns a base64-encoded string
// for future comparison
func Argon2HashString(s string, cfg *Argon2Config) string {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return ""
	}

	if cfg == nil {
		cfg = DefaultArgon2Config()
	}

	hash := argon2.IDKey([]byte(s), salt, cfg.Time, cfg.Memory, cfg.Threads, cfg.KeyLen)

	encodedS := base64.RawStdEncoding.EncodeToString(salt)
	encodedH := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, 64*1024, 1, 4,
		encodedS, encodedH)
}

// Argon2CompareString takes an string an compares it with the sent base64-encoded hash and
// returns a bool (false = not equal, true = equal)
func Argon2CompareString(s, hash string) bool {
	parts := strings.Split(hash, "$")

	c := &Argon2Config{}

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Time, &c.Threads)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	c.KeyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(s), salt, c.Time, c.Memory, c.Threads, c.KeyLen)

	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1)
}
