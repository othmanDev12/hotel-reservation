package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// HandleGetBookings is a function that is available for admin
func (h *BookingHandler) HandleGetBookings(ctx *fiber.Ctx) error {
	bookings, err := h.store.BookingStore.GetBookings(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(bookings)
}

// HandleGetBooking is a function that is available for user
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	return nil
}
