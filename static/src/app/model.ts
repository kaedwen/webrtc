
export interface SignalingMessage {
  type: 'new-ice-candidate' | 'offer' | 'answer';
  data: any;
}

export interface AnswerMessage extends SignalingMessage {
  data: RTCSessionDescriptionInit;
}

export interface OfferMessage extends SignalingMessage {
  data: RTCSessionDescriptionInit;
}

export interface IceCandidateMessage extends SignalingMessage {
  data: RTCIceCandidate;
}

export const IsSignalingMessage = (d: any): d is SignalingMessage => {
  return !!d && typeof(d.type) === 'string';
}

export const IsIceCandidate = (d: any): d is IceCandidateMessage => {
  return IsSignalingMessage(d) && d.type === 'new-ice-candidate';
}

export const IsAnswer = (d: any): d is AnswerMessage => {
  return IsSignalingMessage(d) && d.type === 'answer';
}

export const IsOffer = (d: any): d is OfferMessage => {
  return IsSignalingMessage(d) && d.type === 'offer';
}
