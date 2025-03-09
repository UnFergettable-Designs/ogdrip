package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Generation represents a record of an OpenGraph image generation
type Generation struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	TargetURL     string    `json:"target_url"`
	ImagePath     string    `json:"image_path"`
	HTMLPath      string    `json:"html_path"`
	CreatedAt     time.Time `json:"created_at"`
	ClientIP      string    `json:"client_ip"`
	UserAgent     string    `json:"user_agent"`
	Parameters    string    `json:"parameters"` // JSON string of all parameters
	Status        string    `json:"status"`     // pending, completed, failed
	ErrorMessage  string    `json:"error_message,omitempty"`
	DownloadCount int       `json:"download_count"`
}

// Database struct for SQLite operations
type Database struct {
	db *sql.DB
}

var dbInstance *Database
var dbInitError error

// InitDB initializes the SQLite database
func InitDB() (*Database, error) {
	// If we already have an instance, return it
	if dbInstance != nil {
		return dbInstance, nil
	}

	// If we previously encountered an error, return it
	if dbInitError != nil {
		return nil, dbInitError
	}

	// Ensure the data directory exists
	dataDir := filepath.Join(".", "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		dbInitError = fmt.Errorf("failed to create data directory: %w", err)
		return nil, dbInitError
	}

	dbPath := filepath.Join(dataDir, "generations.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		dbInitError = fmt.Errorf("failed to open database: %w", err)
		return nil, dbInitError
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		dbInitError = fmt.Errorf("failed to connect to database: %w", err)
		return nil, dbInitError
	}

	// Create the generations table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS generations (
		id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		target_url TEXT,
		image_path TEXT,
		html_path TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		client_ip TEXT,
		user_agent TEXT,
		parameters TEXT,
		downloaded BOOLEAN DEFAULT 0,
		cleanup_after TIMESTAMP,
		status TEXT DEFAULT 'pending',
		error_message TEXT,
		download_count INTEGER DEFAULT 0
	);
	`
	if _, err = db.Exec(createTableSQL); err != nil {
		db.Close()
		dbInitError = fmt.Errorf("failed to create table: %w", err)
		return nil, dbInitError
	}

	// Create indexes for faster queries
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_created_at ON generations(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_cleanup_after ON generations(cleanup_after);`,
		`CREATE INDEX IF NOT EXISTS idx_status ON generations(status);`,
	}

	for _, query := range indexQueries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}

	dbInstance = &Database{db: db}
	return dbInstance, nil
}

// CloseDB closes the database connection
func (db *Database) CloseDB() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

// ensureConnection checks if the database connection is valid, and reopens it if closed
func (db *Database) ensureConnection() error {
	if db.db == nil {
		log.Println("Database connection is nil, attempting to reconnect")
		newDB, err := InitDB()
		if err != nil {
			return fmt.Errorf("failed to reconnect to database: %w", err)
		}
		db.db = newDB.db
		return nil
	}

	// Ping the database to check if the connection is still alive
	if err := db.db.Ping(); err != nil {
		log.Println("Database ping failed, attempting to reconnect:", err)
		// Close the existing connection if it's not nil
		if db.db != nil {
			// Just attempt to close, ignore errors as we're reconnecting anyway
			_ = db.db.Close()
		}

		// Clear the instance to ensure a fresh connection
		dbInstance = nil

		// Reopen the connection
		newDB, err := InitDB()
		if err != nil {
			return fmt.Errorf("failed to reconnect to database: %w", err)
		}
		db.db = newDB.db
		// Update the global instance
		dbInstance = db
	}

	return nil
}

