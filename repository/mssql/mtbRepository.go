package mssql

import (
	"context"
	"database/sql"
	"time"
	"tnm-malawi/connectors/callback/models"
)

const (
	MbtQuery = `select MBTtransID from africainv.dbo.MBT where MBTid = @mbtId`
)

func (r *repository) GetMbt(mbtId string) (*models.MbtEntity, error) {
	entitie := new(models.MbtEntity)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, MbtQuery,
		sql.NamedArg{Name: "mbtId", Value: mbtId}).Scan(
		&entitie.Mbt)
	if err != nil {
		return nil, err
	}

	return entitie, nil
}
