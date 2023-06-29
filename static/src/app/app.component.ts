import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { SignalingService } from './services/signaling.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

  @ViewChild('video', { static: true, read: ElementRef })
  public videoElem!: ElementRef<HTMLVideoElement>;

  @ViewChild('audio', { static: true, read: ElementRef })
  public audioElem!: ElementRef<HTMLAudioElement>;

  @ViewChild('start', { static: true, read: ElementRef })
  public startElem!: ElementRef<HTMLButtonElement>;

  private readonly pc: RTCPeerConnection;

  constructor(private signaling: SignalingService) {
    //this.pc = new RTCPeerConnection({iceServers: [{urls: 'stun:stun.l.google.com:19302'}]});
    this.pc = new RTCPeerConnection();

    // Once remote track media arrives, show it in remote element.
    this.pc.ontrack = (event) => {
      const [remoteStream] = event.streams;
      if (event.track.kind === 'video') {
        this.videoElem.nativeElement.srcObject = remoteStream;
        this.videoElem.nativeElement.muted = true;
        this.videoElem.nativeElement.play();
      } else if (event.track.kind === 'audio') {
        this.audioElem.nativeElement.srcObject = remoteStream;
      }
    }

    this.pc.oniceconnectionstatechange = (e) => {
      console.log(e);
    };

    // Send any ice candidates to the other peer.
    this.pc.onicecandidate = ({candidate}) => {
      console.log(candidate);
      if (candidate == null) {
        signaling.next(this.pc.localDescription)
      }
    };

    // Offer to receive 1 audio, and 1 video track
    this.pc.addTransceiver('audio', {direction: 'recvonly'})
    this.pc.addTransceiver('video', {direction: 'recvonly'})
    this.pc.createOffer().then(d => this.pc.setLocalDescription(d)).catch((e) => {
      console.error(e);
    });
    
    // Let the "negotiationneeded" event trigger offer generation.
    this.pc.onnegotiationneeded = async (e) => {
      console.log(e);
      // try {
      //   await this.pc.setLocalDescription();
      //   // Send the offer to the other peer.
      //   this.api.send({desc: this.pc.localDescription});
      // } catch (err) {
      //   console.error(err);
      // }
    };

    this.signaling.subscribe(async (x) => {
      console.log(x);
      if (x.type === 'answer') {
        console.log('Received answer');

        try {
          // the answer
          await this.pc.setRemoteDescription(x);
        } catch(e) {
          console.error(e);
        }

      }
    });
  }

  public async ngOnInit(): Promise<void> {
    
  }

  public run(e: MouseEvent): void {
    this.videoElem.nativeElement.play();
  }

}