// SaveGeneration stores a generation record in the database
func (db *Database) SaveGeneration(gen *Generation) error {
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	// Set cleanup time to 24 hours from now
	cleanupAfter := time.Now().Add(24 * time.Hour)

	// Set default status if not provided
	if gen.Status == "" {
		gen.Status = "pending"
	}

	query := `
	INSERT INTO generations (
		id, title, description, target_url, image_path, html_path,
		created_at, client_ip, user_agent, parameters, cleanup_after,
		status, error_message, download_count
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.db.Exec(
		query,
		gen.ID,
		gen.Title,
		gen.Description,
		gen.TargetURL,
		gen.ImagePath,
		gen.HTMLPath,
		gen.CreatedAt,
		gen.ClientIP,
		gen.UserAgent,
		gen.Parameters,
		cleanupAfter,
		gen.Status,
		gen.ErrorMessage,
		gen.DownloadCount,
	)

	return err
}

// GetGeneration retrieves a generation by ID
func (db *Database) GetGeneration(id string) (*Generation, error) {
	if err := db.ensureConnection(); err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	query := `SELECT id, title, description, target_url, image_path, html_path,
		created_at, client_ip, user_agent, parameters, status, error_message, download_count
		FROM generations WHERE id = ?`

	row := db.db.QueryRow(query, id)

	gen := &Generation{}
	var createdAtStr string

	err := row.Scan(
		&gen.ID,
		&gen.Title,
		&gen.Description,
		&gen.TargetURL,
		&gen.ImagePath,
		&gen.HTMLPath,
		&createdAtStr,
		&gen.ClientIP,
		&gen.UserAgent,
		&gen.Parameters,
		&gen.Status,
		&gen.ErrorMessage,
		&gen.DownloadCount,
	)

	if err != nil {
		return nil, err
	}

	gen.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format: %w", err)
	}

	return gen, nil
}

// ListGenerations returns a list of recent generations
func (db *Database) ListGenerations(limit int) ([]*Generation, error) {
	if err := db.ensureConnection(); err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `SELECT id, title, description, target_url, image_path, created_at,
		status, download_count
		FROM generations ORDER BY created_at DESC LIMIT ?`

	rows, err := db.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var generations []*Generation

	for rows.Next() {
		gen := &Generation{}
		var createdAtStr string

		err := rows.Scan(
			&gen.ID,
			&gen.Title,
			&gen.Description,
			&gen.TargetURL,
			&gen.ImagePath,
			&createdAtStr,
			&gen.Status,
			&gen.DownloadCount,
		)

		if err != nil {
			return nil, err
		}

		gen.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp format: %w", err)
		}

		generations = append(generations, gen)
	}

	return generations, nil
}

// MarkAsDownloaded marks a generation as downloaded and increments download count
func (db *Database) MarkAsDownloaded(id string) error {
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	query := `UPDATE generations SET downloaded = 1, download_count = download_count + 1 WHERE id = ?`
	_, err := db.db.Exec(query, id)
	return err
}

// UpdateStatus updates the status of a generation
func (db *Database) UpdateStatus(id string, status string) error {
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	query := `UPDATE generations SET status = ? WHERE id = ?`
	_, err := db.db.Exec(query, status, id)
	return err
}

// SetErrorMessage sets an error message for a failed generation
func (db *Database) SetErrorMessage(id string, errorMsg string) error {
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	query := `UPDATE generations SET status = 'failed', error_message = ? WHERE id = ?`
	_, err := db.db.Exec(query, errorMsg, id)
	return err
}

// MarkAsCompleted marks a generation as completed
func (db *Database) MarkAsCompleted(id string) error {
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	query := `UPDATE generations SET status = 'completed' WHERE id = ?`
	_, err := db.db.Exec(query, id)
	return err
}

// SetCleanupTime sets when a generation should be cleaned up
func (db *Database) SetCleanupTime(id string, cleanupAfter time.Time) error {
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	query := `UPDATE generations SET cleanup_after = ? WHERE id = ?`
	_, err := db.db.Exec(query, cleanupAfter, id)
	return err
}

// RunCleanup performs cleanup of old generations
func (db *Database) RunCleanup() error {
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	// Get records that are due for cleanup
	query := `SELECT id, image_path, html_path FROM generations WHERE cleanup_after < ?`
	rows, err := db.db.Query(query, time.Now())
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	filesToDelete := make(map[string]struct{})

	for rows.Next() {
		var id, imagePath, htmlPath string

		err := rows.Scan(&id, &imagePath, &htmlPath)
		if err != nil {
			log.Printf("Error scanning cleanup record: %v", err)
			continue
		}

		ids = append(ids, id)

		// Add paths to delete list if they exist
		if imagePath != "" {
			filesToDelete[imagePath] = struct{}{}
		}
		if htmlPath != "" {
			filesToDelete[htmlPath] = struct{}{}
		}
	}

	// Delete files
	for path := range filesToDelete {
		if _, err := os.Stat(path); err == nil {
			err := os.Remove(path)
			if err != nil {
				log.Printf("Error removing file %s: %v", path, err)
			} else {
				log.Printf("Removed file: %s", path)
			}
		}
	}

	// Remove records from database
	if len(ids) > 0 {
		placeholders := make([]string, len(ids))
		args := make([]interface{}, len(ids))

		for i, id := range ids {
			placeholders[i] = "?"
			args[i] = id
		}

		deleteQuery := fmt.Sprintf(
			"DELETE FROM generations WHERE id IN (%s)",
			placeholders,
		)

		_, err = db.db.Exec(deleteQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to delete records: %w", err)
		}

		log.Printf("Cleaned up %d records", len(ids))
	}

	return nil
}

// GenerationParameters captures all parameters for a generation
type GenerationParameters struct {
	WebpageURL  string `json:"webpage_url,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	OgType      string `json:"og_type,omitempty"`
	SiteName    string `json:"site_name,omitempty"`
	TargetURL   string `json:"target_url,omitempty"`
	TwitterCard string `json:"twitter_card,omitempty"`
	ImageWidth  int    `json:"image_width,omitempty"`
	ImageHeight int    `json:"image_height,omitempty"`
	Quality     int    `json:"quality,omitempty"`
	WaitTime    int    `json:"wait_time,omitempty"`
	Debug       bool   `json:"debug,omitempty"`
	Verbose     bool   `json:"verbose,omitempty"`
}

