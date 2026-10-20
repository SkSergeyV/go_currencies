CREATE TABLE exchange_upd (
	id SERIAL PRIMARY KEY,
    exchange_date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    exchange_date_format VARCHAR(255),
    exchange_date_from VARCHAR(255)
)