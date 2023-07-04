"use strict";
(self["webpackChunkwebrtc"] = self["webpackChunkwebrtc"] || []).push([["main"],{

/***/ 5041:
/*!**********************************!*\
  !*** ./src/app/app.component.ts ***!
  \**********************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "AppComponent": () => (/* binding */ AppComponent)
/* harmony export */ });
/* harmony import */ var _projects_PRIVATE_webrtc_gst_static_node_modules_babel_runtime_helpers_esm_asyncToGenerator_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./node_modules/@babel/runtime/helpers/esm/asyncToGenerator.js */ 1670);
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @angular/core */ 2560);
/* harmony import */ var _components_video_video_component__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./components/video/video.component */ 5465);
/* harmony import */ var _components_audio_audio_component__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./components/audio/audio.component */ 210);
/* harmony import */ var _services_signaling_service__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./services/signaling.service */ 1759);






const _c0 = ["vc"];
class AppComponent {
  onClick(_) {
    var _this = this;
    return (0,_projects_PRIVATE_webrtc_gst_static_node_modules_babel_runtime_helpers_esm_asyncToGenerator_js__WEBPACK_IMPORTED_MODULE_0__["default"])(function* () {
      if (!_this.selfAudioRunning) {
        yield _this.startAudio();
      }
      for (const el of _this.audioList) {
        el.instance.toggleMute();
      }
    })();
  }
  startAudio() {
    var _this2 = this;
    return (0,_projects_PRIVATE_webrtc_gst_static_node_modules_babel_runtime_helpers_esm_asyncToGenerator_js__WEBPACK_IMPORTED_MODULE_0__["default"])(function* () {
      const stream = yield navigator.mediaDevices.getUserMedia({
        audio: true,
        video: false
      });
      _this2.pc.addTrack(stream.getAudioTracks()[0]);
      _this2.selfAudioRunning = true;
    })();
  }
  constructor(signaling) {
    var _this3 = this;
    this.signaling = signaling;
    this.audioList = [];
    this.selfAudioRunning = false;
    //this.pc = new RTCPeerConnection({iceServers: [{urls: 'stun:stun.l.google.com:19302'}]});
    this.pc = new RTCPeerConnection();
    // Once remote track media arrives, show it in remote element.
    this.pc.ontrack = event => {
      const [remoteStream] = event.streams;
      if (event.track.kind === 'video') {
        const video = this.vcr.createComponent(_components_video_video_component__WEBPACK_IMPORTED_MODULE_1__.VideoComponent, {});
        video.instance.setStream(remoteStream);
        video.instance.play();
        console.log(video);
      } else if (event.track.kind === 'audio') {
        const audio = this.vcr.createComponent(_components_audio_audio_component__WEBPACK_IMPORTED_MODULE_2__.AudioComponent, {});
        audio.instance.setStream(remoteStream);
        console.log(audio);
        this.audioList.push(audio);
      }
    };
    this.pc.oniceconnectionstatechange = e => {
      console.log(e);
    };
    // Send any ice candidates to the other peer.
    this.pc.onicecandidate = ({
      candidate
    }) => {
      console.log(candidate);
      if (candidate == null) {
        signaling.next(this.pc.localDescription);
      }
    };
    // Offer to receive 1 audio, and 1 video track
    this.pc.addTransceiver('audio', {
      direction: 'sendrecv'
    });
    this.pc.addTransceiver('video', {
      direction: 'recvonly'
    });
    this.pc.createOffer().then(d => this.pc.setLocalDescription(d)).catch(e => {
      console.error(e);
    });
    // Let the "negotiationneeded" event trigger offer generation.
    this.pc.onnegotiationneeded = /*#__PURE__*/function () {
      var _ref = (0,_projects_PRIVATE_webrtc_gst_static_node_modules_babel_runtime_helpers_esm_asyncToGenerator_js__WEBPACK_IMPORTED_MODULE_0__["default"])(function* (e) {
        try {
          yield _this3.pc.setLocalDescription();
          signaling.next(_this3.pc.localDescription);
        } catch (err) {
          console.error(err);
        }
      });
      return function (_x) {
        return _ref.apply(this, arguments);
      };
    }();
    this.signaling.subscribe( /*#__PURE__*/function () {
      var _ref2 = (0,_projects_PRIVATE_webrtc_gst_static_node_modules_babel_runtime_helpers_esm_asyncToGenerator_js__WEBPACK_IMPORTED_MODULE_0__["default"])(function* (x) {
        console.log(x);
        if (x.type === 'answer') {
          console.log('Received answer');
          try {
            // the answer
            yield _this3.pc.setRemoteDescription(x);
          } catch (e) {
            console.error(e);
          }
        }
      });
      return function (_x2) {
        return _ref2.apply(this, arguments);
      };
    }());
  }
  ngOnInit() {
    return (0,_projects_PRIVATE_webrtc_gst_static_node_modules_babel_runtime_helpers_esm_asyncToGenerator_js__WEBPACK_IMPORTED_MODULE_0__["default"])(function* () {})();
  }
}
AppComponent.ɵfac = function AppComponent_Factory(t) {
  return new (t || AppComponent)(_angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵdirectiveInject"](_services_signaling_service__WEBPACK_IMPORTED_MODULE_3__.SignalingService));
};
AppComponent.ɵcmp = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵdefineComponent"]({
  type: AppComponent,
  selectors: [["app-root"]],
  viewQuery: function AppComponent_Query(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵviewQuery"](_c0, 7, _angular_core__WEBPACK_IMPORTED_MODULE_4__.ViewContainerRef);
    }
    if (rf & 2) {
      let _t;
      _angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵqueryRefresh"](_t = _angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵloadQuery"]()) && (ctx.vcr = _t.first);
    }
  },
  hostBindings: function AppComponent_HostBindings(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵlistener"]("click", function AppComponent_click_HostBindingHandler($event) {
        return ctx.onClick($event.target);
      }, false, _angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵresolveDocument"]);
    }
  },
  decls: 2,
  vars: 0,
  consts: [["vc", ""]],
  template: function AppComponent_Template(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵelementContainer"](0, null, 0);
    }
  },
  styles: ["[_nghost-%COMP%] {\n  display: flex;\n  align-items: center;\n  justify-content: center;\n  flex-direction: column;\n  height: 100%;\n  width: 100%;\n}\n/*# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly8uL3NyYy9hcHAvYXBwLmNvbXBvbmVudC5zY3NzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUNBO0VBQ0UsYUFBQTtFQUNBLG1CQUFBO0VBQ0EsdUJBQUE7RUFDQSxzQkFBQTtFQUNBLFlBQUE7RUFDQSxXQUFBO0FBQUYiLCJzb3VyY2VzQ29udGVudCI6WyJcbjpob3N0IHtcbiAgZGlzcGxheTogZmxleDtcbiAgYWxpZ24taXRlbXM6IGNlbnRlcjtcbiAganVzdGlmeS1jb250ZW50OiBjZW50ZXI7XG4gIGZsZXgtZGlyZWN0aW9uOiBjb2x1bW47XG4gIGhlaWdodDogMTAwJTtcbiAgd2lkdGg6IDEwMCU7XG59XG5cbiJdLCJzb3VyY2VSb290IjoiIn0= */"]
});

