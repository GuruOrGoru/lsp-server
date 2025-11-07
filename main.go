package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/guruorgoru/lsp-server/analysis"
	"github.com/guruorgoru/lsp-server/lsp"
	"github.com/guruorgoru/lsp-server/rpc"
)

func main() {
	logger := getNewLogger("/home/guruorgoru/projects/lsp_server/log.txt")
	logger.Println("Logger is initialised")

	writer := bufio.NewWriter(os.Stdout)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("Friendly error:", err)
			continue
		}
		handleMessage(logger, *state, writer, method, content)
	}
}

func getNewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("You should provide a good file :( ")
	}

	return log.New(logfile, "[guruLsp]", log.Ldate|log.Ltime|log.Lshortfile)
}

func unmarshalAndLog(logger *log.Logger, content []byte, v any) error {
	if err := json.Unmarshal(content, v); err != nil {
		logger.Println("Couldn't parse that sorry:", err)
		return err
	}
	return nil
}

func sendResponse(logger *log.Logger, writer *bufio.Writer, msgStruct any, logMsg string) {
	replyMsg := rpc.EncodeMessage(msgStruct)
	if replyMsg == "" {
		logger.Fatalln("Failed to encode " + logMsg + " message")
	}
	_, err := writer.WriteString(replyMsg)
	if err != nil {
		logger.Fatalf("Error writing to stdout: %v", err)
	}
	writer.Flush()
	logger.Println("Sent " + logMsg + " response")
}

func handleMessage(logger *log.Logger, state analysis.State, writer *bufio.Writer, method string, content []byte) {
	logger.Printf("Recieved Method: %v \n", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := unmarshalAndLog(logger, content, &request); err != nil {
			return
		}
		logger.Printf("Connected to: %v \t version: %v", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		msgStruct := lsp.NewInitializeResponse(request.ID)
		sendResponse(logger, writer, msgStruct, "initialize response")
	case "textDocument/didOpen":
		var notification lsp.DidOpenTextDocumentNotifications
		if err := unmarshalAndLog(logger, content, &notification); err != nil {
			return
		}
		diagnostic := state.OpenDocument(notification.Parameters.TextDocument.URI, notification.Parameters.TextDocument.Text)
		diagnosticNotification := lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         notification.Parameters.TextDocument.URI,
				Diagnostics: diagnostic,
			},
		}
		sendResponse(logger, writer, diagnosticNotification, "diagnostic notification")
	case "textDocument/didChange":
		var notification lsp.TextDocumentDidChangeNotification
		if err := unmarshalAndLog(logger, content, &notification); err != nil {
			return
		}
		logger.Printf("Changed file: %v", notification.Parameters.TextDocument.URI)
		for _, change := range notification.Parameters.ContentChanges {
			diagnostics := state.ChangeDocument(notification.Parameters.TextDocument.URI, change.Text)
			diagnosticNotification := lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         notification.Parameters.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			}
			sendResponse(logger, writer, diagnosticNotification, "diagnostic notification")
		}

	case "textDocument/hover":
		var request lsp.TextDocumentHoverRequest
		if err := unmarshalAndLog(logger, content, &request); err != nil {
			return
		}
		logger.Printf("Hover request for file: %v", request.Params.TextDocument.URI)
		hoverInfo := state.GetHoverInformation(request.Params.TextDocument.URI, request.Params.Position.Line, request.Params.Position.Character)
		hoverResponse := lsp.HoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &request.ID,
			},
			Result: lsp.HoverResult{
				Contents: hoverInfo,
			},
		}
		sendResponse(logger, writer, hoverResponse, "hover response")
	case "textDocument/definition":
		var request lsp.TextDocumentDefinitionRequest
		if err := unmarshalAndLog(logger, content, &request); err != nil {
			return
		}
		logger.Printf("Definition request for file: %v", request.Params.TextDocument.URI)
		definitionLocation := state.GetDefinitionLocation(request.Params.TextDocument.URI, request.Params.Position)
		definitionResponse := lsp.DefinitionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &request.ID,
			},
			Result: definitionLocation,
		}
		sendResponse(logger, writer, definitionResponse, "definition response")
	
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
	if err := unmarshalAndLog(logger, content, &request); err != nil {
			logger.Println("Failed to unmarshal code action request:", err)
			return
		}
		logger.Printf("Code Action request for file: %v", request.Params.TextDocument.URI)
		codeActions := state.GetCodeActions(request.Params.TextDocument.URI, request.Params.Range)
		codeActionResponse := lsp.CodeActionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &request.ID,
			},
			Result: codeActions,
		}
		sendResponse(logger, writer, codeActionResponse, "code action response")

	case "textDocument/completion":
		var request lsp.TextDocumentCompletionRequest
		if err := unmarshalAndLog(logger, content, &request); err != nil {
			return
		}
		logger.Printf("Completion request for file: %v", request.Params.TextDocument.URI)
		completions := state.GetCompletions(request.Params.TextDocument.URI, request.Params.Position)
		completionResponse := lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &request.ID,
			},
			Result: completions,
		}
		sendResponse(logger, writer, completionResponse, "completion response")
	
	}
}
