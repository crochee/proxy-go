FROM golang1.16.2:latest as builder
ARG URL=https://OnlyOneFace:lcf13149825.lmy@gitee.com/OnlyOneFace/proxy-go.git
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd proxy-go &&\
    go mod tidy &&\
    go build -ldflags="-s -w" -o proxy -tags jsoniter cmd/proxy

FROM centos:latest
WORKDIR /opt/cloud/proxy-go
RUN mkdir -p /opt/cloud/obs/conf
COPY --from=builder /opt/cloud/proxy-go/proxy .
COPY --from=builder /opt/cloud/obs/conf ./conf
COPY --from=builder /opt/cloud/obs/.build/obs/entrypoint.sh .
RUN chmod +x ./proxy && chmod +x ./entrypoint.sh
VOLUME /obs

EXPOSE 8120
STOPSIGNAL 2

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./proxy"]