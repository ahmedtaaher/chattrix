package mapper

import (
	"chattrix/dto"
	"chattrix/models"
)

func ToMessageResponse(msg *models.Message) dto.MessageResponse {
  return ToMessageResponseWithDepth(msg, 1)
}

func ToMessageResponseWithDepth(msg *models.Message, depth int) dto.MessageResponse {
	var content *string = msg.Content
  if msg.IsDeleted {
    content = nil
  }
  
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

  var reply *dto.MessageResponse
  if depth > 0 && msg.ReplyToMessage != nil {
    r := ToMessageResponseWithDepth(msg.ReplyToMessage, depth-1)
    reply = &r
  }

	return dto.MessageResponse{
		ID:          msg.ID,
		ChatID:      msg.ChatID,
		SenderID:    msg.SenderID,
		Type:        msg.Type,
		Content:     content,
		SentAt:      msg.SentAt,
		EditedAt:    msg.EditedAt,
		IsDeleted:   msg.IsDeleted,
		Attachments: attachments,
		Reactions:   reactions,
    ReplyTo:     reply,
	}
}