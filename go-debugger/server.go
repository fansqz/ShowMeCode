package main

import (
	"bufio"
	"fmt"
	"github.com/fansqz/go-debugger/debugger"
	"github.com/google/go-dap"
	"io"
	"log"
	"net"
	"sync"
)

// handleConnection handles a connection from a single client.
// It reads and decodes the incoming data and dispatches it
// to per-request processing goroutines. It also launches the
// sender goroutine to send resulting messages over the connection
// back to the client.
func handleConnection(conn net.Conn, d debugger.Debugger) {
	// 创建调试session
	debugSession := DebugSession{
		conn:      conn,
		debugger:  d,
		rw:        bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		sendQueue: make(chan dap.Message),
	}

	for {
		err := debugSession.handleRequest()
		if err != nil {
			if err == io.EOF {
				log.Printf("No more data to read:", err)
				break
			}
			log.Printf("Server error: ", err)
		}
	}

	log.Printf("Closing connection from", conn.RemoteAddr())
	debugSession.sendWg.Wait()
	close(debugSession.sendQueue)
	conn.Close()
}

func (d *DebugSession) handleRequest() error {
	request, err := dap.ReadProtocolMessage(d.rw.Reader)
	if err != nil {
		return err
	}
	d.dispatchRequest(request)
	return nil
}

func (d *DebugSession) dispatchRequest(request dap.Message) {
	switch request := request.(type) {
	case *dap.InitializeRequest:
		d.onInitializeRequest(request)
	case *dap.TerminateRequest:
		d.onTerminateRequest(request)
	case *dap.SetBreakpointsRequest:
		d.onSetBreakpointsRequest(request)
	case *dap.ConfigurationDoneRequest:
		d.onConfigurationDoneRequest(request)
	case *dap.ContinueRequest:
		d.onContinueRequest(request)
	case *dap.NextRequest:
		d.onNextRequest(request)
	case *dap.StepInRequest:
		d.onStepInRequest(request)
	case *dap.StepOutRequest:
		d.onStepOutRequest(request)
	case *dap.StackTraceRequest:
		d.onStackTraceRequest(request)
	case *dap.ScopesRequest:
		d.onScopesRequest(request)
	case *dap.VariablesRequest:
		d.onVariablesRequest(request)
	default:
		if baseReq, ok := request.(*dap.Request); ok {
			d.send(newErrorResponse(baseReq.Seq, baseReq.Command, fmt.Sprint("%s is not yet supported", baseReq.Command)))
		}
		fmt.Printf("Unable to process %#v", request)
	}
}

// send Message响应给客户端
func (d *DebugSession) send(message dap.Message) {
	dap.WriteProtocolMessage(d.rw.Writer, message)
	d.rw.Flush()
}

// DebugSession 调试会话
type DebugSession struct {
	conn net.Conn
	// rw is used to read requests and write events/responses
	rw *bufio.ReadWriter

	debugger debugger.Debugger
	// sendQueue is used to capture messages from multiple request
	// processing goroutines while writing them to the client connection
	// from a single goroutine via sendFromQueue. We must keep track of
	// the multiple channel senders with a wait group to make sure we do
	// not close this channel prematurely. Closing this channel will signal
	// the sendFromQueue goroutine that it can exit.
	sendQueue chan dap.Message
	sendWg    sync.WaitGroup
}

// -----------------------------------------------------------------------
// Request Handlers
//
// Below is a dummy implementation of the request handlers.
// They take no action, but just return dummy responses.
// A real debug adaptor would call the debugger methods here
// and use their results to populate each response.

func (d *DebugSession) onInitializeRequest(request *dap.InitializeRequest) {
	response := &dap.InitializeResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	response.Body.SupportsConfigurationDoneRequest = true
	response.Body.SupportsFunctionBreakpoints = false
	response.Body.SupportsConditionalBreakpoints = false
	response.Body.SupportsHitConditionalBreakpoints = false
	response.Body.SupportsEvaluateForHovers = false
	response.Body.ExceptionBreakpointFilters = []dap.ExceptionBreakpointsFilter{}
	response.Body.SupportsStepBack = false
	response.Body.SupportsSetVariable = false
	response.Body.SupportsRestartFrame = false
	response.Body.SupportsGotoTargetsRequest = false
	response.Body.SupportsStepInTargetsRequest = false
	response.Body.SupportsCompletionsRequest = false
	response.Body.CompletionTriggerCharacters = []string{}
	response.Body.SupportsModulesRequest = false
	response.Body.AdditionalModuleColumns = []dap.ColumnDescriptor{}
	response.Body.SupportedChecksumAlgorithms = []dap.ChecksumAlgorithm{}
	response.Body.SupportsRestartRequest = false
	response.Body.SupportsExceptionOptions = false
	response.Body.SupportsValueFormattingOptions = false
	response.Body.SupportsExceptionInfoRequest = false
	response.Body.SupportTerminateDebuggee = false
	response.Body.SupportsDelayedStackTraceLoading = false
	response.Body.SupportsLoadedSourcesRequest = false
	response.Body.SupportsLogPoints = false
	response.Body.SupportsTerminateThreadsRequest = false
	response.Body.SupportsSetExpression = false
	response.Body.SupportsTerminateRequest = false
	response.Body.SupportsDataBreakpoints = false
	response.Body.SupportsReadMemoryRequest = false
	response.Body.SupportsDisassembleRequest = false
	response.Body.SupportsCancelRequest = false
	response.Body.SupportsBreakpointLocationsRequest = false
	// This is a fake set up, so we can start "accepting" configuration
	// requests for setting breakpoints, etc from the client at any time.
	// Notify the client with an 'initialized' event. The client will end
	// the configuration sequence with 'configurationDone' request.
	e := &dap.InitializedEvent{Event: *newEvent("initialized")}
	d.send(e)
	d.send(response)
}

