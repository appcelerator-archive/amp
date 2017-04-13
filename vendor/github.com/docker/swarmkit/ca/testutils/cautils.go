package testutils

import (
	"crypto"
	cryptorand "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	cfcsr "github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/connectionbroker"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/ioutils"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	stateutils "github.com/docker/swarmkit/manager/state/testutils"
	"github.com/docker/swarmkit/remotes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TestCA is a structure that encapsulates everything needed to test a CA Server
type TestCA struct {
	RootCA                      ca.RootCA
	ExternalSigningServer       *ExternalSigningServer
	MemoryStore                 *store.MemoryStore
	Addr, TempDir, Organization string
	Paths                       *ca.SecurityConfigPaths
	Server                      *grpc.Server
	ServingSecurityConfig       *ca.SecurityConfig
	CAServer                    *ca.Server
	Context                     context.Context
	NodeCAClients               []api.NodeCAClient
	CAClients                   []api.CAClient
	Conns                       []*grpc.ClientConn
	WorkerToken                 string
	ManagerToken                string
	ConnBroker                  *connectionbroker.Broker
	KeyReadWriter               *ca.KeyReadWriter
	watchCancel                 func()
}

// Stop cleans up after TestCA
func (tc *TestCA) Stop() {
	tc.watchCancel()
	os.RemoveAll(tc.TempDir)
	for _, conn := range tc.Conns {
		conn.Close()
	}
	if tc.ExternalSigningServer != nil {
		tc.ExternalSigningServer.Stop()
	}
	tc.CAServer.Stop()
	tc.Server.Stop()
	tc.MemoryStore.Close()
}

// NewNodeConfig returns security config for a new node, given a role
func (tc *TestCA) NewNodeConfig(role string) (*ca.SecurityConfig, error) {
	withNonSigningRoot := tc.ExternalSigningServer != nil
	return genSecurityConfig(tc.MemoryStore, tc.RootCA, tc.KeyReadWriter, role, tc.Organization, tc.TempDir, withNonSigningRoot)
}

// WriteNewNodeConfig returns security config for a new node, given a role
// saving the generated key and certificates to disk
func (tc *TestCA) WriteNewNodeConfig(role string) (*ca.SecurityConfig, error) {
	withNonSigningRoot := tc.ExternalSigningServer != nil
	return genSecurityConfig(tc.MemoryStore, tc.RootCA, tc.KeyReadWriter, role, tc.Organization, tc.TempDir, withNonSigningRoot)
}

// NewNodeConfigOrg returns security config for a new node, given a role and an org
func (tc *TestCA) NewNodeConfigOrg(role, org string) (*ca.SecurityConfig, error) {
	withNonSigningRoot := tc.ExternalSigningServer != nil
	return genSecurityConfig(tc.MemoryStore, tc.RootCA, tc.KeyReadWriter, role, org, tc.TempDir, withNonSigningRoot)
}

// WriteNewNodeConfigOrg returns security config for a new node, given a role and an org
// saving the generated key and certificates to disk
func (tc *TestCA) WriteNewNodeConfigOrg(role, org string) (*ca.SecurityConfig, error) {
	withNonSigningRoot := tc.ExternalSigningServer != nil
	return genSecurityConfig(tc.MemoryStore, tc.RootCA, tc.KeyReadWriter, role, org, tc.TempDir, withNonSigningRoot)
}

// External controls whether or not NewTestCA() will create a TestCA server
// configured to use an external signer or not.
var External bool

// NewTestCA is a helper method that creates a TestCA and a bunch of default
// connections and security configs.
func NewTestCA(t *testing.T, krwGenerators ...func(ca.CertPaths) *ca.KeyReadWriter) *TestCA {
	tempdir, err := ioutil.TempDir("", "swarm-ca-test-")
	require.NoError(t, err)
	paths := ca.NewConfigPaths(tempdir)

	rootCA, err := createAndWriteRootCA("swarm-test-CA", paths.RootCA, ca.DefaultNodeCertExpiration)
	require.NoError(t, err)

	return NewTestCAFromRootCA(t, tempdir, rootCA, krwGenerators)
}

// NewTestCAFromRootCA is a helper method that creates a TestCA and a bunch of default
// connections and security configs, given a temp directory and a RootCA to use for signing.
func NewTestCAFromRootCA(t *testing.T, tempBaseDir string, rootCA ca.RootCA, krwGenerators []func(ca.CertPaths) *ca.KeyReadWriter) *TestCA {
	s := store.NewMemoryStore(&stateutils.MockProposer{})

	paths := ca.NewConfigPaths(tempBaseDir)
	organization := identity.NewID()

	var (
		externalSigningServer *ExternalSigningServer
		externalCAs           []*api.ExternalCA
		err                   error
	)

	if External {
		// Start the CA API server.
		externalSigningServer, err = NewExternalSigningServer(rootCA, tempBaseDir)
		assert.NoError(t, err)
		externalCAs = []*api.ExternalCA{
			{
				Protocol: api.ExternalCA_CAProtocolCFSSL,
				URL:      externalSigningServer.URL,
			},
		}
	}

	krw := ca.NewKeyReadWriter(paths.Node, nil, nil)
	if len(krwGenerators) > 0 {
		krw = krwGenerators[0](paths.Node)
	}

	managerConfig, err := genSecurityConfig(s, rootCA, krw, ca.ManagerRole, organization, "", External)
	assert.NoError(t, err)

	managerDiffOrgConfig, err := genSecurityConfig(s, rootCA, krw, ca.ManagerRole, "swarm-test-org-2", "", External)
	assert.NoError(t, err)

	workerConfig, err := genSecurityConfig(s, rootCA, krw, ca.WorkerRole, organization, "", External)
	assert.NoError(t, err)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	baseOpts := []grpc.DialOption{grpc.WithTimeout(10 * time.Second)}
	insecureClientOpts := append(baseOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	clientOpts := append(baseOpts, grpc.WithTransportCredentials(workerConfig.ClientTLSCreds))
	managerOpts := append(baseOpts, grpc.WithTransportCredentials(managerConfig.ClientTLSCreds))
	managerDiffOrgOpts := append(baseOpts, grpc.WithTransportCredentials(managerDiffOrgConfig.ClientTLSCreds))

	conn1, err := grpc.Dial(l.Addr().String(), insecureClientOpts...)
	assert.NoError(t, err)

	conn2, err := grpc.Dial(l.Addr().String(), clientOpts...)
	assert.NoError(t, err)

	conn3, err := grpc.Dial(l.Addr().String(), managerOpts...)
	assert.NoError(t, err)

	conn4, err := grpc.Dial(l.Addr().String(), managerDiffOrgOpts...)
	assert.NoError(t, err)

	serverOpts := []grpc.ServerOption{grpc.Creds(managerConfig.ServerTLSCreds)}
	grpcServer := grpc.NewServer(serverOpts...)

	clusterObj := createClusterObject(t, s, organization, &rootCA, externalCAs...)

	caServer := ca.NewServer(s, managerConfig, paths.RootCA)
	caServer.SetReconciliationRetryInterval(50 * time.Millisecond)
	caServer.SetRootReconciliationInterval(50 * time.Millisecond)
	api.RegisterCAServer(grpcServer, caServer)
	api.RegisterNodeCAServer(grpcServer, caServer)

	ctx := context.Background()

	clusterWatch, clusterWatchCancel, err := store.ViewAndWatch(
		s, func(tx store.ReadTx) error {
			cluster := store.GetCluster(tx, organization)
			caServer.UpdateRootCA(ctx, cluster)
			return nil
		},
		api.EventUpdateCluster{
			Cluster: &api.Cluster{ID: organization},
			Checks:  []api.ClusterCheckFunc{api.ClusterCheckID},
		},
	)
	assert.NoError(t, err)
	go func() {
		for {
			select {
			case event := <-clusterWatch:
				clusterEvent := event.(api.EventUpdateCluster)
				if err := caServer.UpdateRootCA(ctx, clusterEvent.Cluster); err != nil {
					log.G(ctx).WithError(err).Error("ca utils CA server could not update root CA")
				}
			case <-ctx.Done():
				clusterWatchCancel()
				return
			}
		}
	}()

	go grpcServer.Serve(l)
	go caServer.Run(ctx)

	// Wait for caServer to be ready to serve
	<-caServer.Ready()
	remotes := remotes.NewRemotes(api.Peer{Addr: l.Addr().String()})

	caClients := []api.CAClient{api.NewCAClient(conn1), api.NewCAClient(conn2), api.NewCAClient(conn3)}
	nodeCAClients := []api.NodeCAClient{api.NewNodeCAClient(conn1), api.NewNodeCAClient(conn2), api.NewNodeCAClient(conn3), api.NewNodeCAClient(conn4)}
	conns := []*grpc.ClientConn{conn1, conn2, conn3, conn4}

	return &TestCA{
		RootCA:                rootCA,
		ExternalSigningServer: externalSigningServer,
		MemoryStore:           s,
		TempDir:               tempBaseDir,
		Organization:          organization,
		Paths:                 paths,
		Context:               ctx,
		CAClients:             caClients,
		NodeCAClients:         nodeCAClients,
		Conns:                 conns,
		Addr:                  l.Addr().String(),
		Server:                grpcServer,
		ServingSecurityConfig: managerConfig,
		CAServer:              caServer,
		WorkerToken:           clusterObj.RootCA.JoinTokens.Worker,
		ManagerToken:          clusterObj.RootCA.JoinTokens.Manager,
		ConnBroker:            connectionbroker.New(remotes),
		KeyReadWriter:         krw,
		watchCancel:           clusterWatchCancel,
	}
}

func createNode(s *store.MemoryStore, nodeID, role string, csr, cert []byte) error {
	apiRole, _ := ca.FormatRole(role)

	err := s.Update(func(tx store.Tx) error {
		node := &api.Node{
			ID: nodeID,
			Certificate: api.Certificate{
				CSR:  csr,
				CN:   nodeID,
				Role: apiRole,
				Status: api.IssuanceStatus{
					State: api.IssuanceStateIssued,
				},
				Certificate: cert,
			},
			Spec: api.NodeSpec{
				DesiredRole: apiRole,
				Membership:  api.NodeMembershipAccepted,
			},
			Role: apiRole,
		}

		return store.CreateNode(tx, node)
	})

	return err
}

func genSecurityConfig(s *store.MemoryStore, rootCA ca.RootCA, krw *ca.KeyReadWriter, role, org, tmpDir string, nonSigningRoot bool) (*ca.SecurityConfig, error) {
	req := &cfcsr.CertificateRequest{
		KeyRequest: cfcsr.NewBasicKeyRequest(),
	}

	csr, key, err := cfcsr.ParseRequest(req)
	if err != nil {
		return nil, err
	}

	// Obtain a signed Certificate
	nodeID := identity.NewID()

	certChain, err := rootCA.ParseValidateAndSignCSR(csr, nodeID, role, org)
	if err != nil {
		return nil, err
	}

	// If we were instructed to persist the files
	if tmpDir != "" {
		paths := ca.NewConfigPaths(tmpDir)
		if err := ioutil.WriteFile(paths.Node.Cert, certChain, 0644); err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(paths.Node.Key, key, 0600); err != nil {
			return nil, err
		}
	}

	// Load a valid tls.Certificate from the chain and the key
	nodeCert, err := tls.X509KeyPair(certChain, key)
	if err != nil {
		return nil, err
	}

	err = createNode(s, nodeID, role, csr, certChain)
	if err != nil {
		return nil, err
	}

	signingCert := rootCA.Certs
	if len(rootCA.Intermediates) > 0 {
		signingCert = rootCA.Intermediates
	}
	parsedCert, err := helpers.ParseCertificatePEM(signingCert)
	if err != nil {
		return nil, err
	}

	if nonSigningRoot {
		rootCA = ca.RootCA{
			Certs:         rootCA.Certs,
			Digest:        rootCA.Digest,
			Pool:          rootCA.Pool,
			Intermediates: rootCA.Intermediates,
		}
	}

	return ca.NewSecurityConfig(&rootCA, krw, &nodeCert, &ca.IssuerInfo{
		PublicKey: parsedCert.RawSubjectPublicKeyInfo,
		Subject:   parsedCert.RawSubject,
	})
}

func createClusterObject(t *testing.T, s *store.MemoryStore, clusterID string, rootCA *ca.RootCA, externalCAs ...*api.ExternalCA) *api.Cluster {
	cluster := &api.Cluster{
		ID: clusterID,
		Spec: api.ClusterSpec{
			Annotations: api.Annotations{
				Name: store.DefaultClusterName,
			},
			CAConfig: api.CAConfig{
				ExternalCAs: externalCAs,
			},
		},
		RootCA: api.RootCA{
			CACert: rootCA.Certs,
			JoinTokens: api.JoinTokens{
				Worker:  ca.GenerateJoinToken(rootCA),
				Manager: ca.GenerateJoinToken(rootCA),
			},
		},
	}
	if s, err := rootCA.Signer(); err == nil && !External {
		cluster.RootCA.CAKey = s.Key
	}
	assert.NoError(t, s.Update(func(tx store.Tx) error {
		store.CreateCluster(tx, cluster)
		return nil
	}))
	return cluster
}

// CreateRootCertAndKey returns a generated certificate and key for a root CA
func CreateRootCertAndKey(rootCN string) ([]byte, []byte, error) {
	// Create a simple CSR for the CA using the default CA validator and policy
	req := cfcsr.CertificateRequest{
		CN:         rootCN,
		KeyRequest: cfcsr.NewBasicKeyRequest(),
		CA:         &cfcsr.CAConfig{Expiry: ca.RootCAExpiration},
	}

	// Generate the CA and get the certificate and private key
	cert, _, key, err := initca.New(&req)
	return cert, key, err
}

// createAndWriteRootCA creates a Certificate authority for a new Swarm Cluster.
// We're copying ca.CreateRootCA, so we can have smaller key-sizes for tests
func createAndWriteRootCA(rootCN string, paths ca.CertPaths, expiry time.Duration) (ca.RootCA, error) {
	cert, key, err := CreateRootCertAndKey(rootCN)
	if err != nil {
		return ca.RootCA{}, err
	}

	rootCA, err := ca.NewRootCA(cert, cert, key, ca.DefaultNodeCertExpiration, nil)
	if err != nil {
		return ca.RootCA{}, err
	}

	// Ensure directory exists
	err = os.MkdirAll(filepath.Dir(paths.Cert), 0755)
	if err != nil {
		return ca.RootCA{}, err
	}

	// Write the Private Key and Certificate to disk, using decent permissions
	if err := ioutils.AtomicWriteFile(paths.Cert, cert, 0644); err != nil {
		return ca.RootCA{}, err
	}
	if err := ioutils.AtomicWriteFile(paths.Key, key, 0600); err != nil {
		return ca.RootCA{}, err
	}
	return rootCA, nil
}

// ReDateCert takes an existing cert and changes the not before and not after date, to make it easier
// to test expiry
func ReDateCert(t *testing.T, cert, signerCert, signerKey []byte, notBefore, notAfter time.Time) []byte {
	signee, err := helpers.ParseCertificatePEM(cert)
	require.NoError(t, err)
	signer, err := helpers.ParseCertificatePEM(signerCert)
	require.NoError(t, err)
	key, err := helpers.ParsePrivateKeyPEM(signerKey)
	require.NoError(t, err)
	signee.NotBefore = notBefore
	signee.NotAfter = notAfter

	derBytes, err := x509.CreateCertificate(cryptorand.Reader, signee, signer, signee.PublicKey, key)
	require.NoError(t, err)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
}

// CreateCertFromSigner creates a Certificate authority for a new Swarm Cluster given an existing key only.
func CreateCertFromSigner(rootCN string, priv crypto.Signer) ([]byte, error) {
	req := cfcsr.CertificateRequest{
		CN:         rootCN,
		KeyRequest: &cfcsr.BasicKeyRequest{A: ca.RootKeyAlgo, S: ca.RootKeySize},
		CA:         &cfcsr.CAConfig{Expiry: ca.RootCAExpiration},
	}
	cert, _, err := initca.NewFromSigner(&req, priv)
	return cert, err
}
