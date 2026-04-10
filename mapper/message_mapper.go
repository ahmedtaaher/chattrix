package mapper

import (
	"chattrix/dto"
	"chattrix/models"
)

func ToMessageResponse(msg *models.Message) dto.MessageResponse {
	var attachments []dto.AttachmentResponse
	for _, a := range msg.Attachments {
		attachments = append(attachments, dto.AttachmentResponse{
			FileURL:  a.FileURL,
			FileType: a.FileType,
			FileSize: a.FileSize,
		})
	}

	reactionMap := make(map[string]int)
	for _, r := range msg.Reactions {
		reactionMap[r.Reaction]++
	}

	var reactions []dto.ReactionResponse
	for k, v := range reactionMap {
		reactions = append(reactions, dto.ReactionResponse{
			Reaction: k,
			Count:    v,
		})
	}

	return dto.MessageResponse{
		ID:          msg.ID,
		ChatID:      msg.ChatID,
		SenderID:    msg.SenderID,
		Type:        msg.Type,
		Content:     msg.Content,
		SentAt:      msg.SentAt,
		EditedAt:    msg.EditedAt,
		IsDeleted:   msg.IsDeleted,
		Attachments: attachments,
		Reactions:   reactions,
	}
}