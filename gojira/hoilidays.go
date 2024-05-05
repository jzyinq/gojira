package gojira

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

/* create a function that access country code and fetch the holidays for that country from
http endpoint https://date.nager.at/api/v3/PublicHolidays/2024/COUNTRY_CODE prepare a struct with parsed dates
*/

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

func getHolidaysForCountry(countryCode string) (Holidays, error) {
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/2024/%s", countryCode)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	var holidays Holidays
	err = json.NewDecoder(resp.Body).Decode(&holidays)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return holidays, nil
}

func NewHolidays(countryCode string) *Holidays {
	holidays, err := getHolidaysForCountry(countryCode)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return &holidays
}
