package vendors_test

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"napi/app/vendors"
)

func TestSamoSoft_ParsePricesXML(t *testing.T) {
	xmlData := `
<SearchTour_PRICES>
  <prices>
    <price>
      <id>0x001</id>
      <hotel>Grand Hotel</hotel>
      <hotelKey>777</hotelKey>
      <price>1299</price>
      <currency>EUR</currency>
      <checkIn>20250418</checkIn>
      <meal>AI</meal>
      <room>Standard Room</room>
      <nights>7</nights>
    </price>
  </prices>
</SearchTour_PRICES>`

	var parsed vendors.SamoPrices
	err := xml.Unmarshal([]byte(xmlData), &parsed)

	assert.NoError(t, err)
	assert.Len(t, parsed.Prices, 1)
	assert.Equal(t, "Grand Hotel", parsed.Prices[0].Hotel)
	assert.Equal(t, 1299, parsed.Prices[0].Price)
	assert.Equal(t, "EUR", parsed.Prices[0].Currency)
}

func TestSamoSoft_MapToOffer(t *testing.T) {
	p := vendors.SamoPrice{
		Hotel:     "Test Resort",
		HotelKey:  888,
		CheckIn:   "20250418",
		Nights:    5,
		Meal:      "UAI",
		Room:      "Deluxe",
		Price:     950,
		Currency:  "USD",
	}
	offer := vendors.MapSamoToOffer(p, 2700)

	assert.Equal(t, 888, offer.HotelID)
	assert.Equal(t, "UAI", offer.Meal)
	assert.Equal(t, 950*100, offer.Price)
	assert.Equal(t, "Test Resort", offer.HotelName)
	assert.Equal(t, 840, offer.Currency)
}
