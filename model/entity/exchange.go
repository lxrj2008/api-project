package entity

import "database/sql"

// Exchange mirrors the TExchange table schema.
type Exchange struct {
	MQMExchangeCode    string         `db:"mqm_exchange_code"`
	ClearExchangeCode  string         `db:"clear_exchange_code"`
	GlobexExchangeCode string         `db:"globex_exchange_code"`
	Description        sql.NullString `db:"description"`
	SegType            string         `db:"seg_type"`
}
