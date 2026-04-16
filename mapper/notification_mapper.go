package mapper

import (
	"chattrix/dto"
	"chattrix/models"
)


func ToNotificationResponse(n *models.Notification) dto.NotificationResponse {
	if n == nil {
		return dto.NotificationResponse{}
	}

	title := ""
	body := ""

	switch n.Type {

	case "message":
		title = "New Message"
		body = "You received a new message"

	case "mention":
		title = "Mention"
		body = "You were mentioned in a message"

	case "invite":
		title = "Chat Invite"
		body = "You were invited to a chat"
	}

	return dto.NotificationResponse{
		ID:          n.ID,
		Type:        n.Type,
		ReferenceID: n.ReferenceID,
		IsRead:      n.IsRead,
		CreatedAt:   n.CreatedAt,
		Title:       title,
		Body:        body,
	}
}


func ToNotificationResponseList(notifs []models.Notification) []dto.NotificationResponse {
	var result []dto.NotificationResponse

	for _, n := range notifs {
		result = append(result, ToNotificationResponse(&n))
	}

	return result
}