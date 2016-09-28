package bracelet

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/seka17/all/structures"
)

var rePhone = regexp.MustCompile(`[ \-()_]`)

func formatName(name string) string {
	res := make([]byte, 0)

	i := 0
	for _, x := range name {
		tmp, _ := hex.DecodeString(fmt.Sprintf("%x", x))
		if len(tmp) == 1 {
			res = append(res, []byte{0x00, tmp[0]}...)
		} else {
			res = append(res, tmp...)
		}
		i++
	}
	return fmt.Sprintf("%x", string(res))
}

func (this Bracelet) Reset() error {
	_, err := this.Write(this.AddHeader("FACTORY"))
	return err
}

func (this Bracelet) PowerOff() error {
	_, err := this.Write(this.AddHeader("POWEROFF"))
	return err
}

func (this Bracelet) PrintMessage(msg string) error {
	_, err := this.Write(this.AddHeader("MESSAGE," + msg))
	return err
}

func (this Bracelet) SetContacts(contacts []structures.Contact) error {
	// Maximum contacts in contacts book is 10 (first 5 and last 5)
	// Maximum lengh of name of the contact is 10
	contacts = structures.SortContacts(contacts)
	var length int
	if len(contacts) > 10 {
		length = 10
	} else {
		length = len(contacts)
	}
	tmp := make([]string, length)
	for i, v := range contacts {
		if i == length {
			break
		}
		number := rePhone.ReplaceAllString(v.Number, "")
		name := string([]byte(v.Name))
		if len(v.Name) > 10 {
			name = hex.EncodeToString([]byte(name[:10]))
		}

		tmp[i] = fmt.Sprintf("%s,%s", number, formatName(name))
	}

	// Send first 5 entries
	first := length
	if first > 5 {
		first = 5
	}
	_, err := this.Write(this.AddHeader("PHB," + strings.Join(tmp[:first], ",")))
	if err != nil {
		return err
	}
	if length > 5 {
		// Send last 5 entries
		_, err = this.Write(this.AddHeader("PHB2," + strings.Join(tmp[5:], ",")))
	}
	return err
}

func (this Bracelet) Configure(configuration map[string]interface{}) error {
	// Set interval of gps coordinates comming in seconds, default 30
	interval := "30"
	if v, ok := configuration["interval"]; ok {
		interval = strconv.Itoa(v.(int))
	}
	if _, err := this.Write(this.AddHeader("UPLOAD," + interval)); err != nil {
		return err
	}

	// Sms notification on or off, default is off
	sms := "0"
	if v, ok := configuration["sms"]; ok {
		sms = strconv.Itoa(v.(int))
	}
	// Turn alarm sms notifications on/off
	if _, err := this.Write(this.AddHeader("SOSSMS," + sms)); err != nil {
		return err
	}
	// Turn low battery sms notifications on/off
	if _, err := this.Write(this.AddHeader("LOWBAT," + sms)); err != nil {
		return err
	}
	// Turn pedometr sms notifications on/off
	if _, err := this.Write(this.AddHeader("PEDO," + sms)); err != nil {
		return err
	}
	// Turn all sms services on/off
	if _, err := this.Write(this.AddHeader("SMSONOFF," + sms)); err != nil {
		return err
	}

	// Turn belt sensor on/off, default on
	belt := "1"
	if v, ok := configuration["belt"]; ok {
		sms = strconv.Itoa(v.(int))
	}
	if _, err := this.Write(this.AddHeader("REMOVE," + belt)); err != nil {
		return err
	}

	// Wake up terminal GPS module
	if _, err := this.Write(this.AddHeader("CR")); err != nil {
		return err
	}
	return nil

}

func (this Bracelet) Call(number string) error {
	number = rePhone.ReplaceAllString(number, "")
	_, err := this.Write(this.AddHeader("CALL," + number))
	return err
}
