package geoip

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var (
	db *geoip2.Reader
)

// InitDB inits geoip db
func InitDB(dbFile string) error {
	_db, err := geoip2.Open(dbFile)
	if err != nil {
		return err
	}
	db = _db
	return nil
}

// CloseDB closes db
func CloseDB() error {
	var err error
	if db != nil {
		err = db.Close()
		db = nil
	}
	return err
}

// CountryCode returns country isoCode of the ip
func CountryCode(ip string) (string, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return "", errors.New("wrong ip string")
	}
	country, err := db.Country(netIP)
	if err != nil || country == nil {
		return "", err
	}
	return country.Country.IsoCode, nil
}
