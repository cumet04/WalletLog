package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/medalhkr/WalletLog/account"
)

func apiAddTrans(w http.ResponseWriter, r *http.Request) {
	var input account.Transaction
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "failed to unmarshal input")
	}
	if err := account.AddTrans(input); err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed to add transaction")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, input)
}

func apiQueryTrans(w http.ResponseWriter, r *http.Request) {
	res, err := account.QueryTrans("")
	if err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed")
	}
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiQueryTrans(w, r)
	case "POST":
		apiAddTrans(w, r)
	}
}

func main() {
	initdb := *flag.Bool("initialize", false, "performs profile")
	flag.Parse()
	if initdb {
		if err := account.InitializeDB(); err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/transactions/", transactionsHandler)
	http.ListenAndServe(":50000", nil)
}
