package kala

import (
	"os"
	"os/user"
	"path"
)

const (
	Version string = "0.1"
)

type Kala struct {
	Config    *Config
	Container *Container
	Entries   Entries
}

func (kala *Kala) Load(pw []byte) (err error) {
	kala.Config.Passphrase = pw
	err = kala.Container.Load(kala)
	return
}

func New() (kala *Kala, err error) {
	kala = &Kala{}
	kala.Config = &Config{}
	kala.Config.Nonce = &[24]byte{}
	kala.Config.Key = &[32]byte{}
	kala.Config.File, err = mkFile()
	kala.Config.Salt = mkSalt()
	kala.Config.Scrypt_N = 512
	kala.Config.Scrypt_r = 10
	kala.Config.Scrypt_p = 10
	kala.Container = &Container{}
	kala.Entries = Entries{}
	return
}

func (kala *Kala) NewPassphrase(pw []byte) (err error) {
	kala.Config.Passphrase = pw
	err = kdf(kala)
	return
}

type Config struct {
	Nonce      *[24]byte
	Key        *[32]byte
	Salt       []byte
	Passphrase []byte
	File       string
	Scrypt_N   int
	Scrypt_r   int
	Scrypt_p   int
}

func mkFile() (filename string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	filename = path.Join(usr.HomeDir, ".config", "kala")
	if err = os.MkdirAll(filename, 0700); err != nil {
		return
	}
	filename = path.Join(filename, "kala.json")
	return
}
