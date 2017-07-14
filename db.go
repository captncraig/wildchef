package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type Result struct {
	RowID              int    `db:"ID"`
	Raw                []byte `db:"Raw"`
	ID                 string `db:"ItemID"`
	Name               string `db:"Name"`
	Cost               uint32 `db:"Cost"`
	Hearts             uint32 `db:"Hearts"`
	Duration           uint32 `db:"Duration"`
	Effect             uint32 `db:"Effect"`
	EffectString       string `db:"EffectStr"`
	Strength           uint32 `db:"Strength"`
	StrengthString     string `db:"StrengthStr"`
	IngredientsFromMem string `db:"Ingredients"`
	IngredientsIn      string `db:"IngredientsIn"`
}

func init() {
	var err error
	db, err = sqlx.Open("mssql", "server=127.0.0.1;database=botw;connection timeout=30")
	e(err)
	e(db.Ping())
}

func haveRecipe(ings string) bool {
	q := "SELECT COUNT(1) FROM Results WHERE IngredientsIn = ?"
	count := 0
	err := db.Get(&count, q, ings)
	e(err)
	return count > 0
}

func putRecipe(r Result, ingsIn string) {
	q := `INSERT INTO Results
(ItemID, Name, Cost, Hearts, Duration, Effect, EffectStr, Strength, StrengthStr, Ingredients, IngredientsIn, Raw) VALUES
(     ?,    ?,    ?,      ?,        ?,      ?,         ?,        ?,           ?,           ?,             ?,   ?)`
	_, err := db.Exec(q, r.ID, r.Name, r.Cost, r.Hearts, r.Duration,
		r.Effect, r.EffectString, r.Strength, r.StrengthString, r.IngredientsFromMem, ingsIn, r.Raw)
	e(err)
}

func getAll() (r []*Result, err error) {
	err = db.Select(&r, "SELECT * FROM Results ORDER BY ID")
	return
}
