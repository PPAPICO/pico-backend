package parser

import (
	"regexp"
	"strconv"
	"strings"
)

func ParseRegionCode(codeStr, addressStr string) int {
	if strings.Contains(addressStr, "서울") {
		return 11000
	}
	if codeStr != "" {
		cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(codeStr, "")
		if val, err := strconv.Atoi(cleaned); err == nil {
			if val == 11000 || val == 11 || val == 6110000 {
				return 11000
			}
		}
	}
	return 0
}
