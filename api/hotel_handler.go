package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"log"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

type HotelParamsHandler struct {
	Rooms  bool
	Rating int
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {
	var qParams HotelParamsHandler
	if err := ctx.QueryParser(&qParams); err != nil {
		return err
	}
	log.Printf("params are %v", qParams)
	hotels, err := h.hotelStore.GetHotels(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(hotels)
}
