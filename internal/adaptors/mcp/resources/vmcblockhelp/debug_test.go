// Copyright 2025 The MathWorks, Inc.

package vmcblockhelp_test

import (
	"context"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/vmcblockhelp"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
)

func TestPrintOutput(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()
	handler := vmcblockhelp.Handler()

	// Act
	result, err := handler(context.Background(), mockLogger)

	// Assert
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if len(result.Contents) == 0 {
		t.Fatal("No contents returned")
	}

	text := result.Contents[0].Text
	
	// Find the end of the table of contents (where the first --- separator after TOC appears)
	tocEnd := strings.Index(text, "---\n\n#")
	if tocEnd == -1 {
		tocEnd = 3000 // fallback
	}
	
	t.Logf("\n=== TABLE OF CONTENTS ===\n%s\n", text[:tocEnd])
	
	// Find the first AIE block section as an example
	aieStart := strings.Index(text, "# AIE Blocks\n\n## ")
	if aieStart != -1 {
		// Find the next separator
		sectionEnd := strings.Index(text[aieStart:], "\n\n---\n\n")
		if sectionEnd != -1 {
			exampleSection := text[aieStart : aieStart+sectionEnd+7]
			t.Logf("\n=== EXAMPLE SECTION (First AIE Block) ===\n%s\n", exampleSection)
		}
	}
}
