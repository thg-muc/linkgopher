// Package main provides a tool for converting between Windows and Mac file paths.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

const AppTitle = "LinkGopher"

// Determines if a string is a Windows link, a Mac link, or neither.
// It returns an int indicating the type (1 for Windows, -1 for Mac, 0 for neither)
// and the cleaned string.
func checkString(possibleLink string) (int, string) {
	trimmedLink := strings.TrimSpace(possibleLink)

	// If there are both forward and backslashes, we can't determine the type
	if strings.Contains(trimmedLink, `/`) && strings.Contains(trimmedLink, `\\`) {
		return 0, possibleLink
	}

	// Check for Windows link
	if strings.HasPrefix(trimmedLink, `file:\\`) || strings.HasPrefix(trimmedLink, `\\`) {
		return 1, trimmedLink
	}

	// Check for Mac link
	if strings.HasPrefix(trimmedLink, `smb://`) || strings.HasPrefix(trimmedLink, `//`) {
		return -1, trimmedLink
	}

	return 0, possibleLink
}

// Converts a Mac file path to a Windows file path.
func macToWindowsPath(macPath string) string {
	// Convert smb:// to file:\\
	macPath = strings.ReplaceAll(macPath, `smb://`, `file:\\`)
	macPath = strings.ReplaceAll(macPath, `//`, `file:\\`)
	// Split the path into components
	components := strings.Split(macPath, `/`)
	// Combine the components with the Windows path separator
	windowsPath := strings.Join(components, `\`)
	// Remove trailing slash if present
	windowsPath = strings.TrimRight(windowsPath, `\`)

	return windowsPath
}

// Converts a Windows file path to a Mac file path.
func windowsToMacPath(windowsPath string) string {
	// Convert all forward slashes to backslashes
	windowsPath = strings.ReplaceAll(windowsPath, `/`, `\`)
	// Convert file:// to smb://
	windowsPath = strings.ReplaceAll(windowsPath, `file:\\`, `smb://`)
	windowsPath = strings.ReplaceAll(windowsPath, `\\`, `smb://`)
	// Split the path into components
	components := strings.Split(windowsPath, `\`)
	// Combine the components with the Mac path separator
	macPath := strings.Join(components, `/`)

	return macPath
}

// Takes a link string and converts it between Windows and Mac formats.
// It returns a string describing the conversion result and the converted link (if any).
func convertLink(link string) (string, string) {
	linkType, possibleLink := checkString(link)

	if linkType == 1 {
		convertedLink := windowsToMacPath(possibleLink)
		return fmt.Sprintf("%s - Converted Windows to Mac Path: \n%s", AppTitle, convertedLink), convertedLink
	} else if linkType == -1 {
		convertedLink := macToWindowsPath(possibleLink)
		return fmt.Sprintf("%s - Converted Mac to Windows Path: \n%s", AppTitle, convertedLink), convertedLink
	} else {
		// Check if the invalid link contained a "file" prefix
		var extraInfo string
		if strings.Contains(possibleLink, "file:") {
			extraInfo = "\nNote: When passing a windows (file:) path as an argument, please either use single quotes for the entire path ('file:\\\\example.corp\\folder') or escape each backslash with a second one (\\\\)."
		} else {
			extraInfo = ""
		}

		return fmt.Sprintf("%s - No valid link detected!%s", AppTitle, extraInfo), ""
	}
}

// Parses command line arguments to determine the link to convert.
func parseInput(args []string, useClipboard bool) (string, error) {
	if useClipboard {
		return clipboard.ReadAll()
	}
	return strings.Join(args, " "), nil
}

func main() {
	// Always print an empty line at the end
	defer fmt.Println()

	var link string
	var err error
	useClipboard := len(os.Args) < 2

	// Check if the clipboard should be checked
	if useClipboard {
		link, err = parseInput(nil, true)
	} else {
		link, err = parseInput(os.Args[1:], false)
	}

	if err != nil {
		fmt.Printf("%s: Error reading input: %v\n", AppTitle, err)
		return
	}

	// Perform the conversion
	result, convertedLink := convertLink(link)
	fmt.Println(result)

	if useClipboard && convertedLink != "" {
		err := clipboard.WriteAll(convertedLink)
		if err != nil {
			fmt.Printf("%s: Error writing to clipboard: %v\n", AppTitle, err)
		} else {
			fmt.Println("Converted link copied to clipboard.")
		}
	}
}
