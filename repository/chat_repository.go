package repository

import (
	"chattrix/dto"
	"chattrix/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
  return &ChatRepository{db: db}
}

func (c *ChatRepository) GetChatMembers(chatID uuid.UUID) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID

	err := c.db.Table("chat_members").
		Where("chat_id = ?", chatID).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}

func (c *ChatRepository) CreateChat(chat *models.Chat) error {
  return c.db.Create(chat).Error
}

func (c *ChatRepository) AddChatMembers(chatID uuid.UUID, userIDs []uuid.UUID) error {
  var members []models.ChatMember
  for _, uid := range userIDs {
    members = append(members, models.ChatMember{
      ChatID: chatID,
      UserID: uid,
    })
  }
  return c.db.Create(&members).Error
}

func (c *ChatRepository) GetUserChats(userID uuid.UUID) ([]dto.ChatResponse, error) {
	var chats []dto.ChatResponse

	query := `
	SELECT 
		c.id AS chat_id,
		CASE 
			WHEN c.is_group THEN c.name
			ELSE u.username
		END AS name,
		c.is_group,
		m.content AS last_message,
		COUNT(ms.message_id) FILTER (WHERE ms.status != 'seen' AND ms.user_id = ?) AS unread_count

	FROM chats c

	JOIN chat_members cm ON cm.chat_id = c.id

	-- For private chat name (other user)
	LEFT JOIN chat_members cm2 ON cm2.chat_id = c.id AND cm2.user_id != ?
	LEFT JOIN users u ON u.id = cm2.user_id

	-- Last message
	LEFT JOIN LATERAL (
		SELECT content
		FROM messages
		WHERE chat_id = c.id
		ORDER BY sent_at DESC
		LIMIT 1
	) m ON true

	-- Unread count
	LEFT JOIN message_status ms ON ms.message_id IN (
		SELECT id FROM messages WHERE chat_id = c.id
	)

	WHERE cm.user_id = ?

	GROUP BY c.id, c.name, c.is_group, u.username, m.content
	ORDER BY MAX(m.content) DESC;
	`

	err := c.db.Raw(query, userID, userID, userID).Scan(&chats).Error
	return chats, err
}

func (c *ChatRepository) RemoveUser(chatID, userID uuid.UUID) error {
	return c.db.
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Delete(&models.ChatMember{}).Error
}

func (c *ChatRepository) IsMember(chatID, userID uuid.UUID) (bool, error) {
	var count int64

	err := c.db.Model(&models.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Count(&count).Error

	return count > 0, err
}

func (c *ChatRepository) LeaveChat(chatID, userID uuid.UUID) error {
	return c.db.
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Delete(&models.ChatMember{}).Error
}

func (c *ChatRepository) CountMembers(chatID uuid.UUID) (int64, error) {
	var count int64

	err := c.db.Model(&models.ChatMember{}).
		Where("chat_id = ?", chatID).
		Count(&count).Error

	return count, err
}

func (c *ChatRepository) DeleteChat(chatID uuid.UUID) error {
	return c.db.Delete(&models.Chat{}, "id = ?", chatID).Error
}

func (c *ChatRepository) SetPinned(chatID, userID uuid.UUID, pinned bool) error {
	return c.db.Model(&models.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("is_pinned", pinned).Error
}

func (c *ChatRepository) SetMuted(chatID, userID uuid.UUID, muted bool) error {
	return c.db.Model(&models.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("is_muted", muted).Error
}

func (c *ChatRepository) GetUserRole(chatID, userID uuid.UUID) (string, error) {
	var role string

	err := c.db.Model(&models.ChatMember{}).
		Select("role").
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Scan(&role).Error

	return role, err
}

func (c *ChatRepository) UpdateUserRole(chatID, userID uuid.UUID, role string) error {
	return c.db.Model(&models.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("role", role).Error
}

func (c *ChatRepository) GetByID(chatID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat

	err := c.db.First(&chat, "id = ?", chatID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chat not found")
		}
		return nil, err
	}

	return &chat, nil
}

func (c *ChatRepository) SearchChats(userID uuid.UUID, query string) ([]dto.ChatResponse, error) {
  var chats []dto.ChatResponse

	sql := `
	SELECT 
		c.id AS chat_id,
		CASE 
			WHEN c.is_group THEN c.name
			ELSE u.username
		END AS name,
		c.is_group,
		NULL AS last_message,
		0 AS unread_count

	FROM chats c

	JOIN chat_members cm ON cm.chat_id = c.id

	LEFT JOIN chat_members cm2 ON cm2.chat_id = c.id AND cm2.user_id != ?
	LEFT JOIN users u ON u.id = cm2.user_id

	WHERE cm.user_id = ?
	AND (
		c.name ILIKE '%' || ? || '%' OR
		u.username ILIKE '%' || ? || '%'
	)
	`

	err := c.db.Raw(sql, userID, userID, query, query).Scan(&chats).Error

	return chats, err
}

func (c *ChatRepository) CreateInvite(invite *models.ChatInvite) error {
	return c.db.Create(invite).Error
}

func (c *ChatRepository) GetInviteByCode(code string) (*models.ChatInvite, error) {
	var invite models.ChatInvite

	err := c.db.
		Where("invite_code = ?", code).
		First(&invite).Error

	if err != nil {
		return nil, err
	}

	return &invite, nil
}
