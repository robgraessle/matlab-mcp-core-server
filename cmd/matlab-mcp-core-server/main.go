// Copyright 2025 The MathWorks, Inc.

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

func main() {
	// Debug: Log startup
	slog.Info("MATLAB MCP Core Server starting...")
	
	slog.Debug("Initializing mode selector...")
	modeSelector, err := wire.InitializeModeSelector()
	if err != nil {
		// As we failed to even initialize, we cannot use a LoggerFactory,
		// and we can't assume whatever failed had a logger factory to log the error either.
		// In this case, we use the default slog.
		slog.With("error", err).Error("Failed to initialize MATLAB MCP Core Server.")
		os.Exit(1)
	}
	slog.Debug("Mode selector initialized successfully")

	ctx := context.Background()
	slog.Debug("Starting mode selector and waiting for completion...")
	err = modeSelector.StartAndWaitForCompletion(ctx)
	if err != nil {
		slog.With("error", err).Error("Mode selector failed during execution.")
		os.Exit(1)
	}

	slog.Info("MATLAB MCP Core Server exiting normally")
	os.Exit(0)
}
