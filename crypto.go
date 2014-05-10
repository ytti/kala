package kala

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"code.google.com/p/go.crypto/nacl/secretbox"
	"code.google.com/p/go.crypto/scrypt"
)

type decryptError struct {
	err string
}

func (f decryptError) Error() string {
	return fmt.Sprintf("decryption failed: %s", f.err)
}

func crypt(kala *Kala, clear []byte) (crypted []byte) {
	crypted = secretbox.Seal(crypted[:0], clear, kala.Config.Nonce, kala.Config.Key)
	return
}

func decrypt(kala *Kala, crypted []byte) (clear []byte, err error) {
	clear, ok := secretbox.Open(clear[:0], crypted, kala.Config.Nonce, kala.Config.Key)
	if ok != true {
		err = &decryptError{"incorrect password"}
	}
	return
}

func kdfEntry(kala *Kala, salt []byte, key *[32]byte) (err error) {
	keyslice, err := scrypt.Key(kala.Config.Passphrase, salt, kala.Config.Scrypt_N, kala.Config.Scrypt_r, kala.Config.Scrypt_p, 32)
	if err != nil {
		return
	}
	copy(key[:], keyslice)
	return
}
func kdf(kala *Kala) (err error) {
	err = kdfEntry(kala, kala.Config.Salt, kala.Config.Key)
	return
}

func mkSalt() (salt []byte) {
	salt = make([]byte, 32)
	rand.Read(salt)
	return
}

func encode(data interface{}) (encoded []byte, err error) {
	encoded, err = json.MarshalIndent(data, "", "  ")
	return
}

func decode(encoded []byte, v interface{}) (err error) {
	err = json.Unmarshal(encoded, &v)
	return
}
