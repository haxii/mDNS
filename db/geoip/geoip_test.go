package geoip

import (
	"testing"
)

func TestInitAndClosedDB(t *testing.T) {
	err := InitDB("../../file/GeoLite2-Country.mmdb")
	if err != nil {
		t.Error(err)
	}

	if db == nil {
		t.Fail()
	}

	err = CloseDB()
	if err != nil {
		t.Error(err)
	}
}

func TestCountryCode(t *testing.T) {
	err := InitDB("../../file/GeoLite2-Country.mmdb")
	if err != nil {
		t.Error(err)
	}

	if db == nil {
		t.Fail()
	}

	code, err := CountryCode("61.135.169.125")
	if err != nil {
		t.Error(err)
	}
	if code != "CN" {
		t.Fail()
	}
	CloseDB()
}
