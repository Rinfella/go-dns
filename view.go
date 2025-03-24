package main

import (
	"fmt"
	"strings"
)

// View renders the user interface
func (m Model) View() string {
	// Help page
	if m.showingHelp {
		return renderHelpScreen(m.width)
	}

	var s strings.Builder

	// Calculate the main container width based on terminal width
	containerWidth := m.width - 4
	if containerWidth < 20 {
		containerWidth = 80
	}

	// Title
	title := titleStyle.Render("DNS Lookup Tool")
	s.WriteString(title + "\n")

	// Domain input
	domainLabel := inputLabelStyle.Render("Domain: ")
	if m.focus == DomainInput {
		domainLabel = focusedInputLabelStyle.Render("Domain: ")
	}
	s.WriteString(domainLabel + " " + m.domain.View() + "\n")

	// Server input
	serverLabel := inputLabelStyle.Render("Server: ")
	if m.focus == ServerInput {
		serverLabel = focusedInputLabelStyle.Render("Server: ")
	}
	s.WriteString(serverLabel + " " + m.server.View() + "\n")

	if m.server.Value() == "" {
		systemServerInfo := helpStyle.Render(" (using system DNS servers)")
		s.WriteString(systemServerInfo + "\n")
	}

	// Record type selection
	recordTypeLabel := inputLabelStyle.Render("Record Type: ")
	if m.focus == RecordTypeInput {
		recordTypeLabel = focusedInputLabelStyle.Render("Record Type: ")
	}
	s.WriteString(recordTypeLabel + " ")

	for i, rt := range m.recordTypes {
		if i == m.selected {
			if m.focus == RecordTypeInput {
				s.WriteString(focusedSelectedRecordTypeStyle.Render(string(rt)))
			} else {
				s.WriteString(selectedRecordTypeStyle.Render(string(rt)))
			}
		} else {
			s.WriteString(recordTypeStyle.Render(string(rt)))
		}
	}
	s.WriteString("\n\n")

	// Help text
	helpText := helpStyle.Render("Press TAB to switch focus, ENTER or CTRL+SPACE to perform lookup, UP/DOWN for history, Q to quit")
	s.WriteString(helpText + "\n")
	helpText2 := helpStyle.Render("Press H/L to cycle record types, ? for help, Q to quit")
	s.WriteString(helpText2 + "\n\n")

	// Results
	if m.loading {
		s.WriteString(loadingStyle.Render("Querying DNS servers...") + "\n")
	} else if m.result != nil {
		resultHeader := resultHeaderStyle.Render(
			fmt.Sprintf("Results for %s (%s record)",
				m.result.Domain,
				m.result.Type))
		s.WriteString(resultHeader + "\n")

		queryTime := fmt.Sprintf("Query time: %v", m.result.QueryTime)
		s.WriteString(queryTime + "\n\n")

		if m.result.Error != nil {
			s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v ", m.result.Error)) + "\n")
		} else if len(m.result.Records) == 0 {
			s.WriteString(errorStyle.Render("No records found!"))
		} else {
			for i, record := range m.result.Records {
				recordItem := recordItemStyle.Render(fmt.Sprintf("%d. %s", i+1, record))
				s.WriteString(recordItem + "\n")
			}
		}
	}

	// Show history count
	if len(m.history) > 0 {
		historyInfo := historyStyle.Render(
			fmt.Sprintf("%d queries in history. Use UP/DOWN to browse.",
				len(m.history)))
		s.WriteString("\n" + historyInfo + "\n")
	}

	// Wrap everything in the app container
	return appStyle.Width(containerWidth).Render(s.String())
}

// Render separate help screen
func renderHelpScreen(width int) string {
	containerWidth := width - 4
	if containerWidth < 20 {
		containerWidth = 80
	}

	var s strings.Builder

	title := titleStyle.Render("Help - DNS Lookup Tool")
	s.WriteString(title + "\n\n")

	s.WriteString(focusedInputLabelStyle.Render("Navigation:") + "\n")
	s.WriteString("- TAB: Cycle between input fields\n")
	s.WriteString("- ENTER: Perform DNS lookup\n")
	s.WriteString("- LEFT/RIGHT or H/L: Change record type\n")
	s.WriteString("- UP/DOWN: Browse history\n")
	s.WriteString("- ?: Toggle this help screen\n")
	s.WriteString("- Q or Ctrl+c: Quit the application\n\n")

	s.WriteString(focusedInputLabelStyle.Render("Record Types:") + "\n")
	s.WriteString("- A: IPV4 address records\n")
	s.WriteString("- AAAA: IPV6 address records\n")
	s.WriteString("- CNAME: Canonical name records (aliases)\n")
	s.WriteString("- MX: Mail exchange records\n")
	s.WriteString("- TXT: Text records\n")
	s.WriteString("- NS: Name server records\n")

	s.WriteString(focusedInputLabelStyle.Render("Server:") + "\n")
	s.WriteString("- Enter a DNS server address with port (e.g., 8.8.8.8:53)\n")
	s.WriteString("- Leave empty to use your system's default DNS servers\n")

	s.WriteString(focusedInputLabelStyle.Render("History:") + "\n")
	s.WriteString("- Use UP and DOWN arrow keys to navigate through your query histories\n")
	s.WriteString("- History is maintained for the current session only.\n\n")

	s.WriteString(helpStyle.Render("Press ? to return to the main screen"))

	return appStyle.Width(containerWidth).Render(s.String())
}
