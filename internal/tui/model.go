package tui

import (
	"fmt"
	"strings"
	"time"

	"lazymail/internal/app"
	"lazymail/internal/domain"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	panelSidebar = iota
	panelMailList
	panelViewer
)

// Model is the root Bubble Tea model for LazyMail phase-1 shell.
type Model struct {
	width       int
	height      int
	ready       bool
	activePanel int

	svc          *app.Service // nil when running without config
	account      string
	folders      []string
	activeFolder int
	mails        []domain.Mail
	selectedMail int
	openedMail   int

	viewer viewport.Model
	status string
}

// New creates a model with mock folders and messages (no accounts configured).
func New() Model {
	return newWithService(nil, "")
}

// NewWithService creates a model backed by the given app.Service.
func NewWithService(svc *app.Service, account string) Model {
	return newWithService(svc, account)
}

func newWithService(svc *app.Service, account string) Model {
	var folders []string
	if svc != nil {
		folders = svc.Folders(account)
	}
	if len(folders) == 0 {
		folders = []string{"Inbox", "Sent", "Drafts", "Archive"}
	}

	var mails []domain.Mail
	if svc != nil && len(svc.AccountNames()) > 0 {
		mails = svc.Mails(account, folders[0])
	}
	if mails == nil {
		mails = mockMailsForFolder(folders[0])
	}

	m := Model{
		activePanel:  panelMailList,
		svc:          svc,
		account:      account,
		folders:      folders,
		activeFolder: 0,
		mails:        mails,
		selectedMail: 0,
		openedMail:   0,
		viewer:       viewport.New(0, 0),
		status:       "Ready",
	}
	m.refreshViewerContent()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.resizePanels()
		return m, nil

	case tea.KeyMsg:
		if cmd := m.handleGlobalKey(msg); cmd != nil {
			return m, cmd
		}
		return m.handlePanelKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}

	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "Loading LazyMail..."
	}

	sidebarW, listW, viewerW, contentH := m.layoutMetrics()
	statusH := 1

	sidebar := m.renderSidebar(sidebarW, contentH)
	mailList := m.renderMailList(listW, contentH)
	viewer := m.renderViewer(viewerW, contentH)

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mailList, viewer)
	status := lipgloss.NewStyle().
		Width(m.width).
		Height(statusH).
		Padding(0, 1).
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Render(m.statusText())

	return lipgloss.JoinVertical(lipgloss.Left, content, status)
}

func (m *Model) handleGlobalKey(key tea.KeyMsg) tea.Cmd {
	switch key.String() {
	case "q", "ctrl+c":
		return tea.Quit
	case "tab":
		m.activePanel = (m.activePanel + 1) % 3
		m.status = fmt.Sprintf("Focus: %s", panelName(m.activePanel))
	case "shift+tab":
		m.activePanel = (m.activePanel + 2) % 3
		m.status = fmt.Sprintf("Focus: %s", panelName(m.activePanel))
	}
	return nil
}

func (m Model) handlePanelKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.activePanel {
	case panelSidebar:
		switch key.String() {
		case "up", "k":
			if m.activeFolder > 0 {
				m.activeFolder--
				m.loadFolder()
			}
		case "down", "j":
			if m.activeFolder < len(m.folders)-1 {
				m.activeFolder++
				m.loadFolder()
			}
		case "enter":
			m.loadFolder()
		}
	case panelMailList:
		switch key.String() {
		case "up", "k":
			if m.selectedMail > 0 {
				m.selectedMail--
			}
		case "down", "j":
			if m.selectedMail < len(m.mails)-1 {
				m.selectedMail++
			}
		case "g":
			m.selectedMail = 0
		case "G":
			if len(m.mails) > 0 {
				m.selectedMail = len(m.mails) - 1
			}
		case "u":
			m.toggleRead(m.selectedMail)
		case "d":
			m.deleteSelected()
		case "enter":
			m.openedMail = m.selectedMail
			m.status = "Opened selected mail"
		}
		m.refreshViewerContent()
	case panelViewer:
		var cmd tea.Cmd
		m.viewer, cmd = m.viewer.Update(key)
		return m, cmd
	}
	return m, nil
}

func (m Model) handleMouse(mouse tea.MouseMsg) (tea.Model, tea.Cmd) {
	if !m.ready || m.height <= 1 {
		return m, nil
	}

	sidebarW, listW, _, contentH := m.layoutMetrics()
	if mouse.Y >= contentH {
		return m, nil
	}

	panel := panelViewer
	relativeX := mouse.X - sidebarW - listW
	if mouse.X < sidebarW {
		panel = panelSidebar
		relativeX = mouse.X
	} else if mouse.X < sidebarW+listW {
		panel = panelMailList
		relativeX = mouse.X - sidebarW
	}
	_ = relativeX

	m.activePanel = panel

	if mouse.Button == tea.MouseButtonWheelDown {
		if panel == panelViewer {
			m.viewer.LineDown(2)
		}
		return m, nil
	}
	if mouse.Button == tea.MouseButtonWheelUp {
		if panel == panelViewer {
			m.viewer.LineUp(2)
		}
		return m, nil
	}

	if mouse.Action != tea.MouseActionPress {
		return m, nil
	}

	y := mouse.Y - 1
	if y < 0 {
		y = 0
	}

	switch panel {
	case panelSidebar:
		if y >= 0 && y < len(m.folders) {
			m.activeFolder = y
			m.loadFolder()
		}
	case panelMailList:
		if y >= 0 && y < len(m.mails) {
			m.selectedMail = y
			m.openedMail = y
			m.refreshViewerContent()
		}
	}

	return m, nil
}

