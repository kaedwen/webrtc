import { Component, ComponentRef, HostListener, OnInit, ViewChild, ViewContainerRef } from '@angular/core';
import { VideoComponent } from './components/video/video.component';
import { AudioComponent } from './components/audio/audio.component';
import { SignalingService } from './services/signaling.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  @ViewChild("vc", {static: true, read: ViewContainerRef }) 
  public vcr!: ViewContainerRef;

  private audioList: ComponentRef<AudioComponent>[] = [];
  private selfAudioRunning = false;

  private readonly pc: RTCPeerConnection;

  @HostListener('document:click', ['$event.target'])
  async onClick(_: KeyboardEvent) {
    if (!this.selfAudioRunning) {
      await this.startAudio()
    }
  
    for (const el of this.audioList) {
      el.instance.toggleMute();
    }
  }

  private async startAudio() {
    const stream = await navigator.mediaDevices.getUserMedia({
      audio: true,
      video: false
    });

    this.pc.addTrack(stream.getAudioTracks()[0])

    this.selfAudioRunning = true;
  }

  constructor(private signaling: SignalingService) {
    //this.pc = new RTCPeerConnection({iceServers: [{urls: 'stun:stun.l.google.com:19302'}]});
    this.pc = new RTCPeerConnection();

    // Once remote track media arrives, show it in remote element.
    this.pc.ontrack = (event) => {
      const [remoteStream] = event.streams;
      if (event.track.kind === 'video') {
        const video = this.vcr.createComponent(VideoComponent, {});
        video.instance.setStream(remoteStream);
        video.instance.play();
        console.log(video);
      } else if (event.track.kind === 'audio') {
        const audio = this.vcr.createComponent(AudioComponent, {});
        audio.instance.setStream(remoteStream);
        console.log(audio);

        this.audioList.push(audio)
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
    this.pc.addTransceiver('audio', {direction: 'sendrecv'})
    this.pc.addTransceiver('video', {direction: 'recvonly'})
    this.pc.createOffer().then(d => this.pc.setLocalDescription(d)).catch((e) => {
      console.error(e);
    });
    
    // Let the "negotiationneeded" event trigger offer generation.
    this.pc.onnegotiationneeded = async (e) => {
      try {
        await this.pc.setLocalDescription();
        signaling.next(this.pc.localDescription);
      } catch (err) {
        console.error(err);
      }
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

}
