# Power Plan Switcher (Wails + React + Tailwind)

A native Windows utility to view, activate, create, configure, rename,
restore, and delete power plans — a GUI replacement for the earlier `.bat`
version.

## Why these files can't be run in this sandbox

Wails builds a native Windows executable using `cgo` + WebView2, and
`powercfg` only exists on Windows. This project must be built on your own
Windows machine (or cross-compiled from one). Everything below is ready to
drop in — you just need to run it locally.

## Prerequisites (on your Windows machine)

1. Go 1.21+: https://go.dev/dl/
2. Node.js 18+: https://nodejs.org/
3. Wails CLI:
   ```
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
4. Confirm your setup:
   ```
   wails doctor
   ```

## Setting up this project

The cleanest path is to let Wails scaffold a fresh project (so it generates
`go.mod`, `frontend/wailsjs/` bindings, and build config for the right Wails
version), then overwrite it with these files:

```
wails init -n powerplan-wails -t react-ts
cd powerplan-wails
```

Now copy these files into place, overwriting the generated ones:
- `app.go` → replaces the generated one (adds plan management and timeout-setting methods)
- `main.go` → replaces the generated one
- `wails.json` → replaces the generated one
- `frontend/src/App.tsx`, `frontend/src/main.tsx`, `frontend/src/index.css`
- `frontend/src/lib/utils.ts`
- `frontend/src/components/ui/button.tsx`, `badge.tsx`, `confirm-dialog.tsx`
- `frontend/tailwind.config.js`, `frontend/postcss.config.js`, `frontend/vite.config.ts`, `frontend/tsconfig.json`, `frontend/index.html`
- Merge the `dependencies`/`devDependencies` from the provided `frontend/package.json` into the generated one (or replace it — the react-ts template's other scripts are identical).

Then install frontend deps:
```
cd frontend
npm install
npm install -D tailwindcss postcss autoprefixer
cd ..
```

## Run it (dev mode, hot reload)

```
wails dev
```

This opens the app in a dev window and regenerates
`frontend/wailsjs/go/main/App.d.ts` — the typed bindings that `App.tsx`
imports from (`../wailsjs/go/main/App`). You don't write these by hand;
Wails generates them automatically from your `app.go` methods every time
you run `wails dev` or `wails build`.

## Build the final .exe

```
wails build
```

Output lands in `build/bin/PowerPlanSwitcher.exe` — a single portable
executable, no install step for whoever runs it.

## GitHub Actions

The included workflow builds the Windows executable for every pull request
targeting `main` and uploads it as a workflow artifact. Once a PR is merged
and its commit reaches `main`, the same build is published automatically as a
GitHub release with a `build-<run number>` tag. This keeps unmerged PR code
from being released.

## Notes

- **Admin rights**: `powercfg /setactive`, `/delete`, `/changename` work
  without elevation. `/duplicatescheme` (the "Restore" button) also works
  as a standard user in modern Windows.
- **Create and configure**: “New plan” duplicates a selected plan and names
  the copy. “Configure” changes display, disk, and sleep timeouts for AC and
  battery power. Choose milliseconds, seconds, or minutes; Windows stores
  whole seconds, so millisecond values are rounded. Enter `0` for “Never”;
  timeouts are stored per plan without switching your active plan.
- **Real shadcn components**: the `Button`/`Badge`/`ConfirmDialog` here are
  hand-written to match shadcn's API and avoid extra Radix dependencies. If
  you want the real thing later: `npx shadcn@latest init` then
  `npx shadcn@latest add button badge alert-dialog`, and swap the imports.
- **Extending it**: `app.go` is the entire backend surface — add a method
  there (e.g. `GetBatterySaverThreshold`), and it's instantly callable from
  React after the next `wails dev` regenerates bindings.
