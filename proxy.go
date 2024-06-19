package httpv3

import (
	"fmt"
	"strings"
)

func parseSingleProxy(rawProxy string) (OptionProxy, error) {
	split := strings.Split(rawProxy, ":")

	if len(split) != 2 && len(split) != 4 {
		return "", ErrProxyFormatCorrupted
	}

	if len(split) == 2 {
		return OptionProxy(fmt.Sprintf("http://%s:%s", split[0], split[1])), nil
	}

	return OptionProxy(fmt.Sprintf("http://%s:%s@%s:%s", split[2], split[3], split[0], split[1])), nil
}

func parseList(list []string) []OptionProxy {
	var parsedList []OptionProxy

	for _, rawProxy := range list {
		if parsed, err := parseSingleProxy(rawProxy); err == nil {
			parsedList = append(parsedList, parsed)
		}
	}

	return parsedList
}
