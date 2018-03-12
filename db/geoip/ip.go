package geoip

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var (
	geoipDb *geoip2.Reader
)

//InitDB init geoip db
func InitDB(dbFile string) error {
	db, err := geoip2.Open(dbFile)
	if err != nil {
		return err
	}
	geoipDb = db
	return nil
}

//CountryCode return country isoCode of ip
func CountryCode(ip string) (string, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return "", errors.New("wrong ip string")
	}
	country, err := geoipDb.Country(netIP)
	if err != nil || country == nil {
		return "", err
	}
	return country.Country.IsoCode, nil
}
