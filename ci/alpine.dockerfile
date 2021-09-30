FROM envoyproxy/envoy-alpine:v1.19-latest as envoy

FROM golang:1.16-alpine as build
WORKDIR /src
COPY . .
RUN go mod download \
    && go build -o bin/watch

FROM frolvlad/alpine-glibc:alpine-3.14_glibc-2.33

ENV TZ=Asia/Tokyo

WORKDIR /etc/watch
COPY ci/docker-entrypoint.sh /
COPY ci/setting.yaml /etc/watch/
COPY --from=envoy /usr/local/bin/envoy /usr/local/bin/envoy
COPY --from=envoy /etc/envoy/envoy.yaml /etc/watch/envoy.yaml

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && rm  -rf /tmp/* /var/cache/apk/*


COPY --from=build /src/bin/watch /etc/watch/

RUN chmod +x /etc/watch/watch \
    && mv /etc/watch/watch /usr/local/bin/ \
    && chmod +x /docker-entrypoint.sh


ENTRYPOINT [ "/docker-entrypoint.sh" ]
CMD [ "watch", "watch", "-c", "/etc/watch/envoy.yaml", "-s", "/etc/watch/setting.yaml" ]