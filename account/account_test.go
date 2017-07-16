package account

import (
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	if err := PrepareDB("sqlite3", "file:test_account?mode=memory"); err != nil {
		panic(err)
	}
	ret := m.Run()
	CloseDB()

	if err := PrepareDB("mysql", "root@tcp(127.0.0.1:3306)/walletlog_test?parseTime=true"); err != nil {
		panic(err)
	}
	ret += m.Run()
	CloseDB()

	os.Exit(ret)
}

func TestAddPurchaseSingle(t *testing.T) {
	if err := InitializeDB(); err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	input := []Purchase{
		Purchase{
			Source:    "jpbank-card",
			Time:      time.Date(2015, 12, 22, 0, 0, 0, 0, loc),
			Price:     2915,
			Content:   "神戸市水道局 コウベシスイドウキヨク　Ｅ５",
			Raw:       "2015/12/22,神戸市水道局,2915,１,１,2915,コウベシスイドウキヨク　Ｅ５",
			Temporary: true,
		},
	}

	if err := AddPurchase(input); err != nil {
		t.Fatalf("Error by AddPurchase. %v", err)
	}

	results, err := QueryPurchase("")
	if err != nil {
		t.Fatalf("Error by QueryPurchase. %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("result's size is %d instead of 1", len(results))
	}
	if equalPurchase(results[0], input[0]) == false {
		t.Fatalf("unexpected result: %v, expected: %v", results[0], input[0])
	}
}

func TestAddPurchaseMulti(t *testing.T) {
	if err := InitializeDB(); err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	input := []Purchase{
		Purchase{
			Source:    "jpbank-card",
			Time:      time.Date(2015, 12, 22, 0, 0, 0, 0, loc),
			Price:     2915,
			Content:   "神戸市水道局 コウベシスイドウキヨク　Ｅ５",
			Raw:       "2015/12/22,神戸市水道局,2915,１,１,2915,コウベシスイドウキヨク　Ｅ５",
			Temporary: true,
		},
		Purchase{
			Source:    "jpbank-card",
			Time:      time.Date(2015, 12, 22, 0, 0, 0, 0, loc),
			Price:     2915,
			Content:   "神戸市水道局 コウベシスイドウキヨク　Ｅ５",
			Raw:       "2015/12/22,神戸市水道局,2915,１,１,2915,コウベシスイドウキヨク　Ｅ５",
			Temporary: true,
		},
	}

	if err := AddPurchase(input); err != nil {
		t.Fatalf("Error by AddPurchase. %v", err)
	}

	results, err := QueryPurchase("")
	if err != nil {
		t.Fatalf("Error by QueryPurchase. %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("result's size is %d instead of 2", len(results))
	}
	for _, in := range input {
		if (equalPurchase(in, results[0]) == false) && (equalPurchase(in, results[1]) == false) {
			t.Fatalf("expected item is not found: %v in %v", in, results)
		}
	}
}

func TestAddBankTradeSingle(t *testing.T) {
	if err := InitializeDB(); err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	input := []BankTrade{
		BankTrade{
			Source:    "mizuho-bank",
			Time:      time.Date(2017, 6, 26, 0, 0, 0, 0, loc),
			Price:     90206,
			Content:   "ミツイスミトモカ－ド （カ",
			Raw:       "2017.06.26 90,206 円 - ミツイスミトモカ－ド （カ",
			Temporary: true,
		},
	}

	if err := AddTrade(input); err != nil {
		t.Fatalf("Error by AddTrade. %v", err)
	}

	results, err := QueryTrade("")
	if err != nil {
		t.Fatalf("Error by QueryTrade. %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("result's size is %d instead of 1", len(results))
	}
	if equalTrade(results[0], input[0]) == false {
		t.Fatalf("unexpected result: %v, expected: %v", results[0], input[0])
	}
}

func TestAddBankTradeMulti(t *testing.T) {
	if err := InitializeDB(); err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	input := []BankTrade{
		BankTrade{
			Source:    "mizuho-bank",
			Time:      time.Date(2017, 5, 10, 0, 0, 0, 0, loc),
			Price:     211813,
			Content:   "振込 カ）オロ",
			Raw:       "2017.05.10 - 211,813 円 振込 カ）オロ",
			Temporary: true,
		},
		BankTrade{
			Source:    "mizuho-bank",
			Time:      time.Date(2017, 6, 26, 0, 0, 0, 0, loc),
			Price:     90206,
			Content:   "ミツイスミトモカ－ド （カ",
			Raw:       "2017.06.26 90,206 円 - ミツイスミトモカ－ド （カ",
			Temporary: true,
		},
	}

	if err := AddTrade(input); err != nil {
		t.Fatalf("Error by AddTrade. %v", err)
	}

	results, err := QueryTrade("")
	if err != nil {
		t.Fatalf("Error by QueryTrade. %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("result's size is %d instead of 2", len(results))
	}
	for _, in := range input {
		if (equalTrade(in, results[0]) == false) && (equalTrade(in, results[1]) == false) {
			t.Fatalf("expected item is not found: %v in %v", in, results)
		}
	}
}

// equalPurchase checks that two Purchases are same except ID.
func equalPurchase(e1 Purchase, e2 Purchase) bool {
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

// equalTrade checks that two BankTrades are same except ID.
func equalTrade(e1 BankTrade, e2 BankTrade) bool {
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
