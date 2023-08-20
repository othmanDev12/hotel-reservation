package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"github.com/hotel-reservation/util"
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

func (h *RoomHandler) HandleRoomBooking(c *fiber.Ctx) error {
	var params RoomBookingParams
	id := c.Params("id")
	if err := c.BodyParser(&params); err != nil {
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

	book := domain.Book{
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

// 2023-04-15T00:00:00.0Z
