package lsp

type TextDocumentHoverRequest struct {
	Request
	Params TextDocumentHoverParams `json:"params"`
}

type TextDocumentHoverParams struct {
	TextDocumentPositionParams
}

type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}

type HoverResult struct {
	Contents string `json:"contents"`
}
