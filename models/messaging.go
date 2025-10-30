package models

import "time"

// Conversation represents a communication thread between users.
type Conversation struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Members   []ConversationMember
	Messages  []Message
}

// ConversationMember links a user to a conversation.
type ConversationMember struct {
	ID             uint `gorm:"primaryKey"`
	ConversationID uint `gorm:"not null"`
	UserID         uint `gorm:"not null"`
	User           User `gorm:"foreignKey:UserID"`
}

// Message is a single message within a conversation.
type Message struct {
	ID             uint      `gorm:"primaryKey"`
	ConversationID uint      `gorm:"not null"`
	AuthorID       uint      `gorm:"not null"`
	Author         User      `gorm:"foreignKey:AuthorID"`
	Content        string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	Attachments    []MessageAttachment
	ReadBy         []MessageReadBy
}

// MessageAttachment links an attachment to a message.
type MessageAttachment struct {
	ID           uint       `gorm:"primaryKey"`
	MessageID    uint       `gorm:"not null"`
	AttachmentID uint       `gorm:"not null"`
	Attachment   Attachment `gorm:"foreignKey:AttachmentID"`
}

// MessageReadBy tracks which users have read a message.
type MessageReadBy struct {
	ID        uint      `gorm:"primaryKey"`
	MessageID uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	ReadAt    time.Time `gorm:"autoCreateTime"`
}
