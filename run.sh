#!/bin/bash

if [[ -z "${IN_RTSP_URL}" ]]; then
  /app/cam_cp
else
  chmod o+w /dev/stdout
  su vlcuser -c  "/usr/bin/vlc  --intf dummy $IN_RTSP_URL --sout '#transcode{vcodec=MJPG,venc=ffmpeg{strict=1}}:standard{access=http{mime=multipart/x-mixed-replace;boundary=--7b3cc56e5f51db803f790dad720ed50a},mux=mpjpeg,dst=:8080/}' --daemon --file-logging --logfile=/dev/stdout"
  IN_MJPEG_ENABLE=true
  IN_MJPEG_URL="http://localhost:8080/"
  /app/cam_cp
fi
