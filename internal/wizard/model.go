package wizard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fabricatorsltd/go-tools-cli/internal/catalog"
)

type Model struct {
	currentStep int
	Done        bool
	quitting    bool
	Config      Config

	// Step 0, 1: text inputs
	textInputs [2]textinput.Model

	// Step 2: preset selection
	presetCursor int

	// Step 3: module multi-select
	moduleList       []catalog.Module
	moduleCursors    map[int]struct{} // selected indices
	moduleListCursor int

	// Step 4: per-module options
	optionQueue    []optionTask // queue of options to ask
	optionCursor   int          // current selection within a single option's choices
	currentOptTask int          // index into optionQueue

	// Step 5: output depth
	depthCursor int

	// error message
	err string
}

type optionTask struct {
	ModuleName string
	Option     catalog.ModuleOption
}

type Config struct {
	ProjectName     string
	ModulePath      string
	PresetIndex     int
	SelectedModules []string
	ModuleOptions   map[string]map[string]string // module name -> option key -> value
	OutputDepth     OutputDepth
	OutputDir       string
}

func NewModel(outputDir string) Model {
	ti0 := textinput.New()
	ti0.Placeholder = "my-app"
	ti0.Focus()
	ti0.CharLimit = 64
	ti0.Width = 40

	ti1 := textinput.New()
	ti1.Placeholder = "github.com/username/my-app"
	ti1.CharLimit = 128
	ti1.Width = 60

	return Model{
		Config: Config{
			OutputDir:     outputDir,
			ModuleOptions: make(map[string]map[string]string),
		},
		textInputs:    [2]textinput.Model{ti0, ti1},
		moduleCursors: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEsc:
			if m.currentStep > 0 {
				m.currentStep--
			}
			return m, nil
		}

		switch m.currentStep {
		case StepProjectName:
			return m.handleTextInput(msg, 0)
		case StepModulePath:
			return m.handleTextInput(msg, 1)
		case StepPreset:
			return m.handlePresetSelect(msg)
		case StepModules:
			return m.handleModuleSelect(msg)
		case StepModuleOpts:
			return m.handleModuleOpts(msg)
		case StepOutputDepth:
			return m.handleDepthSelect(msg)
		case StepConfirm:
			return m.handleConfirm(msg)
		}
	}

	// Forward to active text input
	if m.currentStep == StepProjectName || m.currentStep == StepModulePath {
		idx := m.currentStep
		var cmd tea.Cmd
		m.textInputs[idx], cmd = m.textInputs[idx].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleTextInput(msg tea.KeyMsg, idx int) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyEnter {
		val := strings.TrimSpace(m.textInputs[idx].Value())
		if idx == 0 {
			if val == "" {
				val = "my-app"
			}
			m.Config.ProjectName = val
			// Pre-fill module path placeholder
			m.textInputs[1].Placeholder = fmt.Sprintf("github.com/username/%s", val)
			m.textInputs[1].Focus()
			m.textInputs[0].Blur()
		} else {
			if val == "" {
				val = fmt.Sprintf("github.com/username/%s", m.Config.ProjectName)
			}
			m.Config.ModulePath = val
		}
		m.currentStep++
		return m, nil
	}
	var cmd tea.Cmd
	m.textInputs[idx], cmd = m.textInputs[idx].Update(msg)
	return m, cmd
}

func (m Model) handlePresetSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.presetCursor > 0 {
			m.presetCursor--
		}
	case tea.KeyDown:
		if m.presetCursor < len(catalog.Presets)-1 {
			m.presetCursor++
		}
	case tea.KeyEnter:
		m.Config.PresetIndex = m.presetCursor
		preset := catalog.Presets[m.presetCursor]
		// Build module list from all modules
		m.moduleList = make([]catalog.Module, 0, len(catalog.Modules))
		for _, mod := range catalog.Modules {
			m.moduleList = append(m.moduleList, mod)
		}
		// Pre-select preset modules
		m.moduleCursors = make(map[int]struct{})
		for i, mod := range m.moduleList {
			for _, pmod := range preset.Modules {
				if mod.Name == pmod {
					m.moduleCursors[i] = struct{}{}
					break
				}
			}
		}
		m.currentStep++
	}
	return m, nil
}

