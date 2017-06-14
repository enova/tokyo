package alert

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/ipv4"

	"github.com/enova/tokyo/src/cfg"
)

/////////////////////////
// Message Packet Format
// ---------------------
//
// Each emitted packet contains three parts:
//
// A. 5-Byte Ascii, Zero-Padded Message-Length (Not Including These 5 Bytes)
// B. Variable-Length Payload
// C. Terminating Newline ('\n')
//
//
// For example:
//
// 00018{"name": "hello"}\n
//
// The message length is 18 = 17 + 1 for both the payload and the
// newline terminator
//
/////////////////////////

// MaxMulticastPerHour limits the number of messages sent to Multicast
const MaxMulticastPerHour int = 100

// MaxTextLength limits the length of text sent in the Text field of the emitted packet
const MaxTextLength int = 128

// Globals: Multicast
var (
	multicastClient   *Multicaster
	multicastThrottle *Throttle
)

// For Testing
var lastMulticastMsg string

// Multicaster provides functionality for emitting multicast messages
type Multicaster struct {
	client *ipv4.PacketConn
	dest   *net.UDPAddr
}

// NewMulticaster returns a new Multicaster object
func NewMulticaster(ip string, port, ttl int) *Multicaster {

	// Group
	group := net.ParseIP(ip)

	// Create Packet-Connection
	conn, err := net.ListenPacket("udp4", "")
	if err != nil {
		Exit(err)
	}

	packetConn := ipv4.NewPacketConn(conn)

	// Multicaster
	multicaster := &Multicaster{
		client: packetConn,
		dest:   &net.UDPAddr{IP: group, Port: port},
	}
	multicaster.SetTTL(ttl)

	return multicaster
}

// Emit sends a message over multicast
func (m *Multicaster) Emit(message []byte) {
	if _, err := m.client.WriteTo(message, nil, m.dest); err != nil {
		Warn(err, Whisper)
	}
}

// SetTTL sets the time to live for future packets
func (m *Multicaster) SetTTL(amount int) {
	m.client.SetTTL(amount)
	m.client.SetMulticastTTL(amount)
}

// Message ...
type multicastMessage struct {
	Meta
	Time  time.Time
	Level string
	Text  string
}

// Configure Multicast
func setMulticast(cfg *cfg.Config) {

	// Set Port
	if !cfg.Has("Port") {
		log.Fatal("Alert.Multicast.Port missing from config")
	}
	port, err := strconv.Atoi(cfg.Get("Port"))

	if err != nil {
		log.Fatal("Alert.Multicast.Port must be an integer")
	}

	// Set Group
	if !cfg.Has("Group") {
		log.Fatal("Alert.Multicast.Group missing from config")
	}
	group := cfg.Get("Group")

	// Set TTL
	ttl := 0
	if cfg.Has("TTL") {
		ttl, _ = strconv.Atoi(cfg.Get("TTL"))
	}

	// Create Multicast-Throttle
	multicastThrottle = NewThrottle(MaxMulticastPerHour)
	multicastClient = NewMulticaster(group, port, ttl)
}

// Send Message Over Multicast
func sendToMulticast(msg *Message) {

	// Limit Text-Length
	var text string
	if len(msg.Text) <= MaxTextLength {
		text = msg.Text
	} else {
		text = msg.Text[:MaxTextLength]
	}

	// For Testing
	lastMulticastMsg = text

	// Multicast enabled?
	if multicastClient == nil {
		return
	}

	// Check Multicast-Throttle
	if ok := multicastThrottle.Update(time.Now()); !ok {
		return
	}

	// Create an easy to marshal version of the message
	multicastMessage := &multicastMessage{
		Meta:  msg.Meta,
		Time:  msg.Now,
		Level: msg.Level.String(),
		Text:  text,
	}

	// Marshal the message
	jsonMsg, err := json.Marshal(multicastMessage)

	// Errored: Whisper A Warning
	if err != nil {
		Warn("Could not marshal the message to json", err, Whisper)
	}

	// Bytes To Emit
	bytes := make([]byte, 5+len(jsonMsg)+1)

	// Copy Message-Length, Payload, Terminating-Newline
	msgLen := fmt.Sprintf("%05d", len(jsonMsg))
	copy(bytes, msgLen)
	copy(bytes[5:], jsonMsg)
	bytes[5+len(jsonMsg)] = byte('\n')

	multicastClient.Emit(bytes)
}
