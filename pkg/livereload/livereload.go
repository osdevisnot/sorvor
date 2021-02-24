// Package livereload implements HTML5 Server Side Events to live reload connected browsers
// Usage: first start the livereload instance
//		liveReload := livereload.New()
//		liveReload.Start()
// then, install an HTTP handler on desired path
//		http.Handle("/livereload", liveReload)
// then, reload the connected browsers
//      liveReload.Reload()
// The target browser must support HTML5 Server Side Events.
// Note: In Chrome and Firefox - there is a limit of 6 active connections per browser + domain.
// Chromium: https://bugs.chromium.org/p/chromium/issues/detail?id=275955
// Firefox: https://bugzilla.mozilla.org/show_bug.cgi?id=906896
package livereload

import (
	"fmt"
	"net/http"
	"time"
)

// JsSnippet is a minimal javascript client for browsers. Embed it in your index.html using a script tag:
//		<script>{{ LiveReload.JsSnippet }}</script>
const JsSnippet = `<script>const source = new EventSource('/livereload');const reload = () => location.reload(true);source.onmessage = reload;source.onerror = () => (source.onopen = reload);console.info('[sørvør] listening for file changes');</script>`

// LiveReload keeps track of connected browser clients and broadcasts messages to them
type LiveReload struct {
	clients  map[chan string]bool // map of connection pool
	incoming chan chan string     // channel to push new clients
	outgoing chan chan string     // channel to push disconnected clients
	messages chan string          // channel to publish new messages
}

// Start manages connections and broadcasts messages to connected clients
func (livereload *LiveReload) Start() {
	go func() {
		for {
			select {
			case s := <-livereload.incoming:
				// attach an incoming client
				livereload.clients[s] = true
			case s := <-livereload.outgoing:
				// detach an outgoing client
				delete(livereload.clients, s)
				close(s)
			case msg := <-livereload.messages:
				// send new message to all connected clients
				for s := range livereload.clients {
					s <- msg
				}
			}
		}
	}()

	go func() {
		for {
			livereload.messages <- "waiting"
			time.Sleep(time.Minute)
		}
	}()
}

// Reload sends a reload signal to all connected browsers
func (livereload *LiveReload) Reload() {
	livereload.messages <- "reload"
}

// sendEvent is a helper to create formatted SSE events based on message data.
func (livereload *LiveReload) sendEvent(res http.ResponseWriter, eventData string) {
	var eventType string

	switch eventData {
	case "ready":
		eventType = "connected"
	case "waiting":
		eventType = "ping"
	default:
		eventType = "message"
	}

	_, _ = res.Write([]byte(fmt.Sprintf("event: %s\nid: 0\ndata: %s\n", eventType, eventData)))
	_, _ = res.Write([]byte("\n\n"))
}

// ServeHTTP is a request handler for http requests
func (livereload *LiveReload) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	connection, ok := res.(http.Flusher)
	if !ok {
		http.Error(res, "Streaming Unsupported !!", http.StatusInternalServerError)
		return
	}

	channel := make(chan string)
	livereload.incoming <- channel

	// close the connection when this client is gone
	ctx := req.Context()
	go func() {
		select {
		case <-ctx.Done():
			livereload.outgoing <- channel
			return
		default:
		}
	}()

	// necessary http headers for SSE streaming...
	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.WriteHeader(200)

	livereload.sendEvent(res, "ready")
	connection.Flush()

	for {
		msg, open := <-channel
		if !open {
			break
		}
		livereload.sendEvent(res, msg)
		connection.Flush()
	}
}

// New creates a new LiveReload struct
func New() *LiveReload {
	livereload := &LiveReload{
		clients:  make(map[chan string]bool),
		incoming: make(chan (chan string)),
		outgoing: make(chan (chan string)),
		messages: make(chan string),
	}

	return livereload
}
