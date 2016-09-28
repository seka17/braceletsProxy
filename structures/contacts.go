package structures

import (
	"sort"
	"strconv"
)

type Contact struct {
	Position string // position in the contact list
	Number   string // phone number
	Name     string // name of contact

	// Additional params
	Params map[string]interface{}
}

func SortContacts(contacts []Contact) []Contact {
	sort.Sort(byPosition(contacts))
	return contacts
}

type byPosition []Contact

func (a byPosition) Len() int      { return len(a) }
func (a byPosition) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPosition) Less(i, j int) bool {
	left, _ := strconv.Atoi(a[i].Position)
	right, _ := strconv.Atoi(a[j].Position)
	return left < right
}
