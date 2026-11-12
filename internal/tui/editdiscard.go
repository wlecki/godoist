package tui

import tea "charm.land/bubbletea/v2"

// discardNoteEditor closes the editor and drops the changes (ctrl+q). esc saves, so this
// is the only way out without the new text. If autosave already sent a changed text,
// the text from before the editor opened goes back to Todoist.
func (m Model) discardNoteEditor() (tea.Model, tea.Cmd) {
	e := m.note
	m.note = nil
	if e.kind == editCommentNew || e.saved == e.orig && !e.saving {
		m.setStatus("changes discarded", false)
		return m, nil
	}
	e.lines = splitRunes(e.orig)
	next := m.saveNote(e, "Changes discarded · the old text is back")
	return m, next
}
