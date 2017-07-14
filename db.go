package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

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
