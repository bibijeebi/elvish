package cli

import (
	"strings"
	"sync"
)

// HistorySuggester provides suggestions based on command history
type HistorySuggester struct {
	mu            sync.RWMutex
	history       []string // Stores command history
	updates       chan struct{}
	maxHistoryLen int
}

// NewHistorySuggester creates a new HistorySuggester
func NewHistorySuggester(maxHistory int) *HistorySuggester {
	return &HistorySuggester{
		history:       make([]string, 0, maxHistory),
		updates:       make(chan struct{}, 1),
		maxHistoryLen: maxHistory,
	}
}

// AddHistory adds a new command to history
func (hs *HistorySuggester) AddHistory(cmd string) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	// Debug print
	println("Adding to history:", cmd)

	// Don't add empty commands or duplicates of the last command
	if cmd == "" || (len(hs.history) > 0 && hs.history[len(hs.history)-1] == cmd) {
		return
	}

	// Add new command
	hs.history = append(hs.history, cmd)

	// Debug print
	println("History length:", len(hs.history))

	// Trim history if it exceeds max length
	if len(hs.history) > hs.maxHistoryLen {
		hs.history = hs.history[1:]
	}

	// Notify of update
	select {
	case hs.updates <- struct{}{}:
	default:
	}
}

// Get implements the Suggester interface
func (hs *HistorySuggester) Get(code string) string {
	if code == "" {
		println("Empty code, no suggestion")
		return ""
	}

	hs.mu.RLock()
	defer hs.mu.RUnlock()

	print("Searching history for:", code)
	print("History length:", len(hs.history))

	// Search history in reverse order (newest first)
	for i := len(hs.history) - 1; i >= 0; i-- {
		print("Checking history entry:", hs.history[i])
		if strings.HasPrefix(hs.history[i], code) {
			suggestion := hs.history[i][len(code):]
			print("Found suggestion:", suggestion)
			return suggestion
		}
	}

	print("No suggestion found")
	return ""
}

// LateUpdates implements the Suggester interface
func (hs *HistorySuggester) LateUpdates() <-chan struct{} {
	return hs.updates
}

// Clear empties the history
func (hs *HistorySuggester) Clear() {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.history = make([]string, 0, hs.maxHistoryLen)
	select {
	case hs.updates <- struct{}{}:
	default:
	}
}
