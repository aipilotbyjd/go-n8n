package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity
type User struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash      string     `json:"-" gorm:"not null"`
	Name              string     `json:"name" gorm:"not null"`
	Role              Role       `json:"role" gorm:"default:'user'"`
	IsActive          bool       `json:"is_active" gorm:"default:true"`
	EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	ProfilePicture    string     `json:"profile_picture,omitempty"`
	Settings          UserSettings `json:"settings" gorm:"serializer:json"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Role represents user role
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleOwner Role = "owner"
)

// UserSettings contains user-specific settings
type UserSettings struct {
	Theme           string `json:"theme"`           // light, dark, auto
	Language        string `json:"language"`        // en, es, de, etc.
	Timezone        string `json:"timezone"`
	DateFormat      string `json:"date_format"`
	FirstDayOfWeek  int    `json:"first_day_of_week"`  // 0=Sunday, 1=Monday
	ExecutionTimeout int   `json:"execution_timeout"`   // in seconds
	SaveExecutions  bool   `json:"save_executions"`
	NotifyOnError   bool   `json:"notify_on_error"`
	NotifyOnSuccess bool   `json:"notify_on_success"`
}

// Session represents a user session
type Session struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Token        string     `json:"token" gorm:"uniqueIndex;not null"`
	RefreshToken string     `json:"refresh_token" gorm:"uniqueIndex"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	ExpiresAt    time.Time  `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
	LastUsedAt   time.Time  `json:"last_used_at"`
}

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Name       string     `json:"name" gorm:"not null"`
	KeyHash    string     `json:"-" gorm:"uniqueIndex;not null"`
	KeyPreview string     `json:"key_preview"` // First 8 chars for identification
	Scopes     []string   `json:"scopes" gorm:"type:text[]"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Team represents a team entity
type Team struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	OwnerID     uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	Settings    TeamSettings   `json:"settings" gorm:"serializer:json"`
	Members     []TeamMember   `json:"members" gorm:"foreignKey:TeamID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TeamMember represents a team member
type TeamMember struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TeamID    uuid.UUID  `json:"team_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Role      TeamRole   `json:"role" gorm:"not null"`
	JoinedAt  time.Time  `json:"joined_at"`
}

// TeamRole represents a role within a team
type TeamRole string

const (
	TeamRoleMember TeamRole = "member"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleOwner  TeamRole = "owner"
)

// TeamSettings contains team-specific settings
type TeamSettings struct {
	MaxWorkflows     int  `json:"max_workflows"`
	MaxExecutions    int  `json:"max_executions"`
	ShareCredentials bool `json:"share_credentials"`
	ShareVariables   bool `json:"share_variables"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	now := time.Now()
	u.PasswordChangedAt = &now
	return nil
}

// CheckPassword verifies the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// VerifyEmail marks the user's email as verified
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	now := time.Now()
	u.EmailVerifiedAt = &now
	u.UpdatedAt = now
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// IsPasswordExpired checks if password needs to be changed (e.g., after 90 days)
func (u *User) IsPasswordExpired(maxAge time.Duration) bool {
	if u.PasswordChangedAt == nil {
		return true
	}
	return time.Since(*u.PasswordChangedAt) > maxAge
}

// CanAccessWorkflow checks if user can access a workflow
func (u *User) CanAccessWorkflow(workflowOwnerID uuid.UUID) bool {
	return u.ID == workflowOwnerID || u.Role == RoleAdmin || u.Role == RoleOwner
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(permission string) bool {
	switch u.Role {
	case RoleOwner:
		return true // Owners have all permissions
	case RoleAdmin:
		// Admins have most permissions except system-level ones
		return permission != "system:manage"
	case RoleUser:
		// Regular users have limited permissions
		allowedPermissions := []string{
			"workflow:read",
			"workflow:create",
			"workflow:update",
			"workflow:delete",
			"workflow:execute",
			"credential:manage",
			"variable:manage",
		}
		for _, allowed := range allowedPermissions {
			if permission == allowed {
				return true
			}
		}
		return false
	default:
		return false
	}
}