/***/ }),

/***/ 6747:
/*!*******************************!*\
  !*** ./src/app/app.module.ts ***!
  \*******************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "AppModule": () => (/* binding */ AppModule)
/* harmony export */ });
/* harmony import */ var _angular_platform_browser__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! @angular/platform-browser */ 4497);
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! @angular/common/http */ 8987);
/* harmony import */ var _app_component__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./app.component */ 5041);
/* harmony import */ var _services_signaling_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./services/signaling.service */ 1759);
/* harmony import */ var _components_video_video_component__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./components/video/video.component */ 5465);
/* harmony import */ var _components_audio_audio_component__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./components/audio/audio.component */ 210);
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @angular/core */ 2560);







class AppModule {}
AppModule.ɵfac = function AppModule_Factory(t) {
  return new (t || AppModule)();
};
AppModule.ɵmod = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵdefineNgModule"]({
  type: AppModule,
  bootstrap: [_app_component__WEBPACK_IMPORTED_MODULE_0__.AppComponent]
});
AppModule.ɵinj = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵdefineInjector"]({
  providers: [_services_signaling_service__WEBPACK_IMPORTED_MODULE_1__.SignalingService],
  imports: [_angular_platform_browser__WEBPACK_IMPORTED_MODULE_5__.BrowserModule, _angular_common_http__WEBPACK_IMPORTED_MODULE_6__.HttpClientModule]
});
(function () {
  (typeof ngJitMode === "undefined" || ngJitMode) && _angular_core__WEBPACK_IMPORTED_MODULE_4__["ɵɵsetNgModuleScope"](AppModule, {
    declarations: [_app_component__WEBPACK_IMPORTED_MODULE_0__.AppComponent, _components_video_video_component__WEBPACK_IMPORTED_MODULE_2__.VideoComponent, _components_audio_audio_component__WEBPACK_IMPORTED_MODULE_3__.AudioComponent],
    imports: [_angular_platform_browser__WEBPACK_IMPORTED_MODULE_5__.BrowserModule, _angular_common_http__WEBPACK_IMPORTED_MODULE_6__.HttpClientModule]
  });
})();

