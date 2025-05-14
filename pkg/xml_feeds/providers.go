package xmlfeeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/rs/zerolog/log"
)

func GetNutrendFeed() (*Shop, error) {
	resp, err := http.Get("https://api.nutrend.eu/feed/heureka")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var shop Shop
	if err := xml.Unmarshal(data, &shop); err != nil {
		return nil, err
	}
	return &shop, nil
}

func GetDafitFeed() (*Shop, error) {
	resp, err := http.Get("https://www.dafit.cz/cron/xmlfeedvo.php?vomail=marie.vaneckova@maximal-limit.cz")
	if err != nil {
		return nil, fmt.Errorf("error triggering XML: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`https://www.dafit.cz/xml/[^"]+\.xml`)
	xmlURL := re.FindString(string(body))

	log.Info().Msgf("XML URL: %s", xmlURL)

	if xmlURL == "" {
		return nil, fmt.Errorf("XML URL nebyla nalezena")
	}

	resp, err = http.Get(xmlURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var shop Shop
	if err := xml.Unmarshal(data, &shop); err != nil {
		return nil, err
	}
	return &shop, nil
}
