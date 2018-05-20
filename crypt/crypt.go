package crypt

// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
// https://crackstation.net/hashing-security.htm
import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	rand := base64.URLEncoding.EncodeToString(b)
	return rand[:s-1], err
}

// // Example: this will give us a 44 byte, base64 encoded output
// token, err := GenerateRandomString(32)
// if err != nil {
// 	// Serve an appropriately vague error to the
//     // user, but log the details internally.
// }

// https://gowebexamples.com/password-hashing/

//HashPassword will hash the provided string using bcrypt and return the result
func HashPassword(password string, saltLen int) (string, error) {
	// Encrypt the password with the salt prepended
	// Append the salt to the front of the password
	// For use later in CheckPasswordHash
	salt, _ := GenerateRandomString(saltLen)
	bytes, err := bcrypt.GenerateFromPassword([]byte(salt+password), 14)
	return salt + string(bytes), err
}

//CheckPasswordHash will hash the provided password and compair that with the original password hash
func CheckPasswordHash(password string, hash string, saltLen int) bool {

	salt := hash[:saltLen-1]
	realHash := hash[saltLen-1:]

	err := bcrypt.CompareHashAndPassword([]byte(realHash), []byte(salt+password))
	return err == nil
}
