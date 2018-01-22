{
  "Dashboard": {
    "title": "Cluster Health - Docker",
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
        "height": 214,
        "panels": [
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
                "hide": false,
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
                "hide": false,
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
          },
          {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "decimals": 0,
            "fill": 1,
            "id": 99,
            "legend": {
              "alignAsTable": true,
              "avg": false,
              "current": true,
              "hideEmpty": true,
              "hideZero": true,
              "max": false,
              "min": false,
              "show": true,
              "total": true,
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
                "expr": "rate(grpc_client_handled_total{grpc_code=\"OK\",grpc_method=\"Heartbeat\",grpc_service=\"docker.swarmkit.v1.Dispatcher\"}[1m])*60",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}::{{grpc_code}}",
                "metric": "grpc_client_handled_total",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Heartbeats per min",
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
            "decimals": 0,
            "fill": 1,
            "id": 68,
            "legend": {
              "alignAsTable": true,
              "avg": false,
              "current": true,
              "hideEmpty": true,
              "hideZero": true,
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
                "expr": "engine_daemon_health_checks_failed_total",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "metric": "grpc_client_handled_total",
                "refId": "A",
                "step": 20
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Healthchecks Failed",
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
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": null,
            "fill": 1,
            "id": 83,
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
                "expr": "swarm_raft_transaction_latency_seconds_bucket{le=\"0.1\"} / ignoring(le) swarm_raft_transaction_latency_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Raft Transaction Latency <= 0.1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": "",
                "logBase": 1,
                "max": "1",
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
            "id": 82,
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
                "expr": "swarm_raft_transaction_latency_seconds_bucket{le=\"0.01\"} / ignoring(le) swarm_raft_transaction_latency_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Raft Transaction Latency <= 0.01s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": "",
                "logBase": 1,
                "max": "1",
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
            "id": 88,
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
                "expr": "swarm_raft_transaction_latency_seconds_bucket{le=\"0.005\"} / ignoring(le) swarm_raft_transaction_latency_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Raft Transaction Latency <= 0.005s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 89,
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
                "expr": "swarm_dispatcher_scheduling_delay_seconds_bucket{le=\"10\"} / ignoring(le) swarm_dispatcher_scheduling_delay_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Dispatcher Scheduling Delay <= 10s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 86,
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
                "expr": "swarm_dispatcher_scheduling_delay_seconds_bucket{le=\"1\"} / ignoring(le) swarm_dispatcher_scheduling_delay_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Dispatcher Scheduling Delay <= 1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 85,
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
                "expr": "swarm_dispatcher_scheduling_delay_seconds_bucket{le=\"0.1\"} / ignoring(le) swarm_dispatcher_scheduling_delay_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Dispatcher Scheduling Delay <= 0.1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 90,
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
                "expr": "swarm_store_memory_store_lock_duration_seconds_bucket{le=\"0.1\"} / ignoring(le) swarm_store_memory_store_lock_duration_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Memory Store Lock Duration <= 0.1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 91,
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
                "expr": "swarm_store_memory_store_lock_duration_seconds_bucket{le=\"0.01\"} / ignoring(le) swarm_store_memory_store_lock_duration_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Memory Store Lock Duration <= 0.01s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 92,
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
                "expr": "swarm_store_memory_store_lock_duration_seconds_bucket{le=\"0.005\"} / ignoring(le) swarm_store_memory_store_lock_duration_seconds_count",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Swarm Memory Store Lock Duration <= 0.005s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 94,
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
                "expr": "engine_daemon_network_actions_seconds_bucket{le=\"5\",action=\"allocate\"} / ignoring(le) engine_daemon_network_actions_seconds_count{action=\"allocate\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Network Allocation Actions <= 5s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 93,
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
                "expr": "engine_daemon_network_actions_seconds_bucket{le=\"1\",action=\"allocate\"} / ignoring(le) engine_daemon_network_actions_seconds_count{action=\"allocate\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Network Allocation Actions <= 1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 95,
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
                "expr": "engine_daemon_network_actions_seconds_bucket{le=\"0.1\",action=\"allocate\"} / ignoring(le) engine_daemon_network_actions_seconds_count{action=\"allocate\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Network Allocation Actions <= 0.1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 96,
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
                "expr": "engine_daemon_network_actions_seconds_bucket{le=\"5\",action=\"connect\"} / ignoring(le) engine_daemon_network_actions_seconds_count{action=\"connect\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Network Connection Actions <= 5s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 97,
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
                "expr": "engine_daemon_network_actions_seconds_bucket{le=\"1\",action=\"connect\"} / ignoring(le) engine_daemon_network_actions_seconds_count{action=\"connect\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Network Connection Actions <= 1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 98,
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
                "expr": "engine_daemon_network_actions_seconds_bucket{le=\"0.1\",action=\"connect\"} / ignoring(le) engine_daemon_network_actions_seconds_count{action=\"connect\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Network Connection Actions <= 0.1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 100,
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
                "expr": "engine_daemon_container_actions_seconds_bucket{le=\"5\",action=\"create\"} / ignoring(le) engine_daemon_container_actions_seconds_count{action=\"create\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Create Actions <= 5s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 101,
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
                "expr": "engine_daemon_container_actions_seconds_bucket{le=\"1\",action=\"create\"} / ignoring(le) engine_daemon_container_actions_seconds_count{action=\"create\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Create Actions <= 1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 102,
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
                "expr": "engine_daemon_container_actions_seconds_bucket{le=\"0.1\",action=\"create\"} / ignoring(le) engine_daemon_container_actions_seconds_count{action=\"create\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Create Actions <= 0.1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 103,
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
                "expr": "engine_daemon_container_actions_seconds_bucket{le=\"5\",action=\"start\"} / ignoring(le) engine_daemon_container_actions_seconds_count{action=\"start\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Start Actions <= 5s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 105,
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
                "expr": "engine_daemon_container_actions_seconds_bucket{le=\"2.5\",action=\"start\"} / ignoring(le) engine_daemon_container_actions_seconds_count{action=\"start\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Start Actions <= 2.5s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 104,
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
                "expr": "engine_daemon_container_actions_seconds_bucket{le=\"1\",action=\"start\"} / ignoring(le) engine_daemon_container_actions_seconds_count{action=\"start\"}",
                "format": "time_series",
                "instant": false,
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Container Start Actions <= 1s",
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
                "decimals": 0,
                "format": "percentunit",
                "label": null,
                "logBase": 1,
                "max": "1",
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
            "id": 80,
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
                "expr": "engine_daemon_events_total",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Docker Engine Events",
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
            "id": 87,
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
                "expr": "rate(engine_daemon_events_total[1m])",
                "format": "time_series",
                "interval": "",
                "intervalFactor": 2,
                "legendFormat": "{{instance}}",
                "refId": "A"
              }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeShift": null,
            "title": "Docker Engine Events Rate",
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
        "showTitle": false,
        "title": "DOCKER",
        "titleSize": "h6"
     }
  ],
  "schemaVersion": 14,
  "style": "dark",
  "tags": [
    "cluster",
    "infrastructure",
    "docker",
    "swarm"
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
