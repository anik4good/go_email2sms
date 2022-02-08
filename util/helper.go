package util

import (
	"log"
	"regexp"
	"strings"

	gonanoid "github.com/matoous/go-nanoid"
)

func GetTransactionID(prefix string) string {
	id, _ := gonanoid.Nanoid(13)
	return prefix + "-" + id
}

func GetValidPhoneNumber(number string) (string, bool) {
	reg := regexp.MustCompilePOSIX("/[^0-9]/")
	msisdn := reg.ReplaceAllString(number, "")

	if len(msisdn) == 11 {
		if strings.HasPrefix(msisdn, "01") {
			log.Println("01")

			return "88" + msisdn, true
		}

		return msisdn, false

	}

	if len(msisdn) == 13 {
		if strings.HasPrefix(msisdn, "8801") {

			return msisdn, true
		}

	}

	return msisdn, false

}

func GetValidPhoneNumberUpdated(number string) bool {
	reg := regexp.MustCompilePOSIX("/[^0-9]/")
	msisdn := reg.ReplaceAllString(number, "")

	if len(msisdn) == 11 {
		if strings.HasPrefix(msisdn, "01") {
			log.Println("01")

			return true
		}

		return false

	}

	if len(msisdn) == 13 {
		if strings.HasPrefix(msisdn, "8801") {

			return true
		}

	}

	if len(msisdn) == 14 {
		if strings.HasPrefix(msisdn, "+8801") {

			return true
		}

	}

	return false

}
