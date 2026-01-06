// Copyright 2025 The MathWorks, Inc.

package queryvmcblockhelp

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/queryvmcblockhelp"
)

const (
	name        = "query_vmc_block_help"
	title       = "Query VMC Block Help"
	description = "Search and retrieve help documentation for specific Vitis Model Composer blocks. Returns detailed documentation including parameters, description, and usage for the requested block."
)

type Usecase interface {
	Execute(ctx context.Context, blockName string) (*queryvmcblockhelp.Result, error)
}

type Tool struct {
	basetool.ToolWithStructuredContentOutput[Args, ReturnArgs]
}

func New(
	loggerFactory basetool.LoggerFactory,
	usecase Usecase,
) *Tool {
	return &Tool{
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredContent(name, title, description, loggerFactory, Handler(usecase)),
	}
}

func (Tool) Name() string {
	return name
}

func (Tool) Description() string {
	return description
}

func Handler(usecase Usecase) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		sessionLogger.Info("Querying VMC block help")
		defer sessionLogger.Info("Done - Querying VMC block help")

		result, err := usecase.Execute(ctx, inputs.BlockName)
		if err != nil {
			return ReturnArgs{}, err
		}

		return ReturnArgs{
			Documentation: result.Documentation,
		}, nil
	}
}
