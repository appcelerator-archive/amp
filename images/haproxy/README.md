### amp-haproxy-controller


HAProxy image including a controller for HAProxy configuration updates on amp stacks events.

# tags

- 1.0.5

# Setting

environment variable: 

- ETCD_ENDPOINTS: should be equal to the ETCD endpoints list, format: "host1:port1, host2:port2, ..."
- NO_DEFAULT_BACKEND: remove all amp default infrastructure backends (for test)


Considering HAProxy doesn't start or update its configuration if one of the declared backend has a DNS name not resolved, the amp-haproxy controller disables it in the HAProxy configuration allowing HAProxy to start without issues. The disabled backend are checked on a regular basis to be enabled when their dns names resolve again.
