FROM envoyproxy/envoy:v1.19-latest as envoy

FROM ncreativesystem/golang-developments:dev

RUN mkdir -p /etc/envoy
COPY --from=envoy /usr/local/bin/envoy /usr/local/bin/envoy
COPY --from=envoy /etc/envoy/envoy.yaml /etc/envoy/envoy.yaml
WORKDIR /src/

CMD [ "sh" ]