package kala

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

var (
	ContainerFile string
)

type Container struct {
	Version  string
	Checksum uint32
	Entries  []byte
}

type ChecksumError struct {
	got    uint32
	wanted uint32
}

func (f ChecksumError) Error() string {
	return fmt.Sprintf("checksum failed, got %d, wanted %d", f.got, f.wanted)
}

func (container *Container) Decrypt(key *[32]byte) (err error) {
	container.Entries, err = Decrypt(container.Entries, key)
	return
}
func (container Container) Decode() (entries Entries, err error) {
	err = Decode(container.Entries, &entries)
	return
}
func (container *Container) AddEntries(entries Entries, key *[32]byte) (err error) {
	encoded, err := Encode(entries)
	if err != nil {
		return
	}
	container.Entries = Crypt([]byte(encoded), key)
	return
}

func (container Container) File() (filename string, err error) {
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

func (container *Container) Load(file string) (err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	if err = Decode(data, &container); err != nil {
		return
	}
	checksum := crc32.ChecksumIEEE(container.Entries)
	if checksum != container.Checksum {
		err = &ChecksumError{checksum, container.Checksum}
	}
	return
}

func (container *Container) Save(file string, key *[32]byte) (err error) {
	container.Version = Version
	container.Checksum = crc32.ChecksumIEEE(container.Entries)
	data, err := Encode(container)
	if err != nil {
		return
	}
	oldfile := file + ".old"
	// dont bailout on failing rename, maybe during runtime original file was deleted, if we bailout, we've lost our data
	err = os.Rename(file, oldfile)
	err = ioutil.WriteFile(file, data, 0600)
	return
}
