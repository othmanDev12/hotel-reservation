package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

// NewUserHandler creates a new constructor for the user handler
func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

// HandleGetUser handle a get user api that accepts a function receiver
func (h *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")
	user, err := h.userStore.GetUserById(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.JSON(map[string]string{"error": "user not found"})
		}
		return err
	}
	return ctx.JSON(user)
}

func (h *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(users)
}

func (h *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")
	err := h.userStore.DeleteUser(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"message": "user has been deleted successfully"})
}

func (h *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params = domain.CreateUserParams{}
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if errV := params.ValidateUser(&params); errV != nil {
		return ctx.JSON(errV)
	}
	user, err := domain.NewCreateUser(params)
	if err != nil {
		return err
	}
	insertedValue, err := h.userStore.CreateUser(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedValue)
}
