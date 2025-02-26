package utils

import (
	"math"
	"strings"
	"time"
)

func GenerateShortCode(name string) string {
	code := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	if len(code) > 12 {
		code = code[:12]
	}
	return code
}

func GetFirstDayOfActualYear() string {
	now := time.Now()
	firstDay := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	return firstDay.Format("2006-01-02")
}

func CalculateTotalPriceWithVat(price float64, vat float64) float64 {
	return math.Ceil(price + (price * vat / 100))
}

func GetCurrentDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}
