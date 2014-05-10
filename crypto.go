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

func crypt(key *[32]byte, nonce *[24]byte, clear []byte) (crypted []byte) {
	crypted = secretbox.Seal(crypted[:0], clear, nonce, key)
	return
}

func decrypt(key *[32]byte, nonce *[24]byte, crypted []byte) (clear []byte, err error) {
	clear, ok := secretbox.Open(clear[:0], crypted, nonce, key)
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

func mkNonce(nonce *[24]byte) {
	non := make([]byte, 24)
	rand.Read(non)
	copy(nonce[:], non)
}

func encode(in interface{}, out *[]byte) (err error) {
	*out, err = json.MarshalIndent(in, "", "  ")
	return
}

func decode(in []byte, out interface{}) (err error) {
	err = json.Unmarshal(in, &out)
	return
}
