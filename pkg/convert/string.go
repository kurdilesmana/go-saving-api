package convert

import (
	"fmt"
	"unicode"
)

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
