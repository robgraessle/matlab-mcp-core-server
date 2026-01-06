// Copyright 2025 The MathWorks, Inc.

package queryvmcblockhelp

import (
	"context"
	"fmt"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/vmcblockhelp"
)

type Result struct {
	Documentation string
}

type Usecase struct{}

func New() *Usecase {
	return &Usecase{}
}

func (u *Usecase) Execute(ctx context.Context, blockName string) (*Result, error) {
	// Use the vmcblockhelp resource's search functionality
	searchTerm := strings.ToLower(strings.TrimSpace(blockName))
	
	blockDocs, err := vmcblockhelp.SearchBlock(searchTerm)
	if err != nil {
		return nil, err
	}

	if len(blockDocs) == 0 {
		return nil, fmt.Errorf("no block found matching '%s'. Try a different search term or check the spelling", blockName)
	}

	// Return the first match (or exact match if found)
	return &Result{
		Documentation: blockDocs,
	}, nil
}
