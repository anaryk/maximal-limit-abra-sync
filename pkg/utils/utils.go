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

// Methon for extract 2025-01-01 from time format 2025-01-16 22:03:23
func ExtractDate(date string) string {
	return date[:10]
}

// Extract EAN code from SumUP product description in format :EAN:12345569:EAN:
func ExtractEANCode(description string) string {
	if strings.Contains(description, ":EAN:") {
		start := strings.Index(description, ":EAN:") + len(":EAN:")
		end := strings.Index(description[start:], ":EAN:")
		if end != -1 {
			return description[start : start+end]
		}
	}
	return ""
}

// Extract BaseCode from Nutrent item id (example VT-023-250-VI -> VT-023-250)
func ExtractBaseCode(itemID string) string {
	if strings.Contains(itemID, "-") {
		start := strings.Index(itemID, "-") + len("-")
		end := strings.Index(itemID[start:], "-")
		if end != -1 {
			return itemID[:start+end]
		}
	}
	return ""
}
