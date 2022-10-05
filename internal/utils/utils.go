package utils

import (
	"coursebench-backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"log"
)

func GetIP(ctx *fiber.Ctx) []string {
	if config.FiberConfig.UseXForwardFor {
		log.Println("UseXForwardFor")
		if len(ctx.IPs()) > 0 {
			log.Println("UseXForwardFor", len(ctx.IPs()))
			return ctx.IPs()
		} else {
			log.Println("UseXForwardFor, but GG")
			return []string{ctx.IP()}
		}
	} else {
		log.Println("Not UseXForwardFor")
		return []string{ctx.IP()}
	}
}
