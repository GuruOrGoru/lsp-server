package lsp

type DidOpenTextDocumentNotifications struct {
	Notification
	Parameters DidOpenTextDocumentParams `json:"params"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}
