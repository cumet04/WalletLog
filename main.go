package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"log"

	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/medalhkr/WalletLog/account"
)

func apiAddTrade(w http.ResponseWriter, r *http.Request) {
	var input []account.BankTrade
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to unmarshal input: %v", err)
		return
	}
	if err := account.AddTrade(input); err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed to add transaction")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, input)
}

func apiQueryTrade(w http.ResponseWriter, r *http.Request) {
	query := r.Form.Get("query")
	res, err := account.QueryTrade(query)
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

func apiAddPurchase(w http.ResponseWriter, r *http.Request) {
	var input []account.Purchase
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to unmarshal input: %v", err)
		return
	}
	if err := account.AddPurchase(input); err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed to add transaction")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, input)
}

func apiQueryPurchase(w http.ResponseWriter, r *http.Request) {
	query := r.Form.Get("query")
	res, err := account.QueryPurchase(query)
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

func purchaseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiQueryPurchase(w, r)
	case "POST":
		apiAddPurchase(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func bankHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiQueryTrade(w, r)
	case "POST":
		apiAddTrade(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	var migrate bool
	flag.BoolVar(&migrate, "migrate", false, "initialize DB")
	flag.Parse()
	// if err := account.PrepareDB("mysql", "root@tcp(127.0.0.1:3306)/walletlog?parseTime=true"); err != nil {
	// if err := account.PrepareDB("sqlite3", "file:main?mode=memory"); err != nil {
	if err := account.PrepareDB("sqlite3", "realdata.db"); err != nil {
		panic(err)
	}
	defer account.CloseDB()

	if migrate {
		log.Println("initialize DB...")
		if err := account.InitializeDB(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	log.Println("start serving")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api/purchase/", purchaseHandler)
	http.HandleFunc("/api/bank_trade/", bankHandler)
	http.ListenAndServe(":50000", nil)
}
