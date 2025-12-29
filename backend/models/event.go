package models

import (
	"fmt"
	"strings"
	"time"

	"shipshipship/utils"

	"gorm.io/gorm"
)

type EventStatus string

// EventStatusDefinition stores metadata for user-defined statuses.
// All statuses are created/managed by admins.
type EventStatusDefinition struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DisplayName string    `json:"display_name" gorm:"not null;uniqueIndex"` // human-friendly name
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`         // URL-friendly identifier
	Order       int       `json:"order" gorm:"default:0"`                   // display ordering
	IsReserved  bool      `json:"is_reserved" gorm:"default:false"`         // kept for backward compatibility
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex"`
	Color     string    `json:"color" gorm:"not null;default:#3B82F6"` // Default blue color
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Events    []Event   `json:"-" gorm:"many2many:event_tags;"`
}

type Event struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	Title        string            `json:"title" gorm:"not null"`
	Slug         string            `json:"slug" gorm:"uniqueIndex"`
	Tags         []Tag             `json:"tags" gorm:"many2many:event_tags;"`
	Media        string            `json:"media"` // JSON string of array
	Status       EventStatus       `json:"status" gorm:"not null"`
	Date         string            `json:"date"`
	Votes        int               `json:"votes" gorm:"default:0"`
	Content      string            `json:"content"` // Markdown content
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `json:"-" gorm:"index"`
	IsPublic     bool              `json:"is_public" gorm:"default:true"`      // Controls if event appears on public page
	HasPublicUrl bool              `json:"has_public_url" gorm:"default:true"` // Controls if event has individual public URL
	Publication  *EventPublication `json:"publication,omitempty" gorm:"foreignKey:EventID"`
}

type EventPublication struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	EventID         uint       `json:"event_id" gorm:"not null;uniqueIndex"`
	EmailSent       bool       `json:"email_sent" gorm:"default:false"`
	EmailSubject    string     `json:"email_subject"`
	EmailContent    string     `json:"email_content"`
	EmailTemplate   string     `json:"email_template"` // "upcoming_feature" or "new_release"
	EmailSentAt     *time.Time `json:"email_sent_at"`
	SubscriberCount int        `json:"subscriber_count" gorm:"default:0"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// EventEmailHistory tracks the history of all emails sent for an event
type EventEmailHistory struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	EventID         uint      `json:"event_id" gorm:"not null;index"`
	EventStatus     string    `json:"event_status"`
	EmailSubject    string    `json:"email_subject"`
	EmailTemplate   string    `json:"email_template"` // "upcoming_feature" or "new_release"
	SubscriberCount int       `json:"subscriber_count" gorm:"default:0"`
	SentAt          time.Time `json:"sent_at"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type UpdateTagRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

type CreateEventRequest struct {
	Title   string      `json:"title" binding:"required"`
	TagIDs  []uint      `json:"tag_ids"` // Array of tag IDs instead of strings
	Media   []string    `json:"media"`
	Status  EventStatus `json:"status" binding:"required"`
	Date    string      `json:"date"`
	Content string      `json:"content"`
}

type UpdateEventRequest struct {
	Title   *string      `json:"title"`
	TagIDs  *[]uint      `json:"tag_ids"` // Pointer to array of tag IDs to distinguish nil from empty
	Media   []string     `json:"media"`
	Status  *EventStatus `json:"status"`
	Date    *string      `json:"date"`
	Content *string      `json:"content"`
}

type VoteRequest struct {
	EventID uint `json:"event_id" binding:"required"`
}

type EventPublishRequest struct {
	IsPublic     *bool `json:"is_public"`
	HasPublicUrl *bool `json:"has_public_url"`
}

