package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/ytti/kala"
)

var (
	container kala.Container
	file      string
	password  []byte
	entries   kala.Entries
)

type questionFunc func() string
type changeFunc func(int) bool

func main() {
	opts := cmdLineParse()
	arg := flag.Arg(0)
	container = kala.Container{}
	file, _ = container.File()
	var key [32]byte
	switch opts["action"] {
	case "show", "add", "edit", "del", "export":
		load()
		getpw()
		kala.KDF(password, &key)
		decryptAndDecode(&key)
	}
	switch opts["action"] {
	case "show":
		show(arg)
	case "add":
		add(arg, &key)
	case "edit":
		edit(arg, &key)
	case "del":
		del(arg, &key)
	case "import":
		getpw()
		kala.KDF(password, &key)
		imp(arg, &key)
	case "export":
		exp(&key)
	case "create":
		create()
	}
}

func show(arg string) {
	for _, i := range findIdx(arg) {
		entries[i].Decrypt(password)
		fmt.Println(entries[i])
	}
}

func add(arg string, key *[32]byte) {
	strs := map[string]string{
		"Name": arg,
		"Host": arg,
	}
	entry := getEntry(strs, []byte{})
	fmt.Printf("Add Entry (YES/no): ")
	if getInput() == "no" {
		fmt.Println("Ignoring entry")
	} else {
		entries.Add(entry)
		save(key)
		fmt.Println("Entry written to file")
	}
}

func edit(arg string, key *[32]byte) {
	strs := map[string]string{
		"act_this":     "Edit this entry (yes/NO): ",
		"act_ack":      "Entry changed\n\n",
		"act_nak":      "Skipping entry\n\n",
		"act_conf":     "entries changed, commit changes (yes/NO): ",
		"act_conf_ack": "Changes written to file\n\n",
		"act_conf_nak": "Changes ignored\n\n",
	}
	change(arg, strs, changeEdit, key)
}

func del(arg string, key *[32]byte) {
	strs := map[string]string{
		"act_this":     "Delete this entry (yes/NO): ",
		"act_ack":      "Marking for deletion\n\n",
		"act_nak":      "Skipping deletion\n\n",
		"act_conf":     "entries marked for deletion, commit changes (yes/NO): ",
		"act_conf_ack": "Deletions written to file\n\n",
		"act_conf_nak": "Deletions ignored\n\n",
	}
	change(arg, strs, changeDel, key)
}

