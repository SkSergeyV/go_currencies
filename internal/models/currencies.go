package models

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"time"
)

type Currencies struct {
	Id            int    `db:"id"`
	CurrencyId    string `db:"currency_id"`
	CurrencyName  string `db:"currency_name"`
	CharCode      string `db:"char_code"`
	UnitRate      string `db:"unit_rate"`
	ExchangeDate  string `db:"exchange_date"`
	Nominal       int    `db:"nominal"`
	NumCode       string `db:"num_code"`
	ExchangeValue string `db:"exchange_value"`
	Difference    string `db:"difference"`
}

type ExportCurrencies struct {
	CurrencyId        string `db:"currency_id" json:"currency_id" xml:"currency_id"`
	CurrencyName      string `db:"currency_name" json:"currency_name" xml:"currency_name"`
	CharCode          string `db:"char_code" json:"char_code" xml:"char_code"`
	NumCode           string `db:"num_code" json:"num_code" xml:"num_code"`
	Nominal           int    `db:"nominal" json:"nominal" xml:"nominal"`
	ExchangeValue     string `db:"exchange_value" json:"exchange_value"  xml:"exchange_value"`
	PrevExchangeValue string `db:"prev_exchange_value" json:"prev_exchange_value" xml:"prev_exchange_value"`
	Difference        string `db:"difference" json:"difference" xml:"difference"`
}

type ExportDates struct {
	ExchangeDate string `db:"exchange_date_format" json:"exchange_date" xml:"exchange_date"`
	ExchangeFrom string `db:"exchange_date_from" json:"exchange_from" xml:"exchange_from"`
}

type ExportResponse struct {
	XMLName    xml.Name     `json:"-" xml:"exchange_rate"`
	Dates      ExportDates  `json:"dates" xml:"dates"`
	Currencies CurrencyList `json:"currencies" xml:"currencies,omitempty"`
}

type CurrencyList []ExportCurrencies

func (l CurrencyList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, c := range l {
		se := xml.StartElement{Name: xml.Name{Local: "currency"}}
		if err := e.EncodeElement(c, se); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

type AccessKey struct {
	ID        int
	KeyName   string
	KeyValue  string
	KeyDesc   string
	CreatedAt time.Time
	ExpiresAt sql.NullTime
	IsActive  bool
}

var ErrKeyNotFound = errors.New("key not found or expired")
