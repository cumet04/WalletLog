package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/medalhkr/WalletLog/account"
)

func apiTransactionsPUT(w http.ResponseWriter, r *http.Request) {
	var input account.Transaction
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "failed to unmarshal input")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, input)
}

func apiTransactionsGET(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(account.TransList)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "failed")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

func main() {
	http.HandleFunc("/transactions/",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				apiTransactionsGET(w, r)
			case "PUT":
				apiTransactionsPUT(w, r)
			}
		})
	http.ListenAndServe(":50000", nil)
}
