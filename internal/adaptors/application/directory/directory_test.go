// Copyright 2025 The MathWorks, Inc.

package directory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-12345"

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(expectedLogDir, nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server-123")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	// Act
	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, directoryInstance, "Directory instance should not be nil")
}

func TestNew_MkdirTempError(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedError := assert.AnError

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return("", expectedError).
		Once()

	// Act
	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)

	// Assert
	require.ErrorIs(t, err, expectedError, "New should return the error from MkdirTemp")
	assert.Nil(t, directoryInstance, "Directory instance should be nil when error occurs")
}

func TestNew_SuppliedBaseDir_MkdirAllErrors(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := "/logs"

	mockConfig.EXPECT().
		BaseDir().
		Return(expectedLogDir).
		Once()

	expectedError := assert.AnError
	mockOSLayer.EXPECT().
		MkdirAll(expectedLogDir, os.FileMode(0o700)).
		Return(expectedError).
		Once()

	// Act
	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)

	//Assert
	require.ErrorIs(t, err, expectedError, "New should return the error from MkdirAll")
	assert.Nil(t, directoryInstance, "Directory instance should be nil when error occurs")
}

func TestDirectory_BaseDir_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-67890"

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(expectedLogDir, nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	baseDir := directoryInstance.BaseDir()

	// Assert
	assert.Equal(t, expectedLogDir, baseDir, "BaseDir should return the expected log directory")
}

func TestDirectory_BaseDir_SuppliedBaseDir_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/logs"

	mockConfig.EXPECT().
		BaseDir().
		Return(expectedLogDir).
		Once()

	mockOSLayer.EXPECT().
		MkdirAll(expectedLogDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	baseDir := directoryInstance.BaseDir()

	// Assert
	assert.Equal(t, expectedLogDir, baseDir, "BaseDir should return the expected log directory")
}

func TestDirectory_ID_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-12345"

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(expectedLogDir, nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	id := directoryInstance.ID()

	// Assert
	require.NoError(t, err, "ID should not return an error")
	assert.Equal(t, expectedSuffix, id, "ID should return the expected ID")
}

func TestDirectory_MkdirTemp_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-11111"
	pattern := "test-pattern-"
	expectedTempDir := "/tmp/matlab-mcp-core-server-11111/test-pattern-22222"

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(expectedLogDir, nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	expectedPattern := pattern + expectedSuffix + "-"
	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedPattern).
		Return(expectedTempDir, nil).
		Once()

	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	tempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.NoError(t, err, "MkdirTemp should not return an error")
	assert.Equal(t, expectedTempDir, tempDir, "MkdirTemp should return the expected temp directory")
}

func TestDirectory_MkdirTemp_EnforcesDashSuffix(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-11111"
	pattern := "test-pattern"
	expectedTempDir := "/tmp/matlab-mcp-core-server-11111/test-pattern-22222"

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(expectedLogDir, nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	expectedPattern := pattern + "-" + expectedSuffix + "-"
	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedPattern).
		Return(expectedTempDir, nil).
		Once()

	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	tempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.NoError(t, err, "MkdirTemp should not return an error")
	assert.Equal(t, expectedTempDir, tempDir, "MkdirTemp should return the expected temp directory")
}

func TestDirectory_MkdirTemp_MkdirTempError(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-33333"
	pattern := "test-pattern-"
	expectedError := assert.AnError

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(expectedLogDir, nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	expectedPattern := pattern + expectedSuffix + "-"
	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedPattern).
		Return("", expectedError).
		Once()

	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	tempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.ErrorIs(t, err, expectedError, "MkdirTemp should return the error from OSLayer.MkdirTemp")
	assert.Empty(t, tempDir, "MkdirTemp should return empty string when error occurs")
}

func TestDirectory_MkdirTemp_SuppliedBaseDir_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/logs"
	pattern := "test-pattern-"
	expectedTempDir := "/logs/test-pattern-22222"

	mockConfig.EXPECT().
		BaseDir().
		Return(expectedLogDir).
		Once()

	mockOSLayer.EXPECT().
		MkdirAll(expectedLogDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	expectedPattern := pattern + expectedSuffix + "-"
	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedPattern).
		Return(expectedTempDir, nil).
		Once()

	directoryInstance, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	tempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.NoError(t, err, "MkdirTemp should not return an error")
	assert.Equal(t, expectedTempDir, tempDir, "MkdirTemp should return the expected temp directory")
}

func TestDirectory_RecordToLogger_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	expectedLogDir := "/tmp/matlab-mcp-core-server-33333"

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(expectedLogDir, nil).
		Once()

	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "123"
	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(filepath.Join(expectedLogDir, directory.MarkerFileName), "").
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	directory, err := directory.New(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	testLogger := testutils.NewInspectableLogger()

	// Act
	directory.RecordToLogger(testLogger)

	// Assert
	infoLogs := testLogger.InfoLogs()
	require.Len(t, infoLogs, 1)

	fields, found := infoLogs["Application directory state"]
	require.True(t, found, "Expected log message not found")

	actualValue, exists := fields["log-dir"]
	require.True(t, exists, "log-dir field not found in log")
	assert.Equal(t, expectedLogDir, actualValue, "log-dir field has incorrect value")
}
