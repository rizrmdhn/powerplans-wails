package main

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// powercfgTimeout bounds how long a single powercfg invocation may run before
// it is cancelled, so a hung process can't freeze the app indefinitely.
const powercfgTimeout = 15 * time.Second

// App struct holds the runtime context Wails injects on startup.
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// PowerPlan mirrors one row of `powercfg /list` output.
type PowerPlan struct {
	GUID   string `json:"guid"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Hidden bool   `json:"hidden"` // true = a known template not currently duplicated onto the system
}

// PlanSettings contains the timeout values, in seconds, stored on a power
// scheme for both mains power (AC) and battery power (DC). A value of 0 means
// "Never", which is how powercfg represents an unlimited timeout.
type PlanSettings struct {
	DisplayAC int `json:"displayAC"`
	DisplayDC int `json:"displayDC"`
	DiskAC    int `json:"diskAC"`
	DiskDC    int `json:"diskDC"`
	SleepAC   int `json:"sleepAC"`
	SleepDC   int `json:"sleepDC"`
}

// knownTemplates are the well-known GUIDs for plans Windows ships but often
// hides from `powercfg /list` unless duplicated onto the system first.
var knownTemplates = []PowerPlan{
	{GUID: "381b4222-f694-41f0-9685-ff5bb260df2e", Name: "Balanced"},
	{GUID: "8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c", Name: "High performance"},
	{GUID: "a1841308-3541-4fab-bc81-f71556f20b4a", Name: "Power saver"},
	{GUID: "e9a42b02-d5df-448d-aa00-03f14749eb61", Name: "Ultimate Performance"},
}

var listLineRe = regexp.MustCompile(`(?i)Power Scheme GUID:\s*([0-9a-fA-F-]{36})\s*\(([^)]*)\)\s*(\*)?`)
var guidRe = regexp.MustCompile(`(?i)[0-9a-fA-F]{8}-(?:[0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}`)
var guidExactRe = regexp.MustCompile(`(?i)^[0-9a-fA-F]{8}-(?:[0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`)
var currentIndexRe = regexp.MustCompile(`(?im)Current (?:AC|DC) Power Setting Index:\s*0x([0-9a-f]+)`)

// notifyPlansChanged keeps the tray menu and any open window in sync after a
// plan-mutating method succeeds, regardless of whether the change came from
// the tray or the window itself.
func (a *App) notifyPlansChanged() {
	refreshTray()
	runtime.EventsEmit(a.ctx, "plans:changed")
}

// validateGUID rejects anything that isn't a well-formed GUID before it is
// passed to powercfg, since powercfg is invoked with the value verbatim.
func validateGUID(guid string) error {
	if !guidExactRe.MatchString(strings.TrimSpace(guid)) {
		return fmt.Errorf("invalid power plan identifier")
	}
	return nil
}

func (a *App) runPowercfg(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(a.ctx, powercfgTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powercfg", args...)
	// powercfg is a console application. Hide its console window when launched
	// from the Wails GUI so operations do not flash a Command Prompt window.
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("powercfg timed out after %s", powercfgTimeout)
		}
		return "", fmt.Errorf("%s: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

// GetPlans returns every plan currently registered on the system, flagging
// the active one, plus any well-known template that isn't installed yet.
func (a *App) GetPlans() ([]PowerPlan, error) {
	out, err := a.runPowercfg("/list")
	if err != nil {
		return nil, err
	}

	matches := listLineRe.FindAllStringSubmatch(out, -1)
	plans := make([]PowerPlan, 0, len(matches))
	seen := map[string]bool{}

	for _, m := range matches {
		guid := strings.ToLower(strings.TrimSpace(m[1]))
		plans = append(plans, PowerPlan{
			GUID:   guid,
			Name:   strings.TrimSpace(m[2]),
			Active: m[3] == "*",
		})
		seen[guid] = true
	}

	// Append known templates that aren't installed, marked Hidden so the
	// UI can offer a "Restore" action instead of Activate/Delete.
	for _, t := range knownTemplates {
		if !seen[t.GUID] {
			t.Hidden = true
			plans = append(plans, t)
		}
	}

	return plans, nil
}

// SetActive switches the active power plan.
func (a *App) SetActive(guid string) error {
	if err := validateGUID(guid); err != nil {
		return err
	}
	_, err := a.runPowercfg("/setactive", guid)
	if err != nil {
		return err
	}
	a.notifyPlansChanged()
	return nil
}

// DeletePlan removes a plan. Windows itself refuses to delete the active
// plan, so we surface that as a clear error rather than letting it fail silently.
func (a *App) DeletePlan(guid string) error {
	if err := validateGUID(guid); err != nil {
		return err
	}
	_, err := a.runPowercfg("/delete", guid)
	if err != nil {
		return fmt.Errorf("couldn't delete plan (it may be the active plan — switch to another plan first): %w", err)
	}
	a.notifyPlansChanged()
	return nil
}

// RestorePlan brings a hidden Windows-default template back via /duplicatescheme.
func (a *App) RestorePlan(guid string) error {
	if err := validateGUID(guid); err != nil {
		return err
	}
	_, err := a.runPowercfg("/duplicatescheme", guid)
	if err != nil {
		return err
	}
	a.notifyPlansChanged()
	return nil
}

// RenamePlan renames a plan in place (name only, keeps existing settings).
func (a *App) RenamePlan(guid string, newName string) error {
	if err := validateGUID(guid); err != nil {
		return err
	}
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return fmt.Errorf("plan name cannot be empty")
	}
	_, err := a.runPowercfg("/changename", guid, newName)
	if err != nil {
		return err
	}
	a.notifyPlansChanged()
	return nil
}

// CreatePlan duplicates an existing scheme and gives the copy a user-supplied name.
// Passing an empty source GUID lets Windows duplicate the active scheme.
func (a *App) CreatePlan(name, sourceGUID string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("plan name cannot be empty")
	}
	sourceGUID = strings.TrimSpace(sourceGUID)
	if sourceGUID != "" {
		if err := validateGUID(sourceGUID); err != nil {
			return err
		}
	}

	args := []string{"/duplicatescheme"}
	if sourceGUID != "" {
		args = append(args, sourceGUID)
	}
	out, err := a.runPowercfg(args...)
	if err != nil {
		return err
	}

	guid := guidRe.FindString(out)
	if guid == "" {
		return fmt.Errorf("Windows created the plan but did not return its identifier")
	}
	if _, err := a.runPowercfg("/changename", guid, name); err != nil {
		return err
	}
	a.notifyPlansChanged()
	return nil
}

var timeoutSettings = []struct {
	subgroup string
	setting  string
}{
	{"SUB_VIDEO", "VIDEOIDLE"},
	{"SUB_DISK", "DISKIDLE"},
	{"SUB_SLEEP", "STANDBYIDLE"},
}

func (a *App) readTimeout(guid, subgroup, setting string) (int, int, error) {
	out, err := a.runPowercfg("/query", guid, subgroup, setting)
	if err != nil {
		return 0, 0, err
	}
	matches := currentIndexRe.FindAllStringSubmatch(out, -1)
	if len(matches) != 2 {
		return 0, 0, fmt.Errorf("could not read timeout values from powercfg")
	}
	var ac, dc int
	if _, err := fmt.Sscanf(matches[0][1], "%x", &ac); err != nil {
		return 0, 0, err
	}
	if _, err := fmt.Sscanf(matches[1][1], "%x", &dc); err != nil {
		return 0, 0, err
	}
	return ac, dc, nil
}

// GetPlanSettings reads the most useful timeouts from the selected scheme.
func (a *App) GetPlanSettings(guid string) (PlanSettings, error) {
	if err := validateGUID(guid); err != nil {
		return PlanSettings{}, err
	}
	var settings PlanSettings
	values := []*int{&settings.DisplayAC, &settings.DisplayDC, &settings.DiskAC, &settings.DiskDC, &settings.SleepAC, &settings.SleepDC}
	for i, item := range timeoutSettings {
		ac, dc, err := a.readTimeout(guid, item.subgroup, item.setting)
		if err != nil {
			return PlanSettings{}, err
		}
		*values[i*2] = ac
		*values[i*2+1] = dc
	}
	return settings, nil
}

// UpdatePlanSettings writes timeout values to the selected scheme without
// switching the user's currently active plan.
func (a *App) UpdatePlanSettings(guid string, settings PlanSettings) error {
	if err := validateGUID(guid); err != nil {
		return err
	}
	values := [][2]int{
		{settings.DisplayAC, settings.DisplayDC},
		{settings.DiskAC, settings.DiskDC},
		{settings.SleepAC, settings.SleepDC},
	}
	for i, item := range timeoutSettings {
		for mode, value := range values[i] {
			if value < 0 || value > 604800 {
				return fmt.Errorf("timeout must be between 0 and 604800 seconds")
			}
			command := "/setacvalueindex"
			if mode == 1 {
				command = "/setdcvalueindex"
			}
			if _, err := a.runPowercfg(command, guid, item.subgroup, item.setting, fmt.Sprint(value)); err != nil {
				return err
			}
		}
	}
	return nil
}