func imp(filename string, key *[32]byte) {
	var err error
	entries, err = kala.Import(filename, password, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d entries found, commit changes to disk, replacing existing data (yes/NO): ", len(entries))
	if getInput() == "yes" {
		save(key)
		fmt.Println("entries written to disk")
	} else {
		fmt.Println("skipping import")
	}
}

func exp(key *[32]byte) {
	str, err := kala.Export(password, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)
}

func create() (ok bool) {
	if _, err := os.Stat(file); err == nil {
		fmt.Printf("%s already exists, not creating\n", file)
		ok = false
		return
	}
	fmt.Printf("Creating new %s, new master passphrase needed\n", file)
	entries = kala.Entries{}
	password = []byte(newpw())
	var key [32]byte
	kala.KDF(password, &key)
	save(&key)
	fmt.Printf("Empty %s created\n", file)
	ok = true
	return
}

func change(arg string, strs map[string]string, changeFn changeFunc, key *[32]byte) {
	var changes uint
	idx := findIdx(arg)
	// reverse the result, as index changes when we delete it
	for i := len(idx) - 1; i >= 0; i-- {
		eidx := idx[i]
		fmt.Println()
		fmt.Print(entries[eidx])
		fmt.Printf("%s", strs["act_this"])
		if getInput() == "yes" {
			changeFn(eidx)
			fmt.Printf("%s", strs["act_ack"])
			changes++
		} else {
			fmt.Printf("%s", strs["ack_nak"])
		}
	}
	if changes > 0 {
		fmt.Printf("%d %s", changes, strs["act_conf"])
		if getInput() == "yes" {
			save(key)
			fmt.Printf("%s", strs["act_conf_ack"])
		} else {
			fmt.Printf("%s", strs["act_conf_nak"])
		}
	}
}

func changeDel(index int) (ok bool) {
	entries.Delete(index)
	ok = true
	return
}
func changeEdit(index int) (ok bool) {
	e := entries[index]
	strs := map[string]string{
		"Name":        e.Name,
		"Host":        e.Host,
		"Username":    e.Username,
		"Information": e.Information,
	}
	entries[index] = getEntry(strs, e.Secret)
	ok = true
	return
}

func getEntry(def map[string]string, oldsecret []byte) (entry kala.Entry) {
	entry.Name = ask("Name", def["Name"], getInput)
	entry.Host = ask("Host", def["Host"], getInput)
	entry.Username = ask("Username", def["Username"], getInput)
	entry.Information = ask("Information", def["Information"], getInput)
	secret := kala.SecretEntry{}
	secret.Passphrase = newpw()
	if secret.Passphrase == "" && len(oldsecret) > 0 {
		entry.Secret = oldsecret
	} else {
		entry.AddSecret(secret, password)
	}
	return
}

func cmdLineParse() (opts map[string]string) {
	opts = make(map[string]string)
	add := flag.Bool("a", false, "add entry")
	del := flag.Bool("d", false, "delete entry")
	edit := flag.Bool("e", false, "edit entry")
	imp := flag.Bool("import", false, "import entries - destructive, not merge")
	exp := flag.Bool("export", false, "export entries")
	create := flag.Bool("create", false, "create file")
	flag.Parse()
	switch {
	case *add:
		opts["action"] = "add"
	case *edit:
		opts["action"] = "edit"
	case *del:
		opts["action"] = "del"
	case *exp:
		opts["action"] = "export"
	case *imp:
		opts["action"] = "import"
	case *create:
		opts["action"] = "create"
	default:
		opts["action"] = "show"
	}
	return
}

func load() {
	if err := container.Load(file); err != nil {
		if create() == true {
			os.Exit(0)
		} else {
			log.Fatal(err)
		}
	}
}

func decryptAndDecode(key *[32]byte) {
	var err error
	if err := container.Decrypt(key); err != nil {
		log.Fatal(err)
	}
	entries, err = container.Decode()
	if err != nil {
		log.Fatal(err)
	}
}

func save(key *[32]byte) {
	if err := container.AddEntries(entries, key); err != nil {
		log.Fatal(err)
	}
	if err := container.Save(file, key); err != nil {
		log.Fatal(err)
	}
}

func findIdx(find string) (found []int) {
	found, err := entries.Find(find)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func ask(question string, def string, fn questionFunc) (answer string) {
	qstr := fmt.Sprintf("%s [%s]", question, def)
	fmt.Printf("%23s : ", qstr)
	astr := fn()
	answer = string(astr)
	if answer == "" {
		answer = def
	}
	return
}

func getInput() (str string) {
	in := bufio.NewReader(os.Stdin)
	str, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	str = strings.TrimRight(str, "\r\n")
	return
}

func getSilent() (str string) {
	a := gopass.GetPasswd()
	str = string(a)
	return
}

func getpw() {
	fmt.Fprintf(os.Stderr, "Passphrase for %s: ", file)
	password = gopass.GetPasswd()
}

func newpw() (pw1 string) {
	pw1, pw2 := newpwask()
	for pw1 != pw2 {
		fmt.Println("Passphrases do not match, try again")
		pw1, pw2 = newpwask()
	}
	return
}

func newpwask() (pw1 string, pw2 string) {
	pw1 = ask("Passphrase", "", getSilent)
	pw2 = ask("Passphrase (again)", "", getSilent)
	return
}