func (m *Model) loadFolder() {
	folder := m.folders[m.activeFolder]
	var mails []domain.Mail
	if m.svc != nil {
		mails = m.svc.Mails(m.account, folder)
	}
	if len(mails) == 0 {
		mails = mockMailsForFolder(folder)
	}
	m.mails = mails
	m.selectedMail = 0
	m.openedMail = 0
	m.status = fmt.Sprintf("Folder: %s", folder)
	m.refreshViewerContent()
}

func (m *Model) toggleRead(idx int) {
	if idx < 0 || idx >= len(m.mails) {
		return
	}
	m.mails[idx].IsRead = !m.mails[idx].IsRead
	if m.svc != nil {
		m.svc.SetRead(m.mails[idx].ID, m.mails[idx].IsRead)
	}
	m.status = "Toggled read state"
}

func (m *Model) deleteSelected() {
	if len(m.mails) == 0 || m.selectedMail < 0 || m.selectedMail >= len(m.mails) {
		return
	}

	deletedSubject := m.mails[m.selectedMail].Subject
	if m.svc != nil {
		m.svc.Delete(m.mails[m.selectedMail].ID)
	}
	m.mails = append(m.mails[:m.selectedMail], m.mails[m.selectedMail+1:]...)
	if m.selectedMail >= len(m.mails) {
		m.selectedMail = len(m.mails) - 1
	}
	if m.selectedMail < 0 {
		m.selectedMail = 0
	}
	m.openedMail = m.selectedMail
	m.status = fmt.Sprintf("Deleted: %s", deletedSubject)
	m.refreshViewerContent()
}

func (m *Model) refreshViewerContent() {
	if len(m.mails) == 0 || m.openedMail < 0 || m.openedMail >= len(m.mails) {
		m.viewer.SetContent("No message selected.")
		return
	}

	mail := m.mails[m.openedMail]
	content := strings.Join([]string{
		fmt.Sprintf("Subject: %s", mail.Subject),
		fmt.Sprintf("From: %s", mail.From),
		fmt.Sprintf("To: %s", strings.Join(mail.To, ", ")),
		fmt.Sprintf("Date: %s", mail.Date.Format(time.RFC1123)),
		fmt.Sprintf("Folder: %s", mail.Folder),
		"",
		mail.Body,
	}, "\n")
	m.viewer.SetContent(content)
}

func (m *Model) resizePanels() {
	sidebarW, listW, viewerW, contentH := m.layoutMetrics()
	m.viewer.Width = max(1, viewerW-2)
	m.viewer.Height = max(1, contentH-2)
	_ = sidebarW
	_ = listW
	m.refreshViewerContent()
}

func (m Model) layoutMetrics() (int, int, int, int) {
	if m.width < 60 {
		return 18, 22, max(20, m.width-40), max(3, m.height-1)
	}
	sidebarW := 22
	listW := max(30, m.width/3)
	viewerW := max(20, m.width-sidebarW-listW)
	contentH := max(3, m.height-1)
	return sidebarW, listW, viewerW, contentH
}

func (m Model) panelStyle(panel int, width, height int) lipgloss.Style {
	style := lipgloss.NewStyle().Width(width).Height(height).Padding(0, 1).Border(lipgloss.NormalBorder())
	if panel == m.activePanel {
		return style.BorderForeground(lipgloss.Color("39"))
	}
	return style.BorderForeground(lipgloss.Color("240"))
}

func (m Model) renderSidebar(width, height int) string {
	lines := make([]string, 0, len(m.folders)+1)
	lines = append(lines, "Folders")
	for i, folder := range m.folders {
		cursor := " "
		if i == m.activeFolder {
			cursor = ">"
		}
		lines = append(lines, fmt.Sprintf("%s %s", cursor, folder))
	}
	return m.panelStyle(panelSidebar, width, height).Render(strings.Join(lines, "\n"))
}

func (m Model) renderMailList(width, height int) string {
	lines := []string{"Mail List"}
	if len(m.mails) == 0 {
		lines = append(lines, "(empty)")
	} else {
		for i, mail := range m.mails {
			cursor := " "
			if i == m.selectedMail {
				cursor = ">"
			}
			readState := "•"
			if mail.IsRead {
				readState = " "
			}
			lines = append(lines, fmt.Sprintf("%s%s %s", cursor, readState, trimRight(mail.Subject, max(8, width-8))))
		}
	}
	return m.panelStyle(panelMailList, width, height).Render(strings.Join(lines, "\n"))
}

func (m Model) renderViewer(width, height int) string {
	header := lipgloss.NewStyle().Bold(true).Render("Mail Viewer")
	body := m.viewer.View()
	return m.panelStyle(panelViewer, width, height).Render(header + "\n" + body)
}

func (m Model) statusText() string {
	folder := m.folders[m.activeFolder]
	position := "0/0"
	if len(m.mails) > 0 {
		position = fmt.Sprintf("%d/%d", m.selectedMail+1, len(m.mails))
	}
	return fmt.Sprintf("[Tab] Focus [Enter] Open [j/k] Move [d] Delete [u] Read [q] Quit | panel:%s folder:%s %s | %s",
		panelName(m.activePanel), folder, position, m.status)
}

func panelName(panel int) string {
	switch panel {
	case panelSidebar:
		return "sidebar"
	case panelMailList:
		return "mail-list"
	default:
		return "viewer"
	}
}

func trimRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) <= width {
		return s
	}
	if width <= 1 {
		return "…"
	}
	return string(runes[:width-1]) + "…"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
