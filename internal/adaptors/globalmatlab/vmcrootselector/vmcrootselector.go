// Copyright 2025 The MathWorks, Inc.

package vmcrootselector

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type Config interface {
	PreferredVMCRoot() string
}

type VMCRootSelector struct {
	config Config
}

func New(config Config) *VMCRootSelector {
	return &VMCRootSelector{
		config: config,
	}
}

func (v *VMCRootSelector) SelectVMCRoot(ctx context.Context, logger entities.Logger) string {
	return v.config.PreferredVMCRoot()
}
