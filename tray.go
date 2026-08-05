package main

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var trayIcon []byte

// trayPlanSlots bounds how many power plans the tray menu can display at
// once. getlantern/systray has no API to remove a menu item once added, so
// slots are pre-created and hidden/shown/relabelled instead of recreated.
const trayPlanSlots = 32

type trayManager struct {
	app *App

	mu        sync.Mutex
	planItems [trayPlanSlots]*systray.MenuItem
	planGUIDs [trayPlanSlots]string
}

var tray *trayManager

// onTrayReady builds the tray icon and menu. It runs on the goroutine
// started by systray.Run from OnStartup.
func (a *App) onTrayReady() {
	systray.SetIcon(trayIcon)
	systray.SetTitle("Power Plan Switcher")
	systray.SetTooltip("Power Plan Switcher")

	tm := &trayManager{app: a}
	tray = tm

	showItem := systray.AddMenuItem("Show window", "Open Power Plan Switcher")
	systray.AddSeparator()

	for i := range tm.planItems {
		item := systray.AddMenuItemCheckbox("", "", false)
		item.Hide()
		tm.planItems[i] = item

		idx := i
		go func() {
			for range item.ClickedCh {
				tm.mu.Lock()
				guid := tm.planGUIDs[idx]
				tm.mu.Unlock()
				if guid == "" {
					continue
				}
				if err := a.SetActive(guid); err == nil {
					runtime.EventsEmit(a.ctx, "plans:changed")
				}
			}
		}()
	}

	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Exit Power Plan Switcher")

	go func() {
		for {
			select {
			case <-showItem.ClickedCh:
				runtime.WindowShow(a.ctx)
				runtime.WindowUnminimise(a.ctx)
			case <-quitItem.ClickedCh:
				runtime.Quit(a.ctx)
				return
			}
		}
	}()

	tm.refresh()
}

func (a *App) onTrayExit() {}

// refresh re-reads the installed power plans and updates the tray's
// checkable plan list so it always mirrors what powercfg reports.
func (tm *trayManager) refresh() {
	plans, err := tm.app.GetPlans()
	if err != nil {
		return
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	slot := 0
	for _, p := range plans {
		if p.Hidden {
			continue // not installed, so it can't be switched to directly
		}
		if slot >= trayPlanSlots {
			break
		}
		item := tm.planItems[slot]
		item.SetTitle(p.Name)
		item.SetTooltip(fmt.Sprintf("Switch to %s", p.Name))
		if p.Active {
			item.Check()
		} else {
			item.Uncheck()
		}
		item.Show()
		tm.planGUIDs[slot] = p.GUID
		slot++
	}
	for ; slot < trayPlanSlots; slot++ {
		tm.planItems[slot].Hide()
		tm.planGUIDs[slot] = ""
	}
}

// refreshTray updates the tray menu after a plan-mutating App method
// succeeds, whether it was triggered from the tray or the app window.
func refreshTray() {
	if tray != nil {
		tray.refresh()
	}
}
