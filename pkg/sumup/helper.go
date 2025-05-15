package sumup

import "regexp"

// ExtractEANFromProductName extracts the EAN from the product name... it will locate the EAN at the end like EAN:1234567890123 and return everything from EAN: to end of string
func ExtractEANFromProductName(productName string) string {
	re := regexp.MustCompile(`EAN:.*$`)
	match := re.FindString(productName)
	return match
}
