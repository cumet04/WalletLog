package account

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type BankTrade struct {
	ID        int       `json:"id"`
	Source    string    `json:"source"`
	Time      time.Time `json:"time"`
	Price     int       `json:"price"`
	Content   string    `json:"content"`
	Raw       string    `json:"raw"`
	Temporary bool      `json:"temporary"`
}

type Purchase struct {
	ID        int       `json:"id"`
	Source    string    `json:"source"`
	Time      time.Time `json:"time"`
	Price     int       `json:"price"`
	Content   string    `json:"content"`
	Raw       string    `json:"raw"`
	Temporary bool      `json:"temporary"`
}

type dbParam struct {
	Type string
	Dsn  string
}

var dbm *gorp.DbMap

func PrepareDB(driver, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	var dbmap gorp.DbMap
	if driver == "sqlite3" {
		dbmap = gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	} else if driver == "mysql" {
		dbmap = gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}}
	} else {
		return errors.New("unknown db driver: " + driver)
	}

	dbmap.AddTableWithName(Purchase{}, "purchases").SetKeys(true, "ID")
	dbmap.AddTableWithName(BankTrade{}, "bank_trades").SetKeys(true, "ID")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return errors.Wrap(err, "failed to create tables")
	}

	dbm = &dbmap
	return nil
}

func CloseDB() {
	dbm.Db.Close()
	dbm = nil
}

func InitializeDB() error {
	if err := dbm.DropTablesIfExists(); err != nil {
		return errors.Wrap(err, "failed to drop tables")
	}
	if err := dbm.CreateTables(); err != nil {
		return errors.Wrap(err, "failed to create tables")
	}
	return nil
}

func AddPurchase(inputs []Purchase) error {
	// MEMO: sliceをまとめて渡すのはうまくいかなかった
	// inputs... -> non pointer
	// []*Purhace... -> 謎エラー
	// []interface{} := map(inputs, return &item) -> 何故か同じ要素を複数挿入になる
	// というかgorpがbulk insert対応してなかった

	for _, item := range inputs {
		err := dbm.Insert(&item)
		if err != nil {
			return errors.Wrap(err, "failed to insert purchases")
		}
	}
	return nil
}

func QueryPurchase(where string) ([]Purchase, error) {
	var res []Purchase
	q := "select * from purchases"
	if where != "" {
		q += " " + where
	}
	_, err := dbm.Select(&res, q)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query purchases")
	}
	return res, nil
}

func AddTrade(inputs []BankTrade) error {
	for _, item := range inputs {
		err := dbm.Insert(&item)
		if err != nil {
			return errors.Wrap(err, "failed to insert purchases")
		}
	}
	return nil
}

func QueryTrade(where string) ([]BankTrade, error) {
	var res []BankTrade
	q := "select * from bank_trades"
	if where != "" {
		q += " " + where
	}
	_, err := dbm.Select(&res, q)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query purchases")
	}
	return res, nil
}
