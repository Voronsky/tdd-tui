package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/common-nighthawk/go-figure"
	"github.com/tdd-tui/internal/uex"
)

type SessionState int

const (
	ChangeAPIStr              = "Change API UEX key"
	ChangeSCUStr              = "Change SCU cargo size"
	ExitStr                   = "Exit"
	GoBackStr                 = "Go Back"
	TradeStr                  = "Trade"
	SettingsStr               = "Change Settings"
	listHeight                = 10
	ViewState    SessionState = iota
	MainState
	SettingsState
	SetApiState
	SetCargoSizeState
	TradeState
)

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	api          string
	cargoSize    uint = 0
	UEXClient    uex.APIClient
)

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

type model struct {
	state     SessionState
	list      list.Model
	choice    string
	styles    styles
	quitting  bool
	textInput textinput.Model
}

type Commodity struct {
	Name string
	Code string
	Kind string
	Buy  float64
	Sell float64
	SCU  float64
}

type item string

type itemDelegate struct {
	styles *styles
}

func (i item) FilterValue() string { return "" }

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func getCommodities() ([]Commodity, error) {
	commodities := []Commodity{}
	//Rest CALL
	return commodities, nil
}

func newStyles(darkBG bool) styles {
	var s styles
	s.title = lipgloss.NewStyle().MarginLeft(2)
	s.item = lipgloss.NewStyle().PaddingLeft(4)
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("86"))
	s.pagination = list.DefaultStyles(darkBG).PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles(darkBG).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return s
}

