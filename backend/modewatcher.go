//go:build windows

package backend

import (
	"context"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ModeWatcherState is pushed to the frontend via Wails events when mode/service status changes.
type ModeWatcherState struct {
	ServiceStatus string `json:"serviceStatus"`
	CurrentMode   string `json:"currentMode"`
	AutoMode      int    `json:"autoMode"`
	ITSCurrentMode string `json:"itsCurrentMode"`
	ITSTargetMode  string `json:"itsTargetMode"`
}

// StartModeWatcher polls service status + dispatcher mode every 200ms in the background.
// It ONLY emits a Wails event when something actually changed, so the frontend
// receives push-based updates with zero polling overhead.
//
// This replaces the frontend setInterval(500ms) approach:
//   - Frontend no longer calls GetServiceAndModeInfo() on a timer
//   - Latency drops to ~200ms (poll interval) + ~1ms (event emit) instead of
//     500ms + 5-15ms (IPC round-trip)
//   - No wasted IPC calls when nothing changed
func StartModeWatcher(ctx context.Context) {
	const pollInterval = 100 * time.Millisecond

	var lastState ModeWatcherState
	firstRun := true

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[ModeWatcher] context cancelled, stopping")
			return
		case <-ticker.C:
		}

		// Gather state — this is pure registry + cached SCM, ~1-3ms
		state := gatherModeState()

		// Only emit if something changed (or first run)
		if firstRun || state != lastState {
			runtime.EventsEmit(ctx, "mode:state-change", state)
			lastState = state
			firstRun = false
		}
	}
}

// gatherModeState reads service status + dispatcher info in one shot.
// Uses cached SCM handle + persistent registry key, so it's very fast (~1-3ms).
func gatherModeState() ModeWatcherState {
	var state ModeWatcherState

	// Service status via cached SCM handle
	scmMu.Lock()
	if scmMgr != nil && scmService != nil {
		status, err := scmService.Query()
		if err == nil {
			state.ServiceStatus = serviceStateToString(status.State)
		} else {
			// Stale handle
			scmService.Close()
			scmService = nil
			state.ServiceStatus = "Error"
		}
	} else {
		state.ServiceStatus = "Not Installed"
	}
	scmMu.Unlock()

	// If SCM handle was stale, try to reconnect (best-effort, non-blocking)
	if state.ServiceStatus == "Error" || state.ServiceStatus == "Not Installed" {
		_, svc, err := getSCMService()
		if err == nil {
			status, err := svc.Query()
			if err == nil {
				state.ServiceStatus = serviceStateToString(status.State)
			}
		}
	}

	// Dispatcher info via cached registry key
	info, _ := GetDispatcherInfo()
	state.CurrentMode = info.CurrentMode
	state.AutoMode = info.AutoMode
	state.ITSCurrentMode = info.ITSCurrentMode
	state.ITSTargetMode = info.ITSTargetMode

	return state
}
