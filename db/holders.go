package db

import (
	"context"
	"database/sql"
)

type Holder struct {
	Address string
	Count   int64
}

func GetHolders(db *sql.DB, policyID string) ([]Holder, error) {
	query := `
	SELECT address.address, SUM(mto.quantity)::BIGINT
	FROM ma_tx_out mto
	JOIN tx_out txo ON mto.tx_out_id = txo.id
	JOIN address ON txo.address_id = address.id
	JOIN multi_asset ma ON mto.ident = ma.id
	WHERE ma.policy = $1
	  AND txo.consumed_by_tx_id IS NULL
	  AND mto.quantity = 1
	  AND address.has_script = FALSE
	GROUP BY address.address;
	`
	rows, err := db.QueryContext(context.Background(), query, policyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holders []Holder
	for rows.Next() {
		var h Holder
		if err := rows.Scan(&h.Address, &h.Count); err != nil {
			return nil, err
		}
		holders = append(holders, h)
	}
	return holders, nil
}
