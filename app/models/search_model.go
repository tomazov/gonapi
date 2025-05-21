package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type SearchRequest struct {
	From          int      `query:"from"`
	To            int      `query:"to"`
	Stars         string   `query:"stars"`
	CheckIn       string   `query:"checkIn"`
	CheckTo       string   `query:"checkTo"`
	Nights        int      `query:"nights"`
	NightsTo      int      `query:"nightsTo"`
	People        string   `query:"people"`
	Food          string   `query:"food"`
	Transport     string   `query:"transport"`
	Price         int      `query:"price"`
	PriceTo       int      `query:"priceTo"`
	Currency      string   `query:"currency"`
	Page          int      `query:"page"`
	CurrencyLocal string   `query:"currencyLocal"`
	ToOperators   []int    `query:"toOperators"`
	AvailableFlight string `query:"availableFlight"`
	StopSale      string   `query:"stopSale"`
	Lang          string   `query:"lang"`
	Group         int      `query:"group"`
	Rating        string   `query:"rating"`
	Number        int      `query:"number"`
	AccessToken   string   `query:"access_token"`
}

// Генерує MD5-хеш на основі ключових параметрів запиту
func (r *SearchRequest) GenerateSearchID() string {
	raw := fmt.Sprintf("f%d%d%s%s%d%d%s%d",
		r.From,
		r.To,
		r.CheckIn,
		r.CheckTo,
		r.Nights,
		r.NightsTo,
		r.People,
		r.Page,
	)
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}
