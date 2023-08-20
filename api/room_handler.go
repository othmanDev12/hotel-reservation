package api

import (
	"fmt"
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
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
}

// NewRoomHandler creates a new constructor for the room handler
func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (h *RoomHandler) HandleRoomBooking(c *fiber.Ctx) error {
	id := c.Params("id")
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
		UserId: userId.Id,
		RoomId: roomId,
	}

	fmt.Println(book)
	return nil
}
