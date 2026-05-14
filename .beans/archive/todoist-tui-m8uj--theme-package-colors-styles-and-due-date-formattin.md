---
# todoist-tui-m8uj
title: 'Theme package: colors, styles, and due date formatting'
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:34:56Z
updated_at: 2026-05-08T17:53:33Z
parent: todoist-tui-0rtm
blocked_by:
    - todoist-tui-m01d
---

## Description

Create the `internal/ui/theme` package with three files: color mapping, reusable lipgloss styles, and smart due date formatting. Used by all M1 panels.

## Requirements

### `internal/ui/theme/colors.go`
- `TodoistColor(name string, fallback string) lipgloss.TerminalColor` — maps all ~30 Todoist color names to hex values
- Color names: `berry_red`, `red`, `orange`, `lime_green`, `green`, `teal`, `blue`, `sky_blue`, `purple`, `pink`, `grey`, `magenta`, `peach`, `yellow`, `ivory`
- Unknown names fall back to `fallback` (typically `config.ThemeConfig.NormalText`)

### `internal/ui/theme/styles.go`
- `Styles` struct holds `*config.ThemeConfig` reference
- `NewStyles(cfg *config.Config) *Styles` constructor
- Methods returning `lipgloss.Style`: `Header()`, `MutedText()`, `SectionSep()`, `PriorityDot()`, `ActiveItem()`, `InactiveItem()`, `ColorDot()`
- Each method reads from the theme config for color values

### `internal/ui/theme/format.go`
- `FormatDueDate(due *model.DueDate, now time.Time) (string, lipgloss.TerminalColor)` — smart date formatting
- Returns relative for near dates: "Today", "Tomorrow", "In 3 days", "Overdue 2d"
- Returns abbreviated weekday+date for this week: "Mon 10"
- Returns full date for further out: "May 10"
- Returns color matching the theme: `TaskOverdue` for overdue, `TaskDueToday` for today, `NormalText` otherwise

## Acceptance Criteria

- All three files exist and compile
- `TodoistColor()` returns correct hex for known names, fallback for unknown
- `FormatDueDate()` handles nil due date (returns empty string)
- `go build ./...` succeeds

## Summary of Changes\n\nCreated  package with three files:\n\n- **colors.go** —  maps 15 Todoist color names to hex values, returns fallback for unknown names\n- **styles.go** —  struct backed by  with 7 style methods: , , , , , , . Helper  converts theme config strings to .\n- **format.go** —  method on  for smart due date formatting: nil → ("", nil), today → "Today" (TaskDueToday), tomorrow → "Tomorrow" (TaskDueToday), 2-7 days → "In N days" (NormalText), overdue → "Overdue Nd" (TaskOverdue), far future → "Jan 2" (NormalText)\n\nAll acceptance criteria met.  and  pass.
