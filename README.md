# LinkGopher

A simple Golang CLI tool which converts network links between Windows and Mac formats.

## Purpose

LinkGopher automatically detects and converts network file paths between Windows format (`file:\\server\path`) and Mac format (`smb://server/path`). Useful for sharing network locations across different operating systems (potentially useful in business environments).

Motivation of this project was to create a simple tool that can be used in command-line environments, such as Git Bash or WSL. Go was chosen because I wanted to find out more about the language (and Python, as an interpreted language, was not the best choice for a cli tool).

## Features

- Automatic format detection (Windows or Mac)
- Supports paths with or without protocol prefixes
- Works with command-line arguments or clipboard content
- Automatically copies converted links to clipboard when no arguments provided

## Usage

```bash
# Convert a link provided as an argument
linkgo smb://example.corp/path/to/folder

# Convert a link from your clipboard (and copy result back to clipboard)
linkgo
```

## Build from Source

```bash
# Clone repository
git clone https://github.com/thg-muc/linkgopher.git
cd linkgopher

# Build executable
go build

# Run tests
go test -v
```
