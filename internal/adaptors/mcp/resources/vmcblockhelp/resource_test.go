// Copyright 2025 The MathWorks, Inc.

package vmcblockhelp_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/vmcblockhelp"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	baseresourcemocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/resources/baseresource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := baseresourcemocks.NewMockLoggerFactory(t)

	// Act
	resource, err := vmcblockhelp.New(mockLoggerFactory)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resource)
	assert.Equal(t, "vmcblockhelp", resource.Name())
	assert.Equal(t, "Vitis Model Composer Block Help", resource.Title())
	assert.Equal(t, "text/markdown", resource.MimeType())
	assert.Equal(t, "vmc-help://blocks", resource.URI())
}

func TestHandler_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	handler := vmcblockhelp.Handler()

	// Act
	result, err := handler(t.Context(), mockLogger)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)
	assert.Equal(t, "text/markdown", result.Contents[0].MIMEType)
	assert.NotEmpty(t, result.Contents[0].Text)
	assert.Contains(t, result.Contents[0].Text, "Vitis Model Composer Block Help")
}
