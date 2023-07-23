package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
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

func (h *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params = domain.CreateUserParams{}
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.ValidateUser(&params); errors != nil {
		return ctx.JSON(errors)
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
