package tools

import (
	"fmt"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"github.com/nyaruka/phonenumbers"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"strings"
)

/*
   Creation Time: 2019 - Oct - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	RegExPhone    *regexp.Regexp
	RegExColour   *regexp.Regexp
	RegExUsername *regexp.Regexp
)

func init() {
	exp1, err := regexp.Compile("^[+\\d]\\d*$")
	if err != nil {
		panic(err)
	}
	RegExPhone = exp1
	exp2, err := regexp.Compile("^#[0-9A-F]{6}$")
	if err != nil {
		panic(err)
	}
	RegExColour = exp2
	exp3, err := regexp.Compile("^[a-zA-Z][\\da-zA-Z]{4,31}$")
	if err != nil {
		panic(err)
	}
	RegExUsername = exp3

}

func SanitizePhone(phoneNumber string, defaultRegion string, testMode bool) string {
	phoneNumber = strings.TrimLeft(phoneNumber, "0+ ")
	if testMode && strings.HasPrefix(phoneNumber, "237400") {
		return phoneNumber
	}
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

func Trim52(num int64) int64 {
	if num > 4503599627370496 {
		return num >> 12
	}
	if num < -4503599627370496 {
		return num >> 12
	}
	return num
}

func TrimU52(num uint64) uint64 {
	if num > 4503599627370496 {
		return num >> 12
	}
	return num
}

func DeleteItemFromArray(slice interface{}, index int) {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	vLength := v.Len()
	if v.Kind() != reflect.Slice {
		panic("slice is not valid")
	}
	if index >= vLength || index < 0 {
		panic("invalid index")
	}
	switch vLength {
	case 1:
		v.SetLen(0)
	default:
		v.Index(index).Set(v.Index(v.Len() - 1))
		v.SetLen(vLength - 1)
	}
}

func TruncateString(str string, l int) string {
	n := 0
	for i := range str {
		n++
		if n > l {
			return str[:i]
		}
	}
	return str
}
