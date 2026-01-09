// Copyright 2025 The MathWorks, Inc.

package queryvmcblockhelp

import (
	"context"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/queryvmcblockhelp"
)

const (
	name        = "query_vmc_block_help"
	title       = "Query VMC Block Help"
	description = "Search and retrieve help documentation for specific Vitis Model Composer blocks. Returns detailed documentation including parameters, description, and WORKING EXAMPLE MODELS with COMPLETE MATLAB CREATION SCRIPTS.\n\nCRITICAL: When creating models with VMC blocks, ALWAYS extract and adapt the embedded example scripts (marked with 'MATLAB CREATION SCRIPT:') rather than creating implementations from scratch. These scripts contain:\n- Verified block library paths (e.g., 'aieDSP/FFT', 'aieUtilities/To Fixed Size')\n- Correct parameter configurations\n- Proven working subsystem structures\n- Proper signal sources and connections\n\nBlock library paths vary and guessing them will cause errors. Use the paths shown in the examples.\n\nPOST-CREATION VALIDATION CHECKLIST:\nAfter creating a model with VMC blocks, always validate with these steps:\n\nREQUIRED:\n1. Set discrete sample time on all source blocks (Constant, Random Source, etc.)\n   - AIE blocks require discrete sample times with offset = 0\n   - Example: set_param([modelName '/SourceBlock'], 'SampleTime', '1')\n2. Update model diagram to check for errors\n   - Command: set_param(modelName, 'SimulationCommand', 'update')\n   - Fix any errors before proceeding\n\nOPTIONAL VERIFICATION:\n- Set appropriate stop time for frame size (for N samples: StopTime = N-1 when SampleTime=1)\n- If comparing with reference, verify error signals are minimal\n- Check signal dimensions match expected sizes\n- Verify Hub block target device matches your hardware\n\nCOMMON ISSUES:\n- Missing sample time → 'AIEImportedIpBlock supports only discrete sample times' error\n- Wrong frame size → Check SSR and input vector size relationship\n- Connection errors → Clear all lines before reconnecting when modifying models"
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

		// Add usage reminder with validation checklist if the documentation contains example scripts
		documentation := result.Documentation
		if containsExampleScripts(documentation) {
			documentation = "⚠️  REMINDER: This block has working example models with complete MATLAB creation scripts.\n\nAFTER CREATING THE MODEL:\n1. Set discrete sample time (e.g., '1') on all source blocks - REQUIRED for AIE blocks\n2. Run: set_param(modelName, 'SimulationCommand', 'update') to validate\n3. Fix any errors before proceeding\n\nExtract scripts from 'MATLAB CREATION SCRIPT:' sections below. Use exact block library paths shown.\n\n---\n\n" + documentation
		}

		return ReturnArgs{
			Documentation: documentation,
		}, nil
	}
}

// containsExampleScripts checks if the documentation contains embedded example scripts
func containsExampleScripts(doc string) bool {
	return strings.Contains(doc, "MATLAB CREATION SCRIPT:")
}
