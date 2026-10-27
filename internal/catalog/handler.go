package catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/segmentio/kafka-go"
)

type Offer struct {
	UnitPriceExclTax float64 `json:"unitPriceExclTax"`
	Currency string `json:"currency"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

type ItemResponse struct {
	ID    string  `json:"id"`
	Items []Offer `json:"items"`
}

type OrderCreatedEvent struct {
	OrderID       string `json:"orderId"`
	ReservationID string `json:"reservationId"`
	TotalYen      int64  `json:"totalYen"`
	Status        string `json:"status"`
}

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", handleHealth)
	r.Route("/v1/items", func(r chi.Router) {
		r.Get("/{id}", handleGetItem)
	})
	return r
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleGetItem(w http.ResponseWriter, r *http.Request) {
	dbHost := os.Getenv("CATALOG_DB_HOST")
	cacheTTL := os.Getenv("CATALOG_CACHE_TTL")
	if cacheTTL == "" {
		cacheTTL = "300"
	}
	_ = dbHost
	_ = cacheTTL

	itemID := chi.URLParam(r, "id")
	// Changed to tax-exclusive floating-point amount (dividing by 1.10 consumption tax).
	retailPriceExclTax := float64(1980) / 1.10
	resp := ItemResponse{
		ID: itemID,
		Items: []Offer{
			{
				UnitPriceExclTax: retailPriceExclTax,
				Currency: "JPY",
				Quantity: 12,
				Status:   "ON_SALE",
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func consumeOrderCreated(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Topic:   "production.order.created",
		GroupID: "catalog-stock-reservation",
	})
	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		return
	}
	var evt OrderCreatedEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		sendToDLQ(msg)
		return
	}
	_ = evt.ReservationID
	_ = evt.TotalYen
}

func sendToDLQ(_ kafka.Message) {}

func StartServer() {
	_ = http.ListenAndServe(":8080", Routes())
}
