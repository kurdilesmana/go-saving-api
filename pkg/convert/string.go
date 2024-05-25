package convert

import (
	"fmt"
	"strings"
	"unicode"
)

func ReplaceSpaceAndDot(str string, replacement rune) string {
	modifiedStr := strings.Map(func(r rune) rune {
		if r == ' ' || r == '.' {
			return replacement
		}
		return unicode.ToLower(r) // Mengubah karakter menjadi huruf kecil
	}, str)
	return modifiedStr
}

func NormalizePhoneNumber(phoneNumber string) string {
	strNo := ""
	for _, c := range phoneNumber {
		if unicode.IsDigit(c) {
			strNo += string(c)
		}
	}

	if string(strNo[0]) == "0" {
		strNo = fmt.Sprintf("62%s", strNo[1:])
	}

	return strNo
}
