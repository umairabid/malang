package ui

import (
    "github.com/charmbracelet/lipgloss"
)

func focusedStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("57")).
        Bold(true)
}

func normalStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("240"))
}

func errorStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true)
}

func applyStyle(s string, focused bool) string {
    if focused {
        return focusedStyle().Render(s)
    }
    return normalStyle().Render(s)
}
