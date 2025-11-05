package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	// --- Public routes ---
	auth := v1.Group("/auth")
	{
		auth.POST("/register", app.registerUser)
		auth.POST("/login", app.login)
	}

	// Publicly accessible routes (if you want GET events public)
	eventsPublic := v1.Group("/events")
	{
		eventsPublic.GET("", app.getAllEvents)
		eventsPublic.GET("/:id", app.getEvent)
		eventsPublic.GET("/:id/attendees", app.getAttendeesForEvent)
	}

	// --- Protected routes (require JWT) ---
	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())

	// Protected event routes
	events := authGroup.Group("/events")
	{
		events.POST("", app.createEvent)
		events.PUT("/:id", app.updateEvent)
		events.DELETE("/:id", app.deleteEvent)

		// attendees under a specific event
		events.POST("/:id/attendees/:userId", app.addAttendeeToEvent)
		events.DELETE("/:id/attendees/:userId", app.deleteAttendeeFromEvent)
	}

	// Protected attendee routes
	attendees := authGroup.Group("/attendees")
	{
		attendees.GET("/:id/events", app.getEventsByAttendee)
	}

	g.GET("/swagger/*any", func(ctx *gin.Context) {
		if ctx.Request.RequestURI == "/swagger/" {
			ctx.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(ctx)
	})

	return g
}
