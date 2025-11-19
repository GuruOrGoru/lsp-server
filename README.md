# LSP Server

A simple Language Server Protocol (LSP) server implemented in Go. This server provides basic LSP features such as diagnostics, hover information, go-to-definition, code actions, and autocompletion for supported languages.

## Features

- Text document synchronization (open, change)
- Diagnostics reporting
- Hover information
- Go-to-definition
- Code actions
- Autocompletion

## Installation

1. Ensure you have Go 1.25.3 or later installed.
2. Clone or download the repository.
3. Run `go mod tidy` to download dependencies.
4. Build the server: `go build -o lsp-server main.go`

## Usage

The LSP server communicates via stdin/stdout using JSON-RPC. It is designed to be integrated with editors that support the Language Server Protocol.

### Running the Server

Execute the built binary:

```bash
./lsp-server
```

The server will log to `log.txt` in the current directory.

### Integrating with an Editor

To use this LSP server with an editor like VS Code:

1. Install the necessary LSP client extension for your editor (e.g., for VS Code, ensure LSP support is available).
2. Configure the editor to launch this server. For example, in VS Code, add to your settings:

```json
{
  "languageserver": {
    "lsp-server": {
      "command": "/path/to/lsp-server",
      "filetypes": ["your-language-extension"]
    }
  }
}
```

Replace `/path/to/lsp-server` with the actual path to the binary and `"your-language-extension"` with the file extensions this server should handle.

## Supported Methods

- `initialize`
- `textDocument/didOpen`
- `textDocument/didChange`
- `textDocument/hover`
- `textDocument/definition`
- `textDocument/codeAction`
- `textDocument/completion`

## Logging

Logs are written to `log.txt` for debugging purposes.