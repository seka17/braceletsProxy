package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

type WiFiData struct {
	MacAddress     string `json:"m"`
	Name           string `json:"i"`
	SignalStrength string `json:"s"`
}

type WifiAccessPoint struct {
	MacAddress         string // The BSSID of the WiFi network.
	SignalStrength     int16  // The received signal strength (RSSI) in dBm.
	Age                uint32 // The number of milliseconds since this network was last detected.
	Channel            uint8  // The WiFi channel, often 1 - 13 for networks in the 2.4GHz range.
	SignalToNoiseRatio uint16 // The current signal to noise ratio measured in dB.
}

func Compare(wd1 []WiFiData, wd2 []WiFiData) (res []WiFiData) {
	same := false
	for _, wifi1 := range wd1 {
		same = false
		for _, wifi2 := range wd2 {
			if wifi1.MacAddress == wifi2.MacAddress {
				same = true
				break
			}
		}
		if same {
			continue
		}
		res = append(res, wifi1)
	}
	return
}

// ParseWiFi разбирает и возвращает список с информацией о WiFi-станциях.
func ParseWiFi(wifis ...WiFiData) ([]WifiAccessPoint, error) {
	wifiAccessPoints := make([]WifiAccessPoint, len(wifis))
	for i, wifi := range wifis {
		signalStrength, err := strconv.ParseUint(wifi.SignalStrength, 16, 16)
		if err != nil {
			return nil, fmt.Errorf("bad WiFi signal strength: %s", wifi.SignalStrength)
		}
		mac, err := hex.DecodeString(wifi.MacAddress)
		if err != nil || len(mac) != 6 {
			return nil, fmt.Errorf("bad WiFi mac address: %s", wifi.MacAddress)
		}
		wifiAccessPoints[i] = WifiAccessPoint{
			MacAddress: fmt.Sprintf("%X:%X:%X:%X:%X:%X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]),
			//SignalStrength: int16(math.Log10(float64(signalStrength)/1000) * 100),
			SignalStrength: int16(signalStrength - 256),
		}
	}
	return wifiAccessPoints, nil
}
