package models

import "time"

// User — учётная запись пользователя сервиса.
type User struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	Email         string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Username      string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	DisplayName   string    `gorm:"size:128" json:"display_name"`
	AvatarURL     string    `gorm:"size:512" json:"avatar_url"`
	NoisePref     int       `gorm:"default:3" json:"noise_pref"`
	LightPref     int       `gorm:"default:3" json:"light_pref"`
	CrowdPref     int       `gorm:"default:3" json:"crowd_pref"`
	SmellPref     int       `gorm:"default:3" json:"smell_pref"`
	VisualPref    int       `gorm:"default:3" json:"visual_pref"`
	EmailVerified bool      `gorm:"default:false;index" json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PlaceRef — упрощённая ссылка на место из соседнего модуля карты.
// В рамках индивидуального стенда модуля пользователей и отзывов
// используется как минимальный stub: содержит только id и название
// места, чтобы у отзывов и избранного был валидный внешний ключ.
// При интеграции с модулем карты Атаханова Н. Р. эта таблица
// заменяется полноценной таблицей places.
type PlaceRef struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

// TableName переопределяет имя таблицы, чтобы при объединении проектов
// схема совпадала с таблицей places модуля карты.
func (PlaceRef) TableName() string { return "places" }

// Review — отзыв пользователя о месте с сенсорными оценками 1..5.
type Review struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	PlaceID   uint64    `gorm:"index;not null" json:"place_id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Text      string    `gorm:"type:text" json:"text"`
	Noise     int       `gorm:"not null" json:"noise"`
	Light     int       `gorm:"not null" json:"light"`
	Crowd     int       `gorm:"not null" json:"crowd"`
	Smell     int       `gorm:"not null" json:"smell"`
	Visual    int       `gorm:"not null" json:"visual"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Favorite — избранное место пользователя.
type Favorite struct {
	UserID    uint64    `gorm:"primaryKey" json:"user_id"`
	PlaceID   uint64    `gorm:"primaryKey" json:"place_id"`
	CreatedAt time.Time `json:"created_at"`
}

// VerificationCode хранит одноразовый код подтверждения email.
type VerificationCode struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Email        string     `gorm:"index;size:255;not null" json:"email"`
	CodeHash     string     `gorm:"size:60;not null" json:"-"`
	PasswordHash string     `gorm:"size:60;not null" json:"-"`
	Username     string     `gorm:"size:64;not null" json:"-"`
	DisplayName  string     `gorm:"size:128" json:"-"`
	Attempts     int        `gorm:"default:0" json:"-"`
	ExpiresAt    time.Time  `gorm:"index;not null" json:"expires_at"`
	UsedAt       *time.Time `json:"used_at"`
	CreatedAt    time.Time  `json:"created_at"`
}