func (m Model) handleModuleSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.moduleListCursor > 0 {
			m.moduleListCursor--
		}
	case tea.KeyDown:
		if m.moduleListCursor < len(m.moduleList)-1 {
			m.moduleListCursor++
		}
	case tea.KeyRunes:
		if msg.String() == " " {
			if _, ok := m.moduleCursors[m.moduleListCursor]; ok {
				delete(m.moduleCursors, m.moduleListCursor)
			} else {
				m.moduleCursors[m.moduleListCursor] = struct{}{}
			}
		}
	case tea.KeyEnter:
		// Collect selected modules
		selected := []string{}
		for i, mod := range m.moduleList {
			if _, ok := m.moduleCursors[i]; ok {
				selected = append(selected, mod.Name)
			}
		}
		m.Config.SelectedModules = selected

		// Build option queue
		m.optionQueue = []optionTask{}
		for _, name := range selected {
			mod := catalog.FindModule(name)
			if mod == nil {
				continue
			}
			for _, opt := range mod.Options {
				m.optionQueue = append(m.optionQueue, optionTask{ModuleName: name, Option: opt})
			}
		}
		m.currentOptTask = 0
		m.optionCursor = 0

		if len(m.optionQueue) == 0 {
			m.currentStep = StepOutputDepth
		} else {
			m.currentStep = StepModuleOpts
		}
	}
	return m, nil
}

func (m Model) handleModuleOpts(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.currentOptTask >= len(m.optionQueue) {
		m.currentStep = StepOutputDepth
		return m, nil
	}
	task := m.optionQueue[m.currentOptTask]
	switch msg.Type {
	case tea.KeyUp:
		if m.optionCursor > 0 {
			m.optionCursor--
		}
	case tea.KeyDown:
		if m.optionCursor < len(task.Option.Choices)-1 {
			m.optionCursor++
		}
	case tea.KeyEnter:
		chosen := task.Option.Choices[m.optionCursor]
		if m.Config.ModuleOptions[task.ModuleName] == nil {
			m.Config.ModuleOptions[task.ModuleName] = make(map[string]string)
		}
		m.Config.ModuleOptions[task.ModuleName][task.Option.Key] = chosen
		m.currentOptTask++
		m.optionCursor = 0
		if m.currentOptTask >= len(m.optionQueue) {
			m.currentStep = StepOutputDepth
		}
	}
	return m, nil
}

func (m Model) handleDepthSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.depthCursor > 0 {
			m.depthCursor--
		}
	case tea.KeyDown:
		if m.depthCursor < len(outputDepths)-1 {
			m.depthCursor++
		}
	case tea.KeyEnter:
		m.Config.OutputDepth = OutputDepth(m.depthCursor)
		m.currentStep = StepConfirm
	}
	return m, nil
}

