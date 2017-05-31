### Agent

The Agent is responsible to send logs and metrics to Elasticsearch using NATS.

there are several parameters to optimize the network usage:

#### logsBufferPeriod

System Variable: LOGS_BUFFER_PERIOD

indicate to the agent to keep the logs in memory and send them in one message every LOGS_BUFFER_PERIOD seconds
if LOGS_BUFFER_PERIOD = 0 then the logs are sent one by one
Default: 0

#### logsBufferSize

System Variable: LOGS_BUFFER_SIZE

It's the maximum number of log messages the agent can keep in memory. If the end of buffer is reached before the end of the time period then the agent send the logs in one message anyway
if LOGS_BUFFER_SIZE = 0 then the logs are sent one by one
Default: 0


#### metricsPeriod

System Variable: METRICS_PERIOD

indicate to the agent to keep the metrics in memory and send them in one message every METRICS_PERIOD seconds.
During the period, the metrics values add added in order to send their average value at the end of the period.
So it sends only one metrics message to Nats per container each period.
Default: 30 seconds
