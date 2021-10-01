# Envoy File Based Dynamic Routing

Config mapを使用してEnvoy File Based Dynamic Routingを実現します。

## 概要

アーキテクチャとしては、

+----------+   +--------------+   +----------------+
|configmap ----> temporary file ----> envoy config |
+----------+   +--------------+   +----------------+

configmapでマウントされているファイルの変更を監視し、変更があった場合
envoyで使用するconfigの一時ファイルを生成します。
envoyのFile based dynamic routingはファイルのmvを監視しているので
一時ファイルから実際のconfigファイルへmvさせてenvoyの自動更新につなげています。

設定ファイル内の`files`キーに設定されている内容をマージして
`output`キーに設定されているファイル名.tmpを一時ファイルで出力し、設定されたファイル名にリネームします。

`files`キーで指定されたyamlファイルについては
`anchors`キーで設定された内容に関してはAnchorとして定義でき、Aliasとして利用できます。
※出力ファイルには記述されません。
また、複数ファイルを定義した際には同じキー項目内容は結合されて出力されます。

## 設定ファイル

| キー          | 型       | 概要               |
| :------------ | :------- | :----------------- |
| settings      | object[] | 設定内容を記述     |
| &ensp; output | string   | 出力するファイル名 |
| &ensp; files  | string[] | 監視するファイル名 |

### Example

setting.yaml

```yaml
settings:
  - output: envoy/config/cds.yaml
    files:
      - data/cds1.yaml
      - data/cds2.yaml
  - output: envoy/config/lds.yaml
    files:
      - data/lds1.yaml
      - data/lds2.yaml
```

以下、envoyの設定ファイル例

/etc/envoy/envoy.yaml

```yaml
node:
  id: example
  cluster: example

dynamic_resources:
  cds_config:
    resource_api_version: V3
    path: envoy/config/cds.yaml
  lds_config:
    resource_api_version: V3
    path: envoy/config/lds.yaml

admin:
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
```

data/cds1.yaml

```yaml
resources:
  - name: backend
    "@type": type.googleapis.com/envoy.config.cluster.v3.Cluster
    lb_policy: ROUND_ROBIN
    type: STRICT_DNS
    dns_lookup_family: V4_ONLY
    load_assignment:
      cluster_name: backend
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: localhost
                    port_value: 80
```

data/cds2.yaml

```yaml
resources:
  - name: backend2
    "@type": type.googleapis.com/envoy.config.cluster.v3.Cluster
    lb_policy: ROUND_ROBIN
    type: STRICT_DNS
    dns_lookup_family: V4_ONLY
    load_assignment:
      cluster_name: backend2
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: localhost
                    port_value: 443
```

data/lds1.yaml

```yaml
anchors:
  upgrade_configs: &upgrade_configs
    - upgrade_type: websocket
  access_log: &access_log
    - name: extensions.access_loggers.file.v3.FileAccessLog.format
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
        path: "/dev/stdout"
        log_format:
          json_format:
            start_time: "%START_TIME%"
            method: "%REQ(:METHOD)%"
            path: "%REQ(:PATH)%"
            status: "%RESPONSE_CODE%"
            flag: "%RESPONSE_FLAGS%"
            message: "%LOCAL_REPLY_BODY%"
resources:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 8080
    "@type": type.googleapis.com/envoy.config.listener.v3.Listener
    name: backend
    filter_chains:
      - filters:
          - name: envoy.filters.network.http_connection_manager
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
              codec_type: AUTO
              access_log: *access_log
              stat_prefix: ingress_http
              upgrade_configs: *upgrade_configs
              route_config:
                name: local_route
                virtual_hosts:
                  - name: local_route
                    domains:
                      - "*"
                    routes:
                      - match:
                          prefix: /
                        route:
                          cluster: backend
              http_filters:
                - name: envoy.filters.http.router
              use_remote_address: true
```

data/lds2.yaml

```yaml
anchors:
  upgrade_configs: &upgrade_configs
    - upgrade_type: websocket
  access_log: &access_log
    - name: extensions.access_loggers.file.v3.FileAccessLog.format
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
        path: "/dev/stdout"
        log_format:
          json_format:
            start_time: "%START_TIME%"
            method: "%REQ(:METHOD)%"
            path: "%REQ(:PATH)%"
            status: "%RESPONSE_CODE%"
            flag: "%RESPONSE_FLAGS%"
            message: "%LOCAL_REPLY_BODY%"

resources:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 8443
    "@type": type.googleapis.com/envoy.config.listener.v3.Listener
    name: grpc_backend
    filter_chains:
      - filters:
          - name: envoy.filters.network.http_connection_manager
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
              codec_type: AUTO
              access_log: *access_log
              stat_prefix: ingress_http
              upgrade_configs: *upgrade_configs
              route_config:
                name: local_route
                virtual_hosts:
                  - name: local_route
                    domains:
                      - "*"
                    routes:
                      - match:
                          prefix: /
                        route:
                          cluster: backend2
              http_filters:
                - name: envoy.filters.http.router
              use_remote_address: true
```
