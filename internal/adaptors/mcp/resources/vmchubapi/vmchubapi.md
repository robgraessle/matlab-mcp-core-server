# Vitis Model Composer Hub API Reference

## Overview

When working with Vitis Model Composer Hub blocks programmatically, you must use the specialized `vmchub_get_param` and `vmchub_set_param` functions instead of standard MATLAB `get_param` and `set_param` functions. These functions provide the correct interface to access and modify Hub block parameters.

## Function Signatures

### vmchub_get_param
```matlab
value = vmchub_get_param(hubBlock, subsystem, parameterName)
```

**Parameters:**
- `hubBlock` - Handle or path to the Hub block (obtained using `xmcFindHubBlock`)
- `subsystem` - Path to the subsystem containing the parameter
- `parameterName` - Name of the parameter to retrieve (string)

**Returns:**
- The current value of the specified parameter

### vmchub_set_param
```matlab
vmchub_set_param(hubBlock, subsystem, parameterName, value)
```

**Parameters:**
- `hubBlock` - Handle or path to the Hub block (obtained using `xmcFindHubBlock`)
- `subsystem` - Path to the subsystem containing the parameter
- `parameterName` - Name of the parameter to set (string)
- `value` - New value for the parameter (string or numeric)

## Complete Working Example

```matlab
model = 'my_model';
open_system(model);

% Find the Hub block in the model
hubBlk = xmcFindHubBlock(model);

% Define the subsystem path
DUT_ss = [model, '/DUT'];

% Setting a numeric parameter (boolean values)
vmchub_set_param(hubBlk, DUT_ss, 'GenerateHwValidationCode', 1);
vmchub_set_param(hubBlk, DUT_ss, 'GenerateHwImage', 1);
vmchub_set_param(hubBlk, DUT_ss, 'CreateTestbench', 1);

% Setting a numeric parameter (integer values)
vmchub_set_param(hubBlk, DUT_ss, 'SelectSubsystem', 1);
vmchub_set_param(hubBlk, DUT_ss, 'HwCosimFifoDepth', 4096);
vmchub_set_param(hubBlk, DUT_ss, 'FPGAClockPeriod', 20);

% Setting string parameters
vmchub_set_param(hubBlk, DUT_ss, 'SelectHardware', 'xcvm1802-vfvc1760-1LHP-i-L');
vmchub_set_param(hubBlk, DUT_ss, 'CodeDirectory', './CodeD');
vmchub_set_param(hubBlk, DUT_ss, 'HwSystemType', 'Linux');
vmchub_set_param(hubBlk, DUT_ss, 'HwTarget', 'hw_emu');
vmchub_set_param(hubBlk, DUT_ss, 'CompilationType', 'IP Catalog');
vmchub_set_param(hubBlk, DUT_ss, 'SynthesisStrategy', 'Flow_RuntimeOptimized');
vmchub_set_param(hubBlk, DUT_ss, 'ImplementationStrategy', 'Performance_Explore');

% Reading parameters back
hardware = vmchub_get_param(hubBlk, DUT_ss, 'SelectHardware');
fifoDepth = vmchub_get_param(hubBlk, DUT_ss, 'HwCosimFifoDepth');
generateImage = vmchub_get_param(hubBlk, DUT_ss, 'GenerateHwImage');

% Verify values match
fprintf('Selected Hardware: %s\n', hardware);
fprintf('FIFO Depth: %d\n', fifoDepth);
fprintf('Generate HW Image: %d\n', generateImage);
```

## Available Hub Block Parameters

The following parameters can be accessed using `vmchub_get_param` and `vmchub_set_param`:

### Hardware Selection
- `SelectHardware` - Target FPGA device (e.g., 'xcvm1802-vfvc1760-1LHP-i-L')
- `SelectSubsystem` - Subsystem selection index

### Code Generation
- `CodeDirectory` - Directory for generated code
- `GenerateHwValidationCode` - Generate hardware validation code (0/1)
- `GenerateHwImage` - Generate hardware image (0/1)
- `HwSystemType` - Hardware system type (e.g., 'Linux', 'Windows')
- `HwTarget` - Hardware target (e.g., 'hw_emu', 'hw')
- `HWCommonSWDir` - Common software directory path
- `TargetSDKDir` - Target SDK directory path

### Compilation Settings
- `CompilationType` - Type of compilation (e.g., 'IP Catalog')
- `HwCosimBurstMode` - Hardware co-simulation burst mode (0/1)
- `HwCosimFifoDepth` - FIFO depth for hardware co-simulation

### IP Packaging
- `IPVendor` - IP vendor name
- `IPLibrary` - IP library name
- `IPName` - IP name
- `IPVersion` - IP version string
- `IPCategory` - IP category
- `IPStatus` - IP status (0/1)
- `IPAutoInferInterface` - Auto-infer interface (0/1)
- `IPUseCommonRepoDir` - Use common repository directory (0/1)
- `IPCommonRepoDir` - Common repository directory path
- `IPUsePlugInProject` - Use plug-in project (0/1)

