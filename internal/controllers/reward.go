package controllers

import (
	"coursebench-backend/internal/controllers/reward"

	"github.com/gofiber/fiber/v2"
)

func RewardRoutes(r fiber.Router) {
	route := r.Group("/reward")
	route.Get("/ranklist", reward.Ranklist)
	route.Post("/set", reward.SetComment)
}
