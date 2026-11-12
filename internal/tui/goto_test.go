package tui

import "testing"

// gg goes to the top, gt opens Today, and a key that is not a jump key works as usual.
func TestGotoPrefix(t *testing.T) {
	m := onRow(t, 2)
	m = send(t, m, key("g"))
	if !m.gPrefix || m.status != gotoHint {
		t.Fatalf("after g: prefix %v status %q", m.gPrefix, m.status)
	}
	m = send(t, m, key("g"))
	if m.gPrefix || m.currentTask() == nil || m.rowCur > 1 {
		t.Errorf("after gg: prefix %v rowCur %d, want the top", m.gPrefix, m.rowCur)
	}
	m = send(t, m, key("g"), key("t"))
	if cur := m.currentNav(); cur == nil || cur.kind != vkToday || m.focus != paneTasks {
		t.Errorf("after gt: nav %+v focus %d, want Today and the task list", cur, m.focus)
	}
	m = send(t, onRow(t, 1), key("g"), key("j"))
	if m.gPrefix || m.rowCur != 2 || m.status == gotoHint {
		t.Errorf("after g j: prefix %v rowCur %d status %q, want j to move down", m.gPrefix, m.rowCur, m.status)
	}
}

// ctrl+k opens the go-to picker. Typing filters it, and enter opens the item.
func TestGotoPicker(t *testing.T) {
	m := send(t, testModel(t), key("ctrl+k"))
	if m.pick == nil || m.pick.kind != pickGoto {
		t.Fatal("ctrl+k did not open the go-to picker")
	}
	m = send(t, m, key("w"), key("o"), key("r"), key("k"), key("enter"))
	if cur := m.currentNav(); m.pick != nil || cur == nil || cur.projectID != "p1" || m.focus != paneTasks {
		t.Errorf("picker open %v nav %+v focus %d, want the work project", m.pick != nil, cur, m.focus)
	}
}
