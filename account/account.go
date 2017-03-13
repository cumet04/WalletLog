package account

import (
	"database/sql"
	"fmt"
	"time"

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

var dbdsn string

const tblTrans = "transactions"

func AddTrans(tr Transaction) error {
	db, err := sql.Open("mysql", dbdsn)
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
	db, err := sql.Open("mysql", dbdsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}
	defer db.Close()

	if q == "" {
		q = "SELECT * FROM " + tblTrans
	}
	rows, err := db.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query")
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

func SetDataSoruceName(dsn string) {
	dbdsn = dsn
}

func ClearDB() error {
	// for unit test

	db, err := sql.Open("mysql", dbdsn)
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	defer db.Close()

	queries := []string{
		"DROP TABLE IF EXISTS " + tblTrans,
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
