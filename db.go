package main

import (
	"log"
	"sort"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type Result struct {
	RowID          int    `db:"ID"`
	Raw            []byte `db:"Raw" json:"-"`
	ID             string `db:"ItemID"`
	Name           string `db:"Name"`
	Cost           uint32 `db:"Cost"`
	Hearts         uint32 `db:"Hearts"`
	Duration       uint32 `db:"Duration"`
	Effect         uint32 `db:"Effect"`
	EffectString   string `db:"EffectStr"`
	Strength       uint32 `db:"Strength"`
	StrengthString string `db:"StrengthStr"`
	Ingredients    string `db:"Ingredients"`
}

func init() {
	var err error
	db, err = sqlx.Open("mssql", "server=127.0.0.1;database=botw;connection timeout=30")
	e(err)
	e(db.Ping())
}

func haveRecipe(ings string) bool {
	q := "SELECT COUNT(1) FROM Results WHERE Ingredients = ?"
	count := 0
	err := db.Get(&count, q, ings)
	e(err)
	return count > 0
}

func putRecipe(r Result) {
	q := `INSERT INTO Results
(ItemID, Name, Cost, Hearts, Duration, Effect, EffectStr, Strength, StrengthStr, Ingredients, Raw) VALUES
(     ?,    ?,    ?,      ?,        ?,      ?,         ?,        ?,           ?,           ?,   ?)`
	_, err := db.Exec(q, r.ID, r.Name, r.Cost, r.Hearts, r.Duration,
		r.Effect, r.EffectString, r.Strength, r.StrengthString, r.Ingredients, r.Raw)
	e(err)
}

func getAll() (r []*Result, err error) {
	err = db.Select(&r, "SELECT * FROM Results ORDER BY ID")
	return
}

func fixData() {
	rows, err := getAll()
	e(err)
	for _, r := range rows {
		strs := strings.Split(r.Ingredients, ",")
		sort.Strings(strs)
		sorted := strings.Join(strs, ",")
		if sorted != r.Ingredients {
			log.Printf("Resorting row %d", r.RowID)
			_, err = db.Exec("UPDATE Results SET Ingredients = ? WHERE ID = ?", sorted, r.RowID)
			e(err)
		}
	}
}