type EventNewsletterRequest struct {
	Subject  string `json:"subject" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Template string `json:"template" binding:"required"`
}

// Requests for status definition management (admin CRUD)
type CreateStatusDefinitionRequest struct {
	DisplayName string  `json:"display_name" binding:"required"`
	Order       *int    `json:"order"`       // optional explicit order
	CategoryID  *string `json:"category_id"` // optional category mapping
}

type UpdateStatusDefinitionRequest struct {
	DisplayName *string `json:"display_name"`
	Order       *int    `json:"order"`
}

// Helper functions for status definitions (logic layer – used by handlers/services)

// GetOrCreateStatusDefinition ensures a status definition exists for a given display name.
func GetOrCreateStatusDefinition(db *gorm.DB, displayName string) (*EventStatusDefinition, error) {
	var existing EventStatusDefinition
	err := db.Where("LOWER(display_name) = ?", strings.ToLower(displayName)).First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Determine order (append at end)
	var maxOrder int
	db.Model(&EventStatusDefinition{}).Select("COALESCE(MAX(`order`),0)").Scan(&maxOrder)

	// Generate unique slug from display name
	slug := utils.GenerateUniqueSlug(db, displayName, "event_status_definitions")

	def := EventStatusDefinition{
		DisplayName: displayName,
		Slug:        slug,
		Order:       maxOrder + 1,
		IsReserved:  false,
	}

	if err := db.Create(&def).Error; err != nil {
		return nil, err
	}
	return &def, nil
}

// SeedDefaultStatuses creates default status definitions if none exist
func SeedDefaultStatuses(db *gorm.DB) error {
	// Check if any statuses exist
	var count int64
	if err := db.Model(&EventStatusDefinition{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing statuses: %w", err)
	}

	// If statuses already exist, don't create defaults
	if count > 0 {
		fmt.Printf("Statuses already exist (%d), skipping default status creation\n", count)
		return nil
	}

	fmt.Printf("No statuses found, creating default statuses\n")

	// Default statuses in order
	defaultStatuses := []struct {
		Name  string
		Order int
	}{
		{"Backlog", 0},
		{"Proposed", 1},
		{"Feedback", 2},
		{"In Progress", 3},
		{"Released", 4},
		{"Archived", 5},
	}

	createdStatuses := []EventStatusDefinition{}

	for _, ds := range defaultStatuses {
		// Check if status already exists
		var existing EventStatusDefinition
		err := db.Where("LOWER(display_name) = ?", strings.ToLower(ds.Name)).First(&existing).Error
		if err == nil {
			fmt.Printf("Status %s already exists, skipping\n", ds.Name)
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check existing status: %w", err)
		}

		// Generate unique slug
		slug := utils.GenerateUniqueSlug(db, ds.Name, "event_status_definitions")

		// Create the status
		status := EventStatusDefinition{
			DisplayName: ds.Name,
			Slug:        slug,
			Order:       ds.Order,
			IsReserved:  false,
		}

		if err := db.Create(&status).Error; err != nil {
			return fmt.Errorf("failed to create status %s: %w", ds.Name, err)
		}

		fmt.Printf("Created default status: %s (order: %d)\n", ds.Name, ds.Order)
		createdStatuses = append(createdStatuses, status)
	}

	// If we created statuses and a theme exists, create mappings
	if len(createdStatuses) > 0 {
		// Check if a theme is applied
		var settings ProjectSettings
		if err := db.First(&settings).Error; err == nil && settings.CurrentThemeID != "" {
			// Load theme manifest
			manifest, err := LoadThemeManifest("./data/themes/current")
			if err == nil && manifest != nil {
				// Create mappings for the newly created statuses
				for _, status := range createdStatuses {
					suggestedCategory := SuggestCategoryForStatus(status.DisplayName, manifest.Categories)

					// Check if mapping already exists
					var existingMapping StatusCategoryMapping
					err := db.Where("status_definition_id = ? AND theme_id = ?", status.ID, settings.CurrentThemeID).First(&existingMapping).Error
					if err == nil {
						continue // Mapping already exists
					}

					// Create mapping
					mapping := StatusCategoryMapping{
						StatusDefinitionID: status.ID,
						ThemeID:            settings.CurrentThemeID,
						CategoryID:         suggestedCategory,
					}

					if err := db.Create(&mapping).Error; err != nil {
						fmt.Printf("Warning: Failed to create mapping for status %s: %v\n", status.DisplayName, err)
					} else {
						fmt.Printf("Created mapping: %s -> %s\n", status.DisplayName, suggestedCategory)
					}
				}
			}
		}
	}

	return nil
}

// SeedStatusDefinitions initializes any legacy statuses found in existing events
func SeedStatusDefinitions(db *gorm.DB) error {
	// First, seed default statuses if database is empty
	if err := SeedDefaultStatuses(db); err != nil {
		fmt.Printf("Warning: Failed to seed default statuses: %v\n", err)
	}

	// Then, detect distinct existing event statuses and seed definitions for them
	var rawStatuses []string
	if err := db.Model(&Event{}).Distinct().Pluck("status", &rawStatuses).Error; err == nil {
		for _, rs := range rawStatuses {
			if rs == "" {
				continue
			}
			_, _ = GetOrCreateStatusDefinition(db, rs) // ignore errors to continue seeding
		}
	}

	return nil
}
