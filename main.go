package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"log"

	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/medalhkr/WalletLog/account"
)

func apiAddTrans(w http.ResponseWriter, r *http.Request) {
	var input account.Transaction
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to unmarshal input: %v", err)
		return
	}
	if err := account.AddTrans(input); err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed to add transaction")
		return
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

func apiGetSource(w http.ResponseWriter, r *http.Request) {
	res, err := account.GetSource(r.URL.Query().Get("name"))
	if err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed to put source")
		return
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

func apiPutSource(w http.ResponseWriter, r *http.Request) {
	var input account.Source
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to unmarshal input: %v", err)
		return
	}
	if err := account.PutSource(input); err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
		fmt.Fprint(w, "failed to put source")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, input)
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiQueryTrans(w, r)
	case "POST":
		apiAddTrans(w, r)
		// default 406
	}
}

func sourceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiGetSource(w, r)
	case "PUT":
		apiPutSource(w, r)
	}
}

func main() {
	var migrate bool
	flag.BoolVar(&migrate, "migrate", false, "initialize DB")
	flag.Parse()
	account.SetDBParam("root", "", "127.0.0.1", 3306, "walletlog", "parseTime=true")
	if migrate {
		log.Println("initialize DB...")
		if err := account.InitializeDB(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	log.Println("start serving")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api/transactions/", transactionsHandler)
	http.HandleFunc("/api/source/", sourceHandler)
	http.ListenAndServe(":50000", nil)
}
