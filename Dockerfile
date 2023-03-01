FROM golang:latest

RUN apt-get update \
  && apt-get install -y wget gcc build-essential mosquitto mosquitto-clients
RUN wget --directory-prefix=/usr/lib https://github.com/project-alvarium/alvarium-sdk-go/raw/main/internal/iota/include/libiota_streams_c.so

COPY . ./app
WORKDIR app

ENV DEVICEID 88906a10-d4e1-4d1e-8e74-b1fdc8239992

RUN make build
CMD ["make" , "run_mqtt"]