{
  "Dashboard": {
    "title": "Container stats",
    "annotations": {
      "list": []
    },
    "editable": true,
    "gnetId": null,
    "graphTooltip": 1,
    "hideControls": false,
    "links": [],
    "refresh": "10s",
    "rows": [
      {
        "collapse": false,
        "height": "250px",
        "panels": [
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "editable": true,
            "error": false,
            "format": "percent",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": true,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 4,
            "interval": null,
            "links": [],
            "mappingType": 1,
            "mappingTypes": [
              {
                "name": "value to text",
                "value": 1
              },
              {
                "name": "range to text",
                "value": 2
              }
            ],
            "maxDataPoints": 100,
            "nullPointMode": "connected",
            "nullText": null,
            "postfix": "",
            "postfixFontSize": "50%",
            "prefix": "",
            "prefixFontSize": "50%",
            "rangeMaps": [
              {
                "from": "null",
                "text": "N/A",
                "to": "null"
              }
            ],
            "span": 4,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": false
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "(sum(node_memory_MemTotal{instance=~\"$node\"}) - sum(node_memory_MemFree{instance=~\"$node\"}+node_memory_Buffers{instance=~\"$node\"}+node_memory_Cached{instance=~\"$node\"}) ) / sum(node_memory_MemTotal{instance=~\"$node\"}) * 100",
                "format": "time_series",
                "interval": "10s",
                "intervalFactor": 1,
                "legendFormat": "",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": "65, 90",
            "title": "Node memory usage",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "current"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "decimals": 2,
            "editable": true,
            "error": false,
            "format": "percent",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": true,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 6,
            "interval": null,
            "links": [],
            "mappingType": 1,
            "mappingTypes": [
              {
                "name": "value to text",
                "value": 1
              },
              {
                "name": "range to text",
                "value": 2
              }
            ],
            "maxDataPoints": 100,
            "nullPointMode": "connected",
            "nullText": null,
            "postfix": "",
            "postfixFontSize": "50%",
            "prefix": "",
            "prefixFontSize": "50%",
            "rangeMaps": [
              {
                "from": "null",
                "text": "N/A",
                "to": "null"
              }
            ],
            "span": 4,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": false
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum(sum by (container_name)( rate(container_cpu_usage_seconds_total{image!=\"\",instance=~\"$node\"}[1m] ) )) / count(node_cpu{mode=\"system\",instance=~\"$node\"}) * 100",
                "format": "time_series",
                "interval": "10s",
                "intervalFactor": 1,
                "legendFormat": "",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": "65, 90",
            "title": "Node CPU usage",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "current"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "decimals": 2,
            "editable": true,
            "error": false,
            "format": "percent",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": true,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 7,
            "interval": null,
            "links": [],
            "mappingType": 1,
            "mappingTypes": [
              {
                "name": "value to text",
                "value": 1
              },
              {
                "name": "range to text",
                "value": 2
              }
            ],
            "maxDataPoints": 100,
            "nullPointMode": "connected",
            "nullText": null,
            "postfix": "",
            "postfixFontSize": "50%",
            "prefix": "",
            "prefixFontSize": "50%",
            "rangeMaps": [
              {
                "from": "null",
                "text": "N/A",
                "to": "null"
              }
            ],
            "span": 4,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": false
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum (container_fs_limit_bytes{instance=~\"$node\"} - container_fs_usage_bytes{instance=~\"$node\"}) / sum(container_fs_limit_bytes{instance=~\"$node\"})",
                "format": "time_series",
                "interval": "10s",
                "intervalFactor": 1,
                "legendFormat": "",
                "metric": "",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": "65, 90",
            "title": "Node filesystem usage",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "current"
          }
        ],
        "repeat": null,
        "repeatIteration": null,
        "repeatRowId": null,
        "showTitle": false,
        "title": "Row",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": "250px",
        "panels": [
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "decimals": 3,
            "editable": true,
            "error": false,
            "fill": 0,
            "grid": {},
            "id": 3,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": false,
              "min": false,
              "rightSide": true,
              "show": true,
              "sort": "current",
              "sortDesc": true,
              "total": false,
              "values": true
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
            "spaceLength": 10,
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "sort_desc(sum(rate(container_cpu_user_seconds_total{image!=\"\",container_label_com_docker_swarm_service_name=~\"$service\",container_label_com_docker_stack_namespace=~\"$stack\",instance=~\"$node\"}[1m])) by (name))",
                "format": "time_series",
                "interval": "10s",
                "intervalFactor": 1,
                "legendFormat": "{{name}}",
                "metric": "container_cpu_user_seconds_total",
                "refId": "A",
                "step": 10
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container CPU usage",
            "tooltip": {
              "msResolution": true,
              "shared": true,
              "sort": 0,
              "value_type": "cumulative"
            },
            "type": "graph",
            "xaxis": {
              "buckets": null,
              "mode": "time",
              "name": null,
              "show": true,
              "values": []
            },
            "yaxes": [
              {
                "decimals": 1,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "label": null,
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              }
            ]
          }
        ],
        "repeat": null,
        "repeatIteration": null,
        "repeatRowId": null,
        "showTitle": false,
        "title": "New row",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": "250px",
        "panels": [
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "decimals": 2,
            "editable": true,
            "error": false,
            "fill": 0,
            "grid": {},
            "id": 2,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": false,
              "min": false,
              "rightSide": true,
              "show": true,
              "sideWidth": 200,
              "sort": "current",
              "sortDesc": true,
              "total": false,
              "values": true
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
            "spaceLength": 10,
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "sort_desc(sum(container_memory_usage_bytes{image!=\"\",container_label_com_docker_swarm_service_name=~\"$service\",container_label_com_docker_stack_namespace=~\"$stack\",instance=~\"$node\"}) by (name))",
                "format": "time_series",
                "interval": "10s",
                "intervalFactor": 1,
                "legendFormat": "{{name}}",
                "metric": "container_memory_usage:sort_desc",
                "refId": "A",
                "step": 10
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Memory Usage",
            "tooltip": {
              "msResolution": false,
              "shared": true,
              "sort": 0,
              "value_type": "cumulative"
            },
            "type": "graph",
            "xaxis": {
              "buckets": null,
              "mode": "time",
              "name": null,
              "show": true,
              "values": []
            },
            "yaxes": [
              {
                "format": "bytes",
                "label": null,
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "label": null,
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
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "decimals": 2,
            "editable": true,
            "error": false,
            "fill": 0,
            "grid": {},
            "id": 8,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": false,
              "min": false,
              "rightSide": true,
              "show": true,
              "sideWidth": 200,
              "sort": "current",
              "sortDesc": true,
              "total": false,
              "values": true
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
            "spaceLength": 10,
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "sort_desc(sum by (name) (rate(container_network_receive_bytes_total{image!=\"\",container_label_com_docker_swarm_service_name=~\"$service\",container_label_com_docker_stack_namespace=~\"$stack\",instance=~\"$node\"}[1m] ) ))",
                "format": "time_series",
                "interval": "10s",
                "intervalFactor": 1,
                "legendFormat": "{{name}}",
                "metric": "container_network_receive_bytes_total",
                "refId": "A",
                "step": 10
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Network Input",
            "tooltip": {
              "msResolution": false,
              "shared": true,
              "sort": 0,
              "value_type": "cumulative"
            },
            "type": "graph",
            "xaxis": {
              "buckets": null,
              "mode": "time",
              "name": null,
              "show": true,
              "values": []
            },
            "yaxes": [
              {
                "format": "bytes",
                "label": null,
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "label": null,
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
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "decimals": 2,
            "editable": true,
            "error": false,
            "fill": 0,
            "grid": {},
            "id": 9,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": false,
              "min": false,
              "rightSide": true,
              "show": true,
              "sideWidth": 200,
              "sort": "current",
              "sortDesc": true,
              "total": false,
              "values": true
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
            "spaceLength": 10,
            "span": 12,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "sort_desc(sum by (name) (rate(container_network_transmit_bytes_total{image!=\"\",container_label_com_docker_swarm_service_name=~\"$service\",container_label_com_docker_stack_namespace=~\"$stack\",instance=~\"$node\"}[1m] ) ))",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{name}}",
                "metric": "container_network_transmit_bytes_total",
                "refId": "B",
                "step": 2
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Network Output",
            "tooltip": {
              "msResolution": false,
              "shared": true,
              "sort": 0,
              "value_type": "cumulative"
            },
            "type": "graph",
            "xaxis": {
              "buckets": null,
              "mode": "time",
              "name": null,
              "show": true,
              "values": []
            },
            "yaxes": [
              {
                "format": "bytes",
                "label": null,
                "logBase": 1,
                "max": null,
                "min": null,
                "show": true
              },
              {
                "format": "short",
                "label": null,
                "logBase": 1,
                "max": null,
                "min": null,
                "show": false
              }
            ]
          }
        ],
        "repeat": null,
        "repeatIteration": null,
        "repeatRowId": null,
        "showTitle": false,
        "title": "New row",
        "titleSize": "h6"
      }
    ],
    "schemaVersion": 14,
    "style": "dark",
    "tags": [
      "docker",
      "service"
    ],
    "templating": {
      "list": [
        {
          "allValue": ".+",
          "current": {
            "text": "All",
            "value": "$__all"
          },
          "datasource": null,
          "hide": 0,
          "includeAll": true,
          "label": "Service",
          "multi": false,
          "name": "service",
          "options": [],
          "query": "label_values(container_cpu_user_seconds_total,container_label_com_docker_swarm_service_name)",
          "refresh": 2,
          "regex": "",
          "sort": 1,
          "tagValuesQuery": "",
          "tags": [],
          "tagsQuery": "",
          "type": "query",
          "useTags": false
        },
        {
          "allValue": ".+",
          "current": {
            "tags": [],
            "text": "All",
            "value": "$__all"
          },
          "datasource": null,
          "hide": 0,
          "includeAll": true,
          "label": "Stack",
          "multi": false,
          "name": "stack",
          "options": [],
          "query": "label_values(container_cpu_user_seconds_total,container_label_com_docker_stack_namespace)",
          "refresh": 2,
          "regex": "",
          "sort": 1,
          "tagValuesQuery": "",
          "tags": [],
          "tagsQuery": "",
          "type": "query",
          "useTags": false
        },
        {
          "allValue": ".+",
          "current": {
            "text": "All",
            "value": "$__all"
          },
          "datasource": null,
          "hide": 0,
          "includeAll": true,
          "label": "Node",
          "multi": false,
          "name": "node",
          "options": [],
          "query": "label_values(container_cpu_user_seconds_total,instance)",
          "refresh": 1,
          "regex": "",
          "sort": 1,
          "tagValuesQuery": "",
          "tags": [],
          "tagsQuery": "",
          "type": "query",
          "useTags": false
        }
      ]
    },
    "time": {
      "from": "now-30m",
      "to": "now"
    },
    "timepicker": {
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
      ]
    },
    "timezone": "browser",
    "title": "Container stats per node, stack and service",
    "version": 1
  }
}
