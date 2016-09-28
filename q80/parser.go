package bracelet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/seka17/all/structures"
)

func (this Bracelet) Parse(command []byte) ([]byte, error) {
	fmt.Println(string(command))
	data := strings.Split(string(command), ",")
	switch data[0] {
	// Ping command from bracelet
	case "LK":
		switch len(data) {
		// Empty body
		case 1:
			return this.AddHeader("LK"), nil
		case 4:
			steps, err := strconv.Atoi(data[1])
			if err != nil {
				return nil, err
			}
			turnover, err := strconv.Atoi(data[2])
			if err != nil {
				return nil, err
			}
			battery, err := strconv.Atoi(data[3])
			if err != nil {
				return nil, err
			}
			_, _, _ = steps, turnover, battery
			return this.AddHeader("LK"), nil
		default:
			return nil, fmt.Errorf("Can't parse LK command -> %s", strings.Join(data[1:], ","))
		}
	case "TKQ", "TKQ2":
		return this.AddHeader(data[0]), nil
	// Bracelet passes information about location. UD2 sends data when bracelet was offline
	case "UD", "UD2":
		// E.g.: UD,220414,134652,A,22.571707,N,113.8613968,E,0.1,0.0,
		// 100,7,60,90,1000,50,0000,4,1,460,0,9360,4082,131,9360,4092,148,9360,4091,143,9360
		// ,4153,141

		// Date: 220414,134652
		// Point: 22.571707,N,113.8613968,E
		// Spped: 0.1
		// Direction: 0.0
		// Altitude: 100
		// Number of satellites: 7
		// GSM signal strength: 60
		// for more check out docs

		point, err := parseLocation(data[1:])
		if err != nil {
			return nil, err
		}
		_ = point
		return nil, nil
	case "AL":
		point, err := parseLocation(data[1:])
		if err != nil {
			return nil, err
		}
		_ = point
		return this.AddHeader("AL"), nil
	default:
		return nil, fmt.Errorf("Unsupported action -> %s, body -> %s", data[0], strings.Join(data[1:], ","))
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
