package vendors

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"napi/app/models"
)

// --- SAMO Adapter Registration ---
func init() {
	for _, recID := range []int{3306, 3357, 3344, 3419, 3438, 3416, 3428} {
		Register(recID, VendorSamoSoft)
	}
}

// --- Main Adapter Entry Point ---
func VendorSamoSoft(req models.SearchRequest, cfg models.TOConfig) (interface{}, error) {
	params := buildSamoQuery(req, cfg)

	rawPrices, err := fetchSamoSoftXML(cfg.APIURL, params, time.Duration(cfg.Timeout)*time.Second, cfg.Proxy)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}

	var offers []models.Offer
	for _, p := range rawPrices {
		offers = append(offers, mapSamoToOffer(p, cfg.RecID))
	}

	log.Printf("✅ SAMOSOFT [%d]: %d offers parsed", cfg.RecID, len(offers))
	return offers, nil
}

// --- XML Structs ---
type SamoPrices struct {
	Prices []SamoPrice `xml:"prices>price"`
}

type SamoPrice struct {
	ID              string `xml:"id"`
	CheckIn         string `xml:"checkIn"`
	CheckOut        string `xml:"checkOut"`
	Nights          int    `xml:"nights"`
	Adult           int    `xml:"adult"`
	Child           int    `xml:"child"`
	Tour            string `xml:"tour"`
	Hotel           string `xml:"hotel"`
	HotelKey        int    `xml:"hotelKey"`
	Star            string `xml:"star"`
	StarKey         int    `xml:"starKey"`
	Meal            string `xml:"meal"`
	MealKey         int    `xml:"mealKey"`
	Room            string `xml:"room"`
	RoomKey         int    `xml:"roomKey"`
	Price           int    `xml:"price"`
	Currency        string `xml:"currency"`
	CurrencyKey     int    `xml:"currencyKey"`
	ConvertedPrice  string `xml:"convertedPrice"`
	ConvertedCurKey int    `xml:"convertedCurrencyKey"`
	URL             string `xml:"hotelUrl"`
	Grouped         int    `xml:"grouped"`
	External        int    `xml:"external"`
}

// --- HTTP Request ---
func fetchSamoSoftXML(apiURL string, params url.Values, timeout time.Duration, proxy string) ([]SamoPrice, error) {
	fullURL := fmt.Sprintf("%s?action=SearchTour_PRICES&type=xml&%s", apiURL, params.Encode())

	client := &http.Client{Timeout: timeout}
	if proxy != "" {
		if proxyURL, err := url.Parse(proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 10 * time.Second,
			}
		}
	}

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed SamoPrices
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	return parsed.Prices, nil
}

// --- Query Builder ---
func buildSamoQuery(req models.SearchRequest, cfg models.TOConfig) url.Values {
	q := url.Values{}
	q.Set("STATEINC", strconv.Itoa(req.To))
	q.Set("TOWNFROMINC", strconv.Itoa(req.From))
	q.Set("CHECKIN_BEG", strings.ReplaceAll(req.CheckIn, "-", ""))
	q.Set("CHECKIN_END", strings.ReplaceAll(req.CheckTo, "-", ""))
	q.Set("NIGHTS_FROM", strconv.Itoa(req.Nights))
	q.Set("NIGHTS_TILL", strconv.Itoa(req.NightsTo))
	q.Set("ADULT", parseAdult(req.People))
	q.Set("CHILD", parseChildCount(req.People))
	if ageStr := parseChildAges(req.People); ageStr != "" {
		q.Set("AGES", ageStr)
	}
	q.Set("COSTMIN", strconv.Itoa(req.Price))
	q.Set("COSTMAX", strconv.Itoa(req.PriceTo))
	q.Set("CURRENCY", mapCurrency(req.CurrencyLocal))
	q.Set("PRICEPAGE", strconv.Itoa(req.Page))
	q.Set("PARTITION_PRICE", "255")
	q.Set("FILTER", "1")
	q.Set("FREIGHT", "1")
	if cfg.Token != "" {
		q.Set("token", cfg.Token)
	}
	return q
}

// --- Mapper: Samo → Offer ---
func mapSamoToOffer(p SamoPrice, recID int) models.Offer {
	return models.Offer{
		OperatorID: recID,
		HotelID:    p.HotelKey,
		HotelName:  p.Hotel,
		CheckIn:    p.CheckIn,
		Nights:     p.Nights,
		RoomName:   p.Room,
		Meal:       p.Meal,
		Price:      p.Price * 100,
		Currency:   mapCurrencyCode(p.Currency),
		BronURL:    p.URL,
		OfferID:    genOfferIDFromSamo(p, recID), // твоя функція генерації
	}
}

// --- Helpers ---
func parseAdult(people string) string {
	if len(people) > 0 {
		return string(people[0])
	}
	return "2"
}

func parseChildCount(people string) string {
	if len(people) < 2 {
		return "0"
	}
	return strconv.Itoa(len(people[1:]) / 2)
}

func parseChildAges(people string) string {
	if len(people) < 3 {
		return ""
	}
	var ages []string
	for i := 1; i+1 < len(people); i += 2 {
		age := people[i : i+2]
		if n, err := strconv.Atoi(age); err == nil {
			ages = append(ages, strconv.Itoa(n))
		}
	}
	return strings.Join(ages, ",")
}

func mapCurrency(code string) string {
	switch strings.ToLower(code) {
	case "eur":
		return "4"
	case "usd":
		return "3"
	case "uah":
		return "1"
	default:
		return "4"
	}
}

func mapCurrencyCode(s string) int {
	switch strings.ToUpper(s) {
	case "UAH":
		return 980
	case "EUR":
		return 978
	case "USD":
		return 840
	default:
		return 980
	}
}
