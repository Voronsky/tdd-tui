package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/tdd-tui/internal/uex"
)

// RenderBloombergGrid computes the side-by-side cell layout matrices.
func RenderBloombergGrid(items []uex.Listing, focusedIndex int, terminalWidth int) string {
	const fixedBoxWidth = 22

	if len(items) == 0 {
		return "Awaiting Market Data from UEX...."
	}

	// 1. Calculate how many boxes can physically fit on screen safely
	maxVisible := terminalWidth / fixedBoxWidth
	if maxVisible < 1 {
		maxVisible = 1 // Ensure at least one box renders even on tiny screens
	}

	// 2. Compute the start and end boundaries of the sliding window
	startIdx := 0
	if focusedIndex >= maxVisible {
		// Shift the window to keep the cursor at the right-most edge
		startIdx = focusedIndex - maxVisible + 1
	}

	endIdx := startIdx + maxVisible
	if endIdx > len(items) {
		endIdx = len(items)
	}

	// 3. Slice the payload to extract only the visible nodes
	visibleSlice := items[startIdx:endIdx]

	var computedCells []string

	for i, item := range visibleSlice {
		actualIndex := startIdx + i
		// Define strict dimension boundaries for each grid element cell
		cellStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2).
			Width(fixedBoxWidth - 2).
			Height(5)

		// Apply highlight matrix formatting if the item matches the active index
		if actualIndex == focusedIndex {
			cellStyle = cellStyle.BorderForeground(lipgloss.Color("86")) // Accent cyan
		}

		// Dynamically assign terminal color weights based on trend boolean evaluations
		// Let's calculate Upward or Downward terend first
		isTrendingUp := item.PriceBuy > item.PriceBuyAvg

		// Default to downward trend
		indicatorColor := "160" // System Red
		directionSign := "▼"

		// Check if line go up
		if isTrendingUp {
			indicatorColor = "46" // System Green
			directionSign = "▲"   // Price is going up
		}

		trendStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(indicatorColor))

		// Compose individual cell string metrics vertically
		cellContent := fmt.Sprintf(
			"%s\n\nB: %.2f\nS: %.2f\n%s",
			lipgloss.NewStyle().Bold(true).Render(item.CommodityName),
			item.PriceBuy,
			item.PriceSell,
			trendStyle.Render(directionSign),
		)

		// Render strings inside boundaries and append to layout sequence slice
		computedCells = append(computedCells, cellStyle.Render(cellContent))
	}

	// Stitch distinct cell columns into a single row matrix block side-by-side
	return lipgloss.JoinHorizontal(lipgloss.Top, computedCells...)
}
