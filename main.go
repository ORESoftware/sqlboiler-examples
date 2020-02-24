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
	"os"
	"regexp"
	"strings"
	q "github.com/volatiletech/sqlboiler/queries/qm"
)

type WhereCondition struct {
	LHS string
	RHS string
	op  string
}

func printRows(rows *sql.Rows){
	defer rows.Close()

	for rows.Next() {

		var a string
		var b string
		var c string
		var d string
		var e string
		var f string

		var row = []interface{}{&a, &b, &c, &d,&e, &f}


		err := rows.Scan(row...)
		if err != nil {
			log.Fatal(err)
		}
		r, err := json.Marshal(row)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("row:", string(r));
	}
}

func printResult(r sql.Result){
	log.Println("result:", r)
}

func getCleanJoin(db *sql.DB) {

	var limit = 500
	var fields = []string{"*"}
	var mappedFields = []string{}

	var space = regexp.MustCompile(`\s+`)

	for _, v := range fields {
		// replace all white space in field name, for sanitation
		str := space.ReplaceAllString(v, "")
		mappedFields = append(mappedFields, str)
	}

	if len(mappedFields) < 1 {
		mappedFields = []string{"*"};
	}

	//rows, err := model.UserTables(
	//	q.Select(mappedFields...),
	//	q.Limit(limit),
	//	q.Where("handle = ?", "dog")).All(context.Background(), db)

	rows, err := model.NewQuery(
		q.Select(mappedFields...),
		q.Limit(limit),
		q.From("user_table u"),
		q.InnerJoin("user_map_table t on u.id = t.user_id"),
		q.Where("handle = ?", "dog")).Query(db)

	if err != nil {
		log.Fatal(err)
	}

	printRows(rows)
}

func getClean(db *sql.DB) {

	var limit = 500
	var fields = []string{"*"}
	var mappedFields = []string{}

	var space = regexp.MustCompile(`\s+`)

	for _, v := range fields {
		// replace all white space in field name, for sanitation
		str := space.ReplaceAllString(v, "")
		mappedFields = append(mappedFields, str)
	}

	if len(mappedFields) < 1 {
		mappedFields = []string{"*"};
	}

	//rows, err := model.UserTables(
	//	q.Select(mappedFields...),
	//	q.Limit(limit),
	//	q.Where("handle = ?", "dog")).All(context.Background(), db)

	rows, err := model.NewQuery(
		q.Select(mappedFields...),
		q.Limit(limit),
		q.From("user_table"),
		q.Where("handle = ?", "dog")).Query(db)

	if err != nil {
		log.Fatal(err)
	}

	printRows(rows)
}

func get(db *sql.DB) {

	var limit = 500
	var fields = []string{"*"}
	var whereClauses = []WhereCondition{}

	var tableName = "user_table"
	var mappedFields = []string{}

	var space = regexp.MustCompile(`\s+`)

	for _, v := range fields {
		// replace all white space in field name, for sanitation
		str := space.ReplaceAllString(v, "")
		mappedFields = append(mappedFields, str)
	}

	var allFields = strings.Join(mappedFields, ",")

	var whereClause = ""
	for _, v := range whereClauses {
		whereClause = whereClause + fmt.Sprintf("'%s' %s '%s'", v.LHS, v.op, v.RHS)
	}

	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s LIMIT %v", allFields, tableName, whereClause, limit)
	log.Println("raw query:", query)
	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	printRows(rows)
}

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

	result, err := stmt.Exec()

	if err != nil {
		log.Fatalln(
			err,
			"Had to rollback due to error, if there was a rollback error, this is it:",
			tx.Rollback(),
		)
		return
	}

	//log.Println("raw query:", query)
	//result, err := db.Query(query)

	printResult(result);

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

	getCleanJoin(db)
	getClean(db)
	get(db)
	insert(db)
	update(db)

	log.Println("Exiting with code:", 0)
	os.Exit(0)

}
