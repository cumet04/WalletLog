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
	Source    string    `json:"type"`
	Time      time.Time `json:"time"`
	Price     int       `json:"price"`
	Content   string    `json:"content"`
	Raw       string    `json:"raw"`
	Temporary bool      `json:"temporary"`
}

type Source struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	Available bool   `json:"available"`
}

type dbParam struct {
	User    string
	Pass    string
	Host    string
	Port    int
	DB      string
	Options string
}

type SourceType int

const (
	SourceBank SourceType = iota
	SourcePayment
)

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
		"CREATE TABLE sources (" +
			"name      VARCHAR(32) PRIMARY KEY," +
			"type      INT NOT NULL," +
			"available BOOL NOT NULL" +
			")",
		"INSERT INTO sources values" +
			"('smbc-bank', '" + fmt.Sprint(SourceBank) + "', true)," +
			"('mizuho-bank', '" + fmt.Sprint(SourceBank) + "', true)," +
			"('smbc-visa', '" + fmt.Sprint(SourcePayment) + "', true)," +
			"('amazon-card', '" + fmt.Sprint(SourcePayment) + "', true)," +
			"('view-card', '" + fmt.Sprint(SourcePayment) + "', true)," +
			"('tokyu-card', '" + fmt.Sprint(SourcePayment) + "', true)," +
			"('jpbank-card', '" + fmt.Sprint(SourcePayment) + "', false)," +
			"('pitapa', '" + fmt.Sprint(SourcePayment) + "', true)",
		// TODO: hash
		"CREATE TABLE transactions (" +
			"id        INT AUTO_INCREMENT PRIMARY KEY," +
			"source    VARCHAR(32) NOT NULL," +
			"FOREIGN KEY (source) REFERENCES sources (name)," +
			"time      DATETIME NOT NULL," +
			"price     INT NOT NULL," +
			"content   VARCHAR(256) NOT NULL," +
			"raw       VARCHAR(512) NOT NULL," +
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

func SetDBParam(u string, pa string, h string, po int, db string, o string) {
	dbparam = &dbParam{u, pa, h, po, db, o}
}

func GetSource(name string) (*Source, error) {
	db, err := sql.Open("mysql", dbparam.dsn())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}
	defer db.Close()

	q := fmt.Sprintf("SELECT name, type, available from sources WHERE name = '%s'", name)
	var res Source
	if err = db.QueryRow(q).Scan(&res.Name, &res.Type, &res.Available); err != nil {
		return nil, errors.Wrap(err, "failed to select query")
	}

	return &res, nil
}

func PutSource(so Source) error {
	db, err := sql.Open("mysql", dbparam.dsn())
	if err != nil {
		return errors.Wrap(err, "failed to connect to db")
	}
	defer db.Close()

	q := fmt.Sprintf("INSERT INTO sources (name, type, available) VALUES('%s', %d, %t) "+
		"ON DUPLICATE KEY UPDATE type=values(type), available=values(available)",
		so.Name, so.Type, so.Available)
	if _, err = db.Exec(q); err != nil {
		return errors.Wrap(err, "failed to insert query")
	}

	return nil
}

func AddTrans(tr Transaction) error {
	db, err := sql.Open("mysql", dbparam.dsn())
	if err != nil {
		return errors.Wrap(err, "failed to connect to db")
	}
	defer db.Close()

	q := fmt.Sprintf("INSERT INTO transactions (source, time, price, content, raw, temporary) "+
		"VALUES('%s', '%s', %d, '%s', '%s', %t)", tr.Source,
		tr.Time.UTC().Format("2006-01-02 15:04:05"), tr.Price, tr.Content, tr.Raw, tr.Temporary)
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
		q = "SELECT * FROM transactions"
	}
	rows, err := db.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query: "+q)
	}
	defer rows.Close()

	var list []Transaction
	for rows.Next() {
		var tr Transaction
		err := rows.Scan(&tr.ID, &tr.Source, &tr.Time, &tr.Price,
			&tr.Content, &tr.Raw, &tr.Temporary)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction")
		}
		tr.Time = tr.Time.Local()
		list = append(list, tr)
	}
	return list, nil
}
