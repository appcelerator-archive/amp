{
  "Dashboard": {
    "title": "Cluster Health - Prometheus",
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
        "height": 250,
        "panels": [
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 69,
            "legend": {
              "alignAsTable": true,
              "avg": true,
              "current": true,
              "hideEmpty": true,
              "hideZero": true,
              "max": true,
              "min": false,
              "rightSide": true,
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
            "span": 8,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "max(scrape_duration_seconds) by (job)",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{job}}",
                "metric": "scrape_duration_seconds",
                "refId": "A",
                "step": 4
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Scrape duration",
            "tooltip": {
              "shared": true,
              "sort": 2,
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
            "id": 106,
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
                "expr": "prometheus_tsdb_head_series",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{job}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Time Series Count",
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
                "decimals": null,
                "format": "short",
                "label": "",
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
          },
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 107,
            "legend": {
              "alignAsTable": true,
              "avg": false,
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
            "span": 6,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "prometheus_tsdb_wal_fsync_duration_seconds",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "quantile-{{quantile}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "TSDB WAL Fsync Duration (seconds)",
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
                "decimals": null,
                "format": "dtdurations",
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
          },
          {
            "columns": [],
            "datasource": null,
            "fontSize": "100%",
            "id": 109,
            "links": [],
            "pageSize": null,
            "scroll": true,
            "showHeader": true,
            "sort": {
              "col": 0,
              "desc": true
            },
            "span": 2,
            "styles": [
              {
                "alias": "Time",
                "dateFormat": "YYYY-MM-DD HH:mm:ss",
                "pattern": "Time",
                "type": "hidden"
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
                "link": false,
                "pattern": "__name__",
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
                "pattern": "hostip",
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
                "pattern": "taskname",
                "thresholds": [],
                "type": "hidden",
                "unit": "short"
              }
            ],
            "targets": [
              {
                "expr": "topk(3, scrape_samples_scraped)",
                "format": "table",
                "instant": true,
                "intervalFactor": 2,
                "legendFormat": "{{job}} {{instance}}",
                "refId": "A"
              }
            ],
            "title": "Top Samples Scraped Series",
            "transform": "table",
            "type": "table"
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": true,
            "colors": [
              "#299c46",
              "rgba(237, 129, 40, 0.89)",
              "#d44a3a"
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
            "id": 113,
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
                "expr": "prometheus_tsdb_wal_corruptions_total",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "refId": "A"
              }
            ],
            "thresholds": "0.5,0.6",
            "title": "TSDB WAL Corruptions Total",
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
              "#299c46",
              "rgba(237, 129, 40, 0.89)",
              "#d44a3a"
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
            "id": 108,
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
                "expr": "prometheus_tsdb_blocks_loaded",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{job}}",
                "refId": "A"
              }
            ],
            "thresholds": "",
            "title": "TSDB Blocks Loaded",
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
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 111,
            "legend": {
              "alignAsTable": true,
              "avg": false,
              "current": true,
              "hideEmpty": true,
              "hideZero": true,
              "max": true,
              "min": false,
              "show": true,
              "sort": "current",
              "sortDesc": true,
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
            "span": 8,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "expr": "topk(5, prometheus_engine_query_duration_seconds)",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{slice}}-quantile-{{quantile}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Top Engine Queries Durations",
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
                "decimals": null,
                "format": "dtdurations",
                "label": "",
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
          },
          {
            "cacheTimeout": null,
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "#299c46",
              "rgba(237, 129, 40, 0.89)",
              "#d44a3a"
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
            "id": 110,
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
                "expr": "prometheus_engine_queries_concurrent_max",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": "",
            "title": "Engine Queries Concurrent Max",
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
            "colorValue": true,
            "colors": [
              "#d44a3a",
              "rgba(237, 129, 40, 0.89)",
              "#299c46"
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
            "id": 114,
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
                "expr": "prometheus_notifications_alertmanagers_discovered",
                "format": "time_series",
                "intervalFactor": 2,
                "legendFormat": "",
                "refId": "A"
              }
            ],
            "thresholds": "1,2",
            "title": "Alertmanagers Discovered",
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
        "title": "PROMETHEUS",
        "titleSize": "h6"
      }
    ],
    "schemaVersion": 14,
    "style": "dark",
    "tags": [
      "cluster",
      "infrastructure",
      "monitoring"
    ],
    "templating": {
      "list": [
        {
          "allValue": null,
          "current": {
            "text": "All",
            "value": "$__all"
          },
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
          "current": {
            "text": "All",
            "value": "$__all"
          },
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
