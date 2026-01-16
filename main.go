package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

type model struct {
	items      []string
	cursor     int
	choice     string
	currentVer string
}

func initialModel(items []string, current string) model {
	m := model{items: items, currentVer: strings.TrimSpace(current)}
	if m.currentVer != "" {
		for i, item := range items {
			if strings.TrimSpace(item) == m.currentVer {
				m.cursor = i
				break
			}
		}
	}
	return m
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			m.choice = m.items[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := ""
	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = "\x1b[32m⚬\x1b[0m "
		}
		s += fmt.Sprintf("%s%s\n", cursor, item)
	}
	return s
}

func printHelp() {
	fmt.Println(`v py
v python
v node
v nodejs
v -v
v version
v --version
↑/↓ 或 k/j 移动
Enter 选择并写入
Esc / Ctrl+C 退出不修改`)
}

func printVersion() {
	version := strings.TrimPrefix(Version, "v")
	fmt.Printf("v %s\n", version)
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	plugin := os.Args[1]
	switch plugin {
	case "-h", "--help", "-help", "help":
		printHelp()
		os.Exit(0)
	case "-v", "--version", "version":
		printVersion()
		os.Exit(0)
	}
	switch plugin {
	case "node":
		plugin = "nodejs"
	case "py":
		plugin = "python"
	}

	root := findProjectRoot()
	if root == "" {
		root, _ = os.Getwd()
	}

	versions := asdfList(plugin)
	if len(versions) == 0 {
		fmt.Fprintln(os.Stderr, "no versions found")
		os.Exit(1)
	}

	current := readCurrentVersion(root, plugin)

	p := tea.NewProgram(initialModel(versions, current))
	m, err := p.Run()
	if err != nil {
		os.Exit(1)
	}

	result := m.(model)
	if result.choice == "" {
		os.Exit(0)
	}

	if strings.TrimSpace(current) != "" && strings.TrimSpace(result.choice) == strings.TrimSpace(current) {
		os.Exit(0)
	}

	if err := writeToolVersions(root, plugin, result.choice); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("\x1b[32mCurrent: %s\x1b[0m\n", result.choice)
}

//
// helpers
//

func asdfList(plugin string) []string {
	out, err := exec.Command("asdf", "list", plugin).Output()
	if err != nil {
		return nil
	}

	var items []string
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		// 去掉 asdf 当前版本标记 *
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)

		if line != "" {
			items = append(items, line)
		}
	}
	return items
}

func findProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if exists(filepath.Join(dir, ".tool-versions")) ||
			exists(filepath.Join(dir, ".git")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readCurrentVersion(root, plugin string) string {
	path := filepath.Join(root, ".tool-versions")
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		var k, v string
		fmt.Sscan(sc.Text(), &k, &v)
		if k == plugin {
			return v
		}
	}
	return ""
}

func writeToolVersions(root, plugin, version string) error {
	path := filepath.Join(root, ".tool-versions")

	m := map[string]string{}

	if b, err := os.ReadFile(path); err == nil {
		sc := bufio.NewScanner(bytes.NewReader(b))
		for sc.Scan() {
			var k, v string
			fmt.Sscan(sc.Text(), &k, &v)
			if k != "" {
				m[k] = v
			}
		}
	}

	m[plugin] = version

	var out bytes.Buffer
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&out, "%s %s\n", k, m[k])
	}

	return os.WriteFile(path, out.Bytes(), 0644)
}