func initialModel() model {
	items := []list.Item{
		item(TradeStr),
		item(SettingsStr),
		item(ExitStr),
	}
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "List of options"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m := model{state: MainState, list: l}
	m.updateStyles(true)
	return m
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func (m *model) updateStyles(isDark bool) {
	m.styles = newStyles(isDark)
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.Styles.HelpStyle = m.styles.help
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}

func (m model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Key Press management for API View
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			m.choice = string(i)
			log.Println("User made choice: ", m.choice)
			if ok && m.choice == ChangeAPIStr {
				m.state = SetApiState
				m = setApiView(m)
				return m, nil
			}
			if ok && m.choice == ChangeSCUStr {
				m.state = SetCargoSizeState
				m = setCargoSizeView(m)
				return m, nil
			}
			if ok && m.choice == GoBackStr {
				m.state = MainState
				m = mainList(m)
				return m, nil
			}
		case "esc", "q":
			m.state = MainState
			m = mainList(m)
			return m, nil
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m model) updateApi(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Key Press management for API View
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			// Get the User's API key
			if m.textInput.Value() == "" {
				m.state = SettingsState
				m = settingsList(m)
				return m, nil
			}

			api = m.textInput.Value()
			log.Println("User set API value to:", api)

			// Update the UEX client
			UEXClient.Token = api

			// Clear the text field for the next time it gets here
			m.textInput.Reset()

			// Now change the state to go back to the Settings State
			m.state = SettingsState
			m = settingsList(m)
			// Go back to the Settings View
			return m, nil
		case "esc":
			m.textInput.Reset()
			m.state = SettingsState
			m = settingsList(m)
			return m, nil
		case "ctrl+c":
			m.textInput.Reset()
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func (m model) updateCargoSize(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Key Press management for API View
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			// Get the User's SCU
			if m.textInput.Value() == "" {
				m.state = SettingsState
				m = settingsList(m)
				return m, nil
			}
			scu, err := strconv.Atoi(m.textInput.Value())
			if err == nil && scu >= 0 {
				cargoSize = uint(scu)
				log.Println("User set CargoSize to:", api)
			}

			// Clear the text field for the next time it gets here
			m.textInput.Reset()

			// Now change the state to go back to the Settings State
			m.state = SettingsState
			m = settingsList(m)
			// Go back to the Settings View
			return m, nil

		case "esc", "q":
			m.textInput.Reset()
			m.state = SettingsState
			m = settingsList(m)
			return m, nil
		case "ctrl+c":
			m.textInput.Reset()
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func (m model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			log.Println("User chose to quit")
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			m.choice = string(i)
			log.Println("User made choice: ", m.choice)
			if ok && m.choice == SettingsStr {
				m.state = SettingsState
				m = settingsList(m)
				return m, nil

			}
			if ok && m.choice == ExitStr {
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//case tea.WindowSizeMsg:
	//	m.list.SetWidth(msg.Width)
	//	return m, nil
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		}

	}

	// Route messages to the right component based on it's state
	switch m.state {

	case SetApiState:
		return m.updateApi(msg)
	case SettingsState:
		return m.updateSettings(msg)
	case SetCargoSizeState:
		return m.updateCargoSize(msg)
	case MainState:
		return m.updateMainMenu(msg)
	default:
		return m, nil

	}

}

func settingsList(m model) model {
	log.Println("Entered the SettingsList() func")
	items := []list.Item{
		item(ChangeAPIStr),
		item(ChangeSCUStr),
		item(GoBackStr),
	}
	const defaultWidth = 45

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What would you like to change?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	//m = model{state: SettingsState, list: l, choice: "Change Settings"}
	m.state = SettingsState
	m.list = l
	m.choice = SettingsStr
	m.updateStyles(true)
	log.Println("Exiting the SettingsList() func")
	return m

}

// Get the API key from UEX
func setApiView(m model) model {
	log.Println("Entered the apiView() func")
	ti := textinput.New()
	ti.Placeholder = "Enter API Key here"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 256
	ti.SetWidth(20)
	log.Println("Exited the apiView() func")
	m.state = SetApiState
	m.textInput = ti
	m.updateStyles(true)
	return m
}

// Set Cargo Size
func setCargoSizeView(m model) model {
	log.Println("Entered the setCargoSizeView() func")
	ti := textinput.New()
	ti.Placeholder = "Enter total cargo size"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 256
	ti.SetWidth(20)
	m.state = SetCargoSizeState
	m.textInput = ti
	m.updateStyles(true)
	log.Println("Exited the setCargoSizeView() func")
	return m
}

func TradeView(m model) model {
	log.Println("Commodities")
	resp, err := uex.GetCommmoddityPrices()
	m.state = TradeState
	if err != nil {
		log.Println("Can't connect to UEX")
		return m
	}
	log.Println("Commodities: ", resp)
	return m
}

func mainList(m model) model {
	log.Println("Entered the MainList() func")
	items := []list.Item{
		item(TradeStr),
		item(SettingsStr),
		item(ExitStr),
	}
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "List of options"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m = model{state: MainState, list: l, choice: "Back to Main Menu"}
	m.updateStyles(true)
	log.Println("Exiting the MainList() func")
	return m

}

// Main View
func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("\n  See you later!\n\n")
	}
	if m.state == MainState {
		myFig := figure.NewFigure("TDD Terminal", "", true).String()
		return tea.NewView("\n\n\n" + myFig + "\n" + m.list.View())
	}
	if m.state == SetApiState {
		var c *tea.Cursor
		if !m.textInput.VirtualCursor() {
			c := m.textInput.Cursor()
			c.Y += lipgloss.Height(m.headerView())

		}
		str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())
		v := tea.NewView(str)
		v.Cursor = c
		return v
	}
	if m.state == SetCargoSizeState {
		var c *tea.Cursor
		if !m.textInput.VirtualCursor() {
			c := m.textInput.Cursor()
			c.Y += lipgloss.Height(m.headerView())

		}
		str := lipgloss.JoinVertical(lipgloss.Top, m.cargoSizeHeader(), m.textInput.View(), m.footerView())
		v := tea.NewView(str)
		v.Cursor = c
		return v
	}
	if m.state == TradeState {
		var c *tea.Cursor
		if !m.textInput.VirtualCursor() {
			c := m.textInput.Cursor()
			c.Y += lipgloss.Height(m.headerView())

		}
		str := lipgloss.JoinVertical(lipgloss.Top, m.cargoSizeHeader(), m.textInput.View(), m.footerView())
		v := tea.NewView(str)
		v.Cursor = c
		return v
	}

	return tea.NewView("\n" + m.list.View())
	//return tea.NewView(mainStyle.Render("\n" + s + "\n"))
	//return tea.NewView(m.settings.View())
}

func (m model) headerView() string { return fmt.Sprintf("\nCurrent API key: %s\n", api) }
func (m model) cargoSizeHeader() string {
	return fmt.Sprintf("\nCurrent Total SCU set to: %d", cargoSize)
}
func (m model) footerView() string { return "\n(q to quit)" }

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	UEXClient = uex.ClientConfig("https://api.uexcorp.space/2.0", "")
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
