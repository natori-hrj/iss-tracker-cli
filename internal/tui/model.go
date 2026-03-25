package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/natori-hrj/iss-tracker-cli/internal/api"
	"github.com/natori-hrj/iss-tracker-cli/internal/geo"
	"github.com/natori-hrj/iss-tracker-cli/internal/ui"
)

type tickMsg time.Time

type dataMsg struct {
	pos       *api.ISSPosition
	astros    *api.AstrosResponse
	location  *geo.Location
	posErr    error
	astroErr  error
	locErr    error
}

type Model struct {
	client    *api.Client
	pos       *api.ISSPosition
	astros    *api.AstrosResponse
	location  *geo.Location
	posErr    error
	astroErr  error
	locErr    error
	quitting  bool
	interval  time.Duration
}

func NewModel(interval time.Duration) Model {
	return Model{
		client:   api.NewClient(),
		interval: interval,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchData(m.client), tickCmd(m.interval))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "r":
			return m, fetchData(m.client)
		}
	case tickMsg:
		return m, tea.Batch(fetchData(m.client), tickCmd(m.interval))
	case dataMsg:
		if msg.pos != nil {
			m.pos = msg.pos
		}
		m.posErr = msg.posErr
		if msg.astros != nil {
			m.astros = msg.astros
		}
		m.astroErr = msg.astroErr
		if msg.location != nil {
			m.location = msg.location
		}
		m.locErr = msg.locErr
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	// Title
	sb.WriteString(ui.TitleStyle.Render(" ISS Tracker "))
	sb.WriteString("\n\n")

	// Map
	if m.pos != nil {
		showUser := m.location != nil
		userLat, userLon := 0.0, 0.0
		if m.location != nil {
			userLat = m.location.Latitude
			userLon = m.location.Longitude
		}
		mapView := ui.RenderMap(m.pos.Latitude, m.pos.Longitude, userLat, userLon, showUser)
		sb.WriteString(ui.MapStyle.Render(mapView))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("  %s ISS  %s You\n",
			ui.ISSMarkerStyle.Render("★"),
			ui.UserMarkerStyle.Render("◉")))
	} else if m.posErr != nil {
		sb.WriteString(ui.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.posErr)))
		sb.WriteString("\n")
	} else {
		sb.WriteString("  Loading ISS position...\n")
	}

	// ISS Info
	var infoLines []string
	if m.pos != nil {
		infoLines = append(infoLines,
			fmt.Sprintf("%s %s",
				ui.LabelStyle.Render("Position:"),
				ui.ValueStyle.Render(fmt.Sprintf("%.4f°N, %.4f°E", m.pos.Latitude, m.pos.Longitude))),
			fmt.Sprintf("%s %s",
				ui.LabelStyle.Render("Updated: "),
				ui.ValueStyle.Render(m.pos.Timestamp.Local().Format("2006-01-02 15:04:05"))),
		)
	}

	if m.location != nil && m.pos != nil {
		dist := geo.DistanceToISS(m.location.Latitude, m.location.Longitude, m.pos.Latitude, m.pos.Longitude)
		infoLines = append(infoLines,
			fmt.Sprintf("%s %s",
				ui.LabelStyle.Render("Your Location:"),
				ui.ValueStyle.Render(fmt.Sprintf("%s, %s (%.4f°N, %.4f°E)",
					m.location.City, m.location.Country, m.location.Latitude, m.location.Longitude))),
			fmt.Sprintf("%s %s",
				ui.LabelStyle.Render("Distance:"),
				ui.ValueStyle.Render(fmt.Sprintf("%.0f km", dist))),
		)

		nextPass := geo.EstimateNextPass(m.location.Latitude, m.location.Longitude, m.pos.Latitude, m.pos.Longitude)
		infoLines = append(infoLines,
			fmt.Sprintf("%s %s",
				ui.LabelStyle.Render("Next Pass (est):"),
				ui.ValueStyle.Render(nextPass.Local().Format("15:04:05"))),
		)
	} else if m.locErr != nil {
		infoLines = append(infoLines,
			ui.ErrorStyle.Render(fmt.Sprintf("Location error: %v", m.locErr)))
	}

	if len(infoLines) > 0 {
		sb.WriteString(ui.InfoStyle.Render(strings.Join(infoLines, "\n")))
		sb.WriteString("\n")
	}

	// Astronauts
	if m.astros != nil {
		var astroLines []string
		astroLines = append(astroLines,
			fmt.Sprintf("%s %s",
				ui.LabelStyle.Render("Astronauts in Space:"),
				ui.ValueStyle.Render(fmt.Sprintf("%d", m.astros.Number))))

		// Group by craft
		craftMap := make(map[string][]string)
		for _, a := range m.astros.People {
			craftMap[a.Craft] = append(craftMap[a.Craft], a.Name)
		}
		for craft, names := range craftMap {
			astroLines = append(astroLines,
				fmt.Sprintf("  %s", ui.LabelStyle.Render(craft+":")))
			for _, name := range names {
				astroLines = append(astroLines,
					fmt.Sprintf("    %s", ui.AstronautStyle.Render("  "+name)))
			}
		}
		sb.WriteString(ui.InfoStyle.Render(strings.Join(astroLines, "\n")))
		sb.WriteString("\n")
	} else if m.astroErr != nil {
		sb.WriteString(ui.ErrorStyle.Render(fmt.Sprintf("Astronaut error: %v", m.astroErr)))
		sb.WriteString("\n")
	}

	// Notice
	sb.WriteString(ui.WarnStyle.Render("  Note: Uses HTTP (not HTTPS). Your IP is sent to ip-api.com for geolocation."))
	sb.WriteString("\n")

	// Help
	sb.WriteString(ui.HelpStyle.Render("  r: refresh  q/esc: quit"))
	sb.WriteString("\n")

	return sb.String()
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchData(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		msg := dataMsg{}

		// Fetch all data concurrently using channels
		type posResult struct {
			pos *api.ISSPosition
			err error
		}
		type astroResult struct {
			astros *api.AstrosResponse
			err    error
		}
		type locResult struct {
			loc *geo.Location
			err error
		}

		posCh := make(chan posResult, 1)
		astroCh := make(chan astroResult, 1)
		locCh := make(chan locResult, 1)

		go func() {
			pos, err := client.GetISSPosition()
			posCh <- posResult{pos, err}
		}()
		go func() {
			astros, err := client.GetAstronauts()
			astroCh <- astroResult{astros, err}
		}()
		go func() {
			loc, err := geo.GetMyLocation()
			locCh <- locResult{loc, err}
		}()

		pr := <-posCh
		msg.pos = pr.pos
		msg.posErr = pr.err

		ar := <-astroCh
		msg.astros = ar.astros
		msg.astroErr = ar.err

		lr := <-locCh
		msg.location = lr.loc
		msg.locErr = lr.err

		return msg
	}
}
