package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// message sending

// SendCommand is used for sending new message from client

type SendCommand struct {
	Message string
}

// NameCommand is used for setting client display name
type NameCommand struct {
	Name string
}

// MessageCommand is used for notifying new messages
type MessageCommand struct {
	Name    string
	Message string
}

//  You can implement a reader to parse the command from the string and write a writer method to convert the command back into a string.
// The Go programming language uses io.Reader and io.Writer as built-in interfaces,
// and hence the implementation of these methods is not necessarily aware that these methods are used for a TCP stream.

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}
func (w *CommandWriter) writeString(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}
func (w *CommandWriter) Write(command interface{}) error {
	// naive implementation ...
	var err error
	switch v := command.(type) {
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v\n", v.Message))
	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v %v\n", v.Name, v.Message))
	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Name))
	default:
		err = UnknownCommand
	}
	return err
}

// reader method
// mostly use for error handling

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}
func (r *CommandReader) Read() (interface{}, error) {
	// Read the first part
	commandName, err := r.reader.ReadString(' ')
	if err != nil {
		return nil, err
	}
	switch commandName {
	case "MESSAGE ":
		user, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}
		message, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return MessageCommand{
			user[:len(user)-1],
			message[:len(message)-1],
		}, nil
		// similar implementation for other commands
	default:
		log.Printf("Unknown command: %v", commandName)
	}
	return nil, UnknownCommand
}

// Setting up the chat server:
// In this section, the server interface can be defined as follows.
// I don't know anything about the interface yet and will research it in the future.
// The interface makes the definition of behaviors clearer.

type ChatServer interface {
	Listen(address string) error
	Broadcast(command interface{}) error
	Start()
	Close()
}

// The server listens to incoming connections using the Listen() method, and the Start() and Close() methods are
//  used to start and stop the server, respectively, and the BroadCast() method is used to send commands to other clients.

// Now we can see the actual implementation of the server.

type TcpChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
}
type client struct {
	conn   net.Conn
	name   string
	writer *protocol.CommandWriter
}

func (s *TcpChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err == nil {
		s.listener = l
	}
	log.Printf("Listening on %v", address)
	return err
}
func (s *TcpChatServer) Close() {
	s.listener.Close()
}
func (s *TcpChatServer) Start() {
	for {
		// need a way to break the loop
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			// handle connection
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

// When the server accepts a connection, it creates a client structure to keep track of clients.
//  At this time we should use mutex to avoid race condition.
//  Goroutine does not solve all problems and we still have to fix all race conditions manually.

func (s *TcpChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}
	s.clients = append(s.clients, client)
	return client
}
func (s *TcpChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

// The main method of serve is to read the message from the client and manage each message accordingly.
//  Since we have two protocol read and write methods,
//   the server only needs to interact with the messages at a high level and there is no need to communicate with the byte stream.
//  If the server receives a SendCommand, all it does is broadcast the message to all other clients.

func (s *TcpChatServer) serve(client *client) {
	cmdReader := protocol.NewCommandReader(client.conn)
	defer s.remove(client)
	for {
		cmd, err := cmdReader.Read()
		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}
		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.name,
				})
			case protocol.NameCommand:
				client.name = v.Name
			}
		}
		if err == io.EOF {
			break
		}
	}
}

func (s *TcpChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {

		client.writer.Write(command)
	}
	return nil
}

// client :

type ChatClient interface {
	Dial(address string) error
	Send(command interface{}) error
	SendMessage(message string) error
	SetName(name string) error
	Start()
	Close()
	Incoming() chan protocol.MessageCommand
}

// The client can connect to the server using the Dial() method, and the Start() and Close() methods are used to start and stop the client.
// Send() are used to send the command to the server, and SetName() and SendMessage() are wrapper methods to set the display name and send the chat message.
//  Finally, the Incoming() method returns a channel for retrieving chat messages from the server.

// The structure of the client and its constructor can be defined as follows.
// The following code has some private variables that consider conn for connections and reader/writer as wrapper constructs for sending commands.

type TcpChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	name      string
	incoming  chan protocol.MessageCommand
}

func NewClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
	}
}

// Most of the methods are quite simple. Dial establishes a connection with the server,
//  and then the reader and writer of the protocol are created.

func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err == nil {
		c.conn = conn
	}
	c.cmdReader = protocol.NewCommandReader(conn)
	c.cmdWriter = protocol.NewCommandWriter(conn)
	return err
}

// Then the Send() method uses the cmdWriter to send the command to the server.

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

// The most important client method is the Start() method, which listens for incoming messages and then returns them to the channel.

func (c *TcpChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error %v", err)
		}
		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}
