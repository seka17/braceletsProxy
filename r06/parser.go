package bracelet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/seka17/all/structures"
)

func (this Bracelet) Parse(command []byte) ([]byte, error) {
	fmt.Println(string(command))
	data := strings.Split(string(command), ";")
	// append([]byte(fmt.Sprintf("ServiceIP:%s,%s;", parsedService[0], parsedService[1])), []byte{0x010, 0x010, 0x01}...)
	switch data[0] {
	// Requesting main service address, it must be the same as gate server
	case "#@H00@#":
		addr, port := this.server.GetAddress()
		return append([]byte(fmt.Sprintf("ServiceIP:%s,%s;", addr, port)), []byte{0x01, 0x01, 0x01}...), nil
	default:
		return nil, fmt.Errorf("Unsupported action -> %s, body -> %s", data[0], strings.Join(data[1:], ";"))
	}
}

func parseLocation(data []string) (*structures.Point, error) {
	// Ignore lbs and wifi, because bracelet sends coordinates
	var lat, lon float64
	var err error
	// Parse longitude
	switch data[4] {
	case "N":
		lon, err = strconv.ParseFloat(data[3], 64)
		break
	case "S":
		lon, err = strconv.ParseFloat("-"+data[3], 64)
		break
	default:
		return nil, fmt.Errorf("Unknown format of longitude zone %s", data[4])
	}
	if err != nil {
		return nil, err
	}

	// Parse latitude
	switch data[6] {
	case "W":
		lat, err = strconv.ParseFloat(data[5], 64)
		break
	case "E":
		lat, err = strconv.ParseFloat("-"+data[5], 64)
		break
	default:
		return nil, fmt.Errorf("Unknown format of latitude zone %s", data[6])
	}
	if err != nil {
		return nil, err
	}

	return &structures.Point{Point: [2]float64{lon, lat}, Accuracy: 0}, nil
}
