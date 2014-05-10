package kala

import "fmt"

type Entry struct {
	Name        string
	Host        string
	Username    string
	Information string
	Secret      []byte
}

func (entry *Entry) AddSecret(kala *Kala, secret SecretEntry) (err error) {
	encoded, err := encode(secret)
	if err != nil {
		return
	}
	var key [32]byte
	kdfEntry(kala, mysalt(kala, entry.Name), &key)
	entry.Secret = crypt(kala, encoded)
	return
}

func (entry Entry) Decode() (secret SecretEntry, err error) {
	err = decode(entry.Secret, &secret)
	return
}

func (entry *Entry) Decrypt(kala *Kala) (err error) {
	var key [32]byte
	kdfEntry(kala, mysalt(kala, entry.Name), &key)
	entry.Secret, err = decrypt(kala, entry.Secret)
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

func mysalt(kala *Kala, entry_name string) (myslt []byte) {
	myslt = []byte(entry_name)
	myslt = append(myslt, kala.Config.Salt...)
	return
}
