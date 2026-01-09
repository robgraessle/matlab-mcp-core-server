// Copyright 2025 The MathWorks, Inc.

package vmcblockhelp

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/baseresource"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

//go:embed vmchelp
var vmcHelpFiles embed.FS

// libraryMetadataRegex matches HTML comments like: <!-- Library: libraryName -->
var libraryMetadataRegex = regexp.MustCompile(`^<!--\s*Library:\s*(\S+)\s*-->`)

// extractLibraryFromContent extracts the library name from HTML comment metadata
// in the first line of the markdown content
func extractLibraryFromContent(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ""
	}
	
	firstLine := strings.TrimSpace(lines[0])
	matches := libraryMetadataRegex.FindStringSubmatch(firstLine)
	if len(matches) >= 2 {
		return matches[1]
	}
	
	return ""
}

type Resource struct {
	*baseresource.Resource
}

func New(loggerFactory baseresource.LoggerFactory) (*Resource, error) {
	baseRes, err := baseresource.New(
		name,
		title,
		description,
		mimeType,
		estimatedSize,
		uri,
		loggerFactory,
		Handler(),
	)
	if err != nil {
		return nil, err
	}

	return &Resource{
		Resource: baseRes,
	}, nil
}

func Handler() baseresource.ResourceHandler {
	return func(_ context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		logger.Info("Returning Vitis Model Composer block help resource")

		// First pass: collect all files and organize by category
		type fileInfo struct {
			path     string
			category string
			name     string
			library  string
			content  []byte
		}
		var files []fileInfo

		err := fs.WalkDir(vmcHelpFiles, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			
			// Skip directories and non-markdown files
			if d.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}

			// Read the file content
			content, readErr := fs.ReadFile(vmcHelpFiles, path)
			if readErr != nil {
				logger.WithError(readErr).With("file", path).Warn("Failed to read embedded file")
				return nil // Continue processing other files
			}

			contentStr := string(content)
			
			// Extract library metadata from first line
			library := extractLibraryFromContent(contentStr)

			// Extract title from first heading (skip metadata line if present)
			title := strings.TrimSuffix(filepath.Base(path), ".md")
			lines := strings.Split(contentStr, "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "#") {
					title = strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
					break
				}
			}

			// Determine category from path
			category := "Other"
			if strings.Contains(path, "AIE/") || strings.HasPrefix(path, "AIE\\") {
				category = "AIE"
			} else if strings.Contains(path, "HDL/") || strings.HasPrefix(path, "HDL\\") {
				category = "HDL"
			} else if strings.Contains(path, "HLS/") || strings.HasPrefix(path, "HLS\\") {
				category = "HLS"
			}

			files = append(files, fileInfo{
				path:     path,
				category: category,
				name:     title,
				library:  library,
				content:  content,
			})

			return nil
		})

		if err != nil {
			logger.WithError(err).Error("Failed to walk embedded filesystem")
			return nil, err
		}

		// Build content with table of contents
		var combinedContent strings.Builder
		combinedContent.WriteString("# Vitis Model Composer Block Help\n\n")
		combinedContent.WriteString("This resource contains comprehensive help documentation for Vitis Model Composer blocks.\n\n")
		combinedContent.WriteString("**Total Blocks:** ")
		combinedContent.WriteString(strings.TrimSpace(strings.Fields(fmt.Sprintf("%d", len(files)))[0]))
		combinedContent.WriteString("\n\n")

		// Build table of contents organized by category
		combinedContent.WriteString("## Table of Contents\n\n")
		
		categories := []string{"AIE", "HDL", "HLS", "Other"}
		for _, cat := range categories {
			var categoryFiles []fileInfo
			for _, f := range files {
				if f.category == cat {
					categoryFiles = append(categoryFiles, f)
				}
			}
			
			if len(categoryFiles) > 0 {
				combinedContent.WriteString("### ")
				combinedContent.WriteString(cat)
				combinedContent.WriteString(" Blocks (")
				combinedContent.WriteString(strings.TrimSpace(strings.Fields(fmt.Sprintf("%d", len(categoryFiles)))[0]))
				combinedContent.WriteString(")\n\n")
				
				for _, f := range categoryFiles {
					combinedContent.WriteString("- [")
					combinedContent.WriteString(f.name)
					combinedContent.WriteString("](#")
					combinedContent.WriteString(strings.ToLower(strings.ReplaceAll(f.name, " ", "-")))
					combinedContent.WriteString(")")
					
					// Add library information if available
					if f.library != "" {
						combinedContent.WriteString(" - **Library:** `")
						combinedContent.WriteString(f.library)
						combinedContent.WriteString("`")
					}
					
					combinedContent.WriteString("\n")
				}
				combinedContent.WriteString("\n")
			}
		}

		combinedContent.WriteString("---\n\n")

		// Second pass: add all block documentation organized by category
		for _, cat := range categories {
			var categoryFiles []fileInfo
			for _, f := range files {
				if f.category == cat {
					categoryFiles = append(categoryFiles, f)
				}
			}
			
			if len(categoryFiles) == 0 {
				continue
			}

			combinedContent.WriteString("# ")
			combinedContent.WriteString(cat)
			combinedContent.WriteString(" Blocks\n\n")

			for _, f := range categoryFiles {
				// Add block header with anchor
				combinedContent.WriteString("## ")
				combinedContent.WriteString(f.name)
				combinedContent.WriteString("\n\n")
				combinedContent.WriteString("**Category:** ")
				combinedContent.WriteString(cat)
				combinedContent.WriteString("  \n")
				
				// Add library information if available
				if f.library != "" {
					combinedContent.WriteString("**Simulink Library:** `")
					combinedContent.WriteString(f.library)
					combinedContent.WriteString("`  \n")
				}
				
				combinedContent.WriteString("**Source File:** `")
				combinedContent.WriteString(f.path)
				combinedContent.WriteString("`\n\n")
				combinedContent.Write(f.content)
				combinedContent.WriteString("\n\n---\n\n")
			}
		}

		return &baseresource.ReadResourceResult{
			Contents: []baseresource.ResourceContents{
				{
					MIMEType: mimeType,
					Text:     combinedContent.String(),
				},
			},
		}, nil
	}
}

