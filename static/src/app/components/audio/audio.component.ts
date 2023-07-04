import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';

@Component({
  selector: 'app-audio',
  template: '<audio #audio></audio>',
  styleUrls: ['./audio.component.scss']
})
export class AudioComponent implements OnInit {
  @ViewChild('audio', { static: true, read: ElementRef })
  public elem!: ElementRef<HTMLAudioElement>;

  ngOnInit(): void {
    this.elem.nativeElement.muted = true;
  }

  public setStream(stream: MediaStream): void {
    this.elem.nativeElement.srcObject = stream;
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

    if (!this.elem.nativeElement.muted) {
      this.play();
    }
  }
}
