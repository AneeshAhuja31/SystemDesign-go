package stats

import "database/sql"

func InsertClick(db *sql.DB, code string) error {
	_, err := db.Exec(
		"INSERT INTO stats(short_code) VALUES($1)",
		code,
	)
	return err
}
