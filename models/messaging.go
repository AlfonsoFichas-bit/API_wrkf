package models

import "time"

type Conversation struct {
	ID            uint   `gorm:"primaryKey"`
	Type          string `gorm:"not null"`
	Name          string
	Description   string
	ProjectID     *uint
	Project       *Project `gorm:"foreignKey:ProjectID"`
	CreatedByID   uint     `gorm:"not null"`
	CreatedBy     User     `gorm:"foreignKey:CreatedByID"`
	LastMessageAt *time.Time
	IsActive      bool      `gorm:"default:true"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	Members       []ConversationMember
	Messages      []Message
}

type ConversationMember struct {
	ID             uint      `gorm:"primaryKey"`
	ConversationID uint      `gorm:"not null"`
	UserID         uint      `gorm:"not null"`
	User           User      `gorm:"foreignKey:UserID"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	LastReadAt     *time.Time
	IsAdmin        bool      `gorm:"default:false"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type Message struct {
	ID             uint      `gorm:"primaryKey"`
	ConversationID uint      `gorm:"not null"`
	SenderID       uint      `gorm:"not null"`
	Sender         User      `gorm:"foreignKey:SenderID"`
	Content        string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	Attachments    []MessageAttachment
	ReadBy         []MessageReadBy
}

type MessageAttachment struct {
	ID        uint      `gorm:"primaryKey"`
	MessageID uint      `gorm:"not null"`
	FileName  string    `gorm:"not null"`
	FileType  string    `gorm:"not null"`
	FileSize  int       `gorm:"not null"`
	URL       string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type MessageReadBy struct {
	ID        uint      `gorm:"primaryKey"`
	MessageID uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	ReadAt    time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
