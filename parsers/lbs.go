package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LBSRequest struct {
	HomeMobileCountryCode uint16            // The mobile country code stored on the SIM card (100-999).
	HomeMobileNetworkCode uint16            // The mobile network code stored on the SIM card (0-32767).
	RadioType             string            // The mobile radio type. Supported values are lte, gsm, cdma, and wcdma.
	Carrier               string            // The clear text name of the cell carrier / operator.
	ConsiderIp            bool              // Should the clients IP address be used to locate it, defaults to true.
	CellTowers            []CellTower       // Array of cell towers
	WifiAccessPoints      []WifiAccessPoint // Array of wifi access points
	IPAddress             string            // Client IP Address
	Fallbacks             *Fallbacks        // The fallback section is a custom addition to the GLS API.
}

type CellTower struct {
	MobileCountryCode uint16 // The mobile country code.
	MobileNetworkCode uint16 // The mobile network code.
	LocationAreaCode  uint16 // The location area code for GSM and WCDMA networks. The tracking area code for LTE networks.
	CellId            uint32 // The cell id or cell identity.
	SignalStrength    int16  // The signal strength for this cell network, either the RSSI or RSCP.
	Age               uint32 // The number of milliseconds since this networks was last detected.
	TimingAdvance     uint8  // The timing advance value for this cell network.
}

type Fallbacks struct {
	LAC bool // If no exact cell match can be found, fall back from exact cell position estimates to more coarse grained cell location area estimates, rather than going directly to an even worse GeoIP based estimate.
	IP  bool // If no position can be estimated based on any of the provided data points, fall back to an estimate based on a GeoIP database based on the senders IP address at the time of the query.
}

// ParseLBS разбирает строку с информацией в формате LBS и возвращает его описание.
// Первым параметром указывается тип радио (gsm, lte, cdma, wcdam и так далее). Вторым — строка
// с данными LBS. В ответ возвращает сформированную для запроса данных структуру.
func ParseLBS(radio, lbsStr string) (*LBSRequest, error) {
	switch radio {
	case "", "gsm", "lte", "cdma", "wcdma":
	default:
		return nil, fmt.Errorf("bad radio type: %s", radio)
	}
	splitted := strings.Split(lbsStr, "-") // разделяем на элементы
	if len(splitted) < 7 {
		return nil, errors.New("agps - wrong data (len < 7)")
	}
	mcc, err := strconv.ParseUint(splitted[3], 16, 16)
	if err != nil {
		return nil, fmt.Errorf("bad MCC: %s", splitted[3])
	}
	mnc, err := strconv.ParseUint(splitted[4], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("bad MNC: %s", splitted[4])
	}
	cellTowers := make([]CellTower, (len(splitted)-5)/3)
	for i := range cellTowers {
		area, err := strconv.ParseUint(splitted[5+i*3], 16, 16)
		if err != nil {
			return nil, fmt.Errorf("bad Area: %s", splitted[5+i*3])
		}
		id, err := strconv.ParseUint(splitted[6+i*3], 16, 32)
		if err != nil {
			return nil, fmt.Errorf("bad Cell ID: %s", splitted[6+i*3])
		}
		dbm, err := strconv.ParseUint(splitted[7+i*3], 16, 16)
		if err != nil {
			return nil, fmt.Errorf("bad DBM: %s", splitted[7+i*3])
		}
		cellTowers[i] = CellTower{
			CellId:            uint32(id),
			LocationAreaCode:  uint16(area),
			MobileCountryCode: uint16(mcc),
			MobileNetworkCode: uint16(mnc),
			SignalStrength:    int16(dbm - 220),
		}
	}
	return &LBSRequest{
		RadioType:             radio,
		HomeMobileCountryCode: uint16(mcc),
		HomeMobileNetworkCode: uint16(mnc),
		ConsiderIp:            false,
		CellTowers:            cellTowers,
	}, nil
}
