package calendar

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: Test this, tener en cuenta que la test db va a haber que randomizarl el
// nombre o algo para que no se pisen entre paquetes
// TODO: Unify event and description words
func ForEach(
	ctx context.Context, conn *pgxpool.Conn, cb func(date, event string),
) error {
	rows, err := conn.Query(ctx, `SELECT date, event FROM calendar;`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var date, event string

		if err := rows.Scan(&date, &event); err != nil {
			return err
		}

		cb(date, event)
	}

	return nil
}
