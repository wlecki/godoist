package tui

import "testing"

// ctrl+q closes the editor without saving. If autosave already saved the new text, the
// old text goes back to Todoist and shows at once.
func TestEditorDiscard(t *testing.T) {
	m := send(t, notesModel(t), key("E"))
	id, old := m.note.taskID, m.note.orig
	m = typeText(t, m, "!")
	m = send(t, m, key("ctrl+q"))
	if m.note != nil || m.pending != 0 {
		t.Fatalf("note open %v pending %d, want closed and no write", m.note != nil, m.pending)
	}

	m = send(t, notesModel(t), key("E"))
	m = typeText(t, m, "!")
	m = send(t, m, noteTickMsg{gen: m.note.gen})      // autosave starts
	m = send(t, m, noteSavedMsg{text: m.note.text()}) // and ends
	if got := m.taskByID(id).Description; got != old+"!" {
		t.Fatalf("description after autosave = %q", got)
	}
	m = send(t, m, key("ctrl+q"))
	if m.note != nil || m.pending != 1 {
		t.Fatalf("note open %v pending %d, want closed and one write", m.note != nil, m.pending)
	}
	if got := m.taskByID(id).Description; got != old {
		t.Errorf("description after discard = %q, want %q", got, old)
	}
}
