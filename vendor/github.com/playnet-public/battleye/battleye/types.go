package battleye

// Type is the representation of a packets type
type Type byte

// Login is the type used when sending logins to the server
var Login Type = 0x00

// LoginOk is being returned by the server on successful login
var LoginOk Type = 0x01

// LoginFail is being returned by the server on invalid credentials
var LoginFail Type = 0x00

// Command is the default packet type on commands
var Command Type = 0x01

// MultiCommand is the packet type when the received packet is not complete yet
// and will be received in multiple pieces
var MultiCommand Type = 0x00

// ServerMessage is the packet type when the server is sending events to clients
var ServerMessage Type = 0x02
