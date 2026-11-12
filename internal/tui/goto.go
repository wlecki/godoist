package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Jump keys, so that a view or a project is one or two keys away from anywhere:
// g is a prefix like in vim and Todoist. gg goes to the top, gt to Today, gu to
// Upcoming, gi to the Inbox, and gp opens the go-to picker. ctrl+k opens the picker too.
// A key after g that is not one of these cancels the prefix and works as usual.

// gotoHint is the status while g waits for the second key.
const gotoHint = "g… · g top · t Today · u Upcoming · i Inbox · p go to…"

// gotoKeys handles the key after g. ok is false if the key is not a jump key: it then
// works as usual.
func (m Model) gotoKeys(key string) (tea.Model, tea.Cmd, bool) {
	m.gPrefix = false
	if m.status == gotoHint {
		m.setStatus("", false)
	}
	switch key {
	case "g":
		switch m.focus {
		case paneNav:
			m.moveNav(-len(m.nav))
		case paneDetail:
			m.detOff, m.comCur = 0, -1
		default:
			m.moveRow(-len(m.rows))
		}
	case "t":
		m.gotoKind(vkToday)
	case "u":
		m.gotoKind(vkUpcoming)
	case "i":
		m.gotoInbox()
	case "p":
		next := m.openGotoPicker()
		return m, next, true
	case "esc":
	default:
		return m, nil, false
	}
	return m, nil, true
}

// gotoNav opens the sidebar item with navKey key and puts the focus on the task list.
func (m *Model) gotoNav(key string) {
	m.pick = nil
	before := m.navCur
	m.selectNav(func(n navItem) bool { return navKey(n) == key })
	if m.navCur != before {
		m.buildRows(true)
	}
	m.detailOpen = false
	m.focus = paneTasks
}

// gotoKind opens the first sidebar item of kind k.
func (m *Model) gotoKind(k viewKind) {
	for _, n := range m.nav {
		if n.header == "" && n.kind == k {
			m.gotoNav(navKey(n))
			return
		}
	}
}

// gotoInbox opens the Inbox project.
func (m *Model) gotoInbox() {
	for _, n := range m.nav {
		if p := m.projects[n.projectID]; n.header == "" && n.kind == vkProject && p != nil && p.InboxProject {
			m.gotoNav(navKey(n))
			return
		}
	}
	m.setStatus("no Inbox project", false)
}

// openGotoPicker opens a filter list of the sidebar items: views, projects, and labels.
func (m *Model) openGotoPicker() tea.Cmd {
	var items []pickItem
	cur := 0
	for i, n := range m.nav {
		if n.header != "" || n.kind == vkLabelHint || n.name == "" {
			continue
		}
		if i == m.navCur {
			cur = len(items)
		}
		text := n.name
		switch n.kind {
		case vkProject:
			text = "# " + n.name
		case vkLabel:
			text = "@ " + n.name
		}
		items = append(items, pickItem{text: text, match: strings.ToLower(n.name), color: n.color, depth: n.depth, nav: navKey(n)})
	}
	m.pick = &picker{kind: pickGoto, title: "Go to", items: items, filter: newPickerFilter(), cur: cur}
	m.pick.filter.Placeholder = "type a view, project, or label"
	m.pick.filter.SetWidth(pickerFilterWidth(m.pickerWidth()))
	return m.pick.filter.Focus()
}
