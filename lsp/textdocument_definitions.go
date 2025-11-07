package lsp

type TextDocumentDefinitionRequest struct {
	Request
	Params TextDocumentDefinitionParams `json:"params"`
}

type TextDocumentDefinitionParams struct {
	TextDocumentPositionParams
}

type DefinitionResponse struct {
	Response
	Result Location `json:"result"`
}
