FROM golang:1.21 as builder

RUN cd /opt && git clone https://github.com/voodoo-sys/calibre-api && cd ./calibre-api && CGO_ENABLED=0 go build -ldflags="-s -w"

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

ENV TZ=Europe/Prague
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install -y --no-install-recommends wget cron supervisor calibre

RUN apt-get clean autoclean; \
    apt-get autoremove --yes; \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

RUN mkdir -p /opt/calibre-api/files
COPY --from=builder /opt/calibre-api/calibre-api /opt/calibre-api/calibre-api
RUN chmod 0777 /opt/calibre-api/calibre-api && chmod 0777 /opt/calibre-api/files

RUN rm -rf /etc/cron.d
COPY cron.d /etc/cron.d
RUN chmod 0644 /etc/cron.d/*

COPY supervisord.conf /etc/supervisor/

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
