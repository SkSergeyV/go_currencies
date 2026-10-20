<p align="center">
</p>


## Курсы валют
GO: go 1.27.0 <br>
<br>

JSON
/api/v1/exchange_rates/current/json?currencies=all&key=abc  Список курсов всех валют на текущий день

/api/v1/exchange_rates/current/json?currencies=usd&key=abc  Курс одной валюты на текущий день

/api/v1/exchange_rates/date/json?date=08.10.2026&key=abc  Список курсов всех валют на дату

<br>


XML
/api/v1/exchange_rates/xml?currencies=all&key=abc  Список курсов всех валют на текущий день

/api/v1/exchange_rates/xml?currencies=usd&key=abc  Курс одной валюты на текущий день

/api/v1/exchange_rates/date/xml?date=08.10.2026&key=abc  Список курсов всех валют на дату

<br>


DEMO (без ключа) (JSON и XML) Ограничение 1 запрос в минуту

/api/demo/v1/exchange_rates/current/json?currencies=all

/api/demo/v1/exchange_rates/current/xml?currencies=all




