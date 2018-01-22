{
  "Dashboard": {
    "title": "Cluster Health Overview",
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
							"maxValue": 100,
							"minValue": 0,
							"show": false,
							"thresholdLabels": false,
							"thresholdMarkers": true
						},
						"id": 72,
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
						"postfix": "UP",
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
						"span": 1,
						"sparkline": {
							"fillColor": "rgba(31, 118, 189, 0.18)",
							"full": false,
							"lineColor": "rgb(31, 120, 193)",
							"show": true
						},
						"tableColumn": "",
						"targets": [
							{
								"expr": "sum(up{job=\"amplifier\"})",
								"format": "time_series",
								"intervalFactor": 2,
								"legendFormat": "",
								"metric": "up",
								"refId": "A",
								"step": 60
							}
						],
						"thresholds": "0.9,1.1",
						"title": "Amplifier",
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
							"rgba(245, 54, 54, 0.9)",
							"rgba(237, 129, 40, 0.89)",
							"rgba(50, 172, 45, 0.97)"
						],
						"datasource": null,
						"decimals": 0,
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
								"text": "",
								"to": "null"
							}
						],
						"span": 1,
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
						"thresholds": "0.1,0.9",
						"title": "HAPROXY",
						"type": "singlestat",
						"valueFontSize": "80%",
						"valueMaps": [
							{
								"op": "=",
								"text": "DOWN",
								"value": "null"
							},
							{
								"op": "=",
								"text": "DOWN",
								"value": "0"
							},
							{
								"op": "=",
								"text": "UP",
								"value": "1"
							}
						],
						"valueName": "current"
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
						"span": 1,
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
						"valueName": "current"
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
						"span": 1,
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
						"valueName": "current"
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
						"span": 1,
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
						"valueName": "current"
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
						"span": 1,
						"sparkline": {
							"fillColor": "rgba(31, 118, 189, 0.18)",
							"full": true,
							"lineColor": "rgb(31, 120, 193)",
							"show": true
						},
						"tableColumn": "Value",
						"targets": [
							{
								"expr": "sum(up{job=\"docker-engine\"})",
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
						"title": "Docker Engines",
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
						"valueName": "current"
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
						"id": 67,
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
						"postfix": " node(s)",
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
						"span": 1,
						"sparkline": {
							"fillColor": "rgba(31, 118, 189, 0.18)",
							"full": true,
							"lineColor": "rgb(31, 120, 193)",
							"show": true
						},
						"tableColumn": "Value",
						"targets": [
							{
								"expr": "avg(swarm_manager_nodes{state=\"ready\"}!=0)",
								"format": "time_series",
								"interval": "",
								"intervalFactor": 2,
								"legendFormat": "",
								"metric": "swarm_manager_nodes",
								"refId": "A",
								"step": 60
							}
						],
						"thresholds": "0.1,0.9",
						"title": "Swarm",
						"type": "singlestat",
						"valueFontSize": "80%",
						"valueMaps": [
							{
								"op": "=",
								"text": "0",
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
						"postfix": " mngr(s)",
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
						"span": 1,
						"sparkline": {
							"fillColor": "rgba(31, 118, 189, 0.18)",
							"full": true,
							"lineColor": "rgb(31, 120, 193)",
							"show": true
						},
						"tableColumn": "Value",
						"targets": [
							{
								"expr": "sum(swarm_node_manager)",
								"format": "time_series",
								"interval": "",
								"intervalFactor": 2,
								"legendFormat": "",
								"metric": "swarm_node_manager",
								"refId": "A",
								"step": 60
							}
						],
						"thresholds": "0.1,0.9",
						"title": "Swarm",
						"type": "singlestat",
						"valueFontSize": "80%",
						"valueMaps": [
							{
								"op": "=",
								"text": "0",
								"value": "null"
							}
						],
						"valueName": "current"
					},
					{
						"aliasColors": {},
						"cacheTimeout": null,
						"combine": {
							"label": "Others",
							"threshold": "0.1"
						},
						"datasource": null,
						"fontSize": "80%",
						"format": "short",
						"id": 66,
						"interval": null,
						"legend": {
							"percentage": false,
							"show": true,
							"values": true
						},
						"legendType": "On graph",
						"links": [],
						"maxDataPoints": 3,
						"nullPointMode": "connected",
						"pieType": "pie",
						"span": 2,
						"strokeWidth": 1,
						"targets": [
							{
								"expr": "avg(swarm_manager_nodes{job=\"docker-engine\"}!=0) by (state)",
								"format": "time_series",
								"instant": true,
								"interval": "",
								"intervalFactor": 2,
								"legendFormat": "{{state}}",
								"metric": "swarm_manager_nodes",
								"refId": "A",
								"step": 2400
							}
						],
						"title": "Swarm Node State",
						"type": "grafana-piechart-panel",
						"valueName": "current"
					},
					{
						"cacheTimeout": null,
						"colorBackground": false,
						"colorValue": false,
						"colors": [
							"#d44a3a",
							"rgba(237, 129, 40, 0.89)",
							"#299c46"
						],
						"datasource": null,
						"decimals": 1,
						"description": "Nodes with monitoring tasks (cAdvisor)",
						"format": "percentunit",
						"gauge": {
							"maxValue": 1,
							"minValue": 0,
							"show": true,
							"thresholdLabels": false,
							"thresholdMarkers": false
						},
						"id": 79,
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
								"expr": "count(cadvisor_version_info) / max(swarm_manager_nodes{state=\"ready\"}!=0)",
								"format": "time_series",
								"intervalFactor": 2,
								"legendFormat": "",
								"refId": "A"
							}
						],
						"thresholds": "0.2,0.99",
						"title": "Monitoring coverage",
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
						"columns": [],
						"datasource": null,
						"fontSize": "100%",
						"id": 81,
						"links": [],
						"pageSize": null,
						"scroll": true,
						"showHeader": true,
						"sort": {
							"col": 0,
							"desc": true
						},
						"span": 8,
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
								"pattern": "taskname",
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
								"pattern": "cadvisorRevision",
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
								"pattern": "job",
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
							},
							{
								"alias": "cAdvisor Version",
								"colorMode": null,
								"colors": [
									"rgba(245, 54, 54, 0.9)",
									"rgba(237, 129, 40, 0.89)",
									"rgba(50, 172, 45, 0.97)"
								],
								"dateFormat": "YYYY-MM-DD HH:mm:ss",
								"decimals": 2,
								"pattern": "cadvisorVersion",
								"thresholds": [],
								"type": "string",
								"unit": "short"
							},
							{
								"alias": "Docker Version",
								"colorMode": null,
								"colors": [
									"rgba(245, 54, 54, 0.9)",
									"rgba(237, 129, 40, 0.89)",
									"rgba(50, 172, 45, 0.97)"
								],
								"dateFormat": "YYYY-MM-DD HH:mm:ss",
								"decimals": 2,
								"pattern": "dockerVersion",
								"thresholds": [],
								"type": "string",
								"unit": "short"
							},
							{
								"alias": "Node",
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
								"alias": "OS Version",
								"colorMode": null,
								"colors": [
									"rgba(245, 54, 54, 0.9)",
									"rgba(237, 129, 40, 0.89)",
									"rgba(50, 172, 45, 0.97)"
								],
								"dateFormat": "YYYY-MM-DD HH:mm:ss",
								"decimals": 2,
								"pattern": "osVersion",
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
								"pattern": "__name__",
								"thresholds": [],
								"type": "hidden",
								"unit": "short"
							},
							{
								"alias": "Kernel Version",
								"colorMode": null,
								"colors": [
									"rgba(245, 54, 54, 0.9)",
									"rgba(237, 129, 40, 0.89)",
									"rgba(50, 172, 45, 0.97)"
								],
								"dateFormat": "YYYY-MM-DD HH:mm:ss",
								"decimals": 2,
								"pattern": "kernelVersion",
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
								"decimals": 2,
								"pattern": "/.*/",
								"thresholds": [],
								"type": "number",
								"unit": "short"
							}
						],
						"targets": [
							{
								"expr": "cadvisor_version_info",
								"format": "table",
								"instant": true,
								"intervalFactor": 2,
								"legendFormat": "",
								"refId": "A"
							}
						],
						"title": "System Information",
						"transform": "table",
						"type": "table"
					},
					{
						"columns": [],
						"datasource": null,
						"fontSize": "100%",
						"id": 80,
						"links": [],
						"pageSize": null,
						"scroll": true,
						"showHeader": true,
						"sort": {
							"col": 0,
							"desc": true
						},
						"span": 4,
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
								"pattern": "branch",
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
								"pattern": "goversion",
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
								"pattern": "revision",
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
								"pattern": "/.*/",
								"thresholds": [],
								"type": "number",
								"unit": "short"
							}
						],
						"targets": [
							{
								"expr": "prometheus_build_info",
								"format": "table",
								"instant": true,
								"intervalFactor": 2,
								"legendFormat": "",
								"refId": "A"
							}
						],
						"title": "Prometheus Version",
						"transform": "table",
						"type": "table"
          },
          {
            "headings": true,
            "id": 82,
            "limit": 10,
            "links": [],
            "query": "",
            "recent": false,
            "search": true,
            "span": 6,
            "starred": true,
            "tags": [
              "infrastructure"
            ],
            "title": "Infrastructure Dashboards",
            "type": "dashlist"
          },
          {
            "headings": true,
            "id": 83,
            "limit": 10,
            "links": [],
            "query": "",
            "recent": false,
            "search": true,
            "span": 6,
            "starred": true,
            "tags": [
              "service"
            ],
            "title": "Service Dashboards",
            "type": "dashlist"
          }
				],
				"repeat": null,
				"repeatIteration": null,
				"repeatRowId": null,
				"showTitle": false,
				"title": "OVERALL STATUS",
				"titleSize": "h6"
			}
		],
    "schemaVersion": 14,
    "style": "dark",
    "tags": [
      "cluster",
      "infrastructure"
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
