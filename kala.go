package kala

import (
	"encoding/json"
	"fmt"

	"code.google.com/p/go.crypto/nacl/secretbox"
	"code.google.com/p/go.crypto/scrypt"
)

const (
	Version  string = "0.1"
	salt     string = "19dd1947a95454ccaf223a731c32db0c" // md5sum(kala)
	scrypt_N int    = 512
	scrypt_r int    = 10
	scrypt_p int    = 10
)

var (
	nonce [24]byte
)

type DecryptError struct {
	err string
}

func (f DecryptError) Error() string {
	return fmt.Sprintf("decryption failed: %s", f.err)
}

func Encode(data interface{}) (encoded []byte, err error) {
	encoded, err = json.MarshalIndent(data, "", "  ")
	return
}

func Decode(encoded []byte, v interface{}) (err error) {
	err = json.Unmarshal(encoded, &v)
	return
}

func Crypt(clear []byte, key *[32]byte) (crypted []byte) {
	crypted = secretbox.Seal(crypted[:0], clear, &nonce, key)
	return
}

func Decrypt(crypted []byte, key *[32]byte) (clear []byte, err error) {
	clear, ok := secretbox.Open(clear[:0], crypted, &nonce, key)
	if ok != true {
		err = &DecryptError{"incorrect password"}
	}
	return
}

func KDFwithSalt(pw []byte, salt []byte, key *[32]byte) (err error) {
	keyslice, err := scrypt.Key(pw, salt, scrypt_N, scrypt_r, scrypt_p, 32)
	copy(key[:], keyslice)
	return
}

func KDF(pw []byte, key *[32]byte) (err error) {
	err = KDFwithSalt(pw, []byte(salt), key)
	return
}
