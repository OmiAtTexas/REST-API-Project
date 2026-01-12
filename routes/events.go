package routes

import (
	"net/http"
	"strconv"

	"example.com/REST-API-Project/models"
	"github.com/gin-gonic/gin"
)

// Get all events
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch events. Try again later.",
		})
		return
	}

	context.JSON(http.StatusOK, events)
}

// Get single event by ID
func getEvent(context *gin.Context) {
	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch event.",
		})
		return
	}

	context.JSON(http.StatusOK, event)
}

// Create a new event
func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return
	}

	// Get userId from context safely
	userIdInterface, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in"})
		return
	}

	userId, ok := userIdInterface.(int64)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID"})
		return
	}

	event.UserID = userId

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create event. Try again later.",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created",
		"event":   event,
	})
}

// Update an existing event
func updateEvent(context *gin.Context) {
	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch the event.",
		})
		return
	}

	// Get userId from context safely
	userIdInterface, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in"})
		return
	}

	userId, ok := userIdInterface.(int64)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID"})
		return
	}

	// Only the owner can update
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event."})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return
	}

	updatedEvent.ID = eventId

	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update event.",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully.",
		"event":   updatedEvent,
	})
}

// Delete an event
func deleteEvent(context *gin.Context) {
	eventIdStr := context.Param("id")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch the event.",
		})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to get event."})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete the event.",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully!",
	})
}
