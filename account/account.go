package account

import (
	"database/sql"
	"fmt"
	"time"

	"log"

	"github.com/pkg/errors"
)

type Transaction struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	Price     int       `json:"price"`
	Content   string    `json:"content"`
	Raw       string    `json:"raw"`
	Temporary bool      `json:"temporary"`
}

type dbParam struct {
	User    string
	Pass    string
	Host    string
	Port    int
	DB      string
	Options string
}

func (param *dbParam) dsn() string {
	if param == nil {
		log.Fatalln("DBParam is not initialized")
	}

	var u string
	if param.Pass == "" {
		u = param.User
	} else {
		u = param.User + ":" + param.Pass
	}

	return fmt.Sprintf("%s@tcp(%s:%d)/%s?%s", u, param.Host, param.Port, param.DB, param.Options)
}

var dbparam *dbParam

const tblTrans = "transactions"

func SetDBParam(u string, pa string, h string, po int, db string, o string) {
	dbparam = &dbParam{u, pa, h, po, db, o}
}

func AddTrans(tr Transaction) error {
	db, err := sql.Open("mysql", dbparam.dsn())
	if err != nil {
		return errors.Wrap(err, "failed to connect to db")
	}
	defer db.Close()

	q := fmt.Sprintf("INSERT INTO %s (type, time, price, content, raw, temporary) "+
		"VALUES('%s', '%s', %d, '%s', '%s', %t)", tblTrans, tr.Type,
		tr.Time.UTC().Format(time.RFC3339), tr.Price, tr.Content, tr.Raw, tr.Temporary)
	if _, err = db.Exec(q); err != nil {
		return errors.Wrap(err, "failed to insert query")
	}

	return nil
}

func QueryTrans(q string) ([]Transaction, error) {
	// FIXME: MySQLのコネクションは毎度貼り直すのか残すべきか
	db, err := sql.Open("mysql", dbparam.dsn())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}
	defer db.Close()

	if q == "" {
		q = "SELECT * FROM " + tblTrans
	}
	rows, err := db.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query: "+q)
	}
	defer rows.Close()

	var list []Transaction
	for rows.Next() {
		var tr Transaction
		err := rows.Scan(&tr.ID, &tr.Type, &tr.Time, &tr.Price,
			&tr.Content, &tr.Raw, &tr.Temporary)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction")
		}
		tr.Time = tr.Time.Local()
		list = append(list, tr)
	}
	return list, nil
}

func InitializeDB() error {
	param := *dbparam
	param.DB = ""
	db, err := sql.Open("mysql", param.dsn())
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	defer db.Close()

	queries := []string{
		"drop database if exists " + dbparam.DB,
		"create database " + dbparam.DB,
		"use " + dbparam.DB,
		"CREATE TABLE " + tblTrans + "(" +
			"id        INT AUTO_INCREMENT, INDEX(id)," +
			"type      TEXT NOT NULL," +
			"time      DATETIME NOT NULL," +
			"price     INT NOT NULL," +
			"content   TEXT NOT NULL," +
			"raw       TEXT NOT NULL," +
			"temporary BOOL NOT NULL" +
			")",
	}
	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "failed to exec query: "+q)
		}
	}
	return nil
}
