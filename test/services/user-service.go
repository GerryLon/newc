package services

import (
	"log"

	"github.com/GerryLon/newc/test/repositories"
)

// UserService a user service for example
//go:generate go run ../../../newc
type UserService struct {
	baseService

	userRepository *repositories.UserRepository
	proRepository  *repositories.ProRepository

	emailService *EmailService

	logger *log.Logger

	debug bool
}