func (d *DebugSession) onTerminateRequest(request *dap.TerminateRequest) {
	err := d.debugger.Terminate()
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	d.send(newErrorResponse(request.Seq, request.Command, "TerminateRequest is not yet supported"))
}

func (d *DebugSession) onSetBreakpointsRequest(request *dap.SetBreakpointsRequest) {
	err := d.debugger.SetBreakpoints(request.Arguments.Source, request.Arguments.Breakpoints)
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.SetBreakpointsResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	response.Body.Breakpoints = make([]dap.Breakpoint, len(request.Arguments.Breakpoints))
	for i, b := range request.Arguments.Breakpoints {
		response.Body.Breakpoints[i].Line = b.Line
		response.Body.Breakpoints[i].Verified = true
	}
	d.send(response)
}

func (d *DebugSession) onConfigurationDoneRequest(request *dap.ConfigurationDoneRequest) {
	err := d.debugger.Run()
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.ConfigurationDoneResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	d.send(response)
}

func (d *DebugSession) onContinueRequest(request *dap.ContinueRequest) {
	err := d.debugger.Continue()
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.ContinueResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	d.send(response)
}

func (d *DebugSession) onNextRequest(request *dap.NextRequest) {
	err := d.debugger.StepOver()
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.NextResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	d.send(response)
}

func (d *DebugSession) onStepInRequest(request *dap.StepInRequest) {
	err := d.debugger.StepIn()
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.StepInResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	d.send(response)
}

func (d *DebugSession) onStepOutRequest(request *dap.StepOutRequest) {
	err := d.debugger.StepOut()
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.StepOutResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	d.send(response)
}

func (d *DebugSession) onStackTraceRequest(request *dap.StackTraceRequest) {
	stacktrace, err := d.debugger.GetStackTrace()
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.StackTraceResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	response.Body = dap.StackTraceResponseBody{
		StackFrames: stacktrace,
		TotalFrames: len(stacktrace),
	}
	d.send(response)
}

func (d *DebugSession) onScopesRequest(request *dap.ScopesRequest) {
	scopes, err := d.debugger.GetScopes(request.Arguments.FrameId)
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.ScopesResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	response.Body = dap.ScopesResponseBody{
		Scopes: scopes,
	}
	d.send(response)
}

func (d *DebugSession) onVariablesRequest(request *dap.VariablesRequest) {
	variables, err := d.debugger.GetVariables(request.Arguments.VariablesReference)
	if err != nil {
		d.send(newErrorResponse(request.Seq, request.Command, err.Error()))
	}
	response := &dap.VariablesResponse{}
	response.Response = *newResponse(request.Seq, request.Command)
	response.Body = dap.VariablesResponseBody{
		Variables: variables,
	}
	d.send(response)
}

// sendErrorResponseWithOpts offers configuration options.
//
//	showUser - if true, the error will be shown to the user (e.g. via a visible pop-up)
func (s *DebugSession) sendErrorResponseWithOpts(request dap.Request, id int, summary, details string, showUser bool) {
	er := &dap.ErrorResponse{}
	er.Type = "response"
	er.Command = request.Command
	er.RequestSeq = request.Seq
	er.Success = false
	er.Message = summary
	er.Body.Error = &dap.ErrorMessage{
		Id:       id,
		Format:   fmt.Sprintf("%s: %s", summary, details),
		ShowUser: showUser,
	}
	s.send(er)
}

func newEvent(event string) *dap.Event {
	return &dap.Event{
		ProtocolMessage: dap.ProtocolMessage{
			Seq:  0,
			Type: "event",
		},
		Event: event,
	}
}

func newResponse(requestSeq int, command string) *dap.Response {
	return &dap.Response{
		ProtocolMessage: dap.ProtocolMessage{
			Seq:  0,
			Type: "response",
		},
		Command:    command,
		RequestSeq: requestSeq,
		Success:    true,
	}
}

func newErrorResponse(requestSeq int, command string, message string) *dap.ErrorResponse {
	er := &dap.ErrorResponse{}
	er.Response = *newResponse(requestSeq, command)
	er.Success = false
	er.Body.Error = &dap.ErrorMessage{}
	er.Body.Error.Format = message
	er.Body.Error.Id = 12345
	return er
}
