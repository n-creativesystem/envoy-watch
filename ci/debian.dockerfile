FROM envoyproxy/envoy:v1.19-latest as envoy

FROM debian

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Tokyo

WORKDIR /etc/watch
COPY ci/docker-entrypoint.sh /
COPY ci/setting.yaml /etc/watch/
COPY --from=envoy /usr/local/bin/envoy /usr/local/bin/envoy
COPY --from=envoy /etc/envoy/envoy.yaml /etc/watch/envoy.yaml
COPY envoy-watch /etc/watch/

RUN apt-get update \
    && apt-get install --no-install-recommends  -y tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && apt-get clean \
    && rm  -rf /tmp/* /var/lib/apt/lists/*

RUN chmod +x /etc/watch/envoy-watch \
    && mv /etc/watch/envoy-watch /usr/local/bin/ \
    && chmod +x /docker-entrypoint.sh


ENTRYPOINT [ "/docker-entrypoint.sh" ]
CMD [ "envoy-watch", "watch", "-c", "/etc/watch/envoy.yaml", "-s", "/etc/watch/setting.yaml" ]