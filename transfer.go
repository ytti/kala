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

func Import(kala *Kala, filename string) (entries Entries, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	imps := TransportEntries{}
	if err = decode(data, &imps); err != nil {
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
		entry.AddSecret(kala, secret)
		entries = append(entries, entry)
	}
	return
}

func Export(kala *Kala) (data string, err error) {
	exports := TransportEntries{}
	if err = kala.Load(kala.Config.Passphrase); err != nil {
		return
	}
	for _, entry := range kala.Entries {
		export := TransportEntry{}
		secret := SecretEntry{}
		if err = entry.Decrypt(kala); err != nil {
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
	dataslice, err := encode(exports)
	data = string(dataslice)
	return
}
