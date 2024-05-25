package convert

import (
	"strings"
)

// convert query params sort (string) to map[string]string,
//	@var	params ex = name,-created_at. field separated by (,) and add (-) for DESC
//	@return	{"name": asc, "created_at": desc}
func SortToMap(paramsSort string) map[string]string {
	res := make(map[string]string)
	sp := strings.Split(paramsSort, ",")

	if sp[0] != "" {
		for _, p := range sp {
			p = strings.TrimSpace(p)
			if p[0] == '-' {
				res[p[1:]] = "DESC"
			} else {
				res[p] = "ASC"
			}
		}
	}

	return res
}
