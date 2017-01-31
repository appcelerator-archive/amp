{
  "Dashboard": {
    "id": null,
    "title": "AMP Load Balancers Realtime",
    "originalTitle": "AMP Load Balancers Realtime",
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
            "hideTimeOverride": false,
            "id": 1,
            "isNew": true,
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
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "span": 6,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "alias": "",
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
                "measurement": "haproxy",
                "policy": "default",
                "query": "SELECT mean(\"scur\"), mean(\"smax\") FROM \"haproxy\" WHERE $timeFilter GROUP BY time($interval) fill(null)",
                "rawQuery": false,
                "refId": "A",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "scur"
                      ],
                      "type": "field"
                    },
                    {
                      "params": [],
                      "type": "mean"
                    }
                  ],
                  [
                    {
                      "params": [
                        "smax"
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
            "title": "Sessions",
            "tooltip": {
              "msResolution": true,
              "shared": true,
              "sort": 0,
              "value_type": "cumulative"
            },
            "type": "graph",
            "xaxis": {
              "show": true
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
            "title": "Queue",
            "error": false,
            "span": 6,
            "editable": true,
            "type": "graph",
            "isNew": true,
            "id": 4,
            "targets": [
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
                        "qcur"
                      ]
                    },
                    {
                      "type": "mean",
                      "params": []
                    }
                  ]
                ],
                "refId": "A",
                "measurement": "haproxy"
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
                        "qmax"
                      ]
                    },
                    {
                      "type": "mean",
                      "params": []
                    }
                  ]
                ],
                "refId": "B",
                "measurement": "haproxy"
              }
            ],
            "datasource": null,
            "renderer": "flot",
            "yaxes": [
              {
                "label": null,
                "show": true,
                "logBase": 1,
                "min": null,
                "max": null,
                "format": "short"
              },
              {
                "label": null,
                "show": true,
                "logBase": 1,
                "min": null,
                "max": null,
                "format": "short"
              }
            ],
            "xaxis": {
              "show": true
            },
            "grid": {
              "threshold1": null,
              "threshold2": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "lines": true,
            "fill": 1,
            "linewidth": 2,
            "points": false,
            "pointradius": 5,
            "bars": false,
            "stack": false,
            "percentage": false,
            "legend": {
              "show": true,
              "values": false,
              "min": false,
              "max": false,
              "current": false,
              "total": false,
              "avg": false
            },
            "nullPointMode": "connected",
            "steppedLine": false,
            "tooltip": {
              "value_type": "cumulative",
              "shared": true,
              "sort": 0,
              "msResolution": false
            },
            "timeFrom": null,
            "timeShift": null,
            "aliasColors": {},
            "seriesOverrides": [],
            "links": []
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
            "id": 3,
            "isNew": true,
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
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "span": 6,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "alias": "",
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
                "measurement": "haproxy",
                "policy": "default",
                "refId": "A",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "bin"
                      ],
                      "type": "field"
                    },
                    {
                      "params": [],
                      "type": "mean"
                    }
                  ],
                  [
                    {
                      "params": [
                        "bout"
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
            "title": "IO",
            "tooltip": {
              "msResolution": false,
              "shared": true,
              "sort": 0,
              "value_type": "cumulative"
            },
            "type": "graph",
            "xaxis": {
              "show": true
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
            "id": 2,
            "isNew": true,
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
            "nullPointMode": "null",
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "span": 6,
            "stack": false,
            "steppedLine": false,
            "targets": [
              {
                "alias": "2xx",
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
                "measurement": "haproxy",
                "policy": "default",
                "refId": "A",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "http_response.2xx"
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
                "alias": "4xx",
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
                "measurement": "haproxy",
                "policy": "default",
                "refId": "B",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "http_response.4xx"
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
                "alias": "5xx",
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
                "measurement": "haproxy",
                "policy": "default",
                "refId": "C",
                "resultFormat": "time_series",
                "select": [
                  [
                    {
                      "params": [
                        "http_response.5xx"
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
            "title": "HTTP Codes",
            "tooltip": {
              "msResolution": false,
              "shared": true,
              "sort": 0,
              "value_type": "cumulative"
            },
            "type": "graph",
            "xaxis": {
              "show": true
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
        "title": "Row"
      },
      {
        "title": "New row",
        "height": "250px",
        "editable": true,
        "collapse": false,
        "panels": [
          {
            "title": "Active servers",
            "error": false,
            "span": 2,
            "editable": true,
            "type": "singlestat",
            "isNew": true,
            "id": 5,
            "targets": [
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
                        "active_servers"
                      ]
                    },
                    {
                      "type": "mean",
                      "params": []
                    }
                  ]
                ],
                "refId": "A",
                "measurement": "haproxy"
              }
            ],
            "links": [],
            "datasource": null,
            "maxDataPoints": 100,
            "interval": null,
            "cacheTimeout": null,
            "format": "none",
            "prefix": "",
            "postfix": "",
            "nullText": null,
            "valueMaps": [
              {
                "value": "null",
                "op": "=",
                "text": "N/A"
              }
            ],
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
            "rangeMaps": [
              {
                "from": "null",
                "to": "null",
                "text": "N/A"
              }
            ],
            "mappingType": 1,
            "nullPointMode": "connected",
            "valueName": "avg",
            "prefixFontSize": "50%",
            "valueFontSize": "80%",
            "postfixFontSize": "50%",
            "thresholds": "",
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "sparkline": {
              "show": true,
              "full": true,
              "lineColor": "rgb(31, 120, 193)",
              "fillColor": "rgba(31, 118, 189, 0.18)"
            },
            "gauge": {
              "show": false,
              "minValue": 0,
              "maxValue": 100,
              "thresholdMarkers": true,
              "thresholdLabels": false
            }
          },
          {
            "title": "Backup servers",
            "error": false,
            "span": 2,
            "editable": true,
            "type": "singlestat",
            "isNew": true,
            "id": 6,
            "targets": [
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
                        "backup_servers"
                      ]
                    },
                    {
                      "type": "mean",
                      "params": []
                    }
                  ]
                ],
                "refId": "A",
                "measurement": "haproxy"
              }
            ],
            "links": [],
            "datasource": null,
            "maxDataPoints": 100,
            "interval": null,
            "cacheTimeout": null,
            "format": "none",
            "prefix": "",
            "postfix": "",
            "nullText": null,
            "valueMaps": [
              {
                "value": "null",
                "op": "=",
                "text": "N/A"
              }
            ],
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
            "rangeMaps": [
              {
                "from": "null",
                "to": "null",
                "text": "N/A"
              }
            ],
            "mappingType": 1,
            "nullPointMode": "connected",
            "valueName": "avg",
            "prefixFontSize": "50%",
            "valueFontSize": "80%",
            "postfixFontSize": "50%",
            "thresholds": "",
            "colorBackground": false,
            "colorValue": false,
            "colors": [
              "rgba(245, 54, 54, 0.9)",
              "rgba(237, 129, 40, 0.89)",
              "rgba(50, 172, 45, 0.97)"
            ],
            "sparkline": {
              "show": true,
              "full": true,
              "lineColor": "rgb(31, 120, 193)",
              "fillColor": "rgba(31, 118, 189, 0.18)"
            },
            "gauge": {
              "show": false,
              "minValue": 0,
              "maxValue": 100,
              "thresholdMarkers": true,
              "thresholdLabels": false
            }
          }
        ]
      }
    ],
    "time": {
      "from": "now-6h",
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
    "templating": {
      "list": []
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
