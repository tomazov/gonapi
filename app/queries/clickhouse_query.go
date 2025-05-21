package queries

import (
	"context"
	"fmt"
	"log"
	"time"

	"napi/app/models"
	"napi/pkg/platform/clickhouse"
)

func InsertOffers(ctx context.Context, offers []models.Offer) error {
	if len(offers) == 0 {
		return nil
	}

	conn := clickhouse.Get()

	batch, err := conn.PrepareBatch(ctx, `
		INSERT INTO offers (
			operatorId, countryId, cityId, hotelId, fromCityId, adultChild,
			ages, checkIn, nights, duration, transportFood,
			transportOption, tourName, stopSale, roomName,
			price, currency, tourOptions, bronUrl, offerId,
			active, updateTime
		)
	`)
	if err != nil {
		return fmt.Errorf("❌ batch prepare: %w", err)
	}

	now := time.Now()

	for _, offer := range offers {
		err := batch.AppendStruct(&models.ClickHouseOffer{
			OperatorID:      offer.OperatorID,
			CountryID:       offer.CountryID,
			CityID:          offer.CityID,
			HotelID:         offer.HotelID,
			FromCityID:      offer.FromCityID,
			AdultChild:      offer.AdultChild,
			Ages:            offer.Ages,
			CheckIn:         offer.CheckIn,
			Nights:          offer.Nights,
			Duration:        offer.Duration,
			TransportFood:   offer.TransportFood,
			TransportOption: offer.TransportOption,
			TourName:        offer.TourName,
			StopSale:        offer.StopSale,
			RoomName:        offer.RoomName,
			Price:           offer.Price,
			Currency:        offer.Currency,
			TourOptions:     offer.TourOptions,
			BronURL:         offer.BronURL,
			OfferID:         offer.OfferID,
			Active:          1,
			UpdateTime:      now,
		})
		if err != nil {
			log.Printf("⚠️ append offer error: %v", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("❌ batch send: %w", err)
	}

	log.Printf("✅ ClickHouse: inserted %d offers", len(offers))
	return nil
}
