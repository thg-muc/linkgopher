package main

import (
	"testing"

	"github.com/atotto/clipboard"
)

// Test cases to check whether Strings are correctly identified
func TestCheckString(t *testing.T) {
	testCases := []struct {
		name             string
		possibleLink     string
		expectedLinkType int
	}{
		{`SMB_link`, `smb://example.corp/_path/folder/etc`, -1},
		{`SMB_link_with_leading_whitespace`, `  smb://example.corp/_path/folder/etc`, -1},
		{`SMB_link_with_trailing_slash`, `smb://example.corp/_path/folder/etc/`, -1},
		{`SMB_link_without_prefix`, `//example.corp/_path/folder/etc`, -1},
		{`SMB_link_with_file`, `smb://example.corp/_path/folder/etc/file.txt`, -1},
		{`SMB_link_broken`, `smb:/example.corppathfolderetc`, 0},
		{`Windows_link`, `file:\\example.corp\_path\folder\etc`, 1},
		{`Windows_link_with_leading_whitespace`, ` file:\\example.corp\_path\folder\etc`, 1},
		{`Windows_link_with_trailing_slash`, `file:\\example.corp\_path\folder\etc\`, 1},
		{`Windows_link_without_prefix`, `\\example.corp\_path\folder\etc`, 1},
		{`Windows_link_with_file`, `file:\\example.corp\_path\folder\etc\file.txt`, 1},
		{`Windows_link_broken`, `file:\example.corppathfolderetc`, 0},
		{`Web_link`, `https://example.com`, 0},
		{`Normal_string`, `normal string`, 0},
		{`Empty_string`, ``, 0},
		{`Mixed_slashes`, `normal string/with/forward\and\backslashes`, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			linkType, _ := checkString(tc.possibleLink)
			if linkType != tc.expectedLinkType {
				t.Errorf("Expected link type %d, got %d", tc.expectedLinkType, linkType)
			}
		})
	}
}

// Test cases to check whether Mac paths are correctly converted to Windows paths
func TestMacToWindowsPath(t *testing.T) {
	testCases := []struct {
		name                string
		macPath             string
		expectedWindowsPath string
	}{
		{`Normal_string_with_whitespace`, ` normal string `, ` normal string `},
		{`SMB_link`, `smb://example.corp/_path/folder/etc`, `file:\\example.corp\_path\folder\etc`},
		{`SMB_link_with_trailing_slash`, `smb://example.corp/_path/folder/etc/`, `file:\\example.corp\_path\folder\etc`},
		{`SMB_link_with_file`, `smb://example.corp/_path/folder/etc/file.txt`, `file:\\example.corp\_path\folder\etc\file.txt`},
		{`SMB_link_with_leading_whitespace`, ` smb://example.corp/_path/folder/etc`, `file:\\example.corp\_path\folder\etc`},
		{`SMB_link_with_leading_and_trailing_whitespace`, `   smb://example.corp/_path/folder/etc     `, `file:\\example.corp\_path\folder\etc`},
		{`SMB_link_without_prefix`, `//example.corp/_path/folder/etc`, `file:\\example.corp\_path\folder\etc`},
		{`SMB_link_with_spaces_in_path`, `smb://example.corp/_path/folder with spaces/etc/file name.txt`, `file:\\example.corp\_path\folder with spaces\etc\file name.txt`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, cleanedLink := checkString(tc.macPath)
			windowsPath := macToWindowsPath(cleanedLink)
			if windowsPath != tc.expectedWindowsPath {
				t.Errorf("Expected Windows path '%s', got '%s'", tc.expectedWindowsPath, windowsPath)
			}
		})
	}
}

// Test cases to check whether Windows paths are correctly converted to Mac paths
func TestWindowsToMacPath(t *testing.T) {
	testCases := []struct {
		name            string
		windowsPath     string
		expectedMacPath string
	}{
		{`Normal_string_with_whitespace`, ` normal string `, ` normal string `},
		{`Windows_link`, `file:\\example.corp\_path\folder\etc`, `smb://example.corp/_path/folder/etc`},
		{`Windows_link_with_trailing_slash`, `file:\\example.corp\_path\folder\etc\`, `smb://example.corp/_path/folder/etc/`},
		{`Windows_link_with_file`, `file:\\example.corp\_path\folder\etc\file.txt`, `smb://example.corp/_path/folder/etc/file.txt`},
		{`Windows_link_with_leading_whitespace`, ` file:\\example.corp\_path\folder\etc`, `smb://example.corp/_path/folder/etc`},
		{`Windows_link_with_leading_and_trailing_whitespace`, `   file:\\example.corp\_path\folder\etc     `, `smb://example.corp/_path/folder/etc`},
		{`Windows_link_without_prefix`, `\\example.corp\_path\folder\etc`, `smb://example.corp/_path/folder/etc`},
		{`Windows_link_with_forward_slashes`, `file://example.corp/_path/folder/etc`, `smb://example.corp/_path/folder/etc`},
		{`Windows_link_with_forward_slashes_and_trailing_slash`, `file://example.corp/_path/folder/etc/`, `smb://example.corp/_path/folder/etc/`},
		{`Windows_link_with_spaces_in_path`, `file:\\example.corp\_path\folder with spaces\etc\file name.txt`, `smb://example.corp/_path/folder with spaces/etc/file name.txt`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, cleanedLink := checkString(tc.windowsPath)
			macPath := windowsToMacPath(cleanedLink)
			if macPath != tc.expectedMacPath {
				t.Errorf("Expected Mac path '%s', got '%s'", tc.expectedMacPath, macPath)
			}
		})
	}
}

// Test cases to check whether the input is correctly parsed
func TestParseInput(t *testing.T) {
	// Helper function to set clipboard content
	setClipboard := func(content string) {
		err := clipboard.WriteAll(content)
		if err != nil {
			t.Fatalf("Failed to set clipboard content: %v", err)
		}
	}

	tests := []struct {
		name          string
		args          []string
		useClipboard  bool
		clipboardData string
		expected      string
		expectError   bool
	}{
		{
			name:         "Command line arguments",
			args:         []string{"file:\\\\example.com\\path\\to\\file"},
			useClipboard: false,
			expected:     "file:\\\\example.com\\path\\to\\file",
			expectError:  false,
		},
		{
			name:         "Multiple command line arguments",
			args:         []string{"file:\\\\example.com\\path", "with spaces\\to\\file"},
			useClipboard: false,
			expected:     "file:\\\\example.com\\path with spaces\\to\\file",
			expectError:  false,
		},
		{
			name:          "Clipboard input",
			args:          []string{},
			useClipboard:  true,
			clipboardData: "smb://example.com/path/to/file",
			expected:      "smb://example.com/path/to/file",
			expectError:   false,
		},
		{
			name:          "Empty clipboard",
			args:          []string{},
			useClipboard:  true,
			clipboardData: "",
			expected:      "",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useClipboard {
				setClipboard(tt.clipboardData)
			}

			result, err := parseInput(tt.args, tt.useClipboard)

			if tt.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected result '%s', but got '%s'", tt.expected, result)
			}
		})
	}
}
