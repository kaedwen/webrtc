package server

import (
	"encoding/json"

	"github.com/pion/webrtc/v3"
)

type SignalingMessageType = string

const (
	MessageTypeIceCandidate SignalingMessageType = "new-ice-candidate"
	MessageTypeAnswer       SignalingMessageType = "answer"
	MessageTypeOffer        SignalingMessageType = "offer"
)

// INCOMING

type IncomingSignalingMessage struct {
	Type SignalingMessageType `json:"type"`
	Data json.RawMessage      `json:"data"`
}

type IceCandidateMessage struct {
	*IncomingSignalingMessage
	Candidate webrtc.ICECandidateInit
}

type OfferMessage struct {
	*IncomingSignalingMessage
	Offer webrtc.SessionDescription
}

type AnswerMessage struct {
	*IncomingSignalingMessage
	Answer webrtc.SessionDescription
}

func (m *IncomingSignalingMessage) IsIceCandidateMessage() bool {
	return m.Type == MessageTypeIceCandidate
}

func (m *IncomingSignalingMessage) IsAnswerMessage() bool {
	return m.Type == MessageTypeAnswer
}

func (m *IncomingSignalingMessage) IsOfferMessage() bool {
	return m.Type == MessageTypeOffer
}

func (m *IncomingSignalingMessage) ToIceCandidateMessage() (*IceCandidateMessage, error) {
	nm := IceCandidateMessage{
		IncomingSignalingMessage: m,
	}

	return &nm, json.Unmarshal(m.Data, &nm.Candidate)
}

func (m *IncomingSignalingMessage) ToAnswerMessage() (*AnswerMessage, error) {
	nm := AnswerMessage{
		IncomingSignalingMessage: m,
	}

	return &nm, json.Unmarshal(m.Data, &nm.Answer)
}

func (m *IncomingSignalingMessage) ToOfferMessage() (*OfferMessage, error) {
	nm := OfferMessage{
		IncomingSignalingMessage: m,
	}

	return &nm, json.Unmarshal(m.Data, &nm.Offer)
}

// OUTGOING

type OutgoingSignalingMessage struct {
	Type SignalingMessageType `json:"type"`
	Data any                  `json:"data"`
}

func NewIceCandidateMessage(candidate webrtc.ICECandidate) *OutgoingSignalingMessage {
	return &OutgoingSignalingMessage{
		Type: MessageTypeIceCandidate,
		Data: candidate,
	}
}

func NewAnswerMessage(answer *webrtc.SessionDescription) *OutgoingSignalingMessage {
	return &OutgoingSignalingMessage{
		Type: MessageTypeAnswer,
		Data: answer,
	}
}

func NewOfferMessage(offer *webrtc.SessionDescription) *OutgoingSignalingMessage {
	return &OutgoingSignalingMessage{
		Type: MessageTypeOffer,
		Data: offer,
	}
}
