package ipc

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/allankerr/freighter/log"
)

const delimeter = byte('\n')

const (
	FifoName = "FIFO_FD"
	InitName = "INIT_FD"
	LogName  = "LOG_FD"
)

type envelope struct {
	Type    MessageType `json:"type"`
	Payload []byte      `json:"payload"`
}

type Pipe struct {
	file     *os.File
	receiver chan *envelope
}

func NewPipe(file *os.File) *Pipe {
	pipe := &Pipe{file, make(chan *envelope)}
	pipe.listen()
	return pipe
}

func (p *Pipe) listen() {

	go func() {
		reader := bufio.NewReader(p.file)
		for {
			line, err := reader.ReadString(delimeter)
			if err != nil {
				if err == io.EOF {
					log.Warn("Pipe listener encountered EOF")
				} else {
					log.WithError(err).Fatal("Failed to read string")
				}
			}
			data := []byte(line)
			envelope := &envelope{}
			if err := json.Unmarshal(data, envelope); err != nil {
				log.WithError(err).Fatal("Failed to unmarshal envelope")
			}
			p.receiver <- envelope
		}
	}()
}

func (p *Pipe) Send(messageType MessageType, payload interface{}) error {

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	envelope := &envelope{
		messageType,
		payloadData,
	}
	envelopeData, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	messageData := append(envelopeData, delimeter)
	_, err = p.file.Write(messageData)
	return err
}

func (p *Pipe) Receive(v interface{}) (MessageType, error) {
	envelope := <-p.receiver
	if err := json.Unmarshal(envelope.Payload, v); err != nil {
		return MessageUnknown, err
	}
	return envelope.Type, nil
}
