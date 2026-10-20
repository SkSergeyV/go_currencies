

CREATE TABLE exchange (
	id SERIAL PRIMARY KEY,
	currency_id VARCHAR(255) NOT NULL,
	currency_value VARCHAR(255),
	unit_rate VARCHAR(255),
    exchange_date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    exchange_date_format VARCHAR(255)
)