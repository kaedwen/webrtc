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
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/core */ 2560);
/* harmony import */ var _services_signaling_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./services/signaling.service */ 1759);




const _c0 = ["video"];
const _c1 = ["audio"];
const _c2 = ["start"];
class AppComponent {
  constructor(signaling) {
    var _this = this;
    this.signaling = signaling;
    //this.pc = new RTCPeerConnection({iceServers: [{urls: 'stun:stun.l.google.com:19302'}]});
    this.pc = new RTCPeerConnection();
    // Once remote track media arrives, show it in remote element.
    this.pc.ontrack = event => {
      const [remoteStream] = event.streams;
      if (event.track.kind === 'video') {
        this.videoElem.nativeElement.srcObject = remoteStream;
        this.videoElem.nativeElement.muted = true;
        this.videoElem.nativeElement.play();
      } else if (event.track.kind === 'audio') {
        this.audioElem.nativeElement.srcObject = remoteStream;
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
      direction: 'recvonly'
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
        console.log(e);
        // try {
        //   await this.pc.setLocalDescription();
        //   // Send the offer to the other peer.
        //   this.api.send({desc: this.pc.localDescription});
        // } catch (err) {
        //   console.error(err);
        // }
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
            yield _this.pc.setRemoteDescription(x);
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
  run(e) {
    this.videoElem.nativeElement.play();
  }
}
AppComponent.ɵfac = function AppComponent_Factory(t) {
  return new (t || AppComponent)(_angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵdirectiveInject"](_services_signaling_service__WEBPACK_IMPORTED_MODULE_1__.SignalingService));
};
AppComponent.ɵcmp = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵdefineComponent"]({
  type: AppComponent,
  selectors: [["app-root"]],
  viewQuery: function AppComponent_Query(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵviewQuery"](_c0, 7, _angular_core__WEBPACK_IMPORTED_MODULE_2__.ElementRef);
      _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵviewQuery"](_c1, 7, _angular_core__WEBPACK_IMPORTED_MODULE_2__.ElementRef);
      _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵviewQuery"](_c2, 7, _angular_core__WEBPACK_IMPORTED_MODULE_2__.ElementRef);
    }
    if (rf & 2) {
      let _t;
      _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵqueryRefresh"](_t = _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵloadQuery"]()) && (ctx.videoElem = _t.first);
      _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵqueryRefresh"](_t = _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵloadQuery"]()) && (ctx.audioElem = _t.first);
      _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵqueryRefresh"](_t = _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵloadQuery"]()) && (ctx.startElem = _t.first);
    }
  },
  decls: 4,
  vars: 0,
  consts: [["autoplay", "", "muted", "", "controls", ""], ["video", ""], ["autoplay", "", "controls", ""], ["audio", ""]],
  template: function AppComponent_Template(rf, ctx) {
    if (rf & 1) {
      _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵelement"](0, "video", 0, 1)(2, "audio", 2, 3);
    }
  },
  styles: ["[_nghost-%COMP%] {\n  display: flex;\n  align-items: center;\n  justify-content: center;\n  flex-direction: column;\n  height: 100%;\n}\n\nvideo[_ngcontent-%COMP%] {\n  max-height: 90%;\n  max-width: 80%;\n  height: 90%;\n  width: auto;\n}\n/*# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly8uL3NyYy9hcHAvYXBwLmNvbXBvbmVudC5zY3NzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUNBO0VBQ0UsYUFBQTtFQUNBLG1CQUFBO0VBQ0EsdUJBQUE7RUFDQSxzQkFBQTtFQUNBLFlBQUE7QUFBRjs7QUFHQTtFQUNFLGVBQUE7RUFDQSxjQUFBO0VBQ0EsV0FBQTtFQUNBLFdBQUE7QUFBRiIsInNvdXJjZXNDb250ZW50IjpbIlxuOmhvc3Qge1xuICBkaXNwbGF5OiBmbGV4O1xuICBhbGlnbi1pdGVtczogY2VudGVyO1xuICBqdXN0aWZ5LWNvbnRlbnQ6IGNlbnRlcjtcbiAgZmxleC1kaXJlY3Rpb246IGNvbHVtbjtcbiAgaGVpZ2h0OiAxMDAlO1xufVxuXG52aWRlbyB7XG4gIG1heC1oZWlnaHQ6IDkwJTtcbiAgbWF4LXdpZHRoOiA4MCU7XG4gIGhlaWdodDogOTAlO1xuICB3aWR0aDogYXV0bztcbn1cbiJdLCJzb3VyY2VSb290IjoiIn0= */"]
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
/* harmony import */ var _angular_platform_browser__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! @angular/platform-browser */ 4497);
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @angular/common/http */ 8987);
/* harmony import */ var _app_component__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./app.component */ 5041);
/* harmony import */ var _services_signaling_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./services/signaling.service */ 1759);
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/core */ 2560);





class AppModule {}
AppModule.ɵfac = function AppModule_Factory(t) {
  return new (t || AppModule)();
};
AppModule.ɵmod = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵdefineNgModule"]({
  type: AppModule,
  bootstrap: [_app_component__WEBPACK_IMPORTED_MODULE_0__.AppComponent]
});
AppModule.ɵinj = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵdefineInjector"]({
  providers: [_services_signaling_service__WEBPACK_IMPORTED_MODULE_1__.SignalingService],
  imports: [_angular_platform_browser__WEBPACK_IMPORTED_MODULE_3__.BrowserModule, _angular_common_http__WEBPACK_IMPORTED_MODULE_4__.HttpClientModule]
});
(function () {
  (typeof ngJitMode === "undefined" || ngJitMode) && _angular_core__WEBPACK_IMPORTED_MODULE_2__["ɵɵsetNgModuleScope"](AppModule, {
    declarations: [_app_component__WEBPACK_IMPORTED_MODULE_0__.AppComponent],
    imports: [_angular_platform_browser__WEBPACK_IMPORTED_MODULE_3__.BrowserModule, _angular_common_http__WEBPACK_IMPORTED_MODULE_4__.HttpClientModule]
  });
})();

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