package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00FFAA")).
	MarginBottom(1)

var panelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#888888")).
	Padding(0, 2).
	Width(40).
	Height(7)

var tickerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#D11141")).
	Padding(0, 1)

type Company struct {
	Name        string
	Cash        int
	MarketShare int
	LastMove    string
	TechLevel   int
	History     []string
	IsScouted   bool
	IsActive    bool
}

type MarketEvent struct {
	Headline string
	Effect   func(m *model)
}

type model struct {
	turnCounter  int
	player       Company
	aiRival      Company
	vulture      Company
	newsTicker   string
	gameOver     bool
	tickerOffset int
	windowWidth  int
	windowHeight int
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialPlayerModel() model {
	return model{
		turnCounter: 0,
		player: Company{
			Name:        "You",
			Cash:        50,
			MarketShare: 50,
			LastMove:    "None",
			TechLevel:   0,
			History:     []string{},
			IsScouted:   true, // See yourself
			IsActive:    true,
		},
		aiRival: Company{
			Name:        "Whale Industries",
			Cash:        50,
			MarketShare: 50,
			LastMove:    "None",
			TechLevel:   0,
			History:     []string{},
			IsScouted:   false, //Hidden because we shouldn't be able to see unless we scout
			IsActive:    true,
		},
		vulture: Company{
			Name:        "Vulture Inc.",
			Cash:        200,
			MarketShare: 0,
			LastMove:    "Lurking...",
			TechLevel:   2,
			History:     []string{},
			IsScouted:   true, // You can always see yourself
			IsActive:    false,
		},
		newsTicker: "",
		gameOver:   false,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.tickerOffset++
		return m, tick()

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			newState := initialPlayerModel()
			newState.windowWidth = m.windowWidth
			newState.windowHeight = m.windowHeight

			return newState, nil
		}

		if m.gameOver {
			return m, nil
		}

		switch msg.String() {
		case "u":
			m.resolveQuarter("Undercut")
		case "m":
			m.resolveQuarter("Match")
		case "p":
			m.resolveQuarter("Premium")
		case "i":
			if !m.gameOver {
				m.handleStore("Invest")
			}
			return m, nil
		case "s":
			if !m.gameOver {
				m.handleStore("Scout")
			}
			return m, nil

		default:
			return m, nil
		}

		m.turnCounter++
		m.tickerOffset = 0
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.gameOver {
		finalScreen := fmt.Sprintf("--- GAME OVER ---\n\nFinal Quarter: %d\nYour Cash: ₹%dM\nAI Cash: ₹%dM\n\nPress 'q' to quit or 'r' to restart", m.turnCounter, m.player.Cash, m.aiRival.Cash)
		centeredGameOver := lipgloss.Place(
			m.windowWidth,
			m.windowHeight,
			lipgloss.Center,
			lipgloss.Center,
			finalScreen,
		)

		return tea.NewView(centeredGameOver)
	}

	title := titleStyle.Render("=== THE COMPUTE CARTEL ===")

	playerText := fmt.Sprintf("%s (Tech: Lvl %d)\n\nCash: ₹%dM\nMarket Share: %d%%\nLast Move: %s",
		m.player.Name, m.player.TechLevel, m.player.Cash, m.player.MarketShare, m.player.LastMove)
	playerBox := panelStyle.BorderForeground(lipgloss.Color("#00FFAA")).Render(playerText)

	// 2. AI Rival Panel (FOG OF WAR LOGIC)
	var aiText string
	if m.aiRival.IsScouted {
		aiText = fmt.Sprintf("%s (Tech: Lvl %d)\n\nCash: ₹%dM\nMarket Share: %d%%\nLast Move: %s",
			m.aiRival.Name, m.aiRival.TechLevel, m.aiRival.Cash, m.aiRival.MarketShare, m.aiRival.LastMove)
	} else {
		// Hide the exact stats!
		aiText = fmt.Sprintf("%s (Tech: ???)\n\nCash: ₹???M\nMarket Share: ???%%\nLast Move: %s",
			m.aiRival.Name, m.aiRival.LastMove)
	}
	aiBox := panelStyle.BorderForeground(lipgloss.Color("#FF0055")).Render(aiText)

	// 3. The Vulture Panel (Only shows if Active)
	var dashboard string
	if m.vulture.IsActive {
		vultureText := fmt.Sprintf("%s\n\nCash: ₹%dM\nMarket Share: %d%%\nLast Move: %s",
			m.vulture.Name, m.vulture.Cash, m.vulture.MarketShare, m.vulture.LastMove)
		vultureBox := panelStyle.BorderForeground(lipgloss.Color("#FFAA00")).Render(vultureText)
		dashboard = lipgloss.JoinHorizontal(lipgloss.Top, playerBox, "  ", aiBox, "  ", vultureBox)
	} else {
		dashboard = lipgloss.JoinHorizontal(lipgloss.Top, playerBox, "  ", aiBox)
	}

	tickerText := m.newsTicker
	if tickerText == "" {
		tickerText = "AWAITING Q1 MARKET OPEN... STAND BY."
	}

	// We add massive empty space at the end so it loops seamlessly
	rawText := fmt.Sprintf(" --- BREAKING NEWS: %s ---                                        ", tickerText)
	runes := []rune(rawText)

	offset := m.tickerOffset % len(runes)
	shifted := append(runes[offset:], runes[:offset]...)

	safeWidth := 60
	if m.windowWidth < 65 {
		safeWidth = m.windowWidth - 10
		safeWidth = max(safeWidth, 10)
	}
	if len(shifted) > safeWidth {
		shifted = shifted[:safeWidth]
	}

	tickerBox := tickerStyle.Render(string(shifted))

	content := fmt.Sprintf("%s\nQuarter: %d\n\n%s\n%s\n[u] Undercut  [m] Match  [p] Premium  [i] Invest(₹10M)  [s] Scout(₹3M)  [q] Quit  [r] Restart ",
		title, m.turnCounter, dashboard, tickerBox)

	centeredScreen := lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	return tea.NewView(centeredScreen)
}

func (m *model) handleStore(purchase string) {
	if purchase == "Invest" {
		if m.player.Cash >= 10 {
			m.player.Cash -= 10
			m.player.TechLevel++
			m.newsTicker = "SUCCESS: R&D funded! Tech Level increased."
		} else {
			m.newsTicker = "❌ ERROR: Insufficient funds for R&D (Needs ₹10M)."
		}
	}

	if purchase == "Scout" {
		if m.player.Cash >= 3 {
			m.player.Cash -= 3
			m.aiRival.IsScouted = true
			m.newsTicker = "SUCCESS: Corporate espionage revealed AI stats!"
		} else {
			m.newsTicker = "❌ ERROR: Insufficient funds for Scouting (Needs ₹3M)."
		}
	}
}

func (m *model) resolveQuarter(playerCurrentMove string) {
	if playerCurrentMove != "Scout" {
		m.aiRival.IsScouted = false
	}

	m.player.History = append(m.player.History, playerCurrentMove)

	if len(m.player.History) > 3 {
		m.player.History = m.player.History[1:]
	}

	if m.turnCounter == 3 && !m.vulture.IsActive {
		m.vulture.IsActive = true
		m.newsTicker = "⚠️ MEGACORP CLOUD ENTERS THE MARKET! Prices destabilize."
	}

	// Vulture Mechanics: Steals 3% share from both players every turn it's active
	if m.vulture.IsActive {
		m.vulture.LastMove = "Market Dump"
		m.vulture.MarketShare += 6
		m.player.MarketShare -= 3
		m.aiRival.MarketShare -= 3
		m.vulture.Cash -= 5
	}

	events := []MarketEvent{
		{
			Headline: "Global GPU Shortage! Hardware maintenance costs spike.",
			Effect:   func(m *model) { m.player.Cash -= 5; m.aiRival.Cash -= 5 },
		},
		{
			Headline: "Govt Subsidizes Domestic Tech! Tax breaks applied.",
			Effect:   func(m *model) { m.player.Cash += 3; m.aiRival.Cash += 3 },
		},
		{
			Headline: "Quiet Quarter. Market remains stable.",
			Effect:   func(m *model) {},
		},
		{
			Headline: "Power Grid Failure in Bengaluru! Servers go dark.",
			Effect:   func(m *model) { m.player.Cash -= 2; m.player.MarketShare -= 5 },
		},
		{
			Headline: "Corporate Espionage! Rival steals your pricing algorithm.",
			Effect:   func(m *model) { m.aiRival.MarketShare += 10; m.player.MarketShare -= 10 },
		},
		{
			Headline: "Viral Open Source Model Released! Compute demand skyrockets.",
			Effect:   func(m *model) { m.player.Cash += 8; m.aiRival.Cash += 8 },
		},
	}
	choices := []string{"Undercut", "Match", "Premium"}

	//random event
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	selectedEvent := events[r.Intn(len(events))]
	m.newsTicker = selectedEvent.Headline
	selectedEvent.Effect(m)

	previousPlayerMove := m.player.LastMove

	if previousPlayerMove == "" || previousPlayerMove == "N/A" {
		m.aiRival.LastMove = "Match"
	} else {
		m.aiRival.LastMove = previousPlayerMove
	}

	if r.Intn(100) < 15 {
		m.aiRival.LastMove = choices[r.Intn(len(choices))]
	}

	m.player.LastMove = playerCurrentMove
	aiMove := m.aiRival.LastMove

	if m.aiRival.Cash >= 30 {
		m.aiRival.Cash -= 10
		m.aiRival.TechLevel++

		if m.aiRival.IsScouted {
			m.newsTicker = "⚠️ ALERT: The AI Rival just upgraded their Tech Level!"
		}
	}

	techBonus := m.player.TechLevel * 2
	aiTechBonus := m.aiRival.TechLevel * 2

	if playerCurrentMove == "Undercut" && aiMove == "Undercut" {
		m.player.Cash += (0 + techBonus)
		m.aiRival.Cash += (0 + aiTechBonus)
	} else if playerCurrentMove == "Undercut" && aiMove == "Match" {
		m.player.Cash += (3 + techBonus)
		m.aiRival.Cash += (1 + aiTechBonus)
		m.player.MarketShare += 5
		m.aiRival.MarketShare -= 5
	} else if playerCurrentMove == "Undercut" && aiMove == "Premium" {
		m.player.Cash += (5 + techBonus)
		m.aiRival.Cash += (0 + aiTechBonus)
		m.player.MarketShare += 10
		m.aiRival.MarketShare -= 10
	} else if playerCurrentMove == "Match" && aiMove == "Undercut" {
		m.player.Cash += (1 + techBonus)
		m.aiRival.Cash += (3 + aiTechBonus)
		m.player.MarketShare -= 5
		m.aiRival.MarketShare += 5
	} else if playerCurrentMove == "Match" && aiMove == "Premium" {
		m.player.Cash += (4 + techBonus)
		m.aiRival.Cash += (1 + aiTechBonus)
		m.player.MarketShare += 5
		m.aiRival.MarketShare -= 5
	} else if playerCurrentMove == "Match" && aiMove == "Match" {
		m.player.Cash += (2 + techBonus)
		m.aiRival.Cash += (2 + aiTechBonus)
	} else if playerCurrentMove == "Premium" && aiMove == "Undercut" {
		m.player.Cash += (0 + techBonus)
		m.aiRival.Cash += (5 + aiTechBonus)
		m.player.MarketShare -= 10
		m.aiRival.MarketShare += 10
	} else if playerCurrentMove == "Premium" && aiMove == "Match" {
		m.player.Cash += (1 + techBonus)
		m.aiRival.Cash += (4 + aiTechBonus)
		m.player.MarketShare -= 5
		m.aiRival.MarketShare += 5
	} else if playerCurrentMove == "Premium" && aiMove == "Premium" {
		m.player.Cash += (4 + techBonus)
		m.aiRival.Cash += (4 + aiTechBonus)
	}

	if len(m.player.History) == 3 &&
		m.player.History[0] == "Match" &&
		m.player.History[1] == "Match" &&
		m.player.History[2] == "Undercut" {

		m.player.Cash += 15
		m.player.MarketShare += 10
		m.aiRival.MarketShare -= 10
		m.newsTicker = "COMBO EXECUTED! The 'Rope-a-Dope' pricing strategy devours market share!"
		m.player.History = []string{} // Reset
	}

	if m.player.Cash <= 0 || m.aiRival.Cash <= 0 || m.player.MarketShare <= 10 || m.aiRival.MarketShare <= 10 || m.turnCounter >= 11 {
		m.gameOver = true
	}
}

func main() {
	m := initialPlayerModel()
	p := tea.NewProgram(m)

	//  Run the program and handle any fatal errors
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
