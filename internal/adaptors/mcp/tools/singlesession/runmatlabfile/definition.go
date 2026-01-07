// Copyright 2025 The MathWorks, Inc.

package runmatlabfile

const (
	name        = "run_matlab_file"
	title       = "Run MATLAB File"
	description = "Execute a MATLAB script file (`script_path`) in an existing MATLAB session and capture its command window output. The script runs with the working directory automatically set to the script's location. The script must exist and be a valid .m file. Returns the command window output or a success message if no output is generated. Note: The Vitis Model Composer Hub block requires specialized APIs instead of standard get_param/set_param. Check available resources before using standard MATLAB functions on the Vitis Model Composer Hub block."
)

type Args struct {
	ScriptPath string `json:"script_path" jsonschema:"The full absolute path to the MATLAB script file to execute - Must be a .m file that exists - Example: C:\\Users\\username\\projects\\analysis.m or /home/user/matlab/simulation.m."`
}
