package home

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