func (m Model) handleConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		m.Done = true
		return m, tea.Quit
	case tea.KeyRunes:
		if msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the current step
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	// Header
	sb.WriteString(titleStyle.Render("◆ go-tools-cli"))
	sb.WriteString("\n")
	sb.WriteString(subtitleStyle.Render("Bootstrap a new Go project with the go-tools ecosystem"))
	sb.WriteString("\n\n")

	// Progress
	sb.WriteString(progressStyle.Render(fmt.Sprintf("Step %d/%d", m.currentStep+1, TotalSteps)))
	sb.WriteString("\n\n")

	switch m.currentStep {
	case StepProjectName:
		sb.WriteString(labelStyle.Render("Project name:"))
		sb.WriteString("\n")
		sb.WriteString(m.textInputs[0].View())
		sb.WriteString("\n\n")
		sb.WriteString(inactiveStyle.Render("Press Enter to continue • Esc to go back"))

	case StepModulePath:
		sb.WriteString(labelStyle.Render("Go module path:"))
		sb.WriteString("\n")
		sb.WriteString(m.textInputs[1].View())
		sb.WriteString("\n\n")
		sb.WriteString(inactiveStyle.Render("Press Enter to continue • Esc to go back"))

	case StepPreset:
		sb.WriteString(labelStyle.Render("Choose a preset:"))
		sb.WriteString("\n\n")
		for i, p := range catalog.Presets {
			cursor := "  "
			name := inactiveStyle.Render(p.Name)
			desc := inactiveStyle.Render("— " + p.Description)
			if i == m.presetCursor {
				cursor = activeStyle.Render("❯ ")
				name = activeStyle.Render(p.Name)
				desc = subtitleStyle.Render("— " + p.Description)
			}
			sb.WriteString(fmt.Sprintf("%s%s %s\n", cursor, name, desc))
		}
		sb.WriteString("\n")
		sb.WriteString(inactiveStyle.Render("↑/↓ navigate • Enter to select • Esc to go back"))

	case StepModules:
		sb.WriteString(labelStyle.Render("Select modules (Space to toggle, Enter to confirm):"))
		sb.WriteString("\n\n")
		// Group by tag
		currentTag := ""
		for i, mod := range m.moduleList {
			tag := ""
			if len(mod.Tags) > 0 {
				tag = mod.Tags[0]
			}
			if tag != currentTag {
				currentTag = tag
				sb.WriteString(subtitleStyle.Render(fmt.Sprintf("  [%s]", tag)))
				sb.WriteString("\n")
			}

			_, isSelected := m.moduleCursors[i]
			check := unselectedStyle.Render("○")
			if isSelected {
				check = selectedStyle.Render("✓")
			}

			cursor := "  "
			name := inactiveStyle.Render(mod.Name)
			desc := inactiveStyle.Render(mod.Description)
			if i == m.moduleListCursor {
				cursor = activeStyle.Render("❯ ")
				name = activeStyle.Render(mod.Name)
				desc = subtitleStyle.Render(mod.Description)
			}
			sb.WriteString(fmt.Sprintf("  %s %s%s  %s\n", check, cursor, name, desc))
		}
		sb.WriteString("\n")
		sb.WriteString(inactiveStyle.Render("↑/↓ navigate • Space toggle • Enter confirm • Esc back"))

	case StepModuleOpts:
		if m.currentOptTask < len(m.optionQueue) {
			task := m.optionQueue[m.currentOptTask]
			sb.WriteString(labelStyle.Render(fmt.Sprintf("[%s] %s:", task.ModuleName, task.Option.Question)))
			sb.WriteString("\n\n")
			for i, choice := range task.Option.Choices {
				cursor := "  "
				name := inactiveStyle.Render(choice)
				if i == m.optionCursor {
					cursor = activeStyle.Render("❯ ")
					name = activeStyle.Render(choice)
				}
				sb.WriteString(fmt.Sprintf("  %s%s\n", cursor, name))
			}
			sb.WriteString("\n")
			sb.WriteString(progressStyle.Render(fmt.Sprintf("Option %d/%d", m.currentOptTask+1, len(m.optionQueue))))
			sb.WriteString("\n")
			sb.WriteString(inactiveStyle.Render("↑/↓ navigate • Enter select • Esc back"))
		}

	case StepOutputDepth:
		sb.WriteString(labelStyle.Render("How much boilerplate to generate?"))
		sb.WriteString("\n\n")
		descriptions := []string{
			"go.mod + main.go + README.md",
			"Minimal + internal/ stubs for each selected module",
			"Boilerplate + wired main.go + Makefile",
		}
		for i, d := range outputDepths {
			cursor := "  "
			name := inactiveStyle.Render(d)
			desc := inactiveStyle.Render("— " + descriptions[i])
			if i == m.depthCursor {
				cursor = activeStyle.Render("❯ ")
				name = activeStyle.Render(d)
				desc = subtitleStyle.Render("— " + descriptions[i])
			}
			sb.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, name, desc))
		}
		sb.WriteString("\n")
		sb.WriteString(inactiveStyle.Render("↑/↓ navigate • Enter select • Esc back"))

	case StepConfirm:
		sb.WriteString(successStyle.Render("◆ Ready to scaffold!"))
		sb.WriteString("\n\n")
		content := strings.Builder{}
		content.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Project:"), valueStyle.Render(m.Config.ProjectName)))
		content.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Module: "), valueStyle.Render(m.Config.ModulePath)))
		preset := catalog.Presets[m.Config.PresetIndex]
		content.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Preset: "), valueStyle.Render(preset.Name)))
		content.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Depth:  "), valueStyle.Render(m.Config.OutputDepth.String())))
		content.WriteString(fmt.Sprintf("  %s\n", labelStyle.Render("Modules:")))
		for _, name := range m.Config.SelectedModules {
			content.WriteString(fmt.Sprintf("    %s %s\n", selectedStyle.Render("✓"), valueStyle.Render(name)))
		}
		if len(m.Config.ModuleOptions) > 0 {
			content.WriteString(fmt.Sprintf("  %s\n", labelStyle.Render("Options:")))
			for mod, opts := range m.Config.ModuleOptions {
				for k, v := range opts {
					content.WriteString(fmt.Sprintf("    %s %s=%s\n", selectedStyle.Render("→"), labelStyle.Render(mod+"."+k), valueStyle.Render(v)))
				}
			}
		}
		content.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Output: "), valueStyle.Render(m.Config.OutputDir)))
		sb.WriteString(boxStyle.Render(content.String()))
		sb.WriteString("\n\n")
		sb.WriteString(activeStyle.Render("Press Enter to generate  •  q to abort"))
	}

	return sb.String()
}
