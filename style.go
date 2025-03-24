package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#2D3748")
	accentColor    = lipgloss.Color("#F56565")
	successColor   = lipgloss.Color("#48BB78")
	warningColor   = lipgloss.Color("#ECC948")
	errorColor     = lipgloss.Color("#F56565")
	textColor      = lipgloss.Color("#E2E8F0")

	// General styles
	appStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Input styles
	inputLabelStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	focusedInputLabelStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	// Record type styles
	recordTypeStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(textColor)

	selectedRecordTypeStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Foreground(primaryColor).
				Bold(true)

	focusedSelectedRecordTypeStyle = lipgloss.NewStyle().
					Padding(0, 1).
					Foreground(primaryColor).
					Background(lipgloss.Color("#1A202C")).
					Bold(true)

	// Result Styles
	resultHeaderStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true).
				MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	loadingStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#718096")).
			Italic(true).
			MarginTop(1)

	recordItemStyle = lipgloss.NewStyle().
			Foreground(textColor)

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0AEC0")).
			Italic(true).
			MarginTop(1)
)
