package service

import (
	"chattrix/dto"
	"chattrix/mapper"
	"chattrix/models"
	"chattrix/repository"
	"chattrix/utils"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type OnlineChecker interface {
	IsOnline(userID uuid.UUID) bool
}

type MessageService struct {
	messageRepo   *repository.MessageRepository
	chatRepo      *repository.ChatRepository
  userRepo      *repository.UserRepository
	onlineChecker OnlineChecker
}

func NewMessageService(messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository, userRepo *repository.UserRepository, onlineChecker OnlineChecker) *MessageService {
	return &MessageService{
		messageRepo:   messageRepo,
		chatRepo:      chatRepo,
		userRepo:      userRepo,
		onlineChecker: onlineChecker,
	}
}

func (s *MessageService) HandleSendMessage(userID uuid.UUID, wsMsg dto.WSMessage) ([]byte, []uuid.UUID, error) {
	var forwardContent *string
  var forwardFiles []models.Attachment
  
  if wsMsg.Content == "" && len(wsMsg.Files) == 0 {
		return nil, nil, errors.New("message can not be empty")
	}

	msgType := wsMsg.Type
	if msgType == "" {
		msgType = "text"
	}

	validTypes := map[string]bool{"text": true, "image": true, "file": true, "voice": true}

	if !validTypes[msgType] {
		return nil, nil, errors.New("invalid message type")
	}

	var content *string
	if wsMsg.Content != "" {
		content = &wsMsg.Content
	}

	if wsMsg.ReplyToID != nil {
    replyMsg, err := s.messageRepo.GetMessageByID(*wsMsg.ReplyToID)
    if err != nil {
      return nil, nil, errors.New("invalid reply message")
    }

    if replyMsg.ChatID != wsMsg.ChatID {
      return nil, nil, errors.New("reply message does not belong to the same chat")
    }
	}

  if wsMsg.ForwardFromMessageID != nil {
    origMsg, err := s.messageRepo.GetMessageWithFullData(*wsMsg.ForwardFromMessageID)
    if err != nil {
      return nil, nil, errors.New("invalid forward message")
    }

    forwardContent = origMsg.Content
    forwardFiles = origMsg.Attachments
  }

	message := models.Message{
		ChatID:           wsMsg.ChatID,
		SenderID:         userID,
		Type:             msgType,
		Content:          content,
		ReplyToMessageID: wsMsg.ReplyToID,
    ForwardFromMessageID: wsMsg.ForwardFromMessageID,
	}

  if wsMsg.ForwardFromMessageID != nil {
    message.Content = forwardContent
  }

	if err := s.messageRepo.CreateMessage(&message); err != nil {
		return nil, nil, err
	}

	for _, f := range wsMsg.Files {
		att := models.Attachment{
			MessageID: message.ID,
			FileURL:   f.FileURL,
			FileType:  f.FileType,
			FileSize:  f.FileSize,
		}

		if err := s.messageRepo.CreateAttachment(&att); err != nil {
			return nil, nil, err
		}
	}

  for _, f := range forwardFiles {
    att := models.Attachment {
      MessageID: message.ID,
      FileURL: f.FileURL,
      FileType: f.FileType,
      FileSize: f.FileSize,
    }

    if err := s.messageRepo.CreateAttachment(&att); err != nil {
      return nil, nil, err
    }
  }

	memberIDs, err := s.chatRepo.GetChatMembers(wsMsg.ChatID)
	if err != nil {
		return nil, nil, err
	}

	onlineMap := make(map[uuid.UUID]bool)

	for _, uid := range memberIDs {
		var status string

    if uid == userID {
      status = "sent"
      onlineMap[uid] = true
    } else {
      if s.onlineChecker.IsOnline(uid) {
        status = "delivered"
        onlineMap[uid] = true
      } else {
        status = "sent"
      }
    }

		err := s.messageRepo.CreateStatus(&models.MessageStatus{
			MessageID: message.ID,
			UserID:    uid,
			Status:    status,
		})

		if err != nil {
			log.Println("create status error:", err)
			return nil, nil, err
		}
	}

  memberSet := make(map[uuid.UUID]bool)
  for _, id := range memberIDs {
    memberSet[id] = true
  }

  if wsMsg.Content != "" {
    usernames := utils.ExtractMentions(wsMsg.Content)
    for _, username := range usernames {
      user, err := s.userRepo.GetByUsername(username)
      if err != nil {
        continue
      }

      if memberSet[user.ID] {
        onlineMap[user.ID] = true
      }
    }
  }

  fullMsg, err := s.messageRepo.GetMessageWithFullData(message.ID)
  if err != nil {
    return nil, nil, err
  }

  responseDto := mapper.ToMessageResponse(fullMsg)

  response, _ := json.Marshal(map[string]interface{}{
    "type":    "message",
    "message": responseDto,
  })

	var onlineUsers []uuid.UUID
	for uid := range onlineMap {
		onlineUsers = append(onlineUsers, uid)
	}


	return response, onlineUsers, nil
}

func (s *MessageService) HandleSeen(userID uuid.UUID, chatID uuid.UUID) ([]byte, []uuid.UUID, error) {
	lastMsgID, err := s.messageRepo.MarkMessagesAsSeen(chatID, userID)
  if err != nil {
    return nil, nil, err
  }
  
	memberIDs, err := s.chatRepo.GetChatMembers(chatID)
	if err != nil {
		return nil, nil, err
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":    "seen",
		"user_id": userID,
		"chat_id": chatID,
    "last_seen_message_id": lastMsgID,
	})

	var receivers []uuid.UUID
	for _, uid := range memberIDs {
		if uid != userID && s.onlineChecker.IsOnline(uid) {
			receivers = append(receivers, uid)
		}
	}

	return response, receivers, nil
}

