package util

import (
	"log"
	"regexp"
	"strings"

	Configuration "github.com/anik4good/go_email2sms/config"
	"github.com/anik4good/go_email2sms/models"
	"github.com/bxcodec/faker/v3"
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

func UserSeed() {
	// queue := new(models.Queue)

	user := new(models.User)

	for i := 0; i < 500; i++ {

		user.ID = 0
		user.Name = faker.Name()
		user.Email = faker.Email()
		user.Status = 0
		//prepare the statement
		//	stmt, _ := s.db.Prepare(`INSERT INTO users(name, email) VALUES (?,?)`)
		// execute query
		//	_, err := stmt.Exec(faker.Name(), faker.Email())

		// res, err := Configuration.GormDBConn.Raw(`INSERT INTO users(name, email,status) VALUES (?,?,?)`, faker.Name(), faker.Email(), 0)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		Configuration.GormDBConn.Create(&user)
		// Print result

	}

	log.Println("User Seeded successfully")

}
