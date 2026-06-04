package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/tdd-tui/internal/uex"
)

// RenderBloombergGrid computes the side-by-side cell layout matrices.
func RenderBloombergGrid(items []uex.Commodity, focusedIndex int, windowOffset int) string {
	// Basic layout
	const cols = 6
	const visibleRows = 3

	const fixedBoxWidth = 22

	currentRow := focusedIndex / cols

	// Now calculate the vertical sliding window to keep the cursor in view
	// Ensuring that , as the trading window scrolls down when moving past the visible floor
	startRow := 0

	if currentRow >= visibleRows {
		startRow = currentRow - visibleRows + 1
	}

	var renderedRows []string

	// 3. Loop through the data to construct only the visible rows
	for r := startRow; r < startRow+visibleRows; r++ {
		var currentColumnBlocks []string

		for c := 0; c < cols; c++ {
			// Translate matrix coordinates back into a flat slice index
			flatIndex := (r * cols) + c

			// Prevent out-of-bounds panics on the final incomplete row
			if flatIndex >= len(items) {
				break
			}

			node := items[flatIndex]

			// Construct the individual cell bounding boxes
			cellStyle := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("68")).
				Width(18).
				Height(4)

			if flatIndex == focusedIndex {
				cellStyle = cellStyle.BorderForeground(lipgloss.Color("86"))
			}

			//TODO: Rethink how to do trends without blowing up Requests
			//isTrendingUp := node.PriceBuy > node.PriceBuyAvg

			//// Default to downward trend
			//indicatorColor := "160" //System Red
			//directionSign := "▼"

			//if isTrendingUp {
			//	indicatorColor = "46" // System Green
			//	directionSign = "▲"   // Price is going up
			//}

			//trendStyle := lipgloss.NewStyle().Foreground((lipgloss.Color(indicatorColor)))

			//Stylize the name of each commodity
			styledName := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFD700")).
				Render(node.Name)

			styledBuy := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("46")).
				Render(fmt.Sprintf("%.2f", node.PriceBuy))

			styledSell := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("46")).
				Render(fmt.Sprintf("%.2f", node.PriceSell))

			content := fmt.Sprintf(
				"Name:%s\n\nB: %s\nS: %s\n", styledName, styledBuy, styledSell)
			currentColumnBlocks = append(currentColumnBlocks, cellStyle.Render(content))
		}

		if len(currentColumnBlocks) > 0 {
			// Stitch the columns together horizontally to form a completed row
			rowString := lipgloss.JoinHorizontal(lipgloss.Top, currentColumnBlocks...)
			renderedRows = append(renderedRows, rowString)
		}
	}

	// 4. Stitch the completed rows together vertically to form the grid matrix
	return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)

}
