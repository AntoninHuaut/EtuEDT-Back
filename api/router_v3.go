package api

import (
	"github.com/AntoninHuaut/EtuEDT-Back/api/v3"

	"github.com/gofiber/fiber/v2"
)

func v3Router(router fiber.Router) {
	univs := router.Group("/univs")

	univs.Get("/", v3.ListUniversities)
	univs.Get("/:univId", v3.GetUniversity)
	univs.Get("/:univId/groups", v3.ListGroups)
	univs.Get("/:univId/groups/:groupId", v3.ListTimetables)
	univs.Get("/:univId/groups/:groupId/:adeResources", v3.GetTimetableMetadata)
	univs.Get("/:univId/groups/:groupId/:adeResources/events", v3.GetTimetableEvents)
	univs.Get("/:univId/rooms", v3.ListRooms)
	univs.Get("/:univId/rooms/:adeResources", v3.GetRoomMetadata)
	univs.Get("/:univId/rooms/:adeResources/events", v3.GetRoomEvents)
}
