package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/middleware/authorize"
	"github.com/newhorizon-tech-vn/tracing-example/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
)

func (h *Handler) CheckClass(c *gin.Context) {

	classId, err := strconv.Atoi(c.Param("classId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	userId, exists := c.Get(authorize.KeyUserId)
	if exists == false {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	roleId, exists := c.Get(authorize.KeyRoleId)
	if exists == false {
		c.JSON(http.StatusNonAuthoritativeInfo, nil)
		return
	}

	startChildSpan(c.Request.Context(), userId.(int))

	attachSpanLog(c.Request.Context(), userId.(int))

	classIds, err := (&services.ClassService{}).GetClassIds(c.Request.Context(), userId.(int), roleId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	index := slices.IndexFunc(classIds, func(id int) bool { return id == classId })
	if index < 0 {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func startChildSpan(ctx context.Context, userId int) {
	// custom child span
	_, span := otel.GetTracerProvider().Tracer("dms-api").Start(ctx, "task-name", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// do something
	time.Sleep(200 * time.Millisecond)

	if span.IsRecording() {
		span.SetAttributes(
			attribute.String("nhz.child.http.method", "GET"),
			attribute.String("nhz.child.http.route", "/projects/:id"),
		)
	}
	span.AddEvent("nhz.child.log", trace.WithAttributes(
		attribute.String("nhz.child.log.severity", "error"),
		attribute.String("nhz.child.log.message", "User not found"),
		attribute.Int("nhz.child.user_id", userId),
	))

	// To mark the entire operation as an error, set the status.
	// Note that recording an error does not automatically change
	// the status.
	span.SetStatus(codes.Error, fmt.Errorf("mock error").Error())
}

func attachSpanLog(ctx context.Context, userId int) {
	// custom child span
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(
			attribute.String("nhz.http.method", "GET"),
			attribute.String("nhz.http.route", "/user/:userId"),
		)
	}
	span.AddEvent("nhz.action1", trace.WithAttributes(
		attribute.String("nhz.action1.severity", "error"),
		attribute.String("nhz.action1.message", "User not found"),
		attribute.Int("nhz.action1.user_id", userId),
	))

	// To mark the entire operation as an error, set the status.
	// Note that recording an error does not automatically change
	// the status.
	span.SetStatus(codes.Error, fmt.Errorf("mock error").Error())
}
