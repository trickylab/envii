// Package tui implements the interactive terminal UI for browsing and
// editing the vault.
package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Trickster-ID/envii/internal/model"
	"github.com/Trickster-ID/envii/internal/store"
)

// level represents the current navigation depth.
type level int

const (
	levelProjects level = iota
	levelEnvs
	levelVars
)

// inputMode is the active text-input action, if any.
type inputMode int

const (
	inputNone inputMode = iota
	inputAddProject
	inputAddEnv
	inputAddVarKey
	inputAddVarValue
	inputEditValue
)

// Model is the root Bubble Tea model.
type Model struct {
	vault      *model.Vault
	store      *store.Store
	passphrase string

	level   level
	pIdx    int // selected project index
	eIdx    int // selected env index
	vIdx    int // selected var index

	reveal map[int]bool // var index -> revealed

	input     textinput.Model
	inputMode inputMode
	pendingKey string // buffer for add-var key

	width  int
	height int

	status string
	errMsg string
	dirty  bool
}

// New constructs the TUI model.
func New(v *model.Vault, s *store.Store, passphrase string) Model {
	ti := textinput.New()
	ti.CharLimit = 4096
	ti.Width = 40

	return Model{
		vault:      v,
		store:      s,
		passphrase: passphrase,
		level:      levelProjects,
		reveal:     map[int]bool{},
		input:      ti,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.inputMode != inputNone {
			return m.updateInput(msg)
		}
		return m.updateNav(msg)

	case saveResult:
		return m.updateSave(msg)
	}
	return m, nil
}

func (m Model) updateNav(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.errMsg = ""
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)

	case "enter", "l", "right":
		m.descend()
	case "esc", "h", "left":
		m.ascend()

	case "a":
		m.startAdd()
	case "e":
		m.startEdit()
	case "d":
		m.deleteCurrent()
	case "r":
		m.toggleReveal()
	case "c":
		m.copyValue()
	case "s":
		return m, m.save()
	}
	return m, nil
}

func (m *Model) moveCursor(delta int) {
	switch m.level {
	case levelProjects:
		m.pIdx = clamp(m.pIdx+delta, len(m.vault.Projects))
	case levelEnvs:
		m.eIdx = clamp(m.eIdx+delta, len(m.currentProject().Envs))
	case levelVars:
		m.vIdx = clamp(m.vIdx+delta, len(m.currentEnv().Vars))
	}
}

func (m *Model) descend() {
	switch m.level {
	case levelProjects:
		if len(m.vault.Projects) == 0 {
			return
		}
		m.level = levelEnvs
		m.eIdx = 0
	case levelEnvs:
		if len(m.currentProject().Envs) == 0 {
			return
		}
		m.level = levelVars
		m.vIdx = 0
		m.reveal = map[int]bool{}
	}
}

func (m *Model) ascend() {
	switch m.level {
	case levelVars:
		m.level = levelEnvs
	case levelEnvs:
		m.level = levelProjects
	}
}

func (m *Model) toggleReveal() {
	if m.level == levelVars && len(m.currentEnv().Vars) > 0 {
		m.reveal[m.vIdx] = !m.reveal[m.vIdx]
	}
}

func (m *Model) copyValue() {
	if m.level != levelVars || len(m.currentEnv().Vars) == 0 {
		return
	}
	v := m.currentEnv().Vars[m.vIdx]
	if err := clipboard.WriteAll(v.Value); err != nil {
		m.errMsg = "clipboard unavailable: " + err.Error()
		return
	}
	m.status = fmt.Sprintf("copied %s to clipboard", v.Key)
}

