# get a certificate
# prepare the CSR with subject alt names

cfg=$(mktemp)
cp /etc/ssl/openssl.cnf $cfg
sed -i '/^\[ req \]$/ a\
req_extensions = v3_req' $cfg
sed -i '/^\[ v3_req \]$/ a\
subjectAltName          = @alternate_names' $cfg
cat >> $cfg << EOF

[ alternate_names ]
DNS.1       = $(hostname -f)
DNS.2       = $(hostname)
IP.1       = 192.168.2.200
EOF
cacfg=$(mktemp)
cat >> $cacfg << EOF

basicConstraints=CA:FALSE
subjectAltName          = @alternate_names
subjectKeyIdentifier = hash

[ alternate_names ]
DNS.1       = $(hostname -f)
DNS.2       = $(hostname)
IP.1       = 192.168.2.200
EOF

openssl genrsa -out {{ ref "/docker/remoteapi/srvkeyfile" }} 2048 || exit 1
openssl req -subj "/CN=$(hostname)" -sha256 -new -key {{ ref "/docker/remoteapi/srvkeyfile" }} -out {{ ref "/docker/remoteapi/srvcertfile" }}.csr || exit 1
curl --data "csr=$(sed 's/+/%2B/g' {{ ref "/docker/remoteapi/srvcertfile" }}.csr);ext=$(cat $cacfg)"  {{ ref "/certificate/ca/service" }}/csr > {{ ref "/docker/remoteapi/srvcertfile" }}
curl {{ ref "/certificate/ca/service" }}/ca > {{ ref "/docker/remoteapi/cafile" }}
rm -f {{ ref "/docker/remoteapi/srvcertfile" }}.csr $cfg $cacfg
