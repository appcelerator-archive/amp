###amp-haproxy-controller Prototype


It's a HAProxy image containing a controller which is able to update HAProxy configuration when amp stacks is created or deleted.

# tags

- 1.0.5

# Setting

environment variable: 

- ETCD_ENDPOINTS: should be equal to the ETCD endpoints list, format: "host1:port1, host2:port2, ..."
- NO_DEFAULT_BACKEND: remove all amp default infrastructure backends (for test)


Concidenring HAProx does'nt start or update its conf if one the declared backend has a DNS name not resolved, the amp-haproxy-controller  disabled them in the HAProxy controller allowing HAProxy to start correcly. The disabled backend are checked on regular basis to be enabled when their dns names are back to resolved.



