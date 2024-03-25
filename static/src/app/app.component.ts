import { Component, ComponentRef, HostListener, OnInit, ViewChild, ViewContainerRef } from '@angular/core';
import { VideoComponent } from './components/video/video.component';
import { AudioComponent } from './components/audio/audio.component';
import { SignalingService } from './services/signaling.service';
import { IsAnswer, IsIceCandidate, IsOffer } from './model';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  @ViewChild("vc", { static: true, read: ViewContainerRef })
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
    this.pc.ontrack = (e) => {
      const [remoteStream] = e.streams;
      if (e.track.kind === 'video') {
        const video = this.vcr.createComponent(VideoComponent, {});
        video.instance.setStream(remoteStream);
        video.instance.play();
        console.log(video);
      } else if (e.track.kind === 'audio') {
        const audio = this.vcr.createComponent(AudioComponent, {});
        audio.instance.setStream(remoteStream);
        console.log(audio);

        this.audioList.push(audio)
      }
    }

    this.pc.oniceconnectionstatechange = (e) => {
      console.log('ICE: state change', e);
    };

    this.pc.onicegatheringstatechange = (e) => {
      switch (this.pc.iceGatheringState) {
        case "new":
          console.log('ICE: gathering is either just starting or has been reset', e);
          break;
        case "gathering":
          console.log('ICE: gathering has begun or is ongoing', e);
          break;
        case "complete":
          console.log('ICE: gathering has ended', e);
          break;
      }
    }

    // Send any ice candidates to the other peer.
    this.pc.onicecandidate = ({ candidate }) => {
      console.log('ICE: candidate', candidate);
      if (candidate !== null) {
        signaling.SendCandidate(candidate);
      } else {
        /* there are no more candidates coming during this negotiation */
      }
    };

    // Offer to receive 1 audio, and 1 video track
    this.pc.addTransceiver('audio', { direction: 'sendrecv' })
    this.pc.addTransceiver('video', { direction: 'recvonly' })

    // Let the "negotiationneeded" event trigger offer generation.
    this.pc.onnegotiationneeded = (e) => {
      this.pc.createOffer()
        .then(async (d) => {
          await this.pc.setLocalDescription(d)
          return this.pc.localDescription!
        })
        .then((d) => {
          console.log('SDP: sending offer', d);
          signaling.SendOffer(d)
        })
        .catch((e) => {
          console.error(e);
        });
    };

    this.signaling.subscribe((m) => {
      switch (true) {
        case IsIceCandidate(m):
          console.log('ICE: received new candidate', m.data);
          break;
        case IsOffer(m):
          console.log('SDP: received offer', m.data);
          break;
        case IsAnswer(m):
          console.log('SDP: received answer', m.data);

          this.pc.setRemoteDescription(m.data).catch((e) => {
            console.error(e);
          });

          break;
        default:
          console.log('WARN: received unknown data', m);
      }
    });
  }

  public async ngOnInit(): Promise<void> {

  }

}
