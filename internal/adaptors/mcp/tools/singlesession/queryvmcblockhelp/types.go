// Copyright 2025 The MathWorks, Inc.

package queryvmcblockhelp

type Args struct {
	BlockName string `json:"block_name" jsonschema:"The name of the Vitis Model Composer block to query. Can be a partial name. The search is case-insensitive and will find the best match."`
}

type ReturnArgs struct {
	Documentation string `json:"documentation" jsonschema:"The complete help documentation for the requested block including description, parameters, and usage examples."`
}