// SerializeParameters converts parameters to a JSON string
func SerializeParameters(params *GenerationParameters) (string, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetRecentGenerations retrieves the most recent generations with pagination
func (db *Database) GetRecentGenerations(limit, offset int) ([]Generation, error) {
	if err := db.ensureConnection(); err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Default to 50 if limit is invalid
	if limit <= 0 {
		limit = 50
	}

	// Default to 0 if offset is invalid
	if offset < 0 {
		offset = 0
	}

	// Get the generations
	query := `
		SELECT id, title, description, target_url, image_path, html_path,
		       created_at, client_ip, user_agent, parameters, status,
		       error_message, download_count
		FROM generations
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query generations: %w", err)
	}
	defer rows.Close()

	var generations []Generation
	for rows.Next() {
		var gen Generation
		var createdAtStr string

		err := rows.Scan(
			&gen.ID,
			&gen.Title,
			&gen.Description,
			&gen.TargetURL,
			&gen.ImagePath,
			&gen.HTMLPath,
			&createdAtStr,
			&gen.ClientIP,
			&gen.UserAgent,
			&gen.Parameters,
			&gen.Status,
			&gen.ErrorMessage,
			&gen.DownloadCount,
		)

		if err != nil {
			log.Printf("Error scanning generation row: %v", err)
			continue
		}

		// Parse the created_at timestamp
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at timestamp: %v", err)
			// Use current time as fallback
			createdAt = time.Now()
		}
		gen.CreatedAt = createdAt

		// Add to results
		generations = append(generations, gen)
	}

	if err := rows.Err(); err != nil {
		return generations, fmt.Errorf("error during row iteration: %w", err)
	}

	return generations, nil
}

// GetGenerationCount returns the total number of generations
func (db *Database) GetGenerationCount() (int, error) {
	if err := db.ensureConnection(); err != nil {
		return 0, fmt.Errorf("database connection error: %w", err)
	}

	query := `SELECT COUNT(*) FROM generations`

	var count int
	err := db.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count generations: %w", err)
	}

	return count, nil
}

// GetGenerationByID retrieves a specific generation by its ID
func (db *Database) GetGenerationByID(id string) (*Generation, error) {
	if err := db.ensureConnection(); err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("empty generation ID")
	}

	// Get the generation from database
	query := `
		SELECT id, title, description, target_url, image_path, html_path,
		       created_at, client_ip, user_agent, parameters, status,
		       error_message, download_count
		FROM generations
		WHERE id = ?
	`

	var gen Generation
	var createdAtStr string

	err := db.db.QueryRow(query, id).Scan(
		&gen.ID,
		&gen.Title,
		&gen.Description,
		&gen.TargetURL,
		&gen.ImagePath,
		&gen.HTMLPath,
		&createdAtStr,
		&gen.ClientIP,
		&gen.UserAgent,
		&gen.Parameters,
		&gen.Status,
		&gen.ErrorMessage,
		&gen.DownloadCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// No generation found with this ID
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query generation by ID: %w", err)
	}

	// Parse the created_at timestamp
	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		log.Printf("Error parsing created_at timestamp: %v", err)
		// Use current time as fallback
		createdAt = time.Now()
	}
	gen.CreatedAt = createdAt

	return &gen, nil
}

// UpdateGenerationStatus updates the status of a generation
func (db *Database) UpdateGenerationStatus(id string, status string, errorMessage string) error {
	// First, ensure we have a valid connection
	if err := db.ensureConnection(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	// Validate ID
	if id == "" {
		return fmt.Errorf("empty generation ID")
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"completed": true,
		"failed":    true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Create query based on whether there's an error message
	var query string
	var args []interface{}

	if errorMessage != "" {
		query = `UPDATE generations SET status = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, errorMessage, id}
	} else {
		query = `UPDATE generations SET status = ? WHERE id = ?`
		args = []interface{}{status, id}
	}

	// Execute the update
	result, err := db.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update generation status: %w", err)
	}

	// Check if the record was found
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no generation found with ID: %s", id)
	}

	return nil
}
