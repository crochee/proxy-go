FROM golang1.16.2:latest as builder
ARG URL=https://OnlyOneFace:lcf13149825.lmy@gitee.com/OnlyOneFace/proxy-go.git
WORKDIR /opt/cloud
RUN git clone ${URL} &&\
    cd proxy-go &&\
    go mod tidy &&\
    go build -ldflags="-s -w" -o proxy -tags jsoniter cmd/proxy
RUN go install github.com/tianon/gosu &&\
    cp ${GO_PATH}/bin/gosu .


FROM centos:latest
WORKDIR /opt/cloud/proxy-go
RUN mkdir -p /opt/cloud/proxy-go/conf
COPY --from=builder /opt/cloud/proxy-go/proxy .
COPY --from=builder /opt/cloud/proxy-go/conf ./conf
COPY --from=builder /opt/cloud/proxy-go/gosu .
COPY --from=builder /opt/cloud/proxy-go/build/package/proxy/entrypoint.sh .

ENV PATH=/opt/cloud/proxy-go:$PATH
RUN chmod +x gosu proxy entrypoint.sh
VOLUME /proxy

EXPOSE 8120
STOPSIGNAL SIGINT

ENTRYPOINT ["entrypoint.sh"]
CMD ["proxy","server"]