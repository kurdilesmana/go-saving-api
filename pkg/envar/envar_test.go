package envar

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		Key    string
		Def    any
		Set    string
		Expect any
	}{
		{
			Key:    "KEY_TEST_INT1",
			Def:    1,
			Set:    "10",
			Expect: 10,
		},
		{
			Key:    "KEY_TEST_INT2",
			Def:    123456789,
			Set:    "",
			Expect: 123456789,
		},
		{
			Key:    "KEY_TEST_INT3",
			Def:    0,
			Set:    "10240985761",
			Expect: 10240985761,
		},
		{
			Key:    "KEY_TEST_BOOL1",
			Def:    false,
			Set:    "True",
			Expect: true,
		},
		{
			Key:    "KEY_TEST_BOOL2",
			Def:    true,
			Set:    "FalSe",
			Expect: false,
		},
		{
			Key:    "KEY_TEST_BOOL3",
			Def:    false,
			Set:    "TRUE",
			Expect: true,
		},
		{
			Key:    "KEY_TEST_BOOL4",
			Def:    false,
			Set:    "",
			Expect: false,
		},
		{
			Key:    "KEY_TEST_STRING1",
			Def:    "",
			Set:    `http://localhost:5000/kk-ddd/wwwdd_1234/rrghj/45rghh\\/.;\][]p[p[\`,
			Expect: `http://localhost:5000/kk-ddd/wwwdd_1234/rrghj/45rghh\\/.;\][]p[p[\`,
		},
		{
			Key:    "KEY_TEST_STRING2",
			Def:    "secret_key",
			Set:    "m5k^^^-++*&`2&!Z(+En{S+FQw@#",
			Expect: "m5k^^^-++*&`2&!Z(+En{S+FQw@#",
		},
		{
			Key:    "KEY_TEST_STRING3",
			Def:    "default-value",
			Set:    "",
			Expect: "default-value",
		},
	}

	for _, test := range tests {
		if strings.TrimSpace(test.Set) != "" {
			os.Setenv(test.Key, test.Set)
		}

		switch test.Def.(type) {
		case int:
			result := GetEnv(test.Key, test.Def.(int))
			assert.Equal(t, test.Expect.(int), result, fmt.Sprintf("%s value not as expected", test.Key))

		case bool:
			result := GetEnv(test.Key, test.Def.(bool))
			assert.Equal(t, test.Expect.(bool), result, fmt.Sprintf("%s value not as expected", test.Key))

		case string:
			result := GetEnv(test.Key, test.Def.(string))
			assert.Equal(t, test.Expect.(string), result, fmt.Sprintf("%s value not as expected", test.Key))

		}
	}
}
