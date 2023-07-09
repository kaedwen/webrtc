# Webrtc streamer based on golang (pion) and gstreamer

This project demonstrates who to stream any source (thanks to gstreamer) in a webrtc session using golang pion to a webrtc client (angular).

## build locally

### Build web client
To build the web client static sources make sure you have `node` and `npm` installed.
```
make build-static
```

### Build service
To build the service make sure you have build the `static` content first.
```
make build
```

## build cross
To build for `arm` or `arm64` just use the sections in the `Makefile`. Make sure you have `podman` installed as it will be used to create the tmp build container.
```
make build-armhf
```