func (m *Model) deleteCurrent() {
	switch m.level {
	case levelProjects:
		if len(m.vault.Projects) == 0 {
			return
		}
		m.vault.Projects = removeAt(m.vault.Projects, m.pIdx)
		m.pIdx = clamp(m.pIdx, len(m.vault.Projects))
	case levelEnvs:
		p := m.currentProject()
		if len(p.Envs) == 0 {
			return
		}
		p.Envs = removeAt(p.Envs, m.eIdx)
		m.eIdx = clamp(m.eIdx, len(p.Envs))
	case levelVars:
		e := m.currentEnv()
		if len(e.Vars) == 0 {
			return
		}
		e.Vars = removeAt(e.Vars, m.vIdx)
		m.vIdx = clamp(m.vIdx, len(e.Vars))
	}
	m.dirty = true
	m.status = "deleted (press s to save)"
}

func (m Model) save() tea.Cmd {
	return func() tea.Msg {
		if err := m.store.Save(m.vault, m.passphrase); err != nil {
			return saveResult{err: err}
		}
		return saveResult{}
	}
}

type saveResult struct{ err error }

// --- input handling ---

func (m *Model) startAdd() {
	m.input.SetValue("")
	m.input.Focus()
	switch m.level {
	case levelProjects:
		m.inputMode = inputAddProject
		m.input.Placeholder = "project name"
	case levelEnvs:
		m.inputMode = inputAddEnv
		m.input.Placeholder = "env name (e.g. dev)"
	case levelVars:
		m.inputMode = inputAddVarKey
		m.input.Placeholder = "KEY"
	}
}

func (m *Model) startEdit() {
	if m.level != levelVars || len(m.currentEnv().Vars) == 0 {
		return
	}
	m.inputMode = inputEditValue
	m.input.SetValue(m.currentEnv().Vars[m.vIdx].Value)
	m.input.Placeholder = "value"
	m.input.Focus()
}

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.inputMode = inputNone
		m.input.Blur()
		return m, nil
	case "enter":
		return m.commitInput()
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) commitInput() (tea.Model, tea.Cmd) {
	val := strings.TrimSpace(m.input.Value())
	switch m.inputMode {
	case inputAddProject:
		if val != "" {
			m.vault.Projects = append(m.vault.Projects, &model.Project{Name: val})
			m.dirty = true
		}
	case inputAddEnv:
		if val != "" {
			p := m.currentProject()
			p.Envs = append(p.Envs, &model.Env{Name: val})
			m.dirty = true
		}
	case inputAddVarKey:
		if val != "" {
			m.pendingKey = val
			m.input.SetValue("")
			m.input.Placeholder = "value"
			m.inputMode = inputAddVarValue
			return m, nil
		}
	case inputAddVarValue:
		e := m.currentEnv()
		e.Vars = append(e.Vars, &model.Var{Key: m.pendingKey, Value: val, Secret: isSecret(m.pendingKey)})
		m.pendingKey = ""
		m.dirty = true
	case inputEditValue:
		m.currentEnv().Vars[m.vIdx].Value = val
		m.dirty = true
	}
	m.inputMode = inputNone
	m.input.Blur()
	m.status = "changed (press s to save)"
	return m, nil
}

// --- helpers ---

func (m *Model) currentProject() *model.Project {
	if len(m.vault.Projects) == 0 {
		return &model.Project{}
	}
	return m.vault.Projects[m.pIdx]
}

func (m *Model) currentEnv() *model.Env {
	p := m.currentProject()
	if len(p.Envs) == 0 {
		return &model.Env{}
	}
	return p.Envs[m.eIdx]
}

func clamp(i, n int) int {
	if n == 0 {
		return 0
	}
	if i < 0 {
		return 0
	}
	if i >= n {
		return n - 1
	}
	return i
}

func removeAt[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}

// isSecret guesses whether a key likely holds a secret value.
func isSecret(key string) bool {
	k := strings.ToUpper(key)
	for _, hint := range []string{"SECRET", "TOKEN", "PASSWORD", "KEY", "PRIVATE", "CREDENTIAL", "API"} {
		if strings.Contains(k, hint) {
			return true
		}
	}
	return false
}
