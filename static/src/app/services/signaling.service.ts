import { WebSocketSubject, WebSocketSubjectConfig } from 'rxjs/webSocket';
import { SignalingMessage } from '../model';
import { Injectable } from '@angular/core';
import { v4 } from 'uuid';

@Injectable({
  providedIn: 'root'
})
export class SignalingService extends WebSocketSubject<SignalingMessage> {
  constructor() {
    const config: WebSocketSubjectConfig<any> = {
      url: `${location.origin.replace("http", "ws")}/signaling/${v4()}`,
    };

    super(config);
  }

  public SendCandidate(candidate: RTCIceCandidate) {
    this.next({
      type: 'new-ice-candidate',
      data: candidate,
    });
  }

  public SendOffer(offer: RTCSessionDescriptionInit) {
    this.next({
      type: 'offer',
      data: offer,
    });
  }
}
