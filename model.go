package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

// InputField represents ehich input is currently focused
type InputField int

const (
	DomainInput InputField = iota
	ServerInput
	RecordTypeInput
)

// Model represents the state of our application
type Model struct {
	domain        textinput.Model
	server        textinput.Model
	recordType    string
	result        *DNSResult
	loading       bool
	err           error
	width         int
	height        int
	focus         InputField
	selected      int
	recordTypes   []RecordType
	history       []DNSResult // Store history of queries
	historyIdx    int         // Current position in history
	showingHelp   bool        // Whether to show the help screen
	lastQueryTime time.Time   // When the ;ast query was performed
}

// NewModel creates a new mofel with default values
func NewModel() Model {
	domainInput := textinput.New()
	domainInput.Placeholder = "Enter a domain (e.g. example.com)"
	domainInput.Focus()

	// Get system DNS servers to display as examples
	systemServers := GetSystemDNSServers()
	defaultServer := ""
	if len(systemServers) > 0 {
		defaultServer = systemServers[0]
	}

	serverInput := textinput.New()
	serverInput.Placeholder = "DNS server (e.g., 8.8.8.8:53, leave empty for system default)"
	serverInput.SetValue(defaultServer)

	return Model{
		domain:      domainInput,
		server:      serverInput,
		recordType:  string(A),
		result:      nil,
		loading:     false,
		err:         nil,
		focus:       DomainInput,
		selected:    0,
		recordTypes: []RecordType{A, AAAA, CNAME, MX, TXT, NS},
		history:     []DNSResult{},
		historyIdx:  -1,
		showingHelp: false,
	}
}

// Initialize the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink // Start the cursor blinking
}

// LookupDNSCmd returns a command to perform a DNS lookup
func lookupDNSCmd(domain string, recordType RecordType, server string) tea.Cmd {
	return func() tea.Msg {
		result := LookupDNS(domain, recordType, server)
		return result
	}
}
