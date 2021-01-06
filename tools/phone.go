package tools

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"github.com/ronaksoft/rony/internal/log"
	"go.uber.org/zap"
	"regexp"
	"strings"
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
	exp1, err := regexp.Compile("^[+\\d]\\d*$")
	if err != nil {
		panic(err)
	}
	RegExPhone = exp1
}

func SanitizePhone(phoneNumber string, defaultRegion string) string {
	phoneNumber = strings.TrimLeft(phoneNumber, "0+ ")
	if !RegExPhone.MatchString(phoneNumber) {
		return ""
	}
	phone, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		phoneNumber = fmt.Sprintf("+%s", phoneNumber)
		phone, err = phonenumbers.Parse(phoneNumber, defaultRegion)
		if err != nil {
			if ce := log.Check(log.DebugLevel, "Error On SanitizePhone"); ce != nil {
				ce.Write(
					zap.Error(err),
					zap.String("Phone", phoneNumber),
					zap.String("DefaultRegion", defaultRegion),
				)
			}
			return ""
		}
	}

	if !phonenumbers.IsValidNumberForRegion(phone, phonenumbers.GetRegionCodeForNumber(phone)) {
		return ""
	}
	return fmt.Sprintf("%d%d", phone.GetCountryCode(), phone.GetNationalNumber())
}

func GetCountryCode(phone string) string {
	phone = fmt.Sprintf("+%s", strings.TrimLeft(phone, "+"))
	ph, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return ""
	}
	return phonenumbers.GetRegionCodeForNumber(ph)
}
