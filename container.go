package kala

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
)

type Container struct {
	Version  string
	Checksum uint32
	Salt     []byte
	Nonce    []byte
	Entries  []byte
}

type checksumError struct {
	got    uint32
	wanted uint32
}

func (f checksumError) Error() string {
	return fmt.Sprintf("checksum failed, got %d, wanted %d", f.got, f.wanted)
}

func (container *Container) addEntries(kala *Kala) (err error) {
	var encoded []byte
	if err = encode(kala.Entries, &encoded); err != nil {
		return
	}
	container.Entries = crypt(kala.Config.Key, kala.Config.Nonce, encoded)
	return
}

func (container *Container) Load(kala *Kala) (err error) {
	data, err := ioutil.ReadFile(kala.Config.File)
	if err != nil {
		return
	}
	if err = decode(data, &container); err != nil {
		return
	}
	checksum := crc32.ChecksumIEEE(container.Entries)
	if checksum != container.Checksum {
		err = &checksumError{checksum, container.Checksum}
		return
	}
	kala.Config.Salt = container.Salt
	copy(kala.Config.Nonce[:], container.Nonce)
	if err = kdf(kala); err != nil {
		return
	}
	if container.Entries, err = decrypt(kala.Config.Key, kala.Config.Nonce, container.Entries); err != nil {
		return
	}
	err = decode(container.Entries, &kala.Entries)
	return
}

func (container *Container) Save(kala *Kala) (err error) {
	mkNonce(kala.Config.Nonce)
	if err = container.addEntries(kala); err != nil {
		return
	}
	container.Version = Version
	container.Salt = kala.Config.Salt
	container.Nonce = kala.Config.Nonce[0:24]
	container.Checksum = crc32.ChecksumIEEE(container.Entries)
	var encoded []byte
	if err = encode(container, &encoded); err != nil {
		return
	}
	oldfile := kala.Config.File + ".old"
	tmpfile := kala.Config.File + ".tmp"
	if err = ioutil.WriteFile(tmpfile, encoded, 0600); err != nil {
		return
	}
	os.Rename(kala.Config.File, oldfile)
	err = os.Rename(tmpfile, kala.Config.File)
	return
}
