package database

import (
	"fmt"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"gorm.io/gorm"
)

// AutoMigrate runs auto migrations for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Message{},
		&models.Contact{},
		&models.Template{},
		&models.APIKey{},
		&models.Call{},
		&models.Transcript{},
		&models.TranscriptSegment{},
	)
}

// DropAllTables drops all tables (use with caution!)
func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.Message{},
		&models.Contact{},
		&models.Template{},
		&models.APIKey{},
		&models.Call{},
		&models.Transcript{},
		&models.TranscriptSegment{},
	)
}

// CreateIndexes creates custom indexes
func CreateIndexes(db *gorm.DB) error {
	// Messages indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_messages_phone_timestamp
		ON messages(from_number, timestamp DESC);
	`).Error; err != nil {
		return fmt.Errorf("failed to create messages index: %w", err)
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_messages_status_created
		ON messages(status, created_at DESC);
	`).Error; err != nil {
		return fmt.Errorf("failed to create messages status index: %w", err)
	}

	// Contacts indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_contacts_last_message
		ON contacts(last_message_at DESC NULLS LAST);
	`).Error; err != nil {
		return fmt.Errorf("failed to create contacts index: %w", err)
	}

	// Full-text search index for message content
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_messages_content_search
		ON messages USING gin(to_tsvector('english', content));
	`).Error; err != nil {
		return fmt.Errorf("failed to create full-text search index: %w", err)
	}

	return nil
}

// CreateTriggers creates database triggers
func CreateTriggers(db *gorm.DB) error {
	// Trigger to automatically update updated_at timestamp
	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';
	`).Error; err != nil {
		return fmt.Errorf("failed to create trigger function: %w", err)
	}

	// Apply trigger to all tables
	tables := []string{"messages", "contacts", "templates", "api_keys", "calls", "transcripts", "transcript_segments"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf(`
			DROP TRIGGER IF EXISTS update_%s_updated_at ON %s;
			CREATE TRIGGER update_%s_updated_at
			BEFORE UPDATE ON %s
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
		`, table, table, table, table)).Error; err != nil {
			return fmt.Errorf("failed to create trigger for %s: %w", table, err)
		}
	}

	return nil
}
