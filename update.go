package main

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

// Update handles messages and user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// fmt.Printf("Key Pressed: %+v\n", msg)
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "?":
			m.showingHelp = !m.showingHelp
			return m, nil

		case "enter":
			if m.focus == DomainInput && m.domain.Value() != "" {
				// Perform DNS lookup when Enter is pressed in domain input
				if time.Since(m.lastQueryTime) < 500*time.Millisecond {
					return m, nil
				}

				// Perform DNS lookup when Enter is pressed in Domain Input
				m.loading = true
				m.lastQueryTime = time.Now()
				return m, lookupDNSCmd(m.domain.Value(), RecordType(m.recordType), m.server.Value())
			}

		case "ctrl+@":
			// Perform DNS lookup when Enter is pressed in Domain Input
			m.loading = true
			m.lastQueryTime = time.Now()
			return m, lookupDNSCmd(m.domain.Value(), RecordType(m.recordType), m.server.Value())

		case "tab":
			if !m.showingHelp {
				// Cycle focus between input fields
				switch m.focus {
				case DomainInput:
					m.focus = ServerInput
					m.domain.Blur()
					m.server.Focus()
				case ServerInput:
					m.focus = RecordTypeInput
					m.server.Blur()
				case RecordTypeInput:
					m.focus = DomainInput
					m.domain.Focus()
				}
			}
			return m, nil

		case "h", "left":
			if m.focus == RecordTypeInput {
				m.selected = (m.selected - 1 + len(m.recordTypes)) % len(m.recordTypes)
				m.recordType = string(m.recordTypes[m.selected])
			}

		case "l", "right":
			if m.focus == RecordTypeInput {
				m.selected = (m.selected + 1) % len(m.recordTypes)
				m.recordType = string(m.recordTypes[m.selected])
			}

		case "up", "down":
			if m.showingHelp {
				return m, nil
			}

			if m.focus == RecordTypeInput {
				if msg.String() == "up" && len(m.history) > 0 {
					// Cycle through history
					if m.historyIdx == -1 {
						m.historyIdx = len(m.history) - 1
					} else if m.historyIdx > 0 {
						m.historyIdx--
					}

					if m.historyIdx >= 0 && m.historyIdx < len(m.history) {
						historyItem := m.history[m.historyIdx]
						m.domain.SetValue(historyItem.Domain)
						m.recordType = string(historyItem.Type)

						// Find the index of the record type
						for i, rt := range m.recordTypes {
							if rt == historyItem.Type {
								m.selected = i
								break
							}
						}
					}
				} else if msg.String() == "down" && m.historyIdx > -1 {
					// Go forward in history
					if m.historyIdx < len(m.history)-1 {
						m.historyIdx++

						historyItem := m.history[m.historyIdx]
						m.domain.SetValue(historyItem.Domain)
						m.recordType = string(historyItem.Type)

						// Find the index of the record type
						for i, rt := range m.recordTypes {
							if rt == historyItem.Type {
								m.selected = i
								break
							}
						}
					} else {
						m.historyIdx = -1
						m.domain.SetValue("")
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case DNSResult:
		m.result = &msg
		m.loading = false

		// Add to history if successful and not already in history
		if msg.Error == nil {
			// Check if already in history to avoid duplicates
			isDuplicate := false
			for _, hist := range m.history {
				if hist.Domain == msg.Domain && hist.Type == msg.Type {
					isDuplicate = true
					break
				}
			}

			if !isDuplicate {
				m.history = append(m.history, msg)
			}
		}

		return m, nil
	}

	// Handle updates to the text inputs
	if !m.showingHelp {
		switch m.focus {
		case DomainInput:
			m.domain, cmd = m.domain.Update(msg)
		case ServerInput:
			m.server, cmd = m.server.Update(msg)
		}
	}

	return m, cmd
}
