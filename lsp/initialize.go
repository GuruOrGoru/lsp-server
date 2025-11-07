package lsp

type InitializeRequest struct {
	Request
	Params InitializeParams `json:"params"`
}

type InitializeParams struct {
	ClientInfo *ClientInfo `json:"clientInfo,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type (
	ServerCapabilities struct {
		TextDocumentSync int  `json:"textDocumentSync"`
		HoverProvider    bool `json:"hoverProvider,omitempty"`
		DefinitionProvider bool `json:"definitionProvider,omitempty"`
		CodeActionProvider bool `json:"codeActionProvider,omitempty"`
	 	CompletionProvider  *CompletionOptions `json:"completionProvider,omitempty"`
	}
	CompletionOptions struct {
		ResolveProvider   bool     `json:"resolveProvider,omitempty"`
		TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	}
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version,omitempty"`
	}
)

func NewInitializeResponse(id int) *InitializeResponse {
	return &InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 1,
				HoverProvider:    true,
				DefinitionProvider: true,
				CodeActionProvider: true,
				CompletionProvider: &CompletionOptions{
					ResolveProvider:   false,
					TriggerCharacters: []string{"."},
				},
			},
			ServerInfo: &ServerInfo{
				Name:    "MyLSPServer",
				Version: "0.0.0.0.0.0-beta1.preview",
			},
		},
	}
}
