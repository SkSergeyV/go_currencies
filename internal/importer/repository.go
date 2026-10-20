package importer

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImportRepository struct {
	Dbpool *pgxpool.Pool
}

func NewImportRepository(dbpool *pgxpool.Pool) *ImportRepository {
	return &ImportRepository{
		Dbpool: dbpool,
	}
}

func calcDifference(value, previous float64) string {
	diff := value - previous
	rounded := math.Round(diff*100) / 100
	s := strconv.FormatFloat(rounded, 'f', 2, 64)
	if diff > 0 {
		s = "+" + s
	}
	return s
}

func (r *ImportRepository) InsertExchangeRates(ctx context.Context, resp JsonImport) error {
	const query = `INSERT INTO exchange (
		currency_id,
		currency_value,
		unit_rate,
		exchange_date_time,
		exchange_date_format,
		prev_currency_value,
		difference
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tx, err := r.Dbpool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	batch := &pgx.Batch{}
	dateFormatted := resp.Date.Format("02.01.2006")

	for _, c := range resp.Valute {
		batch.Queue(
			query,
			c.ID,
			strconv.FormatFloat(c.Value, 'f', 2, 64),
			strconv.Itoa(c.Nominal),
			resp.Date,
			dateFormatted,
			strconv.FormatFloat(c.Previous, 'f', 2, 64),
			calcDifference(c.Value, c.Previous),
		)
	}

	br := tx.SendBatch(ctx, batch)

	for range resp.Valute {
		_, err := br.Exec()
		if err != nil {
			_ = br.Close()
			return fmt.Errorf("batch exec: %w", err)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("batch close: %w", err)
	}

	// Теперь соединение свободно — можно коммитить
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	committed = true
	return nil
}

func (r *ImportRepository) UpdateExchangeUpd(ctx context.Context, updateDate time.Time) error {
	now := time.Now()
	formattedDateTime := now.Format("2006-01-02 15:04:05.999")
	dateFormatted := now.Format("02.01.2006")
	dateFrom := updateDate.Format("02.01.2006")

	query := `
		UPDATE exchange_upd 
		SET 
			exchange_date_time = $1,
			exchange_date_format = $2,
			exchange_date_from = $3
		WHERE id = 1
	`

	_, err := r.Dbpool.Exec(ctx, query, formattedDateTime, dateFormatted, dateFrom)
	if err != nil {
		return fmt.Errorf("update exchange_upd: %w", err)
	}

	return nil
}
