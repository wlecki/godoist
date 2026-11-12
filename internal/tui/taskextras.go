package tui

import (
	"context"
	"fmt"
	"net/url"

	tea "charm.land/bubbletea/v2"
)

// Small task keys: O opens the task in Todoist on the web, Y copies its web address,
// ctrl+a selects all tasks of the view, and D makes a copy of the task.

// taskURL is the address of a task in the Todoist web app.
func taskURL(id string) string { return "https://app.todoist.com/app/task/" + url.PathEscape(id) }

// openInTodoist opens the task under the cursor in the Todoist web app (O).
func (m Model) openInTodoist() (tea.Model, tea.Cmd) {
	t := m.currentTask()
	if t == nil {
		return m, nil
	}
	if err := openURL(taskURL(t.ID)); err != nil {
		m.setStatus("could not open Todoist: "+err.Error(), true)
		return m, nil
	}
	m.setStatus("opened “"+plain(t.Content)+"” in Todoist", false)
	return m, nil
}

// yankURL copies the web address of the task under the cursor (Y).
func (m Model) yankURL() (tea.Model, tea.Cmd) {
	t := m.currentTask()
	if t == nil {
		return m, nil
	}
	m.setStatus("copied the Todoist link", false)
	return m, tea.SetClipboard(taskURL(t.ID))
}

// selectAll selects all open tasks of the task list (ctrl+a). A second ctrl+a, with all
// of them selected, clears the selection.
func (m *Model) selectAll() {
	var ids []string
	for _, r := range m.rows {
		if r.task != nil && !r.done && !r.doneSub {
			ids = append(ids, r.task.ID)
		}
	}
	all := len(ids) > 0
	for _, id := range ids {
		all = all && m.selected[id]
	}
	if all {
		m.clearSelection()
		m.setStatus("selection cleared", false)
		return
	}
	if m.selected == nil {
		m.selected = map[string]bool{}
	}
	for _, id := range ids {
		m.selected[id] = true
	}
	m.setStatus(fmt.Sprintf("%d tasks selected", len(ids)), false)
}

// duplicateTask adds a copy of the task under the cursor (D): the same name, description,
// project, section, parent, priority, labels, and due date. Sub-tasks and comments do not
// come along.
func (m *Model) duplicateTask() tea.Cmd {
	t := m.currentTask()
	if t == nil {
		return nil
	}
	src, client := *t, m.client
	fields := map[string]any{"priority": src.Priority}
	if src.Description != "" {
		fields["description"] = src.Description
	}
	if len(src.Labels) > 0 {
		fields["labels"] = src.Labels
	}
	switch {
	case src.Due != nil && src.Due.String != "":
		fields["due_string"] = src.Due.String // keeps a repeat
	case src.Due != nil:
		fields["due_date"] = src.Due.Date
	}
	return m.write("", func(ctx context.Context) (string, error) {
		nt, err := client.AddTask(ctx, src.Content, src.ProjectID, src.Section())
		if err != nil {
			return "", err
		}
		if _, err := client.UpdateTask(ctx, nt.ID, fields); err != nil {
			return "", fmt.Errorf("copied, but not all fields: %w", err)
		}
		if p := src.Parent(); p != "" {
			if err := client.MoveToParent(ctx, nt.ID, p); err != nil {
				return "", fmt.Errorf("copied, but not under its parent task: %w", err)
			}
		}
		return "Copied “" + plain(src.Content) + "”", nil
	})
}
