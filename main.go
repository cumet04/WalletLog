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
	res, err := account.QueryTrade("")
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
	// refs := make([]*account.Purchase, len(input))
	// for i, item := range input {
	// 	refs[i] = &item
	// }
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
	res, err := account.QueryPurchase("")
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
	// account.SetDBParam("mysql", "root@tcp(127.0.0.1:3306)/walletlog_test?parseTime=true")
	// account.SetDBParam("sqlite3", "main.db")
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
