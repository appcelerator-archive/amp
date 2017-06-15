{
  "Dashboard": {
    "title": "Cluster Health",
    "annotations": {
      "list": []
    },
    "editable": true,
    "gnetId": null,
    "graphTooltip": 0,
    "hideControls": false,
    "links": [],
    "refresh": false,
    "rows": [
      {
        "collapse": false,
        "height": 174,
        "panels": [
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 1,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": false
            },
            "id": 15,
            "interval": null,
            "links": [],
            "mappingType": 2,
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
              },
              {
                "from": "0",
                "text": "DOWN",
                "to": "0"
              },
              {
                "from": "1",
                "text": "Up",
                "to": "999"
              }
            ],
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "up{job=\"haproxy\"}",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "up",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1",
            "title": "HAPROXY",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 14,
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
            "postfix": " UP",
            "postfixFontSize": "50%",
            "prefix": "",
            "prefixFontSize": "50%",
            "rangeMaps": [
              {
                "from": "null",
                "text": "N/A",
                "to": "null"
              },
              {
                "from": "0",
                "text": "Down",
                "to": "0"
              },
              {
                "from": "1",
                "text": "Up",
                "to": "999"
              }
            ],
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "Value",
            "targets": [
              {
                "expr": "sum(etcd_server_has_leader)",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "etcd_server_has_leader",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1",
            "title": "ETCD",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 12,
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
            "postfix": " UP",
            "postfixFontSize": "50%",
            "prefix": "",
            "prefixFontSize": "50%",
            "rangeMaps": [
              {
                "from": "null",
                "text": "N/A",
                "to": "null"
              },
              {
                "from": "0",
                "text": "Down",
                "to": "0"
              },
              {
                "from": "1",
                "text": "Up",
                "to": "999"
              }
            ],
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "count(up{job=\"nats\"}==1)",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "up",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1",
            "title": "NATS",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 13,
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
            "postfix": " UP",
            "postfixFontSize": "50%",
            "prefix": "",
            "prefixFontSize": "50%",
            "rangeMaps": [
              {
                "from": "null",
                "text": "N/A",
                "to": "null"
              },
              {
                "from": "0",
                "text": "Down",
                "to": "0"
              },
              {
                "from": "1",
                "text": "Up",
                "to": "999"
              }
            ],
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "count(up{job=\"elasticsearch\"}==1)",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "es_cluster_status",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1",
            "title": "ES",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "0",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 16,
            "interval": null,
            "links": [],
            "mappingType": 2,
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
            "postfix": " UP",
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
            "repeat": null,
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "Value",
            "targets": [
              {
                "expr": "count(up{job=\"docker-engine\"}==1)",
                "format": "table",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "up",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1",
            "title": "DOCKER",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              },
              {
                "op": "=",
                "text": "Down",
                "value": "0"
              },
              {
                "op": "=",
                "text": "Up",
                "value": "1"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 17,
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
            "postfix": " KO",
            "postfixFontSize": "50%",
            "prefix": "",
            "prefixFontSize": "50%",
            "rangeMaps": [
              {
                "from": "null",
                "text": "0",
                "to": "null"
              }
            ],
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "Value",
            "targets": [
              {
                "expr": "count(up{job=\"docker-engine\"}) - count(up{job=\"docker-engine\"}==1)",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "up",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1",
            "title": "DOCKER",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "0",
                "value": "null"
              }
            ],
            "valueName": "avg"
          }
        ],
        "repeat": null,
        "repeatIteration": null,
        "repeatRowId": null,
        "showTitle": true,
        "title": "OVERALL STATUS",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": 250,
        "panels": [
          {
            "columns": [
              {
                "text": "Current",
                "value": "current"
              }
            ],
            "fontSize": "100%",
            "id": 65,
            "links": [],
            "pageSize": null,
            "scroll": true,
            "showHeader": true,
            "sort": {
              "col": 0,
              "desc": true
            },
            "span": 6,
            "styles": [
              {
                "alias": "Time",
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "pattern": "Time",
                "type": "date"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "decimals": 2,
                "pattern": "graph_driver",
                "thresholds": [],
                "type": "string",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "instance",
                "thresholds": [],
                "type": "string",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "os",
                "thresholds": [],
                "type": "string",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "version",
                "thresholds": [],
                "type": "string",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "Current",
                "thresholds": [],
                "type": "hidden",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "_name_",
                "thresholds": [],
                "type": "hidden",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "architecture",
                "thresholds": [],
                "type": "hidden",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "commit",
                "thresholds": [],
                "type": "hidden",
                "unit": "short"
              },
              {
                "alias": "",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "Value",
                "thresholds": [],
                "type": "hidden",
                "unit": "short"
              }
            ],
            "targets": [
              {
                "expr": "engine_daemon_engine_info",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "OS={{os}} driver={{graph_driver}} version={{version}} instance={{instance}}",
                "metric": "engine_daemon_engine_info",
                "refId": "A",
                "step": 10
              }
            ],
            "title": "Docker Engine Information",
            "transform": "timeseries_aggregations",
            "type": "table"
          },
          {
            "columns": [
              {
                "text": "Current",
                "value": "current"
              }
            ],
            "fontSize": "100%",
            "hideTimeOverride": false,
            "id": 62,
            "links": [],
            "pageSize": 12,
            "scroll": true,
            "showHeader": true,
            "sort": {
              "col": 0,
              "desc": true
            },
            "span": 3,
            "styles": [
              {
                "alias": "Time",
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "pattern": "Time",
                "type": "date"
              },
              {
                "alias": "Status",
                "colorMode": "value",
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "decimals": 0,
                "pattern": "Current",
                "thresholds": [
                  "0.1",
                  "0.9"
                ],
                "type": "number",
                "unit": "short"
              },
              {
                "alias": "Instance",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "Metric",
                "thresholds": [],
                "type": "number",
                "unit": "short"
              }
            ],
            "targets": [
              {
                "expr": "up{job=\"docker-engine\"}",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "metric": "up",
                "refId": "A",
                "step": 20
              }
            ],
            "title": "Docker Engine List",
            "transform": "timeseries_aggregations",
            "type": "table"
          },
          {
            "columns": [
              {
                "text": "Current",
                "value": "current"
              }
            ],
            "fontSize": "100%",
            "id": 63,
            "links": [],
            "pageSize": null,
            "scroll": true,
            "showHeader": true,
            "sort": {
              "col": 0,
              "desc": true
            },
            "span": 3,
            "styles": [
              {
                "alias": "Time",
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "pattern": "Time",
                "type": "date"
              },
              {
                "alias": "Status",
                "colorMode": "value",
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "decimals": 0,
                "pattern": "Current",
                "thresholds": [
                  "0.1",
                  "0.9"
                ],
                "type": "number",
                "unit": "short"
              },
              {
                "alias": "Instance",
                "colorMode": null,
                "colors": [
                  "rgba(245, 54, 54, 0.9)",
                  "rgba(237, 129, 40, 0.89)",
                  "rgba(50, 172, 45, 0.97)"
                ],
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "decimals": 2,
                "pattern": "Metric",
                "thresholds": [],
                "type": "string",
                "unit": "short"
              }
            ],
            "targets": [
              {
                "expr": "up{job=\"nodes\"}",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "metric": "up",
                "refId": "A",
                "step": 20
              }
            ],
            "title": "Node List",
            "transform": "timeseries_aggregations",
            "type": "table"
          }
        ],
        "repeat": null,
        "repeatIteration": null,
        "repeatRowId": null,
        "showTitle": false,
        "title": "Dashboard Row",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": 250,
        "panels": [
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 36,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "hideEmpty": true,
              "max": true,
              "min": false,
              "rightSide": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "node_load1",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "metric": "node_load1",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Load 1 Minute",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 37,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "hideEmpty": true,
              "max": true,
              "min": false,
              "rightSide": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "node_load5",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "metric": "node_load5",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Load 5 Minutes",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 38,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "hideEmpty": true,
              "max": true,
              "min": false,
              "rightSide": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "node_load15",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "metric": "node_load15",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Load 15 Minutes",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 42,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "100 - 100 * node_memory_MemFree / node_memory_MemTotal",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Memory Usage",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "percent",
                "label": "",
                "logBase": 1,
                "max": "100",
                "min": "0",
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
            "fill": 1,
            "id": 43,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "100 - 100 * node_filesystem_free / node_filesystem_size{fstype=~\"xfs|ext4\",mountpoint=\"/rootfs\"}",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "/ Usage",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "percent",
                "label": "",
                "logBase": 1,
                "max": "100",
                "min": "0",
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
            "decimals": 1,
            "fill": 1,
            "id": 44,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "100 - 100 * node_filesystem_free / node_filesystem_size{fstype=~\"xfs|ext4\",mountpoint=\"/rootfs/var/lib/docker\"}",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "/var/lib/docker Usage",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "percent",
                "label": "",
                "logBase": 1,
                "max": "100",
                "min": "0",
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
        "showTitle": true,
        "title": "SYSTEM",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": 250,
        "panels": [
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 53,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "haproxy_server_current_queue",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{backend}}",
                "metric": "haproxy_server_current_queue",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Server Current Queue",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 54,
            "legend": {
              "alignAsTable": true,
              "avg": false,
              "current": true,
              "max": false,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "haproxy_frontend_current_sessions",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{frontend}}",
                "metric": "haproxy_frontend_current_sessions",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Frontend Http Sessions",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 55,
            "legend": {
              "alignAsTable": true,
              "avg": false,
              "current": true,
              "max": false,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "rate(haproxy_frontend_http_responses_total[1m])",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{frontend}}-{{code}}",
                "metric": "haproxy_frontend_http_responses_total",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Frontend Http Responses (per min)",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
        "showTitle": true,
        "title": "HAPROXY",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": 140,
        "panels": [
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 1,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum(etcd_server_leader_changes_seen_total)",
                "format": "time_series",
                "hide": false,
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "etcd_server_leader_changes_seen_total",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "5,100",
            "title": "Leader Changes Seen",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 2,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum(etcd_server_proposals_committed_total)",
                "format": "time_series",
                "intervalFactor": 2,
                "metric": "etcd_server_proposals_committed_total",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "",
            "title": "Proposals Committed",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 3,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum(etcd_server_proposals_applied_total)",
                "format": "time_series",
                "intervalFactor": 2,
                "metric": "etcd_server_proposals_applied_total",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "",
            "title": "Proposal Applied",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum(etcd_server_proposals_failed_total)",
                "format": "time_series",
                "intervalFactor": 2,
                "metric": "etcd_server_proposals_failed_total",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,10",
            "title": "Proposal Failed",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 5,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum(etcd_server_proposals_pending)",
                "format": "time_series",
                "intervalFactor": 2,
                "metric": "etcd_server_proposals_pending",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "50,100",
            "title": "Proposal Pending",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 45,
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
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 3,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "sum(rate({grpc_type=\"unary\",grpc_code!=\"OK\"} [1m]))",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} RPC Rate",
                "metric": "grpc_server_started_total",
                "refId": "A",
                "step": 20
              },
              {
                "expr": "sum(rate(grpc_server_started_total{grpc_type=\"unary\",grpc_code!=\"OK\"} [1m])) - sum(rate(grpc_server_handled_total{grpc_type=\"unary\"} [1m]))",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} RPC Failed Rate",
                "metric": "grpc_server_handled_total",
                "refId": "B",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "RPC Rate",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "ops",
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
            "fill": 1,
            "id": 47,
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
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 3,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "etcd_debugging_mvcc_db_total_size_in_bytes",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} DB Size",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "DB Size",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "decbytes",
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
            "fill": 1,
            "id": 49,
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
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 3,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "process_resident_memory_bytes{job=\"etcd\"}",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} Resident Memory",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Memory",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "decbytes",
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
            "fill": 1,
            "id": 46,
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
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 3,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "sum(grpc_server_started_total {grpc_service=\"etcdserverpb.Watch\",grpc_type=\"bidi_stream\",grpc_code!=\"OK\"}) - sum(grpc_server_handled_total {grpc_service=\"etcdserverpb.Watch\",grpc_type=\"bidi_stream\"})",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "Watch Streams",
                "metric": "grpc_server_handled_total",
                "refId": "A",
                "step": 20
              },
              {
                "expr": "sum(grpc_server_started_total {grpc_service=\"etcdserverpb.Lease\",grpc_type=\"bidi_stream\"}) - sum(grpc_server_handled_total {grpc_service=\"etcdserverpb.Lease\",grpc_type=\"bidi_stream\"})",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "Lease Streams",
                "metric": "grpc_server_handled_total",
                "refId": "B",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Active Streams",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 48,
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
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "histogram_quantile(0.99, sum(rate(etcd_disk_wal_fsync_duration_seconds_bucket{job=\"etcd\"} [5m])) by (instance, le))",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} WAL fsync",
                "metric": "etcd_disk_wal_fsync_duration_seconds_bucket",
                "refId": "A",
                "step": 20
              },
              {
                "expr": "histogram_quantile(0.99, sum(rate(etcd_disk_backend_commit_duration_seconds_bucket [5m])) by (instance, le))",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} DB fsync",
                "metric": "etcd_disk_backend_commit_duration_seconds_bucket",
                "refId": "B",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Disk Sync Duration",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "s",
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
            "fill": 1,
            "id": 50,
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
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "rate(etcd_network_client_grpc_received_bytes_total [1m])",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} Client Traffic In",
                "metric": "etcd_network_client_grpc_received_bytes_total",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Client Traffic In",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "Bps",
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
            "fill": 1,
            "id": 51,
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
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "rate(etcd_network_client_grpc_sent_bytes_total [1m])",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} Client Traffic Out",
                "metric": "etcd_network_client_grpc_sent_bytes_total",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Client Traffic Out",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "Bps",
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
        "showTitle": true,
        "title": "ETCD",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": 232,
        "panels": [
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 9,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "min(es_cluster_status)",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{cluster}}",
                "metric": "es_cluster_status",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "",
            "title": "Cluster Status",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              },
              {
                "op": "=",
                "text": "UP",
                "value": "1"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 8,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "avg(es_cluster_nodes_number{job=\"elasticsearch\"})",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1",
            "title": "Nodes Number",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "format": "none",
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "max(es_jvm_mem_heap_used_percent)",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "es_jvm_mem_heap_used_percent",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "70,90",
            "title": "Mem Heap Percent Used",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "format": "none",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 10,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": false,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "sum(es_cluster_pending_tasks_number)",
                "format": "time_series",
                "intervalFactor": 2,
                "metric": "es_cluster_pending_tasks_number",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "50,500",
            "title": "Pending Tasks",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 11,
            "legend": {
              "alignAsTable": true,
              "avg": false,
              "current": true,
              "max": false,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "es_cluster_task_max_waiting_time_seconds",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{cluster}}",
                "metric": "es_cluster_task_max_waiting_time_seconds",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Task Max Waiting Time",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
        "showTitle": true,
        "title": "ELASTICSEARCH",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": 250,
        "panels": [
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 33,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "gnatsd_varz_connections",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{server_id}}",
                "metric": "gnatsd_varz_connections",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Connections",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 34,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "gnatsd_varz_subscriptions",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{server_id}}",
                "metric": "gnatsd_varz_subscriptions",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Subscriptions",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 35,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "rightSide": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "gnatsd_varz_slow_consumers",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{server_id}}",
                "metric": "gnatsd_varz_slow_consumers",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Slow Consumers",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
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
            "fill": 1,
            "id": 56,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 3,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "rate(gnatsd_varz_in_msgs[1m])",
                "format": "time_series",
                "interval": "500ms",
                "intervalFactor": 2,
                "legendFormat": "{{server_id}}",
                "metric": "gnatsd_varz_in_msgs",
                "refId": "A",
                "step": 10
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "NATS Msgs In Rate",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "wpm",
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
            "fill": 1,
            "id": 60,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 3,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "rate(gnatsd_varz_out_msgs[1m])",
                "format": "time_series",
                "interval": "500ms",
                "intervalFactor": 2,
                "legendFormat": "{{server_id}}",
                "metric": "gnatsd_varz_out_msgs",
                "refId": "A",
                "step": 10
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "NATS Msgs Out Rate",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "wpm",
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
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "decimals": null,
            "format": "short",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 57,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": true,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "gnatsd_varz_in_msgs",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "gnatsd_varz_in_msgs",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "",
            "title": "Msgs In",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "datasource": null,
            "decimals": null,
            "format": "short",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 58,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": true,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "gnatsd_varz_out_msgs",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "gnatsd_varz_out_msgs",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "",
            "title": "Msgs Out",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "rgba(50, 172, 45, 0.97)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(245, 54, 54, 0.9)"
            ],
            "datasource": null,
            "decimals": null,
            "format": "short",
            "gauge": {
              "maxValue": 100,
              "minValue": 0,
              "show": false,
              "thresholdLabels": false,
              "thresholdMarkers": true
            },
            "id": 59,
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
            "span": 2,
            "sparkline": {
              "fillColor": "rgba(31, 118, 189, 0.18)",
              "full": true,
              "lineColor": "rgb(31, 120, 193)",
              "show": true
            },
            "tableColumn": "",
            "targets": [
              {
                "expr": "gnatsd_varz_in_msgs - gnatsd_varz_out_msgs",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "metric": "gnatsd_varz_in_msgs",
                "refId": "A",
                "step": 60
              }
            ],
            "thresholds": "1,1000",
            "title": "Msgs In - Out",
            "type": "singlestat",
            "valueFontSize": "80%",
            "valueMaps": [
              {
                "op": "=",
                "text": "N/A",
                "value": "null"
              }
            ],
            "valueName": "avg"
          }
        ],
        "repeat": null,
        "repeatIteration": null,
        "repeatRowId": null,
        "showTitle": true,
        "title": "NATS",
        "titleSize": "h6"
      },
      {
        "collapse": false,
        "height": 214,
        "panels": [
          {
            "cards": {
              "cardPadding": null,
              "cardRound": null
            },
            "color": {
              "cardColor": "#b4ff00",
              "colorScale": "sqrt",
              "colorScheme": "interpolateOranges",
              "exponent": 0.5,
              "mode": "spectrum"
            },
            "dataFormat": "timeseries",
            "heatmap": {},
            "highlightCards": true,
            "id": 40,
            "links": [],
            "minSpan": 4,
            "span": 4,
            "targets": [
              {
                "expr": "engine_daemon_container_actions_seconds_bucket",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{job}}",
                "metric": "engine_daemon_container_actions_seconds_bucket",
                "refId": "A",
                "step": 20
              }
            ],
            "title": "Container Actions Seconds",
            "tooltip": {
              "show": true,
              "showHistogram": false
            },
            "type": "heatmap",
            "xAxis": {
              "show": true
            },
            "xBucketNumber": null,
            "xBucketSize": null,
            "yAxis": {
              "decimals": null,
              "format": "short",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true,
              "splitFactor": null
            },
            "yBucketNumber": null,
            "yBucketSize": null
          },
          {
            "cards": {
              "cardPadding": null,
              "cardRound": null
            },
            "color": {
              "cardColor": "#b4ff00",
              "colorScale": "sqrt",
              "colorScheme": "interpolateOranges",
              "exponent": 0.5,
              "mode": "spectrum"
            },
            "dataFormat": "timeseries",
            "heatmap": {},
            "highlightCards": true,
            "id": 41,
            "links": [],
            "minSpan": 4,
            "span": 4,
            "targets": [
              {
                "expr": "engine_daemon_network_actions_seconds_bucket",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{job}}",
                "metric": "engine_daemon_network_actions_seconds_bucket",
                "refId": "A",
                "step": 20
              }
            ],
            "title": "Network Actions Seconds",
            "tooltip": {
              "show": true,
              "showHistogram": false
            },
            "type": "heatmap",
            "xAxis": {
              "show": true
            },
            "xBucketNumber": null,
            "xBucketSize": null,
            "yAxis": {
              "decimals": null,
              "format": "short",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true,
              "splitFactor": null
            },
            "yBucketNumber": null,
            "yBucketSize": null
          },
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "decimals": 0,
            "fill": 0,
            "id": 64,
            "legend": {
              "alignAsTable": true,
              "avg": false,
              "current": true,
              "hideEmpty": true,
              "hideZero": true,
              "max": true,
              "min": false,
              "show": true,
              "total": false,
              "values": true
            },
            "lines": true,
            "linewidth": 1,
            "links": [],
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "span": 4,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "grpc_server_handled_total{grpc_service=\"docker.swarmkit.v1.RaftMembership\"}!=0",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}} {{grpc_method}}={{grpc_code}}",
                "metric": "grpc_server_handled_total",
                "refId": "A",
                "step": 20
              },
              {
                "expr": "sum(grpc_server_handled_total{grpc_method=\"Join\",grpc_code=\"OK\",grpc_service=\"docker.swarmkit.v1.RaftMembership\"})-sum(grpc_server_handled_total{grpc_method=\"Leave\",grpc_code=\"OK\",grpc_service=\"docker.swarmkit.v1.RaftMembership\"})",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "Join - Leave",
                "metric": "grpc_server_handled_total",
                "refId": "B",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Membership Operations",
            "tooltip": {
              "shared": true,
              "sort": 0,
              "value_type": "individual"
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
                "format": "short",
                "label": null,
                "logBase": 1,
                "max": null,
                "min": "0",
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
        "showTitle": true,
        "title": "Docker Engine",
        "titleSize": "h6"
      }
    ],
    "schemaVersion": 14,
    "style": "dark",
    "tags": [],
    "templating": {
      "list": [
        {
          "allValue": null,
          "current": {},
          "datasource": null,
          "hide": 2,
          "includeAll": true,
          "label": "Node",
          "multi": true,
          "name": "node",
          "options": [],
          "query": "label_values(node_load1, instance)",
          "refresh": 1,
          "regex": "",
          "sort": 0,
          "tagValuesQuery": "",
          "tags": [],
          "tagsQuery": "",
          "type": "query",
          "useTags": false
        },
        {
          "allValue": null,
          "current": {},
          "datasource": null,
          "hide": 2,
          "includeAll": true,
          "label": "Docker Engine",
          "multi": true,
          "name": "engine",
          "options": [],
          "query": "label_values(engine_daemon_engine_info, instance)",
          "refresh": 1,
          "regex": "",
          "sort": 0,
          "tagValuesQuery": "",
          "tags": [],
          "tagsQuery": "",
          "type": "query",
          "useTags": false
        }
      ]
    },
    "time": {
      "from": "now-1h",
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
    "version": 1
  }
}
