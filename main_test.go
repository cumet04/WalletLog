package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"io/ioutil"

	"github.com/medalhkr/WalletLog/account"
)

func TestMain(m *testing.M) {
	if err := account.PrepareDB("sqlite3", "file:test_main?mode=memory"); err != nil {
		panic(err)
	}
	ret := m.Run()
	account.CloseDB()

	os.Exit(ret)
}

func TestBankGetPost(t *testing.T) {
	account.InitializeDB()
	ts := httptest.NewServer(http.HandlerFunc(bankHandler))

	// insert
	input := strings.NewReader(`[
		{
			"source": "mizuho-bank",
			"time": "2017-05-10T00:00:00+09:00",
			"price": 211813,
			"content": "振込 カ）オロ",
			"raw": "2017.05.10 - 211,813 円 振込 カ）オロ",
			"temporary": true
		}
	]`)
	r, err := http.Post(ts.URL, "application/json", input)
	if err != nil {
		t.Fatalf("Error by http.Put(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}

	// get
	r, err = http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}
	var list []account.BankTrade
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		t.Fatalf("Error by json decode. %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	expected := account.BankTrade{
		Source:    "mizuho-bank",
		Time:      time.Date(2017, 5, 10, 0, 0, 0, 0, loc),
		Price:     211813,
		Content:   "振込 カ）オロ",
		Raw:       "2017.05.10 - 211,813 円 振込 カ）オロ",
		Temporary: true,
	}
	if equalBank(list[0], expected) == false {
		t.Fatalf("unexpected result: %v, expected: %v", list[0], expected)
	}
}

func TestPurchaseGetPost(t *testing.T) {
	account.InitializeDB()
	ts := httptest.NewServer(http.HandlerFunc(purchaseHandler))
	defer ts.Close()

	// insert
	input := strings.NewReader(`[
		{
			"source": "jpbank-card",
			"time": "2015-12-22T00:00:00+09:00",
			"price": 2915,
			"content": "神戸市水道局 コウベシスイドウキヨク　Ｅ５",
			"raw": "2015/12/22,神戸市水道局,2915,１,１,2915,コウベシスイドウキヨク　Ｅ５",
			"temporary": true
		}
	]`)
	r, err := http.Post(ts.URL, "application/json", input)
	if err != nil {
		t.Fatalf("Error by http.Put(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}

	// get
	r, err = http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}
	var list []account.Purchase
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		t.Fatalf("Error by json decode. %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	expected := account.Purchase{
		Source:    "jpbank-card",
		Time:      time.Date(2015, 12, 22, 0, 0, 0, 0, loc),
		Price:     2915,
		Content:   "神戸市水道局 コウベシスイドウキヨク　Ｅ５",
		Raw:       "2015/12/22,神戸市水道局,2915,１,１,2915,コウベシスイドウキヨク　Ｅ５",
		Temporary: true,
	}
	if equalPurchase(list[0], expected) == false {
		t.Fatalf("unexpected result: %v, expected: %v", list[0], expected)
	}
}

func equalPurchase(e1 account.Purchase, e2 account.Purchase) bool {
	if e1.Source != e2.Source {
		return false
	}
	if e1.Time.Equal(e2.Time) == false {
		return false
	}
	if e1.Price != e2.Price {
		return false
	}
	if e1.Content != e2.Content {
		return false
	}
	if e1.Raw != e2.Raw {
		return false
	}
	if e1.Temporary != e2.Temporary {
		return false
	}
	return true
}

func equalBank(e1 account.BankTrade, e2 account.BankTrade) bool {
	if e1.Source != e2.Source {
		return false
	}
	if e1.Time.Equal(e2.Time) == false {
		return false
	}
	if e1.Price != e2.Price {
		return false
	}
	if e1.Content != e2.Content {
		return false
	}
	if e1.Raw != e2.Raw {
		return false
	}
	if e1.Temporary != e2.Temporary {
		return false
	}
	return true
}