### Hardware Description
- `HardwareDescription` - Hardware description language ('Verilog' or 'VHDL')
- `VHDLLibrary` - VHDL library name
- `UseSTDLogic` - Use STD_LOGIC (0/1)

### Synthesis and Implementation
- `SynthesisStrategy` - Synthesis strategy (e.g., 'Flow_RuntimeOptimized')
- `ImplementationStrategy` - Implementation strategy (e.g., 'Performance_Explore')

### Testing
- `CreateTestbench` - Create testbench (0/1)

### Clocking
- `EnableMultipleClocks` - Enable multiple clocks (0/1)
- `FPGAClockPeriod` - FPGA clock period in nanoseconds
- `SimulinkSystemPeriod` - Simulink system period
- `ClockPinLocation` - Clock pin location
- `ProvideClockEnableClearPin` - Provide clock enable/clear pin (0/1)

### Display and Analysis
- `BlockIconDisplay` - Block icon display mode
- `PerformAnalysis` - Perform analysis (0/1)
- `AnalyzerType` - Type of analyzer (e.g., 'Resource')

### Remote Caching
- `RemoteIPCache` - Remote IP cache (0/1)
- `CreateInterfaceDocument` - Create interface document (0/1)

## Important Notes

### Data Type Handling

When setting parameters with string values '0' or '1', convert them to numeric:
```matlab
if strcmp(value, '1') || strcmp(value, '0')
    vmchub_set_param(hubBlk, DUT_ss, paramName, str2num(value));
else
    vmchub_set_param(hubBlk, DUT_ss, paramName, value);
end
```

When comparing retrieved values:
```matlab
retrieved_value = vmchub_get_param(hubBlk, DUT_ss, paramName);
if strcmp(expected, '1') || strcmp(expected, '0')
    % Compare as numeric
    assert(retrieved_value == str2num(expected));
else
    % Compare as string
    assert(strcmp(retrieved_value, expected));
end
```

### Finding the Hub Block

Always use `xmcFindHubBlock` to locate the Hub block:
```matlab
hubBlk = xmcFindHubBlock(modelName);
```

### Subsystem Path

Construct the subsystem path by combining the model name and subsystem:
```matlab
DUT_ss = [modelName, '/DUT'];
```

### Performance Considerations

- Getting parameters is typically fast (< 2 seconds for 38 parameters)
- Setting parameters may take longer (< 50 seconds for 38 parameters)
- Consider batching parameter changes when possible

## Common Mistakes to Avoid

1. **Don't use standard `get_param`/`set_param`** - These will not work correctly with Hub block parameters
2. **Don't forget to convert string '0'/'1' to numeric** - Boolean parameters expect numeric values
3. **Don't modify parameters without finding the Hub block first** - Always call `xmcFindHubBlock`
4. **Don't use incorrect subsystem paths** - Ensure the path matches your model structure

## Full Reference Example

This complete example demonstrates proper usage patterns:

```matlab
model = 'vmchub_set_get_example';
open_system(model);

% Step 1: Find the Hub block
hubBlk = xmcFindHubBlock(model);
DUT_ss = [model, '/DUT'];

% Step 2: Define parameters and values
params = {
    'SelectHardware';
    'CodeDirectory';
    'GenerateHwValidationCode';
    'HwCosimFifoDepth';
    'IPVendor';
    'SynthesisStrategy';
};

values = {
    'xcvm1802-vfvc1760-1LHP-i-L';
    './CodeD';
    '1';
    '4096';
    'SAI';
    'Flow_RuntimeOptimized';
};

% Step 3: Set parameters with proper type handling
for i = 1:length(params)
    if strcmp(values{i}, '1') || strcmp(values{i}, '0')
        vmchub_set_param(hubBlk, DUT_ss, params{i}, str2num(values{i}));
    else
        vmchub_set_param(hubBlk, DUT_ss, params{i}, values{i});
    end
end

% Step 4: Read parameters back and verify
for i = 1:length(params)
    retrieved = vmchub_get_param(hubBlk, DUT_ss, params{i});
    
    if strcmp(values{i}, '1') || strcmp(values{i}, '0')
        assert(retrieved == str2num(values{i}), ...
            sprintf('Parameter %s mismatch', params{i}));
    else
        assert(strcmp(retrieved, values{i}), ...
            sprintf('Parameter %s mismatch', params{i}));
    end
    
    fprintf('✓ %s verified\n', params{i});
end

fprintf('All parameters set and verified successfully!\n');
```

---

**Source:** Based on reference implementation from AMD Vitis Model Composer test suite
