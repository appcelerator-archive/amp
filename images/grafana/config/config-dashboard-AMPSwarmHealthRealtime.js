{
  "Dashboard": {
    "id": null,
    "title": "AMP Swarm Health Realtime",
    "tags": [],
    "style": "dark",
    "timezone": "browser",
    "editable": true,
    "hideControls": false,
    "sharedCrosshair": false,
    "rows": [
      {
        "collapse": false,
        "editable": true,
        "height": "",
        "panels": [],
        "showTitle": false,
        "title": "DockerLevel"
      },
      {
        "collapse": false,
        "editable": true,
        "height": "250px",
        "panels": [
          {
            "aliasColors": {},
            "bars": false,
            "datasource": null,
            "editable": true,
            "error": false,
            "fill": 1,
            "grid": {
              "threshold1": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2": null,
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "id": 1,
            "legend": {
              "avg": false,
              "current": false,
              "max": false,
              "min": false,
              "show": true,
              "total": false,
              "values": false
            },
            "lines": true,
            "linewidth": 2,
            "links": [],
            "nullPointMode": "connected",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "alias": "[[tag_com.docker.swarm.service.name]]",
                "dsType": "influxdb",
                "fields": [
                  {
                    "func": "mean",
                    "name": "usage_percent"
                  }
                ],
                "groupBy": [
                  {
                    "interval": "auto",
                    "params": [
                      "auto"
                    ],
                    "type": "time"
                  },
                  {
                    "key": "datacenter",
                    "params": [
                      "datacenter"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "host",
                    "params": [
                      "host"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "container_name",
                    "params": [
                      "/^$ServiceName$/"
                    ],
                    "type": "tag"
                  }
                ],
                "hide": false,
                "measurement": "docker_container_mem",
                "policy": "default",
                "query": "SELECT mean(\"usage_percent\") as usage FROM \"docker_container_mem\" WHERE \"com.docker.swarm.service.name\" =~ /$ServiceName/  and \"datacenter\" =~ /$DataCenter/ and \"host\" =~ /$HostName/ and $timeFilter GROUP BY time($interval), \"datacenter\", \"host\", \"com.docker.swarm.service.name\"",
                "rawQuery": true,
                "refId": "A",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "usage_percent"
                      ],
                      "type": "field"
                    },
                    {
                      "params": [],
                      "type": "mean"
                    }
                  ]
                ],
                "tags": []
              },
              {
                "policy": "default",
                "dsType": "influxdb",
                "resultFormat": "time_series",
                "tags": [],
                "groupBy": [
                  {
                    "type": "time",
                    "params": [
                      "$interval"
                    ]
                  },
                  {
                    "type": "tag",
                    "params": [
                      "engine_host"
                    ]
                  }
                ],
                "select": [
                  [
                    {
                      "type": "field",
                      "params": [
                        "usage_percent"
                      ]
                    },
                    {
                      "type": "sum",
                      "params": []
                    }
                  ]
                ],
                "refId": "B",
                "measurement": "docker_container_mem",
                "alias": "total_[[tag_engine_host]]",
                "query": "SELECT sum(\"usage_percent\") FROM \"docker_container_mem\" WHERE \"datacenter\" =~ /$DataCenter/ and \"host\" =~ /$HostName/ and $timeFilter GROUP BY time($interval), \"engine_host\"",
                "rawQuery": true
              }
            ],
            "timeFrom": null,
            "timeShift": null,
            "title": "AMP Memory Utilization",
            "tooltip": {
              "msResolution": true,
              "shared": true,
              "value_type": "cumulative",
              "sort": 0
            },
            "type": "graph",
            "xaxis": {
              "show": true
            },
            "yaxes": [
              {
                "format": "short",
                "label": "Percentage",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              }
            ]
          },
          {
            "aliasColors": {},
            "bars": false,
            "datasource": null,
            "editable": true,
            "error": false,
            "fill": 1,
            "grid": {
              "threshold1": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2": null,
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "id": 5,
            "legend": {
              "avg": false,
              "current": false,
              "max": false,
              "min": false,
              "show": true,
              "total": false,
              "values": false
            },
            "lines": true,
            "linewidth": 2,
            "links": [],
            "nullPointMode": "connected",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "alias": "[[tag_com.docker.swarm.service.name]]",
                "dsType": "influxdb",
                "fields": [
                  {
                    "func": "mean",
                    "name": "usage_percent"
                  }
                ],
                "groupBy": [
                  {
                    "interval": "auto",
                    "params": [
                      "auto"
                    ],
                    "type": "time"
                  },
                  {
                    "key": "datacenter",
                    "params": [
                      "datacenter"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "host",
                    "params": [
                      "host"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "container_name",
                    "params": [
                      "com.docker.swarm.service.name"
                    ],
                    "type": "tag"
                  }
                ],
                "hide": false,
                "measurement": "docker_container_cpu",
                "policy": "default",
                "query": "SELECT max(\"usage_percent\") FROM \"docker_container_cpu\" WHERE \"com.docker.swarm.service.name\" =~ /$ServiceName/  and \"datacenter\" =~ /$DataCenter/ and \"host\" =~ /$HostName/ and $timeFilter GROUP BY time($interval), \"datacenter\", \"host\", \"com.docker.swarm.service.name\"",
                "rawQuery": true,
                "refId": "A",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "type": "field",
                      "params": [
                        "usage_percent"
                      ]
                    },
                    {
                      "type": "max",
                      "params": []
                    }
                  ]
                ],
                "tags": []
              },
              {
                "policy": "default",
                "dsType": "influxdb",
                "resultFormat": "time_series",
                "tags": [],
                "groupBy": [
                  {
                    "type": "time",
                    "params": [
                      "$interval"
                    ]
                  },
                  {
                    "type": "tag",
                    "params": [
                      "engine_host"
                    ]
                  },
                  {
                    "type": "fill",
                    "params": [
                      "null"
                    ]
                  }
                ],
                "select": [
                  [
                    {
                      "type": "field",
                      "params": [
                        "usage_percent"
                      ]
                    },
                    {
                      "type": "sum",
                      "params": []
                    }
                  ]
                ],
                "refId": "B",
                "measurement": "docker_container_cpu",
                "alias": "total_[[tag_engine_host]]",
                "query": "SELECT sum(\"usage_percent\") FROM \"docker_container_cpu\" WHERE  \"datacenter\" =~ /$DataCenter/ and \"host\" =~ /$HostName/ and $timeFilter GROUP BY time($interval), \"engine_host\" fill(null)",
                "rawQuery": true
              }
            ],
            "timeFrom": null,
            "timeShift": null,
            "title": "AMP CPU Utilization",
            "tooltip": {
              "msResolution": true,
              "shared": true,
              "value_type": "cumulative",
              "sort": 0
            },
            "type": "graph",
            "xaxis": {
              "show": true
            },
            "yaxes": [
              {
                "format": "short",
                "label": "Percentage",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              }
            ]
          },
          {
            "aliasColors": {},
            "bars": false,
            "datasource": null,
            "editable": true,
            "error": false,
            "fill": 1,
            "grid": {
              "threshold1": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2": null,
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "id": 6,
            "legend": {
              "avg": false,
              "current": false,
              "max": false,
              "min": false,
              "show": true,
              "total": false,
              "values": false
            },
            "lines": true,
            "linewidth": 2,
            "links": [],
            "nullPointMode": "connected",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "alias": "[[tag_com.docker.swarm.service.name]]",
                "dsType": "influxdb",
                "fields": [
                  {
                    "func": "mean",
                    "name": "usage_percent"
                  }
                ],
                "groupBy": [
                  {
                    "interval": "auto",
                    "params": [
                      "auto"
                    ],
                    "type": "time"
                  },
                  {
                    "key": "datacenter",
                    "params": [
                      "datacenter"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "host",
                    "params": [
                      "host"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "container_name",
                    "params": [
                      "container_name"
                    ],
                    "type": "tag"
                  }
                ],
                "hide": false,
                "measurement": "docker_container_blkio",
                "policy": "default",
                "query": "SELECT non_negative_derivative(last(\"io_service_bytes_recursive_total\"))/1000 FROM \"docker_container_blkio\" WHERE  \"com.docker.swarm.service.name\" =~ /$ServiceName/ and \"datacenter\" =~ /$DataCenter/  and \"host\" =~ /$HostName/ and $timeFilter GROUP BY time($interval), \"datacenter\", \"host\", \"com.docker.swarm.service.name\"",
                "rawQuery": true,
                "refId": "A",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "usage_percent"
                      ],
                      "type": "field"
                    },
                    {
                      "params": [],
                      "type": "mean"
                    }
                  ]
                ],
                "tags": []
              },
              {
                "policy": "default",
                "dsType": "influxdb",
                "resultFormat": "time_series",
                "tags": [],
                "groupBy": [
                  {
                    "type": "time",
                    "params": [
                      "$interval"
                    ]
                  },
                  {
                    "type": "fill",
                    "params": [
                      "null"
                    ]
                  }
                ],
                "select": [
                  [
                    {
                      "type": "field",
                      "params": [
                        "value"
                      ]
                    },
                    {
                      "type": "mean",
                      "params": []
                    }
                  ]
                ],
                "refId": "B",
                "query": "SELECT non_negative_derivative(last(\"io_service_bytes_recursive_total\"))/1000 FROM \"docker_container_blkio\" WHERE  \"datacenter\" =~ /$DataCenter/  and \"host\" =~ /$HostName/ and $timeFilter GROUP BY time($interval), \"datacenter\", \"engine_host\"",
                "rawQuery": true,
                "alias": "total_[[tag_engine_host]]"
              }
            ],
            "timeFrom": null,
            "timeShift": null,
            "title": "AMP Block I/O Utilization",
            "tooltip": {
              "msResolution": false,
              "shared": true,
              "value_type": "cumulative",
              "sort": 0
            },
            "type": "graph",
            "xaxis": {
              "show": true
            },
            "yaxes": [
              {
                "format": "short",
                "label": "Mega Bytes",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              }
            ]
          },
          {
            "aliasColors": {},
            "bars": false,
            "datasource": null,
            "editable": true,
            "error": false,
            "fill": 1,
            "grid": {
              "threshold1": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2": null,
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "id": 7,
            "legend": {
              "avg": false,
              "current": false,
              "max": false,
              "min": false,
              "show": true,
              "total": false,
              "values": false
            },
            "lines": true,
            "linewidth": 2,
            "links": [],
            "nullPointMode": "connected",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "alias": "[[tag_com.docker.swarm.service.name]]:rx_bytes",
                "dsType": "influxdb",
                "fields": [
                  {
                    "func": "mean",
                    "name": "usage_percent"
                  }
                ],
                "groupBy": [
                  {
                    "interval": "auto",
                    "params": [
                      "auto"
                    ],
                    "type": "time"
                  },
                  {
                    "key": "datacenter",
                    "params": [
                      "datacenter"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "host",
                    "params": [
                      "host"
                    ],
                    "type": "tag"
                  },
                  {
                    "key": "container_name",
                    "params": [
                      "container_name"
                    ],
                    "type": "tag"
                  }
                ],
                "hide": false,
                "measurement": "docker_container_net",
                "policy": "default",
                "query": "SELECT non_negative_derivative(last(\"rx_bytes\"))/1000 FROM \"docker_container_net\" WHERE \"com.docker.swarm.service.name\" =~ /$ServiceName/ and \"datacenter\" =~ /$DataCenter/  and \"host\" =~ /$HostName/ and $timeFilter GROUP BY time($interval), \"datacenter\", \"engine_host\", \"com.docker.swarm.service.name\"",
                "rawQuery": true,
                "refId": "A",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "usage_percent"
                      ],
                      "type": "field"
                    },
                    {
                      "params": [],
                      "type": "mean"
                    }
                  ]
                ],
                "tags": []
              },
              {
                "alias": "[[tag_com.docker.swarm.service.name]]:tx_bytes",
                "dsType": "influxdb",
                "groupBy": [
                  {
                    "params": [
                      "$interval"
                    ],
                    "type": "time"
                  },
                  {
                    "params": [
                      "null"
                    ],
                    "type": "fill"
                  }
                ],
                "policy": "default",
                "query": "SELECT non_negative_derivative(last(\"tx_bytes\"))/1000 FROM \"docker_container_net\" WHERE \"com.docker.swarm.service.name\" =~ /$ServiceName/ and \"datacenter\" =~ /$DataCenter/  and \"host\" =~ /$HostName/ and  $timeFilter GROUP BY time($interval), \"datacenter\", \"engine_host\", \"com.docker.swarm.service.name\"",
                "rawQuery": true,
                "refId": "B",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "value"
                      ],
                      "type": "field"
                    },
                    {
                      "params": [],
                      "type": "mean"
                    }
                  ]
                ],
                "tags": []
              }
            ],
            "timeFrom": null,
            "timeShift": null,
            "title": "AMP Network Utilization",
            "tooltip": {
              "msResolution": false,
              "shared": true,
              "value_type": "cumulative",
              "sort": 0
            },
            "type": "graph",
            "xaxis": {
              "show": true
            },
            "yaxes": [
              {
                "format": "short",
                "label": "Mega Bytes",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              }
            ]
          }
        ],
        "title": "New row"
      }
    ],
    "time": {
      "from": "now-5m",
      "to": "now"
    },
    "timepicker": {
      "collapse": false,
      "enable": true,
      "notice": false,
      "now": true,
      "refresh_intervals": [
        "5s",
        "10s",
        "30s",
        "1m",
        "5m",
        "15m",
        "30m",
        "1h",
        "2h",
        "1d"
      ],
      "status": "Stable",
      "time_options": [
        "5m",
        "15m",
        "1h",
        "6h",
        "12h",
        "24h",
        "2d",
        "7d",
        "30d"
      ],
      "type": "timepicker"
    },
    "templating": {
      "list": [
        {
          "type": "query",
          "datasource": null,
          "refresh": 1,
          "name": "ServiceName",
          "hide": 0,
          "options": [
            {
              "text": "All",
              "value": "$__all",
              "selected": true
            },
            {
              "text": "amp-agent",
              "value": "amp-agent",
              "selected": false
            },
            {
              "text": "amp-log-worker",
              "value": "amp-log-worker",
              "selected": false
            },
            {
              "text": "amp-ui",
              "value": "amp-ui",
              "selected": false
            },
            {
              "text": "amplifier",
              "value": "amplifier",
              "selected": false
            },
            {
              "text": "elasticsearch",
              "value": "elasticsearch",
              "selected": false
            },
            {
              "text": "etcd",
              "value": "etcd",
              "selected": false
            },
            {
              "text": "grafana",
              "value": "grafana",
              "selected": false
            },
            {
              "text": "haproxy",
              "value": "haproxy",
              "selected": false
            },
            {
              "text": "influxdb",
              "value": "influxdb",
              "selected": false
            },
            {
              "text": "kapacitor",
              "value": "kapacitor",
              "selected": false
            },
            {
              "text": "nats",
              "value": "nats",
              "selected": false
            },
            {
              "text": "registry",
              "value": "registry",
              "selected": false
            },
            {
              "text": "telegraf-agent",
              "value": "telegraf-agent",
              "selected": false
            },
            {
              "text": "telegraf-haproxy",
              "value": "telegraf-haproxy",
              "selected": false
            }
          ],
          "includeAll": true,
          "multi": true,
          "query": "SHOW TAG VALUES FROM \"docker_container_mem\" WITH KEY = \"com.docker.swarm.service.name\"",
          "current": {
            "text": "All",
            "value": [
              "$__all"
            ],
            "tags": []
          },
          "regex": "",
          "label": "ServiceName"
        },
        {
          "current": {
            "text": "All",
            "value": "$__all",
            "selected": false
          },
          "datasource": null,
          "hide": 0,
          "includeAll": true,
          "label": "DataCenter",
          "multi": true,
          "name": "DataCenter",
          "options": [
            {
              "text": "All",
              "value": "$__all",
              "selected": true
            },
            {
              "text": "dc1",
              "value": "dc1",
              "selected": false
            }
          ],
          "query": "SHOW TAG VALUES FROM \"docker_container_mem\" WITH KEY = \"datacenter\"",
          "refresh": 1,
          "regex": "/([^/]*$)/",
          "type": "query"
        },
        {
          "current": {
            "value": [
              "$__all"
            ],
            "text": "All"
          },
          "datasource": null,
          "hide": 0,
          "includeAll": true,
          "label": "HostName",
          "multi": true,
          "name": "HostName",
          "options": [
            {
              "text": "All",
              "value": "$__all",
              "selected": true
            },
            {
              "text": "f0181ec14413",
              "value": "f0181ec14413",
              "selected": false
            }
          ],
          "query": "SHOW TAG VALUES FROM \"docker_container_mem\" WITH KEY = \"host\"",
          "refresh": 1,
          "regex": "/([^/]*$)/",
          "type": "query"
        }
      ]
    },
    "annotations": {
      "list": []
    },
    "schemaVersion": 12,
    "version": 2,
    "links": [],
    "gnetId": null
  }
}