func (s *MessageService) HandleTyping(userID uuid.UUID, chatID uuid.UUID, isTyping bool) ([]byte, []uuid.UUID, error) {
	memberIDs, err := s.chatRepo.GetChatMembers(chatID)
	if err != nil {
		return nil, nil, err
	}

	eventType := "typing"
	if !isTyping {
		eventType = "stop_typing"
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":    eventType,
		"user_id": userID,
		"chat_id": chatID,
	})

	var receivers []uuid.UUID
	for _, uid := range memberIDs {
		if uid != userID && s.onlineChecker.IsOnline(uid) {
			receivers = append(receivers, uid)
		}
	}

	return response, receivers, nil
}

func (s *MessageService) HandleUnreadMessages(userID uuid.UUID) ([]byte, error) {
	messages, err := s.messageRepo.GetUnreadMessages(userID)
	if err != nil {
		return nil, err
	}

	_ = s.messageRepo.MarkAsDelivered(userID)

	response, err := json.Marshal(map[string]interface{}{
		"type":     "unread_messages",
		"messages": messages,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *MessageService) GetMessages(chatID uuid.UUID, limit, offset int) ([]models.Message, error) {
	return s.messageRepo.GetMessages(chatID, limit, offset)
}

func (s *MessageService) EditMessage(userID, messageID uuid.UUID, content string) error {
	msg, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	if msg.SenderID != userID {
		return errors.New("not allowed")
	}

	return s.messageRepo.UpdateMessageContent(messageID, content)
}

func (s *MessageService) DeleteMessage(userID, messageID uuid.UUID) error {
	msg, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	if msg.SenderID != userID {
		return errors.New("not allowed")
	}

	return s.messageRepo.SoftDelete(messageID)
}

func (s *MessageService) ToggleReaction(userID, messageID uuid.UUID, reaction string) error {
	_, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

  exists, err := s.messageRepo.ReactionExists(messageID, userID, reaction)
  if err != nil {
    return err
  }

  if exists {
    return s.messageRepo.RemoveReaction(messageID, userID, reaction)
  }

	return s.messageRepo.AddReaction(&models.MessageReaction{
		MessageID: messageID,
		UserID:    userID,
		Reaction:  reaction,
	})
}

func (s *MessageService) GetMembersByMessage(messageID uuid.UUID) ([]uuid.UUID, error) {
	msg, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	return s.chatRepo.GetChatMembers(msg.ChatID)
}

func (s *MessageService) EditMessageRealtime(userID, messageID uuid.UUID, content string) ([]byte, []uuid.UUID, error) {
	err := s.EditMessage(userID, messageID, content)
	if err != nil {
		return nil, nil, err
	}

	msg, err := s.messageRepo.GetMessageWithFullData(messageID)
  if err != nil {
    return nil, nil, err
  }

	memberIDs, err := s.chatRepo.GetChatMembers(msg.ChatID)
  if err != nil {
    return nil, nil, err
  }

  var receivers []uuid.UUID
  for _, uid := range memberIDs {
    if uid != userID && s.onlineChecker.IsOnline(uid) {
      receivers = append(receivers, uid)
    }
  }

  responseDTO := mapper.ToMessageResponse(msg)

	response, _ := json.Marshal(map[string]interface{}{
		"type":       "edit",
		"message": responseDTO,
	})

	return response, receivers, nil
}

func (s *MessageService) DeleteMessageRealtime(userID, messageID uuid.UUID) ([]byte, []uuid.UUID, error) {
  msg, err := s.messageRepo.GetMessageByID(messageID)
	if err != nil {
    return nil, nil, err
  }

  if err := s.DeleteMessage(userID, messageID); err != nil {
		return nil, nil, err
	}

	memberIDs, err := s.chatRepo.GetChatMembers(msg.ChatID)
	if err != nil {
		return nil, nil, err
	}

  var receivers []uuid.UUID
  for _, uid := range memberIDs {
    if uid != userID && s.onlineChecker.IsOnline(uid) {
      receivers = append(receivers, uid)
    }
  }

  updatedMsg, err := s.messageRepo.GetMessageWithFullData(messageID)
  if err != nil {
    return nil, nil, err
  }

  responseDTO := mapper.ToMessageResponse(updatedMsg)

	response, _ := json.Marshal(map[string]interface{}{
		"type":       "delete",
		"message": responseDTO,
	})

	return response, receivers, nil
}

func (s *MessageService) HandleReactionRealtime(userID uuid.UUID, dto dto.WSReaction) ([]byte, []uuid.UUID, error) {
	err := s.ToggleReaction(userID, dto.MessageID, dto.Reaction)
	if err != nil {
		return nil, nil, err
	}

	msg, err := s.messageRepo.GetMessageWithFullData(dto.MessageID)
	if err != nil {
		return nil, nil, err
	}

	memberIDs, _ := s.chatRepo.GetChatMembers(msg.ChatID)

	var receivers []uuid.UUID
	for _, uid := range memberIDs {
		if s.onlineChecker.IsOnline(uid) {
			receivers = append(receivers, uid)
		}
	}

	responseDTO := mapper.ToMessageResponse(msg)

	response, _ := json.Marshal(map[string]interface{}{
		"type":    "reaction",
		"message": responseDTO,
	})

	return response, receivers, nil
}

func (s *MessageService) GetPaginatedMessages(chatID uuid.UUID, before *time.Time, limit int) ([]dto.MessageResponse, error) {
  msgs, err := s.messageRepo.GetMessagesByChat(chatID, before, limit)
  if err != nil {
    return nil, err
  }

  var responses []dto.MessageResponse

  for _, m := range msgs {
    responses = append(responses, mapper.ToMessageResponse(&m))
  }

  return responses, nil
}