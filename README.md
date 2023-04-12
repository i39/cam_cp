#GoCam

GoCam is an open-source project written in Go that aims to provide an easy-to-use solution for processing images obtained from CCTV cameras. It supports retrieving images via RTSP protocol or retrieving JPEG images through URLs. The images are then filtered through different filters, such as YOLOv4 or DeepStack, to detect certain objects. The filtered images can then be sent through different output methods, such as email, Telegram, or saved to a folder.

##Installation

To install GoCam, you should first install Go on your system. Then, you can clone the repository:

bash
Copy code
git clone https://github.com/<username>/GoCam.git
Next, you should build the project by running:

'''
make docker
'''

For using deppstack or CPU based filter, or

'''
make docker_yolo_gpu
'''

for using GPU based Yolo detection filter



GoCam can be configured with flags or environment wariables, that specifies the input, filter, and output parameters. Here is an example of docker-compose.yml:

'''
services:
  cam_cp:
    image: i39.in/cam_cp_gpu:master
    restart: always
    deploy:
      resources:
        reservations:
          devices:
            - capabilities: [gpu]
    container_name: some_camera
    hostname: some_host
    logging:
      options:
        max-size: "10m"
        max-file: "1"
    volumes:
      - ./tmp:/tmp
      - ./yolo:/yolo
    environment:
      - TZ=Europe/Warsaw
      - IN_RTSP_URL=rtsp://user:password@192.168.1.1/video2
      - OUT_FILE_ENABLE=false
      - OUT_FILE_DIR=/tmp/cam1
      - FILTER_YOLO_ENABLE=true
      - FILTER_YOLO_CONFIG=/yolo/yolov4-tiny.cfg
      - FILTER_YOLO_WEIGHTS=/yolo/yolov4-tiny.weights
      - FILTER_YOLO_LABELS=person,car
      - FILTER_YOLO_PROBABILITY=75
      - OUT_TELEGRAM_ENABLE=false
      - OUT_TELEGRAM_TOKEN=5125940686:AAFL4gswL9S--0yoloA9odpeiyDURO_c68U
      - OUT_TELEGRAM_CHAT_ID=2778603
      - OUT_EMAIL_ENABLE=true
      - OUT_EMAIL_HOST=192.168.1.1
      - OUT_EMAIL_PORT=1025
      - OUT_EMAIL_FROM=some_camera@some_host.local
      - OUT_EMAIL_TO=some_camera@some_host.local
      - OUT_EMAIL_SUBJECT=Alarm motion
      - DEBUG=true

'''

You can customize the file to specify the input, filter, and output methods that you want to use. Once you have created the configuration file, you can run GoCam using the following command:

You can check all possible flags by use --help flag.

Contributing

If you find a bug or want to contribute to GoCam, feel free to open an issue or submit a pull request. We welcome any suggestions, bug fixes, or improvements.

License

GoCam is licensed under the MIT License. See the LICENSE file for more information.
