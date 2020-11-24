package livereload

import (
	"fmt"
	"net/http"
	"time"
)

// Snippet is a minimal javascript client for reloading browse
// Embed in your `index.html` using a script tag: <script>{{ LiveReload.Snippet }}</script>
const Snippet = `
	const source = new EventSource('/livereload');
	const reload = () => location.reload(true);
	source.onmessage = reload;
	source.onerror = () => (source.onopen = reload);
	console.log('[sørvør] listening for file changes');
`

// LiveReload keeps track of connected browser clients and
// broadcasts messages to those connected browser clients.
type LiveReload struct {
	// map of connection pool, keys are channels over which we push an event,
	// value is always boolean true (kinda useless but map is easier)
	clients map[chan string]bool
	// channel to push new client to
	incoming chan chan string
	// channel to push disconnected clients
	outgoing chan chan string
	// channel into which messages are pushed
	messages chan string
}

// Start manages connections and broadcasts messages to current connected browser clients
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

// SendEvent is helper to create formatted SSE events based on event type and data.
func (livereload LiveReload) SendEvent(res http.ResponseWriter, eventType string, eventData string) {
	res.Write([]byte(fmt.Sprintf("event: %s\nid: 0\ndata: %s\n", eventType, eventData)))
	res.Write([]byte("\n\n"))
}

// ServeHTTP handler for `/livereload` urls
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

	livereload.SendEvent(res, "connected", "ready")
	connection.Flush()

	for {
		msg, open := <-channel
		if !open {
			break
		}
		livereload.SendEvent(res, "ping", msg)
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
