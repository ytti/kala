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
	encoded, err := encode(kala.Entries)
	if err != nil {
		return
	}
	container.Entries = crypt(kala, []byte(encoded))
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
	if err = kdf(kala); err != nil {
		return
	}
	if container.Entries, err = decrypt(kala, container.Entries); err != nil {
		return
	}
	if err = decode(container.Entries, &kala.Entries); err != nil {
		return
	}
	return
}

func (container *Container) Save(kala *Kala) (err error) {
	if err = container.addEntries(kala); err != nil {
		return
	}
	container.Version = Version
	container.Checksum = crc32.ChecksumIEEE(container.Entries)
	container.Salt = kala.Config.Salt
	data, err := encode(container)
	if err != nil {
		return
	}
	oldfile := kala.Config.File + ".old"
	// dont bailout on failing rename, maybe during runtime original file was deleted, if we bailout, we've lost our data
	err = os.Rename(kala.Config.File, oldfile)
	err = ioutil.WriteFile(kala.Config.File, data, 0600)
	return
}
