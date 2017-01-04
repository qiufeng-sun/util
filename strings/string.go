package strings

import (
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
//
func Split(s, sep string, funcPre func(src string) string) []string {
	//
	if funcPre != nil && s != "" {
		s = funcPre(s)
	}

	if "" == s {
		return nil
	}

	return strings.Split(s, sep)
}
