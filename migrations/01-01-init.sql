CREATE DATABASE currencies

CREATE TABLE currencies (
	id SERIAL PRIMARY KEY,
	currency_id VARCHAR(255) NOT NULL,
	currency_name VARCHAR(255),
	char_code VARCHAR(255),
	num_code VARCHAR(255),
	nominal INTEGER DEFAULT 1,
	sort_order INTEGER DEFAULT 0, 
	type_status INTEGER DEFAULT 1
)
