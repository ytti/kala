package kala

import "regexp"

type Entries []Entry

func (entries *Entries) Add(entry ...Entry) {
	*entries = append(*entries, entry...)
}

func (entries *Entries) Delete(index int) {
	a := *entries
	*entries = append(a[:index], a[index+1:]...)
}

func (entries Entries) Find(str string) (found []int, err error) {
	re, err := regexp.Compile("(?i).*" + str + ".*")
	if err != nil {
		return
	}
	for i, entry := range entries {
		switch {
		case re.MatchString(entry.Name):
		case re.MatchString(entry.Host):
		case re.MatchString(entry.Username):
		case re.MatchString(entry.Information):
		default:
			continue
		}
		found = append(found, i)
	}
	return
}
