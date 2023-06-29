import { WebSocketSubject, WebSocketSubjectConfig } from 'rxjs/webSocket';
import { Injectable } from '@angular/core';
import { v4 } from 'uuid';

@Injectable({
  providedIn: 'root'
})
export class SignalingService extends WebSocketSubject<any> {
  constructor() {
    const config: WebSocketSubjectConfig<any> = {
      url: `${location.origin.replace("http", "ws")}/signaling/${v4()}`,
    };

    super(config);
  }
}
