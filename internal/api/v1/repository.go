package v1

import (
	"context"
	"currencies/internal/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type V1Repository struct {
	Dbpool *pgxpool.Pool
}

func NewV1Repository(dbpool *pgxpool.Pool) *V1Repository {
	return &V1Repository{
		Dbpool: dbpool,
	}
}

func (r *V1Repository) GetExportData(currency string) ([]models.ExportCurrencies, error) {

	var query string
	var args []any
	currencyUpper := strings.ToUpper(currency)

	if currencyUpper == "ALL" {

		query = `
			SELECT c.currency_id, c.currency_name, c.char_code, c.num_code, c.nominal,
				exch.currency_value AS exchange_value,
				exch.prev_currency_value AS prev_exchange_value,
				exch.difference
			FROM currencies AS c
			LEFT JOIN LATERAL (
				SELECT ex.currency_value, ex.prev_currency_value, ex.difference
				FROM exchange AS ex
				WHERE ex.currency_id = c.currency_id
				ORDER BY ex.exchange_date_time DESC
				LIMIT 1
			) AS exch ON true
			WHERE c.type_status = 1
			ORDER BY c.sort_order ASC;
		`
	} else {

		query = `
			SELECT c.currency_id, c.currency_name, c.char_code, c.num_code, c.nominal,
				exch.currency_value AS exchange_value,
				exch.prev_currency_value AS prev_exchange_value,
				exch.difference
			FROM currencies AS c
			LEFT JOIN LATERAL (
				SELECT ex.currency_value, ex.prev_currency_value, ex.difference
				FROM exchange AS ex
				WHERE ex.currency_id = c.currency_id
				ORDER BY ex.exchange_date_time DESC
				LIMIT 1
			) AS exch ON true
			WHERE c.type_status = 1
			  AND c.char_code = $1
			ORDER BY c.sort_order ASC;
		`
		args = []any{currencyUpper}
	}

	rows, err := r.Dbpool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exportData, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ExportCurrencies])
	if err != nil {
		return nil, err
	}

	return exportData, nil
}

func (r *V1Repository) GetUpdateExportDates() (models.ExportDates, error) {
	var d models.ExportDates

	query := `
		SELECT exchange_date_format, exchange_date_from
		FROM exchange_upd
		WHERE id = 1
	`

	err := r.Dbpool.QueryRow(context.Background(), query).Scan(
		&d.ExchangeDate,
		&d.ExchangeFrom,
	)
	if err != nil {
		return models.ExportDates{}, err
	}

	return d, nil
}

func (r *V1Repository) GetExportDataDates(date string) ([]models.ExportCurrencies, error) {
	var query string
	var args []any

	baseQuery := `
		SELECT c.currency_id, c.currency_name, c.char_code, c.num_code, c.nominal,
			exch.currency_value AS exchange_value,
			exch.prev_currency_value AS prev_exchange_value,
			exch.difference
		FROM currencies AS c
		LEFT JOIN LATERAL (
			SELECT ex.currency_value, ex.prev_currency_value, ex.difference
			FROM exchange AS ex
			WHERE ex.currency_id = c.currency_id
	`

	if date != "" {
		parsedDate, err := time.Parse("02.01.2006", date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, expected DD.MM.YYYY: %w", err)
		}

		startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1)

		baseQuery += `
				AND ex.exchange_date_time >= $1
				AND ex.exchange_date_time < $2
		`
		args = []any{startOfDay, endOfDay}
	}

	query = baseQuery + `
			ORDER BY ex.exchange_date_time DESC
			LIMIT 1
		) AS exch ON true
		WHERE c.type_status = 1
		ORDER BY c.sort_order ASC;
	`

	rows, err := r.Dbpool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exportData, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ExportCurrencies])
	if err != nil {
		return nil, err
	}

	return exportData, nil
}

func (r *V1Repository) GetUpdateExportDatesPlug(date string) (models.ExportDates, error) {
	var d models.ExportDates

	query := `
     SELECT '` + date + `' as exchange_date_format, '` + date + `' as exchange_date_from
	`

	err := r.Dbpool.QueryRow(context.Background(), query).Scan(
		&d.ExchangeDate,
		&d.ExchangeFrom,
	)
	if err != nil {
		return models.ExportDates{}, err
	}

	return d, nil
}

func (r *V1Repository) GetActiveByValue(ctx context.Context, keyValue string) (*models.AccessKey, error) {
	var row pgx.Row
	row = r.Dbpool.QueryRow(ctx, `
        SELECT id, key_name, key_value, key_desc, created_at, expires_at, is_active
        FROM api_keys
        WHERE key_value = $1 AND is_active = TRUE AND (expires_at IS NULL OR expires_at > NOW())
        LIMIT 1
    `, strings.TrimSpace(keyValue))

	var key models.AccessKey
	err := row.Scan(&key.ID,
		&key.KeyName,
		&key.KeyValue,
		&key.KeyDesc,
		&key.CreatedAt,
		&key.ExpiresAt,
		&key.IsActive)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return nil, models.ErrKeyNotFound
		}
		return nil, err
	}
	return &key, nil
}
