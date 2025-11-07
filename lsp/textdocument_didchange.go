package lsp

type TextDocumentDidChangeNotification struct {
	Notification
	Parameters TextDocumentDidChangeParams `json:"params"`
}

type TextDocumentDidChangeParams struct {
	TextDocument   VersionTextDocumentIdentifier    `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}
