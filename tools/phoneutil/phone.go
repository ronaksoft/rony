package phoneutil

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

/*
   Creation Time: 2019 - Oct - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	RegExPhone *regexp.Regexp
)

func init() {
	exp1, err := regexp.Compile(`^[+\d]\d*$`)
	if err != nil {
		panic(err)
	}
	RegExPhone = exp1
}

func SanitizePhone(phoneNumber string, defaultRegion string) (string, error) {
	if defaultRegion == "" {
		defaultRegion = GetCountryCode(phoneNumber)
	}
	phoneNumber = strings.TrimLeft(phoneNumber, "0+ ")
	if !RegExPhone.MatchString(phoneNumber) {
		return "", fmt.Errorf("did not match phone number regex")
	}
	phone, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		phoneNumber = fmt.Sprintf("+%s", phoneNumber)
		phone, err = phonenumbers.Parse(phoneNumber, defaultRegion)
		if err != nil {
			return "", err
		}
	}

	if !phonenumbers.IsValidNumberForRegion(phone, phonenumbers.GetRegionCodeForNumber(phone)) {
		return "", fmt.Errorf("not valid for region")
	}

	return fmt.Sprintf("%d%d", phone.GetCountryCode(), phone.GetNationalNumber()), nil
}

func GetCountryCode(phone string) string {
	phone = fmt.Sprintf("+%s", strings.TrimLeft(phone, "+"))
	ph, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return ""
	}

	return phonenumbers.GetRegionCodeForNumber(ph)
}
