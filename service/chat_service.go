package service

import (
	"chattrix/dto"
	"chattrix/models"
	"chattrix/repository"
	"errors"

	"github.com/google/uuid"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
  return &ChatService{chatRepo: chatRepo}
}

func (s *ChatService) CreateChat(userID uuid.UUID, chatRequest dto.CreateChatRequest) (*models.Chat, error) {
	if chatRequest.IsGroup {

		if len(chatRequest.UserIDs) < 1 {
			return nil, errors.New("group must have at least 1 user")
		}

		if chatRequest.Name == nil || *chatRequest.Name == "" {
			return nil, errors.New("group name is required")
		}

	} else {
		if len(chatRequest.UserIDs) != 1 {
			return nil, errors.New("private chat must have exactly 1 user")
		}
	}

	chat := models.Chat{
		IsGroup:   chatRequest.IsGroup,
		Name:      chatRequest.Name,
		CreatedBy: &userID,
	}

	if err := s.chatRepo.CreateChat(&chat); err != nil {
		return nil, err
	}

	memberIDs := append(chatRequest.UserIDs, userID)

	if err := s.chatRepo.AddChatMembers(chat.ID, memberIDs); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (s *ChatService) GetUserChats(userID uuid.UUID) ([]dto.ChatResponse, error) {
	return s.chatRepo.GetUserChats(userID)
}

func (s *ChatService) AddUsers(requesterID, chatID uuid.UUID, userIDs []uuid.UUID) error {
	isAdmin, err := s.isAdmin(chatID, requesterID)
	if err != nil || !isAdmin {
		return errors.New("only admin can add users")
	}

	return s.chatRepo.AddChatMembers(chatID, userIDs)
}

func (s *ChatService) RemoveUser(requesterID, chatID, userID uuid.UUID) error {
	isAdmin, err := s.isAdmin(chatID, requesterID)
	if err != nil || !isAdmin {
		return errors.New("only admin can remove users")
	}

	return s.chatRepo.RemoveUser(chatID, userID)
}

func (s *ChatService) LeaveChat(userID, chatID uuid.UUID) error {
	isMember, err := s.chatRepo.IsMember(chatID, userID)
	if err != nil || !isMember {
		return errors.New("not a member of this chat")
	}

	if err := s.chatRepo.LeaveChat(chatID, userID); err != nil {
		return err
	}

	count, err := s.chatRepo.CountMembers(chatID)
	if err != nil {
		return err
	}

	if count == 0 {
		return s.chatRepo.DeleteChat(chatID)
	}

	return nil
}

func (s *ChatService) PinChat(userID, chatID uuid.UUID, pinned bool) error {
	isMember, err := s.chatRepo.IsMember(chatID, userID)
	if err != nil || !isMember {
		return errors.New("not authorized")
	}

	return s.chatRepo.SetPinned(chatID, userID, pinned)
}

func (s *ChatService) MuteChat(userID, chatID uuid.UUID, muted bool) error {
	isMember, err := s.chatRepo.IsMember(chatID, userID)
	if err != nil || !isMember {
		return errors.New("not authorized")
	}

	return s.chatRepo.SetMuted(chatID, userID, muted)
}

func (s *ChatService) isAdmin(chatID, userID uuid.UUID) (bool, error) {
	role, err := s.chatRepo.GetUserRole(chatID, userID)
	if err != nil {
		return false, err
	}

	return role == "admin", nil
}

func (s *ChatService) ChangeUserRole(requesterID, chatID, targetUserID uuid.UUID, role string) error {
	isAdmin, err := s.isAdmin(chatID, requesterID)
	if err != nil || !isAdmin {
		return errors.New("only admin can change roles")
	}

	if role != "admin" && role != "member" {
		return errors.New("invalid role")
	}

	return s.chatRepo.UpdateUserRole(chatID, targetUserID, role)
}

func (s *ChatService) DeleteChat(requesterID, chatID uuid.UUID) error {
	_, err := s.chatRepo.GetByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

  isAdmin, err := s.isAdmin(chatID, requesterID)
	if err != nil || !isAdmin {
		return errors.New("only admin can delete chat")
	}

	return s.chatRepo.DeleteChat(chatID)
}

func (s *ChatService) SearchChats(userID uuid.UUID, query string) ([]dto.ChatResponse, error) {
	return s.chatRepo.SearchChats(userID, query)
}