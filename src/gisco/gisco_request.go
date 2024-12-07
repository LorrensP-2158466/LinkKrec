package gisco

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const GiscoAPIEndpoint = "https://gisco-services.ec.europa.eu/addressapi/search?"

var europeanCountryCodes = map[string]string{
	"albania":                "al",
	"andorra":                "ad",
	"austria":                "at",
	"belarus":                "by",
	"belgium":                "be",
	"bosnia and herzegovina": "ba",
	"bulgaria":               "bg",
	"croatia":                "hr",
	"czech republic":         "cz",
	"denmark":                "dk",
	"estonia":                "ee",
	"finland":                "fi",
	"france":                 "fr",
	"germany":                "de",
	"greece":                 "gr",
	"hungary":                "hu",
	"iceland":                "is",
	"ireland":                "ie",
	"italy":                  "it",
	"latvia":                 "lv",
	"liechtenstein":          "li",
	"lithuania":              "lt",
	"luxembourg":             "lu",
	"malta":                  "mt",
	"moldova":                "md",
	"monaco":                 "mc",
	"montenegro":             "me",
	"netherlands":            "nl",
	"north macedonia":        "mk",
	"norway":                 "no",
	"poland":                 "pl",
	"portugal":               "pt",
	"romania":                "ro",
	"russia":                 "ru",
	"san marino":             "sm",
	"serbia":                 "rs",
	"slovakia":               "sk",
	"slovenia":               "si",
	"spain":                  "es",
	"sweden":                 "se",
	"switzerland":            "ch",
	"ukraine":                "ua",
	"united kingdom":         "gb",
	"vatican city":           "va",
}

type Coordinates struct {
	Long float64
	Lat  float64
}

// Pretty prints with cardinal directions
func (c *Coordinates) Pretty() string {
	ns := "N"
	if c.Lat < 0 {
		ns = "S"
	}

	ew := "E"
	if c.Long < 0 {
		ew = "W"
	}

	return fmt.Sprintf(
		"%.6f°%s, %.6f°%s",
		math.Abs(c.Lat),
		ns,
		math.Abs(c.Long),
		ew,
	)
}

type GiscoResponse struct {
	Count   int `json:"count"`
	Results []struct {
		XY []float64 `json:"XY"`
	} `json:"results"`
}

func CoordinatesFromAddress(country string, city string, road *string, housenumber *string) *Coordinates {

	params := url.Values{}
	params.Add("coutnry", europeanCountryCodes[strings.ToLower(country)])
	params.Add("city", city)
	if road != nil {
		params.Add("road", *road)
	}
	if housenumber != nil {
		params.Add("housenumber", *housenumber)
	}
	endpoint := GiscoAPIEndpoint + params.Encode()
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}
	clnt := &http.Client{Transport: tr}

	res, err := clnt.Get(endpoint)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var giscoResp GiscoResponse
	if err := json.NewDecoder(res.Body).Decode(&giscoResp); err != nil {
		return nil
	}

	if giscoResp.Count == 0 || len(giscoResp.Results) == 0 {
		return nil
	}

	// take the first because why not
	xy := giscoResp.Results[0].XY
	if len(xy) != 2 {
		return nil
	}

	return &Coordinates{
		Long: xy[0],
		Lat:  xy[1],
	}

}
