package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/util"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(s *db.Store) *HotelHandler {
	return &HotelHandler{
		store: s,
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {
	hotels, err := h.store.HotelStore.GetHotels(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(hotels)
}

func (h *HotelHandler) HandleGetRoomsByHotelId(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, _ := util.ObjectIdParser(id)
	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.RoomStore.GetRooms(ctx.Context(), filter)
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	hotel, err := h.store.HotelStore.GetHotelById(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(hotel)
}

func (h *HotelHandler) HandlePutHotel(ctx *fiber.Ctx) error {
	var (
		id     = ctx.Params("id")
		values = bson.M{}
	)
	if err := ctx.BodyParser(&values); err != nil {
		return err
	}
	oid, _ := util.ObjectIdParser(id)

	update := bson.M{"$set": values}

	err := h.store.HotelStore.UpdateHotel(ctx.Context(), bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"message": "updated hotel are: " + id})
}
