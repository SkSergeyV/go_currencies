package home

import (
	"context"

	"currencies/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CurrenciesRepository struct {
	Dbpool *pgxpool.Pool
}

func NewCurrenciesRepository(dbpool *pgxpool.Pool) *CurrenciesRepository {
	return &CurrenciesRepository{
		Dbpool: dbpool,
	}
}

func (r *CurrenciesRepository) getCurrencies() ([]models.Currencies, error) {

	query := "SELECT c.id, c.currency_id, c.currency_name, c.char_code, c.nominal, c.num_code, exch.unit_rate, exch.exchange_date_format AS exchange_date, exch.currency_value AS exchange_value, exch.difference FROM currencies AS c LEFT JOIN LATERAL (SELECT ex.unit_rate, ex.exchange_date_format, ex.currency_value, ex.difference FROM exchange AS ex WHERE ex.currency_id = c.currency_id ORDER BY ex.exchange_date_time DESC LIMIT 1) AS exch ON true WHERE c.type_status = 1 ORDER BY c.sort_order ASC;"

	rows, err := r.Dbpool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	currencies, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Currencies])
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

func (r *CurrenciesRepository) getUpdateExchangeDate() (string, error) {
	query := "SELECT exchange_date_format FROM exchange_upd WHERE id=1"
	var exchange_date_format string
	r.Dbpool.QueryRow(context.Background(), query).Scan(&exchange_date_format)
	return exchange_date_format, nil
}

func (r *CurrenciesRepository) GetUpdateDates() (dateFormat, dateFrom string, err error) {
	query := `
        SELECT exchange_date_format, exchange_date_from
        FROM exchange_upd
        WHERE id = 1
    `
	err = r.Dbpool.QueryRow(context.Background(), query).Scan(&dateFormat, &dateFrom)
	if err != nil {
		return "", "", err
	}

	return dateFormat, dateFrom, nil
}
