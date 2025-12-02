package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// HistoryEntry represents a terraform command execution history
type HistoryEntry struct {
	ID          int64
	Directory   string
	Action      string // "apply", "destroy"
	Timestamp   time.Time
	ConfigFile  string // tfvars file used
	ConfigData  string // content of tfvars at that time
	Success     bool
	ErrorMsg    string
}

// HistoryDB manages terraform execution history
type HistoryDB struct {
	db *sql.DB
}

// NewHistoryDB creates a new history database
func NewHistoryDB() (*HistoryDB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}

	dbDir := filepath.Join(home, ".t9s")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db dir: %w", err)
	}

	dbPath := filepath.Join(dbDir, "history.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	hdb := &HistoryDB{db: db}
	if err := hdb.init(); err != nil {
		db.Close()
		return nil, err
	}

	return hdb, nil
}

// init creates the history table
func (h *HistoryDB) init() error {
	query := `
	CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		directory TEXT NOT NULL,
		action TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		config_file TEXT,
		config_data TEXT,
		success INTEGER NOT NULL,
		error_msg TEXT,
		INDEX idx_directory (directory),
		INDEX idx_timestamp (timestamp)
	);
	`
	_, err := h.db.Exec(query)
	return err
}

// AddEntry adds a new history entry
func (h *HistoryDB) AddEntry(entry *HistoryEntry) error {
	query := `
	INSERT INTO history (directory, action, timestamp, config_file, config_data, success, error_msg)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := h.db.Exec(query,
		entry.Directory,
		entry.Action,
		entry.Timestamp,
		entry.ConfigFile,
		entry.ConfigData,
		entry.Success,
		entry.ErrorMsg,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	entry.ID = id
	return nil
}

// GetByDirectory retrieves history entries for a specific directory
func (h *HistoryDB) GetByDirectory(directory string, limit int) ([]*HistoryEntry, error) {
	query := `
	SELECT id, directory, action, timestamp, config_file, config_data, success, error_msg
	FROM history
	WHERE directory = ?
	ORDER BY timestamp DESC
	LIMIT ?
	`
	rows, err := h.db.Query(query, directory, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*HistoryEntry
	for rows.Next() {
		entry := &HistoryEntry{}
		var timestamp string
		err := rows.Scan(
			&entry.ID,
			&entry.Directory,
			&entry.Action,
			&timestamp,
			&entry.ConfigFile,
			&entry.ConfigData,
			&entry.Success,
			&entry.ErrorMsg,
		)
		if err != nil {
			return nil, err
		}
		entry.Timestamp, _ = time.Parse("2006-01-02 15:04:05", timestamp)
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetRecent retrieves recent history entries across all directories
func (h *HistoryDB) GetRecent(limit int) ([]*HistoryEntry, error) {
	query := `
	SELECT id, directory, action, timestamp, config_file, config_data, success, error_msg
	FROM history
	ORDER BY timestamp DESC
	LIMIT ?
	`
	rows, err := h.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*HistoryEntry
	for rows.Next() {
		entry := &HistoryEntry{}
		var timestamp string
		err := rows.Scan(
			&entry.ID,
			&entry.Directory,
			&entry.Action,
			&timestamp,
			&entry.ConfigFile,
			&entry.ConfigData,
			&entry.Success,
			&entry.ErrorMsg,
		)
		if err != nil {
			return nil, err
		}
		entry.Timestamp, _ = time.Parse("2006-01-02 15:04:05", timestamp)
		entries = append(entries, entry)
	}

	return entries, nil
}

// Close closes the database connection
func (h *HistoryDB) Close() error {
	return h.db.Close()
}

