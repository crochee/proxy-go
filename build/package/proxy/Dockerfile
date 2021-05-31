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
WORKDIR /opt/cloud/proxy
RUN groupadd -r -g 1000 service && useradd -r -g service -u 1000 serviceUser
RUN mkdir -p /opt/cloud/proxy/conf /opt/cloud/proxy/data
COPY --from=builder /opt/cloud/proxy-go/proxy .
COPY --from=builder /opt/cloud/proxy-go/conf ./conf
COPY --from=builder /opt/cloud/proxy-go/gosu .
COPY --from=builder /opt/cloud/proxy-go/build/package/proxy/entrypoint.sh .

# 加入用户组
RUN chown -R service:serviceUser .
# 更改权限
RUN chmod a-rwx,u+rx ./proxy ./gosu ./entrypoint.sh &&\
    chmod -R a-rwx,u+r ./conf &&\
    chmod -R a-rwx,u+rw ./data

VOLUME /opt/cloud/proxy/data

ENTRYPOINT ["entrypoint.sh"]

EXPOSE 8120
STOPSIGNAL SIGINT
CMD ["proxy"]