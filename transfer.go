package kala

import "io/ioutil"

const ()

var ()

type TransportEntry struct {
	Name        string
	Host        string
	Username    string
	Information string
	Secret      SecretEntry
}
type TransportEntries []TransportEntry

func Import(filename string, pw []byte, key *[32]byte) (entries Entries, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	imps := TransportEntries{}
	if err = Decode(data, &imps); err != nil {
		return
	}
	entries = Entries{}
	for _, imp := range imps {
		entry := Entry{}
		entry.Name = imp.Name
		entry.Host = imp.Host
		entry.Username = imp.Username
		entry.Information = imp.Information
		secret := SecretEntry{}
		secret.Passphrase = imp.Secret.Passphrase
		entry.AddSecret(secret, pw)
		entries = append(entries, entry)
	}
	return
}

func Export(pw []byte, key *[32]byte) (data string, err error) {
	exports := TransportEntries{}
	container := Container{}
	file, err := container.File()
	if err != nil {
		return
	}
	if err = container.Load(file); err != nil {
		return
	}
	if err = container.Decrypt(key); err != nil {
		return
	}
	entries, err := container.Decode()
	if err != nil {
		return
	}
	for _, entry := range entries {
		export := TransportEntry{}
		secret := SecretEntry{}
		if err = entry.Decrypt(pw); err != nil {
			return
		}
		if secret, err = entry.Decode(); err != nil {
			return
		}
		export.Name = entry.Name
		export.Host = entry.Host
		export.Username = entry.Username
		export.Information = entry.Information
		export.Secret.Passphrase = secret.Passphrase
		exports = append(exports, export)
	}
	dataslice, err := Encode(exports)
	data = string(dataslice)
	return
}
