# Vitis Model Composer MCP Core Server

Run Vitis Model Composer from AMD® using AI applications with this MCP Server. The Vitis Model Composer MCP Core Server allows your AI applications to:

- Start and quit Vitis Model Composer (which runs MATLAB®).
- Create Vitis Model Composer models by interacting with an AI agent.
- Analyze and debug Vitis Model Composer models.
  
## Table of Contents
  - [Building from Source](#building-from-source)
  - [Setup](#setup)
  - [Arguments](#arguments)
  - [Tools](#tools)
  - [Resources](#resources)

## Building from Source

To build the Vitis Model Composer MCP Core Server from source code on Linux:

1. **Install Go** (version 1.21 or later)
   - Download and install Go from [https://go.dev/doc/install](https://go.dev/doc/install)
   - Verify installation:
     ```sh
     go version
     ```

2. **Install build dependencies**
   ```sh
   make install
   ```
   This installs required tools: wire (dependency injection), mockery (mock generation), gotestsum (test runner), and golangci-lint (linter).

3. **Download Vitis Model Composer help documentation**
   ```sh
   make update-vmc-help
   ```
   This downloads the VMC block help documentation from the [VMC_Help GitHub repository](https://github.com/Xilinx/VMC_Help) and prepares it for embedding in the server binary.

4. **Build the server**
   ```sh
   make build-for-glnxa64
   ```

5. **Make the binary executable**
   ```sh
   chmod +x .bin/glnxa64/matlab-mcp-core-server
   ```

The built binary will be located at `.bin/glnxa64/matlab-mcp-core-server`.

## Setup for GitHub Copilot in Visual Studio Code

To add the Vitis Model Composer MCP Core Server to a workspace `mcp.json` file:

1. In VS Code, open your workspace folder.
2. Create or open the `.vscode/mcp.json` file in your workspace. If the `.vscode` folder doesn't exist, create it first.
3. Add the MCP server configuration to the `mcp.json` file:
   ```json
   {
       "servers": {
           "matlab": {
               "type": "stdio",
               "command": "/fullpath/to/matlab-mcp-core-server-binary",
               "args": [
                   "--vmc-root=/tools/Xilinx/2025.2/Model_Composer",
                   "--matlab-root=/usr/local/MATLAB/R2024b"
               ]
           }
       }
   }
   ```
   Remember to:
   - Replace `/fullpath/to/matlab-mcp-core-server-binary` with the actual path to the server binary you built.
   - Adjust the `--vmc-root` and `--matlab-root` paths to match your installation locations
   - Add any additional [Arguments](#arguments) as needed
   - On Windows, use double backslashes in paths (e.g., `"C:\\tools\\Xilinx\\2025.2\\Model_Composer"`)

4. Save the `mcp.json` file.
5. Start the MCP server in VS Code:
   - Reload VS Code: Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac), type "Developer: Reload Window", and press Enter.
   - When VS Code reloads, you will be prompted to trust the MCP server configuration. Review the configuration and click "Trust" to enable the server.
   - Verify the server is running: Open the GitHub Copilot Chat view and check that the Vitis Model Composer tools are available in the tools picker.

For more information about adding MCP servers in VS Code, see [Add an MCP Server (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers)

## Arguments

Customize the behavior of the server by providing arguments in the `args` array when configuring your AI application.

| Argument | Description | Example |
| ------------- | ------------- | ------------- |
| vmc-root | **Required for Vitis Model Composer.** Full path specifying which Vitis Model Composer installation to use. Do not include `/bin` in the path. When specified, the server launches Vitis Model Composer instead of MATLAB directly. | `"--vmc-root=/tools/Xilinx/2025.2/Model_Composer"` |
| matlab-root | Full path specifying which MATLAB to use. Do not include `/bin` in the path. Required when using `--vmc-root`. By default, the server tries to find the first MATLAB on the system PATH. | `"--matlab-root=/home/usr/MATLAB/R2025a"` |
| initialize-matlab-on-startup | To initialize Vitis Model Composer (or MATLAB) as soon as you start the server, set this argument to `true`. By default, it only starts when the first tool is called. | `"--initialize-matlab-on-startup=true"` |
| initial-working-folder | Specify the folder where MATLAB starts and where the server generates any MATLAB scripts. If you do not provide the argument, MATLAB starts in these locations: <br><br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | `"--initial-working-folder=C:\\Users\\name\\MyProject"` |

## Tools

1. `detect_matlab_toolboxes`
   - Lists installed MATLAB toolboxes with version information.
 
2. `check_matlab_code`
   - Performs static code analysis on a MATLAB script. Returns warnings about coding style, potential errors, deprecated functions, performance issues, and best practice violations. This is a non-destructive, read-only operation that helps identify code quality issues without executing the script.
   - Inputs:
     - `script_path` (string): Absolute path to the MATLAB script file to analyze. Must be a `.m` file within an allowed directory. The file is not modified during analysis. Example: `C:\Users\username\matlab\myFunction.m` or `/home/user/scripts/analysis.m`.
 
3. `evaluate_matlab_code`
   - Evaluates a string of MATLAB code and returns the output.
   - Inputs:
     - `code` (string): MATLAB code to evaluate.
     - `project_path` (string): Absolute path to an allowed project directory. MATLAB sets this directory as the current working folder. Example: `C:\Users\username\matlab-project` or `/home/user/research`.
 
4. `run_matlab_file`
   - Executes a MATLAB script and returns the output. The script must be a valid `.m file`.
   - Inputs:
     - `script_path` (string): Absolute path to the MATLAB script file to execute. Must be a valid `.m` file within an allowed directory. Example: `C:\Users\username\projects\analysis.m` or `/home/user/matlab/simulation.m`.
 
5. `run_matlab_test_file`
   - Executes a MATLAB test script and returns comprehensive test results. Designed specifically for MATLAB unit test files that follow MATLAB testing framework conventions.
   - Inputs:
     - `script_path` (string): Absolute path to the MATLAB test script file. Must be a valid `.m` file containing MATLAB unit tests, within an allowed directory. Example: `C:\Users\username\tests\testMyFunction.m` or `/home/user/matlab/tests/test_analysis.m`.

6. `query_vmc_block_help`
   - Search and retrieve help documentation for specific Vitis Model Composer blocks. Returns detailed documentation including parameters, description, and usage examples for the requested block. This tool searches through all available block documentation and returns the best match.
   - Inputs:
     - `block_name` (string): The name of the Vitis Model Composer block to query. Can be a partial name (e.g., 'Abs', 'FFT', 'FIR'). The search is case-insensitive and will find the best match.
   - Example usage: "Query help for the HLS Abs block" or "What are the parameters for the FFT block?"

## Resources
The MCP server provides [Resources (MCP)](https://modelcontextprotocol.io/specification/2025-03-26/server/resources) to help your AI application write better code and understand Vitis Model Composer blocks. To see instructions for using these resources, refer to the documentation of your AI application that explains how to use resources. 

1. `matlab_coding_guidelines`
   - Provides comprehensive MATLAB coding standards for improving code readability, maintainability, and collaboration. The guidelines encompass naming conventions, formatting, commenting, performance optimization, and error handling.
   - URI: `guidelines://coding`
   - MIME Type: `text/markdown`
   - Source: [MATLAB Coding Standards (GitHub)](https://github.com/matlab/rules/blob/main/matlab-coding-standards.md)

2. `vmcblockhelp`
   - Provides comprehensive help documentation for Vitis Model Composer blocks from the VMC_Help GitHub repository. Includes block descriptions, parameters, usage examples, and best practices.
   - URI: `vmc-help://blocks`
   - MIME Type: `text/markdown`
   - Source: [VMC_Help (GitHub)](https://github.com/Xilinx/VMC_Help)

# 
When using the Vitis Model Composer MCP Core Server, you should thoroughly review and validate all tool calls before you run them. Always keep a human in the loop for important actions and only proceed once you are confident the call will do exactly what you expect. For more information, see [User Interaction Model (MCP)](https://modelcontextprotocol.io/specification/2025-06-18/server/tools#user-interaction-model) and [Security Considerations (MCP)](https://modelcontextprotocol.io/specification/2025-06-18/server/tools#security-considerations).

---

Copyright 2026 Advanced Micro Devices, Inc.

Based on [MATLAB MCP Core Server](https://github.com/matlab/matlab-mcp-core-server): Copyright 2025 The MathWorks, Inc.

----