/***/ }),

/***/ 210:
/*!*****************************************************!*\
  !*** ./src/app/components/audio/audio.component.ts ***!
  \*****************************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "AudioComponent": () => (/* binding */ AudioComponent)
/* harmony export */ });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ 2560);


const _c0 = ["audio"];
class AudioComponent {
  ngOnInit() {
    this.elem.nativeElement.muted = true;
  }
  setStream(stream) {
    this.elem.nativeElement.srcObject = stream;
  }
  play() {
    this.elem.nativeElement.play();
  }
  pause() {
    this.elem.nativeElement.pause();
  }
  toggleMute() {
    this.elem.nativeElement.muted = !this.elem.nativeElement.muted;
    console.log('toggle mute to', this.elem.nativeElement.muted);
    if (!this.elem.nativeElement.muted) {
      this.play();
    }
  }
}
AudioComponent.ɵfac = function AudioComponent_Factory(t) {
  return new (t || AudioComponent)();
};
AudioComponent.ɵcmp = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵdefineComponent"]({
  type: AudioComponent,
  selectors: [["app-audio"]],
  viewQuery: function AudioComponent_Query(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵviewQuery"](_c0, 7, _angular_core__WEBPACK_IMPORTED_MODULE_0__.ElementRef);
    }
    if (rf & 2) {
      let _t;
      _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵqueryRefresh"](_t = _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵloadQuery"]()) && (ctx.elem = _t.first);
    }
  },
  decls: 2,
  vars: 0,
  consts: [["audio", ""]],
  template: function AudioComponent_Template(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵelement"](0, "audio", null, 0);
    }
  },
  styles: ["\n/*# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IiIsInNvdXJjZVJvb3QiOiIifQ== */"]
});

/***/ }),

/***/ 5465:
/*!*****************************************************!*\
  !*** ./src/app/components/video/video.component.ts ***!
  \*****************************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "VideoComponent": () => (/* binding */ VideoComponent)
/* harmony export */ });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ 2560);


