import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { SignalingService } from './services/signaling.service';
import { VideoComponent } from './components/video/video.component';
import { AudioComponent } from './components/audio/audio.component';

@NgModule({
  declarations: [
    AppComponent,
    VideoComponent,
    AudioComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
  ],
  providers: [SignalingService],
  bootstrap: [AppComponent]
})
export class AppModule { }