// SearchBlock searches for a block by name and returns its documentation
func SearchBlock(searchTerm string) (string, error) {
	searchLower := strings.ToLower(strings.TrimSpace(searchTerm))
	
	type blockInfo struct {
		path       string
		title      string
		category   string
		library    string
		content    []byte
		exactMatch bool
	}
	
	var matches []blockInfo

	err := fs.WalkDir(vmcHelpFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Read file content
		content, readErr := fs.ReadFile(vmcHelpFiles, path)
		if readErr != nil {
			return nil // Skip files we can't read
		}

		contentStr := string(content)
		
		// Extract library metadata
		library := extractLibraryFromContent(contentStr)

		// Extract title from first heading
		title := strings.TrimSuffix(filepath.Base(path), ".md")
		lines := strings.Split(contentStr, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "#") {
				title = strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
				break
			}
		}

		// Determine category
		category := "Other"
		if strings.Contains(path, "AIE/") || strings.HasPrefix(path, "AIE\\") {
			category = "AIE"
		} else if strings.Contains(path, "HDL/") || strings.HasPrefix(path, "HDL\\") {
			category = "HDL"
		} else if strings.Contains(path, "HLS/") || strings.HasPrefix(path, "HLS\\") {
			category = "HLS"
		}

		// Check for match
		titleLower := strings.ToLower(title)
		if strings.Contains(titleLower, searchLower) {
			matches = append(matches, blockInfo{
				path:       path,
				title:      title,
				category:   category,
				library:    library,
				content:    content,
				exactMatch: titleLower == searchLower,
			})
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to search block help: %w", err)
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no block found matching '%s'", searchTerm)
	}

	// Prefer exact matches
	var selected blockInfo
	for _, match := range matches {
		if match.exactMatch {
			selected = match
			break
		}
	}

	// If no exact match, use the first match
	if selected.title == "" {
		selected = matches[0]
	}

	// Build documentation
	var doc strings.Builder
	
	doc.WriteString("# ")
	doc.WriteString(selected.title)
	doc.WriteString("\n\n")
	doc.WriteString("**Category:** ")
	doc.WriteString(selected.category)
	doc.WriteString("  \n")
	
	// Add library information if available
	if selected.library != "" {
		doc.WriteString("**Simulink Library:** `")
		doc.WriteString(selected.library)
		doc.WriteString("`  \n")
	}
	
	doc.WriteString("**Source File:** `")
	doc.WriteString(selected.path)
	doc.WriteString("`\n\n")
	
	if len(matches) > 1 {
		doc.WriteString("**Note:** Found ")
		doc.WriteString(fmt.Sprintf("%d", len(matches)))
		doc.WriteString(" matching blocks. Showing: ")
		doc.WriteString(selected.title)
		doc.WriteString(". Other matches: ")
		for i, m := range matches {
			if m.title != selected.title {
				if i > 0 {
					doc.WriteString(", ")
				}
				doc.WriteString(m.title)
				doc.WriteString(" (")
				doc.WriteString(m.category)
				doc.WriteString(")")
			}
		}
		doc.WriteString("\n\n")
	}
	
	doc.WriteString("---\n\n")
	doc.Write(selected.content)

	return doc.String(), nil
}