const _c0 = ["video"];
class VideoComponent {
  setStream(stream) {
    this.elem.nativeElement.srcObject = stream;
    this.elem.nativeElement.muted = true;
  }
  play() {
    this.elem.nativeElement.play();
  }
  pause() {
    this.elem.nativeElement.pause();
  }
  toggleMute() {
    this.elem.nativeElement.muted = !this.elem.nativeElement.muted;
    console.log('toggle mute to', this.elem.nativeElement.muted);
  }
}
VideoComponent.ɵfac = function VideoComponent_Factory(t) {
  return new (t || VideoComponent)();
};
VideoComponent.ɵcmp = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵdefineComponent"]({
  type: VideoComponent,
  selectors: [["app-video"]],
  viewQuery: function VideoComponent_Query(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵviewQuery"](_c0, 7, _angular_core__WEBPACK_IMPORTED_MODULE_0__.ElementRef);
    }
    if (rf & 2) {
      let _t;
      _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵqueryRefresh"](_t = _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵloadQuery"]()) && (ctx.elem = _t.first);
    }
  },
  decls: 2,
  vars: 0,
  consts: [["video", ""]],
  template: function VideoComponent_Template(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_0__["ɵɵelement"](0, "video", null, 0);
    }
  },
  styles: ["[_nghost-%COMP%] {\n  max-height: 90%;\n  max-width: 80%;\n  height: 90%;\n  width: auto;\n}\n\nvideo[_ngcontent-%COMP%] {\n  height: 100%;\n  width: 100%;\n}\n/*# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly8uL3NyYy9hcHAvY29tcG9uZW50cy92aWRlby92aWRlby5jb21wb25lbnQuc2NzcyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFDQTtFQUNFLGVBQUE7RUFDQSxjQUFBO0VBQ0EsV0FBQTtFQUNBLFdBQUE7QUFBRjs7QUFHQTtFQUNFLFlBQUE7RUFDQSxXQUFBO0FBQUYiLCJzb3VyY2VzQ29udGVudCI6WyJcbjpob3N0IHtcbiAgbWF4LWhlaWdodDogOTAlO1xuICBtYXgtd2lkdGg6IDgwJTtcbiAgaGVpZ2h0OiA5MCU7XG4gIHdpZHRoOiBhdXRvO1xufVxuXG52aWRlbyB7XG4gIGhlaWdodDogMTAwJTtcbiAgd2lkdGg6IDEwMCU7XG59Il0sInNvdXJjZVJvb3QiOiIifQ== */"]
});

/***/ }),

/***/ 1759:
/*!***********************************************!*\
  !*** ./src/app/services/signaling.service.ts ***!
  \***********************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "SignalingService": () => (/* binding */ SignalingService)
/* harmony export */ });
/* harmony import */ var rxjs_webSocket__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! rxjs/webSocket */ 3526);
/* harmony import */ var uuid__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! uuid */ 2535);
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/core */ 2560);



class SignalingService extends rxjs_webSocket__WEBPACK_IMPORTED_MODULE_0__.WebSocketSubject {
  constructor() {
    const config = {
      url: `${location.origin.replace("http", "ws")}/signaling/${(0,uuid__WEBPACK_IMPORTED_MODULE_1__["default"])()}`
    };
    super(config);
  }
}
SignalingService.ɵfac = function SignalingService_Factory(t) {
  return new (t || SignalingService)();
};
SignalingService.ɵprov = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵdefineInjectable"]({
  token: SignalingService,
  factory: SignalingService.ɵfac,
  providedIn: 'root'
});

/***/ }),

/***/ 2340:
/*!*****************************************!*\
  !*** ./src/environments/environment.ts ***!
  \*****************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "environment": () => (/* binding */ environment)
/* harmony export */ });
// This file can be replaced during build by using the `fileReplacements` array.
// `ng build` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.
const environment = {
  production: false
};
/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.

/***/ }),

/***/ 4431:
/*!*********************!*\
  !*** ./src/main.ts ***!
  \*********************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _angular_platform_browser__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! @angular/platform-browser */ 4497);
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/core */ 2560);
/* harmony import */ var _app_app_module__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./app/app.module */ 6747);
/* harmony import */ var _environments_environment__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./environments/environment */ 2340);




if (_environments_environment__WEBPACK_IMPORTED_MODULE_1__.environment.production) {
  (0,_angular_core__WEBPACK_IMPORTED_MODULE_2__.enableProdMode)();
}
_angular_platform_browser__WEBPACK_IMPORTED_MODULE_3__.platformBrowser().bootstrapModule(_app_app_module__WEBPACK_IMPORTED_MODULE_0__.AppModule).catch(err => console.error(err));

/***/ })

},
/******/ __webpack_require__ => { // webpackRuntimeModules
/******/ var __webpack_exec__ = (moduleId) => (__webpack_require__(__webpack_require__.s = moduleId))
/******/ __webpack_require__.O(0, ["vendor"], () => (__webpack_exec__(4431)));
/******/ var __webpack_exports__ = __webpack_require__.O();
/******/ }
]);
//# sourceMappingURL=main.js.map