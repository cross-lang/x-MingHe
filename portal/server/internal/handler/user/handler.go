package user

import (
	"portal/internal/service/user"
)

type Handler struct {
	UserService *user.Service
}
