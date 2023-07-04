import { Component, ElementRef, ViewChild } from '@angular/core';

@Component({
  selector: 'app-video',
  template: '<video #video></video>',
  styleUrls: ['./video.component.scss']
})
export class VideoComponent {
  @ViewChild('video', { static: true, read: ElementRef })
  public elem!: ElementRef<HTMLVideoElement>;

  public setStream(stream: MediaStream): void {
    this.elem.nativeElement.srcObject = stream;
    this.elem.nativeElement.muted = true;
  }

  public play(): void {
    this.elem.nativeElement.play()
  }

  public pause(): void {
    this.elem.nativeElement.pause()
  }

  public toggleMute(): void {
    this.elem.nativeElement.muted = !this.elem.nativeElement.muted;
    console.log('toggle mute to', this.elem.nativeElement.muted);
  }
}
