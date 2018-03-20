// microArch/frontend/socket/socket.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package socket

import (
	"errors"
	"fmt"

	"github.com/antha-lang/antha/microArch/logger"
)

type MessageType int

const (
	NONE = iota
	DISPLAY
	ACKNOWLEDGE
	INPUT
	SELECT
)

func (m MessageType) String() string {
	switch m {
	case NONE:
		return ""
	case DISPLAY:
		return "display"
	case ACKNOWLEDGE:
		return "acknowledge"
	case INPUT:
		return "input"
	case SELECT:
		return "select"
	}
	return ""
}
func (m *MessageType) UnmarshalJSON(value []byte) error {
	switch string(value) {
	case "display":
		*m = DISPLAY
	case "acknowledge":
		*m = ACKNOWLEDGE
	case "input":
		*m = INPUT
	case "select":
		*m = SELECT
	case "":
		*m = NONE
	default:
		return errors.New("Unknown data type")
	}
	return nil
}

func (m MessageType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, m.String())), nil
}

type Socket struct {
	ID string
	//TODO put here the reference to the monolith socket library
}

func NewSocket(id string) *Socket {
	ret := new(Socket)
	ret.ID = id
	return ret
}

type LocalizedMessage struct {
	Locale  string
	Message string
}

func NewLocalizedMessage(locale, message string) *LocalizedMessage {
	ret := new(LocalizedMessage)
	ret.Locale = locale
	ret.Message = message
	return ret
}

type LocalizedOption struct {
	ID       string
	Messages []LocalizedMessage
}

func NewLocalizedOption(id string, messages []LocalizedMessage) *LocalizedOption {
	ret := new(LocalizedOption)
	ret.ID = id
	ret.Messages = messages
	return ret
}

type Request struct {
	Type    MessageType
	Options []LocalizedOption
}

func NewRequest(messageType MessageType, options []LocalizedOption) *Request {
	ret := new(Request)
	ret.Type = messageType
	ret.Options = options
	return ret
}

type SocketMessage struct {
	Display []LocalizedMessage
	Request Request
}

func NewSocketMessage(display []LocalizedMessage, request Request) *SocketMessage {
	ret := new(SocketMessage)
	ret.Display = display
	ret.Request = request
	return ret
}

//buildMessageSocketMessage builds a socket message from the parameters of the message call
func buildMessageSocketMessage(message string) *SocketMessage {
	mes := make([]LocalizedMessage, 0)
	mes = append(mes, *NewLocalizedMessage("", message))
	return NewSocketMessage(mes, *NewRequest(DISPLAY, nil))
}

//Message sends a message through the socket and expects an ack answer.
func (sock *Socket) Message(message string) error {
	sock.sockSend(*buildMessageSocketMessage(message))
	return nil
	//TODO handle timeout errors
}

//buildOptionSocketMessage builds a socket message with the parameters of the option function
func buildOptionSocketMessage(message string, options map[string]string) *SocketMessage {
	opts := make([]LocalizedOption, 0)
	for k, v := range options {
		mess := make([]LocalizedMessage, 0)
		mess = append(mess, *NewLocalizedMessage("", v))
		opts = append(opts, *NewLocalizedOption(k, mess))
	}
	req := *NewRequest(SELECT, opts)
	mes := make([]LocalizedMessage, 0)
	mes = append(mes, *NewLocalizedMessage("", message))
	optionRequest := NewSocketMessage(mes, req)
	return optionRequest
}

//RequestOption sends a message through the socket and expects an answer, being that one
// of the given options as parameters.
// Options map key will be id for that option selection and the value will be the referred option text
func (sock *Socket) RequestOption(message string, options map[string]string) (string, error) {
	return sock.sockSend(*buildOptionSocketMessage(message, options))
}

//Message sends a message through the socket and expects an ack answer.
func (sock *Socket) RequestText(message string) (string, error) {
	req := *NewRequest(INPUT, nil)
	mes := make([]LocalizedMessage, 0)
	mes = append(mes, *NewLocalizedMessage("", message))
	return sock.sockSend(*NewSocketMessage(mes, req))
}

//sockSend sends the socketMessage on the underlying socket io connection and waits for a response
// It is not this func duty to unmarshal the response, it is however expected to be a json string
// that depends on the other endpoint. Errors reported by this method should be derived
// from communication problems, thus, not from the data being collected.
func (sock *Socket) sockSend(mess SocketMessage) (string, error) {
	//TODO implement with underlying socket.io connection
	ll := fmt.Sprintf("socket %s... Sending request: %v.", sock.ID, mess)
	logger.Info(ll)
	return "", nil
}
