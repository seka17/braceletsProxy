package bracelet

import (
	"regexp"

	"github.com/seka17/all/structures"
)

var rePhone = regexp.MustCompile("[ -()_]")

func (this Bracelet) Reset() error {
	return nil
}

func (this Bracelet) PowerOff() error {
	return nil
}

func (this Bracelet) PrintMessage(msg string) error {
	return nil
}

func (this Bracelet) SetContacts(contacts []structures.Contact) error {
	return nil
}

func (this Bracelet) Configure(configuration map[string]interface{}) error {
	return nil
}

func (this Bracelet) Call(number string) error {
	return nil
}
