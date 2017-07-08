package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"io/ioutil"

	"github.com/medalhkr/WalletLog/account"
)

func setup() {
	account.SetDBParam("root", "", "127.0.0.1", 3306, "walletlog_test", "parseTime=true")
	if err := account.InitializeDB(); err != nil {
		log.Fatalf("ERROR in setup: %v", err)
	}
}

func teardown() {
	// Nothing to do
}

func TestSinglePutSource(t *testing.T) {
	if err := account.InitializeDB(); err != nil {
		log.Fatalf("ERROR in setup: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(sourceHandler))
	defer ts.Close()

	// insert
	input := strings.NewReader(`{
		"name": "test-card",
		"type": 0,
		"available": true
	}`)
	req, _ := http.NewRequest(http.MethodPut, ts.URL, input)
	req.Header.Add("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error by http.Put(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}

	// get
	r, err = http.Get(ts.URL + "/?name=test-card")
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}
	var res account.Source
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		t.Fatalf("Error by json decode. %v", err)
	}
	expected := account.Source{
		Name:      "test-card",
		Type:      0,
		Available: true,
	}
	if equalSource(res, expected) == false {
		t.Fatalf("unexpected result: %v, expected: %v", res, expected)
	}
}

func TestDuplicatePutSource(t *testing.T) {
	if err := account.InitializeDB(); err != nil {
		log.Fatalf("ERROR in setup: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(sourceHandler))
	defer ts.Close()

	for i := range []int{1, 2} {
		// insert
		input := strings.NewReader(`{
			"name": "test-card",
			"type": 0,
			"available": true
		}`)
		req, _ := http.NewRequest(http.MethodPut, ts.URL, input)
		req.Header.Add("Content-Type", "application/json")
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error by http.Put() at try %d. %v", i, err)
		}
		if r.StatusCode != 200 {
			msg, _ := ioutil.ReadAll(r.Body)
			t.Fatalf("The status code is %d at try %d\nmsg: %s", i, r.StatusCode, string(msg))
		}

		// get
		r, err = http.Get(ts.URL + "/?name=test-card")
		if err != nil {
			t.Fatalf("Error by http.Get(). %v", err)
		}
		if r.StatusCode != 200 {
			msg, _ := ioutil.ReadAll(r.Body)
			t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
		}
		var res account.Source
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			t.Fatalf("Error by json decode. %v", err)
		}
		expected := account.Source{
			Name:      "test-card",
			Type:      0,
			Available: true,
		}
		if equalSource(res, expected) == false {
			t.Fatalf("unexpected result: %v, expected: %v", res, expected)
		}
	}
}

func TestSingleAddAndList(t *testing.T) {
	if err := account.InitializeDB(); err != nil {
		log.Fatalf("ERROR in setup: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(transactionsHandler))
	defer ts.Close()

	// test add transaction
	input := strings.NewReader(`{
		"type": "smbc-visa",
		"time": "2017-03-12T00:00:00+09:00",
		"price": 2000,
		"content": "kabe",
		"raw": "20170310, 1000, kabe",
		"temporary": true
	}`)
	r, err := http.Post(ts.URL, "application/json", input)
	if err != nil {
		t.Fatalf("Error by http.Post(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}

	// test get transactions
	r, err = http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}
	var list []account.Transaction
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		t.Fatalf("Error by json decode. %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("result's size is wrong: %d", len(list))
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	expected := account.Transaction{
		Source:    "smbc-visa",
		Time:      time.Date(2017, 3, 12, 0, 0, 0, 0, loc),
		Price:     2000,
		Content:   "kabe",
		Raw:       "20170310, 1000, kabe",
		Temporary: true,
	}
	if equalTrans(list[0], expected) == false {
		t.Fatalf("unexpected result: %v, expected: %v", list[0], expected)
	}
}

func TestSingleAddAndListJP(t *testing.T) {
	if err := account.InitializeDB(); err != nil {
		log.Fatalf("ERROR in setup: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(transactionsHandler))
	defer ts.Close()

	// test add transaction
	input := strings.NewReader(`{
		"type": "jpbank-card",
		"time": "2015-12-22T00:00:00+09:00",
		"price": 2915,
		"content": "神戸市水道局 コウベシスイドウキヨク　Ｅ５",
		"raw": "2015/12/22,神戸市水道局,2915,１,１,2915,コウベシスイドウキヨク　Ｅ５",
		"temporary": true
	}`)
	r, err := http.Post(ts.URL, "application/json", input)
	if err != nil {
		t.Fatalf("Error by http.Post(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}

	// test get transactions
	r, err = http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}
	if r.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(r.Body)
		t.Fatalf("The status code is %d\nmsg: %s", r.StatusCode, string(msg))
	}
	var list []account.Transaction
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		t.Fatalf("Error by json decode. %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("result's size is wrong: %d", len(list))
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	expected := account.Transaction{
		Source:    "jpbank-card",
		Time:      time.Date(2015, 12, 22, 0, 0, 0, 0, loc),
		Price:     2915,
		Content:   "神戸市水道局 コウベシスイドウキヨク　Ｅ５",
		Raw:       "2015/12/22,神戸市水道局,2915,１,１,2915,コウベシスイドウキヨク　Ｅ５",
		Temporary: true,
	}
	if equalTrans(list[0], expected) == false {
		t.Fatalf("unexpected result: %v, expected: %v", list[0], expected)
	}
}
func equalSource(e1 account.Source, e2 account.Source) bool {
	if e1.Name != e2.Name {
		return false
	}
	if e1.Type != e2.Type {
		return false
	}
	if e1.Available != e2.Available {
		return false
	}
	return true
}

func equalTrans(e1 account.Transaction, e2 account.Transaction) bool {
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

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	if ret == 0 {
		teardown()
	}
	os.Exit(ret)
}
