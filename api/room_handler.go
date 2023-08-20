package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"github.com/hotel-reservation/util"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type RoomHandler struct {
	store *db.Store
}

type RoomBookingParams struct {
	NumPersons int       `json:"numPersons"`
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
}

// NewRoomHandler creates a new constructor for the room handler
func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (p RoomBookingParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("the current date that you pass is on the past")
	}
	return nil
}

func (h *RoomHandler) HandleRoomBooking(c *fiber.Ctx) error {
	var params RoomBookingParams
	id := c.Params("id")
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	roomId, err := util.ObjectIdParser(id)
	if err != nil {
		return err
	}
	userId, ok := c.Context().Value("user").(*domain.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(GenericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}
	ok, err = h.roomIsBooked(c, params)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(GenericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s is already booked", id),
		})
	}

	book := domain.Booking{
		UserId:        userId.Id,
		RoomId:        roomId,
		NumberPersons: params.NumPersons,
		FromDate:      params.FromDate,
		TillDate:      params.TillDate,
	}

	inserted, err := h.store.BookingStore.InsertBooking(c.Context(), &book)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) roomIsBooked(c *fiber.Ctx, params RoomBookingParams) (bool, error) {
	filter := bson.M{
		"fromDate": bson.M{"$gte": params.FromDate},
		"tillDate": bson.M{"$lte": params.TillDate},
	}

	bookings, err := h.store.BookingStore.GetBookings(c.Context(), filter)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}
