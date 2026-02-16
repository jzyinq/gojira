package gojira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var countryCodeRegex = regexp.MustCompile("(?i)_([a-z]{2})")

type Holiday struct {
	Date        string `json:"date"`
	LocalName   string `json:"localName"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type Holidays []Holiday

func (h *Holidays) GetHolidaysForMonth(month time.Month) Holidays {
	var holidaysForMonth Holidays
	for _, holiday := range *h {
		t, err := time.Parse("2006-01-02", holiday.Date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if t.Month() == month {
			holidaysForMonth = append(holidaysForMonth, holiday)
		}
	}
	return holidaysForMonth
}

func (h *Holidays) IsHoliday(t *time.Time) bool {
	for _, holiday := range *h {
		t2, err := time.Parse("2006-01-02", holiday.Date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if t2.Equal(*t) {
			return true
		}
	}
	return false
}

func (h *Holiday) GetTime() (*time.Time, error) {
	t, err := time.Parse("2006-01-02", h.Date)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func getHolidaysForCountry(countryCode string) (*Holidays, error) {
	currentYear := time.Now().Year()
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/%s", currentYear, countryCode)
	logrus.Infof("fetching holidays from url: %s", url)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		logrus.Error("could not fetch holidays from url: ", url)
		return nil, err
	}
	defer resp.Body.Close()
	var holidays Holidays
	err = json.NewDecoder(resp.Body).Decode(&holidays)
	if err != nil {
		return nil, err
	}
	return &holidays, nil
}

func NewHolidays(countryCode string) (*Holidays, error) {
	holidays, err := getHolidaysForCountry(countryCode)
	if err != nil {
		logrus.Error(err)
		return &Holidays{}, err
	}

	return holidays, nil
}

func GetCountryFromLCTime(timeString string) (string, error) {
	match := countryCodeRegex.FindStringSubmatch(timeString)
	if match == nil {
		return "", fmt.Errorf("could not parse country from LC_TIME")
	}
	return strings.ToUpper(match[1]), nil
}
