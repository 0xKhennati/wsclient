package wsClient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client for JSON-RPC communication
type Client struct {
	conn    *websocket.Conn
	counter int32 // Counter for pending messages
	url     string
}

// NewClient creates a new WebSocket client
func NewClient(wsURL string) (*Client, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn: conn,
		url:  wsURL,
	}, nil
}

// Send sends a request without waiting for response and increments counter
func (c *Client) Send(request *Request) error {
	if c.conn == nil {
		return fmt.Errorf("connection is closed")
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	err = c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	atomic.AddInt32(&c.counter, 1)
	return nil
}

// Receive receives a response and decrements counter
func (c *Client) Receive(response *Response) error {
	if c.conn == nil {
		return fmt.Errorf("connection is closed")
	}

	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	atomic.AddInt32(&c.counter, -1)

	if err := json.Unmarshal(data, response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// SendAndReceive sends a request and receives response
func (c *Client) SendAndReceive(request *Request, response *Response) error {
	if err := c.Send(request); err != nil {
		return err
	}
	return c.Receive(response)
}

// PendingCounter returns the number of pending messages
func (c *Client) PendingCounter() int32 {
	return atomic.LoadInt32(&c.counter)
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	atomic.StoreInt32(&c.counter, 0)
	return err
}

// GetURL returns the WebSocket URL
func (c *Client) GetURL() string {
	return c.url
}

// IsConnected checks if the client is connected
func (c *Client) IsConnected() bool {
	return c.conn != nil
}

// CheckAndReopenConnection checks if counter > 0 and reopens connection if needed
func (c *Client) CheckAndReopenConnection() error {
	if c.PendingCounter() > 0 {
		// Close existing connection
		if c.conn != nil {
			c.conn.Close()
		}

		// Reset counter
		atomic.StoreInt32(&c.counter, 0)

		// Reopen connection
		conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
		if err != nil {
			return fmt.Errorf("failed to reopen connection: %w", err)
		}

		c.conn = conn
	}
	return nil
}
