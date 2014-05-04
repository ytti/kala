package kala

import "fmt"

type Entry struct {
	Name        string
	Host        string
	Username    string
	Information string
	Secret      []byte
}

func (entry *Entry) AddSecret(secret SecretEntry, pw []byte) (err error) {
	var key [32]byte
	KDFwithSalt(pw, []byte(entry.Name), &key)
	encoded, err := Encode(secret)
	if err != nil {
		return
	}
	entry.Secret = Crypt([]byte(encoded), &key)
	return
}

func (entry Entry) Decode() (secret SecretEntry, err error) {
	err = Decode(entry.Secret, &secret)
	return
}

func (entry *Entry) Decrypt(pw []byte) (err error) {
	var key [32]byte
	KDFwithSalt(pw, []byte(entry.Name), &key)
	entry.Secret, err = Decrypt(entry.Secret, &key)
	return
}

func (entry Entry) String() (str string) {
	show_secret := true
	secret, err := entry.Decode()
	if err != nil {
		show_secret = false
	}
	str = fmt.Sprintf("%s\n", entry.Name)
	str += fmt.Sprintf("%12s : %s\n", "Host", entry.Host)
	str += fmt.Sprintf("%12s : %s\n", "Username", entry.Username)
	str += fmt.Sprintf("%12s : %s\n", "Information", entry.Information)
	if show_secret {
		str += fmt.Sprintf("%12s : %s\n", "Passphrase", secret.Passphrase)
	}
	return
}
