package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/oresoftware/sqlboiler-examples/v1/model"
	"github.com/volatiletech/sqlboiler/boil"
	"log"
	_ "github.com/lib/pq"
)

func insert(db *sql.DB) {

	var m = model.UserMapTable{

	}
	var err error
	m.UserID = 1
	m.ValueType = "string"
	m.Key = "foo"

	m.Value, err = json.Marshal(struct{ Val string }{"star"})
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Insert(context.Background(), tx, boil.Infer()); err != nil {
		log.Fatal("Could not insert model:", err)
	}

	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}

}

func mustGetProp(mp map[string]interface{}, key string, f1 func(k string, v interface{}) bool, f2 func(v interface{})) {
	if val, ok := mp[key]; !ok {
		log.Fatal("missing user id")
	} else {
		if f1(key, val) {
			f2(val)
		}
	}
}

func getProp(mp map[string]string, key string, f func(v interface{})) {
	if val, ok := mp["value_type"]; ok {
		f(val)
	}
}

func getString(mp map[string]string, mpKey string, dbField string) string {
	if val, ok := mp[mpKey]; ok {
		return fmt.Sprintf(" %s = '%s'", dbField, val)
	}
	return ""
}

func update(db *sql.DB) {

	var mp = map[string]string{}

	mp["user_id"] = "1"
	mp["value"] = `{"bob":3}`
	mp["key"] = "couple3"
	mp["value_type"] = "string"
	mp["added"] = "false"

	var userId = mp["user_id"]

	if userId == "" {
		log.Fatal("Missing user_id")
	}

	var setString = ""

	for p, v := range []string{"user_id", "value", "key", "value_type", "added"} {
		if p < 4 {
			setString = setString + getString(mp, v, v) + ","
		} else {
			setString = setString + getString(mp, v, v)
		}
	}

	log.Println("set string:", setString)

	//var err error

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf("UPDATE user_map_table SET %s WHERE user_id = %s ;", setString, userId)
	stmt, err := tx.Prepare(query)

	if err != nil {
		log.Fatalln(
			err,
			"Had to rollback due to error, if there was a rollback error, this is it:",
			tx.Rollback(),
		)
	}

	defer stmt.Close()

	rows, err := stmt.Exec()

	if err != nil {
		log.Fatalln(
			err,
			"Had to rollback due to error, if there was a rollback error, this is it:",
			tx.Rollback(),
		)
		return;
	}

	//log.Println("raw query:", query)
	//rows, err := db.Query(query)

	log.Println("rows:", rows)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres")
	if err != nil {
		log.Fatal(err)
	}

	insert(db)
	update(db)

}
