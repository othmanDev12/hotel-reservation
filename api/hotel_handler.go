package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/util"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(hotels)
}

func (h *HotelHandler) HandleGetRoomsByHotelId(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, _ := util.ObjectIdParser(id)
	filter := bson.M{"hotelID": oid}
	rooms, err := h.roomStore.GetRooms(ctx.Context(), filter)
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}
