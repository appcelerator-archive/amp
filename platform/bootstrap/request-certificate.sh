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

openssl genrsa -out {{ var "/docker/remoteapi/srvkeyfile" }} 2048 || exit 1
openssl req -subj "/CN=$(hostname)" -sha256 -new -key {{ var "/docker/remoteapi/srvkeyfile" }} -out {{ var "/docker/remoteapi/srvcertfile" }}.csr || exit 1
curl --data "csr=$(sed 's/+/%2B/g' {{ var "/docker/remoteapi/srvcertfile" }}.csr);ext=$(cat $cacfg)"  {{ var "/certificate/ca/service" }}/csr > {{ var "/docker/remoteapi/srvcertfile" }}
curl {{ var "/certificate/ca/service" }}/ca > {{ var "/docker/remoteapi/cafile" }}
rm -f {{ var "/docker/remoteapi/srvcertfile" }}.csr $cfg $cacfg
