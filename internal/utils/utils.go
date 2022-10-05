package utils

import (
	"coursebench-backend/internal/config"
	"github.com/gofiber/fiber/v2"
)

func GetIP(ctx *fiber.Ctx) []string {
	if config.FiberConfig.UseXForwardFor {
		if len(ctx.IPs()) > 0 {
			return ctx.IPs()
		} else {
			return []string{ctx.IP()}
		}
	} else {
		return []string{ctx.IP()}
	}
}
