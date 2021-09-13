package common

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseNumber(message string) (float64, error) {
	var kilo = strings.Count(message, "k")
	message = strings.Replace(message, "k", "", -1)
	message = strings.Replace(message, ",", ".", -1)
	num, err := strconv.ParseFloat(message, 64)
	if err != nil {
		return num, fmt.Errorf("🚫 Couldn't parse <b>%v</b> to number: %v", message, err)
	}
	for i := 0; i < kilo; i++ {
		num *= 1000
	}
	return num, err
}
