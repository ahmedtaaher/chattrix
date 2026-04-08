package service

import (
	"chattrix/dto"
	"chattrix/models"
	"chattrix/repository"
	"encoding/json"

	"github.com/google/uuid"
)

type OnlineChecker interface {
  IsOnline(userID uuid.UUID) bool
}

type MessageService struct {
	messageRepo    *repository.MessageRepository
  chatRepo       *repository.ChatRepository
  onlineChecker  OnlineChecker
}

func NewMessageService(messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository, onlineChecker OnlineChecker) *MessageService {
	return &MessageService{
		messageRepo:    messageRepo,
		chatRepo:       chatRepo,
		onlineChecker:  onlineChecker,
	}
}

func (s *MessageService) HandleSendMessage(userID uuid.UUID, wsMsg dto.WSMessage) ([]byte, []uuid.UUID, error) {
	message := models.Message{
		ChatID:   wsMsg.ChatID,
		SenderID: userID,
		Type:     "text",
		Content:  &wsMsg.Content,
	}

	if err := s.messageRepo.CreateMessage(&message); err != nil {
		return nil, nil, err
	}

	memberIDs, err := s.chatRepo.GetChatMembers(wsMsg.ChatID)
	if err != nil {
		return nil, nil, err
	}

  var onlineUsers []uuid.UUID

	for _, uid := range memberIDs {

		status := "sent"

		if uid != userID && s.onlineChecker.IsOnline(uid) {
			status = "delivered"
			onlineUsers = append(onlineUsers, uid)
		}

		_ = s.messageRepo.CreateStatus(&models.MessageStatus{
			MessageID: message.ID,
			UserID:    uid,
			Status:    status,
		})
	}

  if s.onlineChecker.IsOnline(userID) {
		onlineUsers = append(onlineUsers, userID)
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":    "message",
		"message": message,
	})

	return response, onlineUsers, nil
}

func (s *MessageService) HandleSeen(userID uuid.UUID,chatID uuid.UUID) ([]byte, []uuid.UUID, error) {
	if err := s.messageRepo.MarkMessagesAsSeen(chatID, userID); err != nil {
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

