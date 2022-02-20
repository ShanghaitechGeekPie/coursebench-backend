package tracing

import (
	"context"
	"coursebench-backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer(config.GlobalConf.ServiceName)
}

func GetSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func GetRouterSpan(c *fiber.Ctx) (context.Context, trace.Span) {
	return GetSpan(c.UserContext(), c.Method()+" "+utils.CopyString(c.OriginalURL()))
}
