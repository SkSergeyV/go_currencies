SELECT c.currency_id, c.currency_name, c.char_code, c.nominal, c.num_code, 
(SELECT SUBSTRING(exch.currency_value FROM 1 FOR CHAR_LENGTH(exch.currency_value) - 2) FROM exchange AS exch WHERE exch.currency_id=c.currency_id 
ORDER BY exch.exchange_date_time DESC LIMIT 1) AS exchange_value, 
(SELECT exch.exchange_date_format FROM exchange AS exch WHERE exch.currency_id=c.currency_id ORDER BY exch.exchange_date_time DESC LIMIT 1) AS exchange_date,
(SELECT exch.unit_rate FROM exchange AS exch WHERE exch.currency_id=c.currency_id ORDER BY exch.exchange_date_time DESC LIMIT 1) AS unit_rate  
FROM currencies AS c WHERE c.type_status=1 ORDER BY c.sort_order ASC

SELECT
    c.id,
    c.currency_id,
    c.currency_name,
    c.char_code,
    c.nominal,
    c.num_code,
    exch.unit_rate,
    exch.exchange_date_format AS exchange_date,
    exch.currency_value AS exchange_value
FROM currencies AS c
LEFT JOIN LATERAL (
    SELECT ex.unit_rate, ex.exchange_date_format, ex.currency_value
    FROM exchange AS ex
    WHERE ex.currency_id = c.currency_id
    ORDER BY ex.exchange_date_time DESC
    LIMIT 1
) AS exch ON true
WHERE c.type_status = 1
ORDER BY c.sort_order ASC;


SELECT
    c.id,
    c.currency_id,
    c.currency_name,
    c.char_code,
    c.nominal,
    c.num_code,
    exch.unit_rate,
    exch.exchange_date_format AS exchange_date,
    left(exch.currency_value, length(exch.currency_value) - 2) AS exchange_value
FROM currencies AS c
LEFT JOIN LATERAL (
    SELECT ex.unit_rate, ex.exchange_date_format, ex.currency_value
    FROM exchange AS ex
    WHERE ex.currency_id = c.currency_id
    ORDER BY ex.exchange_date_time DESC
    LIMIT 1
) AS exch ON true
WHERE c.type_status = 1
ORDER BY c.sort_order ASC;