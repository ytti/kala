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
	fish *kala.Kala
	err  error
)

type questionFunc func() string
type changeFunc func(int)

func main() {
	opts, arg := cmdLineParse()
	if fish, err = kala.New(); err != nil {
		log.Fatal(err)
	}
	switch opts["action"] {
	case "show", "add", "edit", "del", "export":
		if haveFile() {
			if err = fish.Load(getpw()); err != nil {
				log.Fatal(err)
			}
		} else {
			create()
		}
	}
	switch opts["action"] {
	case "show":
		show(arg)
	case "add":
		add(arg)
	case "edit":
		edit(arg)
	case "del":
		del(arg)
	case "import":
		imp(arg)
	case "export":
		exp()
	}
}

func show(arg string) {
	for _, i := range findIdx(arg) {
		fish.Entries[i].Decrypt(fish)
		fmt.Println(fish.Entries[i])
	}
}

func add(arg string) {
	strs := map[string]string{
		"Name": arg,
		"Host": arg,
	}
	entry := getEntry(strs, []byte{})
	fmt.Printf("Add Entry (YES/no): ")
	if getInput() == "no" {
		fmt.Println("Ignoring entry")
	} else {
		fish.Entries.Add(entry)
		save()
		fmt.Println("Entry written to file")
	}
}

func edit(arg string) {
	strs := map[string]string{
		"act_this":     "Edit this entry (yes/NO): ",
		"act_ack":      "Entry changed\n\n",
		"act_nak":      "Skipping entry\n\n",
		"act_conf":     "entries changed, commit changes (yes/NO): ",
		"act_conf_ack": "Changes written to file\n\n",
		"act_conf_nak": "Changes ignored\n\n",
	}
	change(arg, strs, changeEdit)
}

func del(arg string) {
	strs := map[string]string{
		"act_this":     "Delete this entry (yes/NO): ",
		"act_ack":      "Marking for deletion\n\n",
		"act_nak":      "Skipping deletion\n\n",
		"act_conf":     "entries marked for deletion, commit changes (yes/NO): ",
		"act_conf_ack": "Deletions written to file\n\n",
		"act_conf_nak": "Deletions ignored\n\n",
	}
	change(arg, strs, changeDel)
}

func imp(filename string) {
	fish.NewPassphrase([]byte(newpw()))
	if fish.Entries, err = kala.Import(fish, filename); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d entries found, commit changes to disk, replacing existing data (yes/NO): ", len(fish.Entries))
	if getInput() == "yes" {
		save()
		fmt.Println("entries written to disk")
	} else {
		fmt.Println("skipping import")
	}
}

func exp() {
	if str, err := kala.Export(fish); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(str)
	}
}

func haveFile() (file bool) {
	file = false
	if _, err := os.Stat(fish.Config.File); err == nil {
		file = true
	}
	return
}

func create() {
	fmt.Printf("Creating new %s, new master passphrase needed\n", fish.Config.File)
	fish.NewPassphrase([]byte(newpw()))
	save()
	fmt.Printf("Empty %s created\n", fish.Config.File)
}

func change(arg string, strs map[string]string, changeFn changeFunc) {
	var changes uint
	idx := findIdx(arg)
	// reverse the result, as index changes when we delete it
	for i := len(idx) - 1; i >= 0; i-- {
		eidx := idx[i]
		fmt.Println()
		fmt.Print(fish.Entries[eidx])
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
			save()
			fmt.Printf("%s", strs["act_conf_ack"])
		} else {
			fmt.Printf("%s", strs["act_conf_nak"])
		}
	}
}

func changeDel(index int) {
	fish.Entries.Delete(index)
}

func changeEdit(index int) {
	e := fish.Entries[index]
	strs := map[string]string{
		"Name":        e.Name,
		"Host":        e.Host,
		"Username":    e.Username,
		"Information": e.Information,
	}
	fish.Entries[index] = getEntry(strs, e.Secret)
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
		entry.AddSecret(fish, secret)
	}
	return
}

func cmdLineParse() (opts map[string]string, arg string) {
	opts = make(map[string]string)
	add := flag.Bool("a", false, "add entry")
	del := flag.Bool("d", false, "delete entry")
	edit := flag.Bool("e", false, "edit entry")
	imp := flag.Bool("import", false, "import entries - destructive, not merge")
	exp := flag.Bool("export", false, "export entries")
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
	default:
		opts["action"] = "show"
	}
	arg = flag.Arg(0)
	return
}

func save() {
	if err := fish.Container.Save(fish); err != nil {
		log.Fatal(err)
	}
}

func findIdx(find string) (found []int) {
	found, err := fish.Entries.Find(find)
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

func getpw() (pw []byte) {
	fmt.Fprintf(os.Stderr, "Passphrase for %s: ", fish.Config.File)
	pw = []byte(gopass.GetPasswd())
	return
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
