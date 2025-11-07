package analysis

import (
	"strings"

	"github.com/guruorgoru/lsp-server/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() *State {
	return &State{
		Documents: make(map[string]string),
	}
}

func getDiagnostics(content string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}

	for row, line := range splitLines(content) {
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source: "common sense",
				Message: "Be respectful while coding",
			})
		}

		if strings.Contains(line, "Windows") {
			idx := strings.Index(line, "Windows")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: LineRange(row, idx, idx+len("Windows")),
				Severity: 1,
				Source: "common sense",
				Message: "Now you are just being mean",
			})
		}
	}

	return diagnostics
}

func (s *State) OpenDocument(document, content string) []lsp.Diagnostic {
	s.Documents[document] = content
	return getDiagnostics(content)
}

func (s *State) ChangeDocument(document, content string) []lsp.Diagnostic {
	s.Documents[document] = content
	return getDiagnostics(content)
}

func (s *State) GetHoverInformation(document string, positionLine int, positionCharacter int) string {
	content, exists := s.Documents[document]
	if !exists {
		return "Document not found"
	}

	lines := splitLines(content)
	if positionLine < 0 || positionLine >= len(lines) {
		return "Line out of range"
	}

	line := lines[positionLine]
	words := splitWords(line)
	charCount := 0
	for _, word := range words {
		if positionCharacter >= charCount && positionCharacter <= charCount+len(word) {
			return "Hover info for word: " + word
		}
		charCount += len(word) + 1 // +1 for space
	}

	return "No word found at the given position"
}

func (s *State) GetDefinitionLocation(document string, position lsp.Position) lsp.Location {
	content, exists := s.Documents[document]
	if !exists {
		return lsp.Location{}
	}

	_ = content

	return lsp.Location{
		URI: document,
		Range: lsp.Range{
			Start: lsp.Position{Line: 0, Character: 0},
			End:   lsp.Position{Line: 0, Character: 0},
		},
	}
}

func (s *State) GetCodeActions(document string, position lsp.Range) []lsp.CodeAction {
	content, exists := s.Documents[document]
	if !exists {
		return nil
	}

	actions := []lsp.CodeAction{}

	for row, line := range splitLines(content) {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[document] = []lsp.TextEdit{
				{
					Range: LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},

			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[document] = []lsp.TextEdit{
				{
					Range: LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Censor VS Code",
				Edit: &lsp.WorkspaceEdit{
					Changes: censorChange,
				},
			})
		}
	}
	return actions
}

func (s *State) GetCompletions(document string, position lsp.Position) []lsp.CompletionItem {
	content, exists := s.Documents[document]
	if !exists {
		return nil
	}
	_ = content

	items  := []lsp.CompletionItem{
		{
			Label: "HelloWorld",
			Kind: 2,
			Detail: "Sample completion item",
			Documentation: "This is a sample completion item for demonstration purposes.",
		},
		{
			Label: "GuruOrGoru",
			Kind: 2,
			Detail: "Completion for the guru",
			Documentation: "This completion item represents the guru himself.",
		},
		{
			Label: "AvhiIsGay",
			Kind: 2,
			Detail: "Just a fun completion",
			Documentation: "This completion item is just for fun.",
		},
		{
			Label: "GoLang",
			Kind: 2,
			Detail: "Programming Language",
			Documentation: "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
		},
	}
	return items

}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{Line: line, Character: start},
		End:   lsp.Position{Line: line, Character: end},
	}
}

func splitLines(content string) []string {
	var lines []string
	currentLine := ""
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(char)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

func splitWords(line string) []string {
	var words []string
	currentWord := ""
	for _, char := range line {
		if char == ' ' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}
	if currentWord != "" {
		words = append(words, currentWord)
	}
	return words
}

