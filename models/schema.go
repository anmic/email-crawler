package models

import "gopkg.in/pg.v4"

func CreateSchema(db *pg.DB) error {
	queries := []string{
		`CREATE SCHEMA IF NOT EXISTS public;`,
		`CREATE TABLE IF NOT EXISTS leads (id serial, email varchar(500) UNIQUE, created_at timestamptz)`,
		`CREATE TABLE IF NOT EXISTS tokens (id serial, access_token varchar(500), token_type varchar(500),
		refresh_token varchar(500), expiry timestamptz)`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func DropSchema(db *pg.DB) error {
	_, err := db.Exec("DROP SCHEMA IF EXISTS public CASCADE")
	if err != nil {
		return err
	}
	return nil
}
