package account

import "time"

type Transaction struct {
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	Price     int       `json:"price"`
	Content   string    `json:"content"`
	Raw       string    `json:"raw"`
	Temporary bool      `json:"temporary"`
}

var TransList = []Transaction{
	Transaction{"visa", time.Now(), 2000, "kabe", "20170310, 1000, kabe", true},
	Transaction{"suica", time.Now(), 100, "zihanki", "20170312, 100, zihanki", true},
}
