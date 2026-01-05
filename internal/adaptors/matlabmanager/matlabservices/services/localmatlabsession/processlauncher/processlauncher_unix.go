// Copyright 2025 The MathWorks, Inc.
//go:build !windows

package processlauncher

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"golang.org/x/sys/unix"
)

func startMatlab(_ entities.Logger, matlabRoot string, vmcRoot string, workingDir string, args []string, env []string, stdIO *stdIO) (*os.Process, error) {
	// If vmcRoot is specified, launch model_composer instead of MATLAB
	var executablePath string
	var processArgs []string
	
	if vmcRoot != "" {
		// Launch model_composer from VMC root
		executablePath = filepath.Join(vmcRoot, "bin", "model_composer")
		if _, err := os.Stat(executablePath); err != nil {
			return nil, fmt.Errorf("model_composer executable not found at %s: %w", executablePath, err)
		}
		
		// model_composer requires -matlab <matlab-root> followed by MATLAB arguments
		// Format: model_composer -matlab <matlab-root> <matlab-args>
		processArgs = []string{executablePath, "-matlab", matlabRoot}
		processArgs = append(processArgs, args...)
	} else {
		// Launch MATLAB directly
		executablePath = filepath.Join(matlabRoot, "bin", "matlab")
		if _, err := os.Stat(executablePath); err != nil {
			return nil, err
		}
		
		// Careful here, for start process, we need the path first. From the doc:
		//   > StartProcess starts a new process with the program, arguments and attributes specified by name, argv and attr.
		//   > The argv slice will become os.Args in the new process, so it normally starts with the program name.
		processArgs = append([]string{executablePath}, args...)
	}

	attr := &os.ProcAttr{
		Dir:   workingDir,
		Env:   env,
		Files: []*os.File{stdIO.stdIn, stdIO.stdOut, stdIO.stdErr},
		Sys: &unix.SysProcAttr{
			Setsid: true, // Create a new session
		},
	}

	process, err := os.StartProcess(executablePath, processArgs, attr)
	if err != nil {
		if vmcRoot != "" {
			return nil, fmt.Errorf("error starting Model Composer: %w", err)
		}
		return nil, fmt.Errorf("error starting MATLAB: %w", err)
	}

	return process, nil
}
