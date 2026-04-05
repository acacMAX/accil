package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/accil/accil/internal/ai"
)

// Session represents a conversation session
type Session struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Messages  []ai.Message `json:"messages"`
}

// Manager manages sessions
type Manager struct {
	sessionsDir string
}

// NewManager creates a new session manager
func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sessionsDir := filepath.Join(home, ".ai-cli", "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return nil, err
	}

	return &Manager{sessionsDir: sessionsDir}, nil
}

// NewSession creates a new session
func (m *Manager) NewSession(name string) *Session {
	now := time.Now()
	return &Session{
		ID:        fmt.Sprintf("%d", now.UnixNano()),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		Messages:  []ai.Message{},
	}
}

// Save saves a session
func (m *Manager) Save(session *Session) error {
	session.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(m.sessionsDir, session.ID+".json")
	return os.WriteFile(path, data, 0644)
}

// Load loads a session by ID
func (m *Manager) Load(id string) (*Session, error) {
	path := filepath.Join(m.sessionsDir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Delete deletes a session
func (m *Manager) Delete(id string) error {
	path := filepath.Join(m.sessionsDir, id+".json")
	return os.Remove(path)
}

// List returns all sessions sorted by updated time
func (m *Manager) List() ([]*Session, error) {
	entries, err := os.ReadDir(m.sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Session{}, nil
		}
		return nil, err
	}

	var sessions []*Session
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5]
		session, err := m.Load(id)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}

	// Sort by updated time (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})

	return sessions, nil
}

// GetLastSession returns the most recently updated session
func (m *Manager) GetLastSession() (*Session, error) {
	sessions, err := m.List()
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, nil
	}

	return sessions[0], nil
}

// AddMessage adds a message to a session
func (s *Session) AddMessage(role, content string) {
	s.Messages = append(s.Messages, ai.Message{
		Role:    role,
		Content: content,
	})
}

// AddToolResult adds a tool result message
func (s *Session) AddToolResult(toolCallID, name, result string) {
	s.Messages = append(s.Messages, ai.Message{
		Role:       "tool",
		Content:    result,
		ToolCallID: toolCallID,
		Name:       name,
	})
}
