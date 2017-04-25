package ca_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	cfconfig "github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/ca/testutils"
	"github.com/docker/swarmkit/manager/state/store"
	"github.com/docker/swarmkit/watch"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadRootCASuccess(t *testing.T) {
	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	// Remove the CA cert
	os.RemoveAll(tc.Paths.RootCA.Cert)

	rootCA, err := ca.DownloadRootCA(tc.Context, tc.Paths.RootCA, tc.WorkerToken, tc.ConnBroker)
	require.NoError(t, err)
	require.NotNil(t, rootCA.Pool)
	require.NotNil(t, rootCA.Certs)
	_, err = rootCA.Signer()
	require.Equal(t, err, ca.ErrNoValidSigner)
	require.Equal(t, tc.RootCA.Certs, rootCA.Certs)

	// Remove the CA cert
	os.RemoveAll(tc.Paths.RootCA.Cert)

	// downloading without a join token also succeeds
	rootCA, err = ca.DownloadRootCA(tc.Context, tc.Paths.RootCA, "", tc.ConnBroker)
	require.NoError(t, err)
	require.NotNil(t, rootCA.Pool)
	require.NotNil(t, rootCA.Certs)
	_, err = rootCA.Signer()
	require.Equal(t, err, ca.ErrNoValidSigner)
	require.Equal(t, tc.RootCA.Certs, rootCA.Certs)
}

func TestDownloadRootCAWrongCAHash(t *testing.T) {
	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	// Remove the CA cert
	os.RemoveAll(tc.Paths.RootCA.Cert)

	// invalid token
	for _, invalid := range []string{
		"invalidtoken", // completely invalid
		"SWMTKN-1-3wkodtpeoipd1u1hi0ykdcdwhw16dk73ulqqtn14b3indz68rf-4myj5xihyto11dg1cn55w8p6", // mistyped
	} {
		_, err := ca.DownloadRootCA(tc.Context, tc.Paths.RootCA, invalid, tc.ConnBroker)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid join token")
	}

	// invalid hash token
	splitToken := strings.Split(tc.ManagerToken, "-")
	splitToken[2] = "1kxftv4ofnc6mt30lmgipg6ngf9luhwqopfk1tz6bdmnkubg0e"
	replacementToken := strings.Join(splitToken, "-")

	os.RemoveAll(tc.Paths.RootCA.Cert)

	_, err := ca.DownloadRootCA(tc.Context, tc.Paths.RootCA, replacementToken, tc.ConnBroker)
	require.Error(t, err)
	require.Contains(t, err.Error(), "remote CA does not match fingerprint.")
}

func TestCreateSecurityConfigEmptyDir(t *testing.T) {
	if testutils.External {
		return // this doesn't require any servers at all
	}
	tc := testutils.NewTestCA(t)
	defer tc.Stop()
	assert.NoError(t, tc.CAServer.Stop())

	// Remove all the contents from the temp dir and try again with a new node
	os.RemoveAll(tc.TempDir)
	krw := ca.NewKeyReadWriter(tc.Paths.Node, nil, nil)
	nodeConfig, err := tc.RootCA.CreateSecurityConfig(tc.Context, krw,
		ca.CertificateRequestConfig{
			Token:      tc.WorkerToken,
			ConnBroker: tc.ConnBroker,
		})
	assert.NoError(t, err)
	assert.NotNil(t, nodeConfig)
	assert.NotNil(t, nodeConfig.ClientTLSCreds)
	assert.NotNil(t, nodeConfig.ServerTLSCreds)
	assert.Equal(t, tc.RootCA, *nodeConfig.RootCA())

	root, err := helpers.ParseCertificatePEM(tc.RootCA.Certs)
	assert.NoError(t, err)

	issuerInfo := nodeConfig.IssuerInfo()
	assert.NotNil(t, issuerInfo)
	assert.Equal(t, root.RawSubjectPublicKeyInfo, issuerInfo.PublicKey)
	assert.Equal(t, root.RawSubject, issuerInfo.Subject)
}

func TestCreateSecurityConfigNoCerts(t *testing.T) {
	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	krw := ca.NewKeyReadWriter(tc.Paths.Node, nil, nil)
	root, err := helpers.ParseCertificatePEM(tc.RootCA.Certs)
	assert.NoError(t, err)

	validateNodeConfig := func(rootCA *ca.RootCA) {
		nodeConfig, err := rootCA.CreateSecurityConfig(tc.Context, krw,
			ca.CertificateRequestConfig{
				Token:      tc.WorkerToken,
				ConnBroker: tc.ConnBroker,
			})
		assert.NoError(t, err)
		assert.NotNil(t, nodeConfig)
		assert.NotNil(t, nodeConfig.ClientTLSCreds)
		assert.NotNil(t, nodeConfig.ServerTLSCreds)
		assert.Equal(t, tc.RootCA, *nodeConfig.RootCA())

		issuerInfo := nodeConfig.IssuerInfo()
		assert.NotNil(t, issuerInfo)
		assert.Equal(t, root.RawSubjectPublicKeyInfo, issuerInfo.PublicKey)
		assert.Equal(t, root.RawSubject, issuerInfo.Subject)
	}

	// Remove only the node certificates form the directory, and attest that we get
	// new certificates that are locally signed
	os.RemoveAll(tc.Paths.Node.Cert)
	validateNodeConfig(&tc.RootCA)

	// Remove only the node certificates form the directory, get a new rootCA, and attest that we get
	// new certificates that are issued by the remote CA
	os.RemoveAll(tc.Paths.Node.Cert)
	rootCA, err := ca.GetLocalRootCA(tc.Paths.RootCA)
	assert.NoError(t, err)
	validateNodeConfig(&rootCA)
}

func TestLoadSecurityConfigExpiredCert(t *testing.T) {
	if testutils.External {
		return // this doesn't require any servers at all
	}
	tc := testutils.NewTestCA(t)
	defer tc.Stop()
	s, err := tc.RootCA.Signer()
	require.NoError(t, err)

	krw := ca.NewKeyReadWriter(tc.Paths.Node, nil, nil)
	now := time.Now()

	_, _, err = tc.RootCA.IssueAndSaveNewCertificates(krw, "cn", "ou", "org")
	require.NoError(t, err)
	certBytes, _, err := krw.Read()
	require.NoError(t, err)

	// A cert that is not yet valid is not valid even if expiry is allowed
	invalidCert := testutils.ReDateCert(t, certBytes, tc.RootCA.Certs, s.Key, now.Add(time.Hour), now.Add(time.Hour*2))
	require.NoError(t, ioutil.WriteFile(tc.Paths.Node.Cert, invalidCert, 0700))

	_, err = ca.LoadSecurityConfig(tc.Context, tc.RootCA, krw, false)
	require.Error(t, err)
	require.IsType(t, x509.CertificateInvalidError{}, errors.Cause(err))

	_, err = ca.LoadSecurityConfig(tc.Context, tc.RootCA, krw, true)
	require.Error(t, err)
	require.IsType(t, x509.CertificateInvalidError{}, errors.Cause(err))

	// a cert that is expired is not valid if expiry is not allowed
	invalidCert = testutils.ReDateCert(t, certBytes, tc.RootCA.Certs, s.Key, now.Add(-2*time.Minute), now.Add(-1*time.Minute))
	require.NoError(t, ioutil.WriteFile(tc.Paths.Node.Cert, invalidCert, 0700))

	_, err = ca.LoadSecurityConfig(tc.Context, tc.RootCA, krw, false)
	require.Error(t, err)
	require.IsType(t, x509.CertificateInvalidError{}, errors.Cause(err))

	// but it is valid if expiry is allowed
	_, err = ca.LoadSecurityConfig(tc.Context, tc.RootCA, krw, true)
	require.NoError(t, err)
}

func TestLoadSecurityConfigInvalidCert(t *testing.T) {
	if testutils.External {
		return // this doesn't require any servers at all
	}
	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	// Write some garbage to the cert
	ioutil.WriteFile(tc.Paths.Node.Cert, []byte(`-----BEGIN CERTIFICATE-----\n
some random garbage\n
-----END CERTIFICATE-----`), 0644)

	krw := ca.NewKeyReadWriter(tc.Paths.Node, nil, nil)

	_, err := ca.LoadSecurityConfig(tc.Context, tc.RootCA, krw, false)
	assert.Error(t, err)
}

func TestLoadSecurityConfigInvalidKey(t *testing.T) {
	if testutils.External {
		return // this doesn't require any servers at all
	}
	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	// Write some garbage to the Key
	ioutil.WriteFile(tc.Paths.Node.Key, []byte(`-----BEGIN EC PRIVATE KEY-----\n
some random garbage\n
-----END EC PRIVATE KEY-----`), 0644)

	krw := ca.NewKeyReadWriter(tc.Paths.Node, nil, nil)

	_, err := ca.LoadSecurityConfig(tc.Context, tc.RootCA, krw, false)
	assert.Error(t, err)
}

func TestLoadSecurityConfigIncorrectPassphrase(t *testing.T) {
	if testutils.External {
		return // this doesn't require any servers at all
	}
	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	paths := ca.NewConfigPaths(tc.TempDir)
	_, _, err := tc.RootCA.IssueAndSaveNewCertificates(ca.NewKeyReadWriter(paths.Node, []byte("kek"), nil),
		"nodeID", ca.WorkerRole, tc.Organization)
	require.NoError(t, err)

	_, err = ca.LoadSecurityConfig(tc.Context, tc.RootCA, ca.NewKeyReadWriter(paths.Node, nil, nil), false)
	require.IsType(t, ca.ErrInvalidKEK{}, err)
}

func TestLoadSecurityConfigIntermediates(t *testing.T) {
	if testutils.External {
		return // this doesn't require any servers at all
	}
	tempdir, err := ioutil.TempDir("", "test-load-config-with-intermediates")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)
	paths := ca.NewConfigPaths(tempdir)
	krw := ca.NewKeyReadWriter(paths.Node, nil, nil)

	rootCA, err := ca.NewRootCA(testutils.ECDSACertChain[2], nil, nil, ca.DefaultNodeCertExpiration, nil)
	require.NoError(t, err)

	// loading the incomplete chain fails
	require.NoError(t, krw.Write(testutils.ECDSACertChain[0], testutils.ECDSACertChainKeys[0], nil))
	_, err = ca.LoadSecurityConfig(context.Background(), rootCA, krw, false)
	require.Error(t, err)

	intermediate, err := helpers.ParseCertificatePEM(testutils.ECDSACertChain[1])
	require.NoError(t, err)

	// loading the complete chain succeeds
	require.NoError(t, krw.Write(append(testutils.ECDSACertChain[0], testutils.ECDSACertChain[1]...), testutils.ECDSACertChainKeys[0], nil))
	secConfig, err := ca.LoadSecurityConfig(context.Background(), rootCA, krw, false)
	require.NoError(t, err)
	require.NotNil(t, secConfig)
	issuerInfo := secConfig.IssuerInfo()
	require.NotNil(t, issuerInfo)
	require.Equal(t, intermediate.RawSubjectPublicKeyInfo, issuerInfo.PublicKey)
	require.Equal(t, intermediate.RawSubject, issuerInfo.Subject)

	// set up a GRPC server using these credentials
	secConfig.ServerTLSCreds.Config().ClientAuth = tls.RequireAndVerifyClientCert
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	serverOpts := []grpc.ServerOption{grpc.Creds(secConfig.ServerTLSCreds)}
	grpcServer := grpc.NewServer(serverOpts...)
	go grpcServer.Serve(l)
	defer grpcServer.Stop()

	// we should be able to connect to the server using the client credentials
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithTransportCredentials(secConfig.ClientTLSCreds),
	}
	conn, err := grpc.Dial(l.Addr().String(), dialOpts...)
	require.NoError(t, err)
	conn.Close()
}

func TestSecurityConfigUpdateRootCA(t *testing.T) {
	tc := testutils.NewTestCA(t)
	defer tc.Stop()
	tcConfig, err := tc.NewNodeConfig("worker")
	require.NoError(t, err)

	// create the "original" security config, and we'll update it to trust the test server's
	cert, key, err := testutils.CreateRootCertAndKey("root1")
	require.NoError(t, err)
	rootCA, err := ca.NewRootCA(cert, cert, key, ca.DefaultNodeCertExpiration, nil)
	require.NoError(t, err)

	tempdir, err := ioutil.TempDir("", "test-security-config-update")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)
	configPaths := ca.NewConfigPaths(tempdir)

	secConfig, err := rootCA.CreateSecurityConfig(context.Background(),
		ca.NewKeyReadWriter(configPaths.Node, nil, nil), ca.CertificateRequestConfig{})
	require.NoError(t, err)
	// update the server TLS to require certificates, otherwise this will all pass
	// even if the root pools aren't updated
	secConfig.ServerTLSCreds.Config().ClientAuth = tls.RequireAndVerifyClientCert

	// set up a GRPC server using these credentials
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	serverOpts := []grpc.ServerOption{grpc.Creds(secConfig.ServerTLSCreds)}
	grpcServer := grpc.NewServer(serverOpts...)
	go grpcServer.Serve(l)
	defer grpcServer.Stop()

	// we should not be able to connect to the test CA server using the original security config, and should not
	// be able to connect to new server using the test CA's client credentials
	dialOptsBase := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(10 * time.Second),
	}
	dialOpts := append(dialOptsBase, grpc.WithTransportCredentials(secConfig.ClientTLSCreds))
	_, err = grpc.Dial(tc.Addr, dialOpts...)
	require.Error(t, err)
	require.IsType(t, x509.UnknownAuthorityError{}, err)

	dialOpts = append(dialOptsBase, grpc.WithTransportCredentials(tcConfig.ClientTLSCreds))
	_, err = grpc.Dial(l.Addr().String(), dialOpts...)
	require.Error(t, err)
	require.IsType(t, x509.UnknownAuthorityError{}, err)

	// we can't connect to the test CA's external server either
	csr, _, err := ca.GenerateNewCSR()
	require.NoError(t, err)
	req := ca.PrepareCSR(csr, "cn", ca.ManagerRole, secConfig.ClientTLSCreds.Organization())

	externalServer := tc.ExternalSigningServer
	tcSigner, err := tc.RootCA.Signer()
	require.NoError(t, err)
	if testutils.External {
		// stop the external server and create a new one because the external server actually has to trust our client certs as well.
		updatedRoot, err := ca.NewRootCA(append(tc.RootCA.Certs, cert...), tcSigner.Cert, tcSigner.Key, ca.DefaultNodeCertExpiration, nil)
		require.NoError(t, err)
		externalServer, err = testutils.NewExternalSigningServer(updatedRoot, tc.TempDir)
		require.NoError(t, err)
		defer externalServer.Stop()

		secConfig.ExternalCA().UpdateURLs(externalServer.URL)
		_, err = secConfig.ExternalCA().Sign(context.Background(), req)
		require.Error(t, err)
		// the type is weird (it's wrapped in a bunch of other things in ctxhttp), so just compare strings
		require.Contains(t, err.Error(), x509.UnknownAuthorityError{}.Error())
	}

	// update the root CA on the "original"" security config to support both the old root
	// and the "new root" (the testing CA root)
	rSigner, err := rootCA.Signer()
	require.NoError(t, err)
	updatedRootCA, err := ca.NewRootCA(append(rootCA.Certs, tc.RootCA.Certs...), rSigner.Cert, rSigner.Key, ca.DefaultNodeCertExpiration, nil)
	require.NoError(t, err)
	err = secConfig.UpdateRootCA(&updatedRootCA, updatedRootCA.Pool)
	require.NoError(t, err)

	// can now connect to the test CA using our modified security config, and can cannect to our server using
	// the test CA config
	conn, err := grpc.Dial(tc.Addr, dialOpts...)
	require.NoError(t, err)
	conn.Close()

	dialOpts = append(dialOptsBase, grpc.WithTransportCredentials(secConfig.ClientTLSCreds))
	conn, err = grpc.Dial(tc.Addr, dialOpts...)
	require.NoError(t, err)
	conn.Close()

	// we can also now connect to the test CA's external signing server
	if testutils.External {
		secConfig.ExternalCA().UpdateURLs(externalServer.URL)
		_, err := secConfig.ExternalCA().Sign(context.Background(), req)
		require.NoError(t, err)
	}
}

func TestSecurityConfigSetWatch(t *testing.T) {
	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	secConfig, err := tc.NewNodeConfig(ca.ManagerRole)
	require.NoError(t, err)
	issuer := secConfig.IssuerInfo()

	w := watch.NewQueue()
	defer w.Close()
	secConfig.SetWatch(w)

	configWatch, configCancel := w.Watch()
	defer configCancel()

	require.NoError(t, ca.RenewTLSConfigNow(context.Background(), secConfig, tc.ConnBroker))
	select {
	case ev := <-configWatch:
		nodeTLSInfo, ok := ev.(*api.NodeTLSInfo)
		require.True(t, ok)
		require.Equal(t, &api.NodeTLSInfo{
			TrustRoot:           tc.RootCA.Certs,
			CertIssuerPublicKey: issuer.PublicKey,
			CertIssuerSubject:   issuer.Subject,
		}, nodeTLSInfo)
	case <-time.After(time.Second):
		require.FailNow(t, "on TLS certificate update, we should have gotten a security config update")
	}

	require.NoError(t, secConfig.UpdateRootCA(&tc.RootCA, tc.RootCA.Pool))
	select {
	case ev := <-configWatch:
		nodeTLSInfo, ok := ev.(*api.NodeTLSInfo)
		require.True(t, ok)
		require.Equal(t, &api.NodeTLSInfo{
			TrustRoot:           tc.RootCA.Certs,
			CertIssuerPublicKey: issuer.PublicKey,
			CertIssuerSubject:   issuer.Subject,
		}, nodeTLSInfo)
	case <-time.After(time.Second):
		require.FailNow(t, "on TLS certificate update, we should have gotten a security config update")
	}

	configCancel()
	w.Close()

	// ensure that we can still update tls certs and roots without error even though the watch is closed
	require.NoError(t, secConfig.UpdateRootCA(&tc.RootCA, tc.RootCA.Pool))
	require.NoError(t, ca.RenewTLSConfigNow(context.Background(), secConfig, tc.ConnBroker))
}

// enforce that no matter what order updating the root CA and updating TLS credential happens, we
// end up with a security config that has updated certs, and an updated root pool
func TestRenewTLSConfigUpdateRootCARace(t *testing.T) {
	tc := testutils.NewTestCA(t)
	defer tc.Stop()
	paths := ca.NewConfigPaths(tc.TempDir)

	secConfig, err := tc.WriteNewNodeConfig(ca.ManagerRole)
	require.NoError(t, err)

	leafCert, err := ioutil.ReadFile(paths.Node.Cert)
	require.NoError(t, err)

	cert, key, err := testutils.CreateRootCertAndKey("extra root cert for external CA")
	require.NoError(t, err)
	extraExternalRootCA, err := ca.NewRootCA(append(cert, tc.RootCA.Certs...), cert, key, ca.DefaultNodeCertExpiration, nil)
	require.NoError(t, err)
	extraExternalServer, err := testutils.NewExternalSigningServer(extraExternalRootCA, tc.TempDir)
	require.NoError(t, err)
	defer extraExternalServer.Stop()
	secConfig.ExternalCA().UpdateURLs(extraExternalServer.URL)

	externalPool := x509.NewCertPool()
	externalPool.AppendCertsFromPEM(tc.RootCA.Certs)
	externalPool.AppendCertsFromPEM(cert)

	csr, _, err := ca.GenerateNewCSR()
	require.NoError(t, err)
	signReq := ca.PrepareCSR(csr, "cn", ca.WorkerRole, tc.Organization)

	for i := 0; i < 5; i++ {
		cert, _, err := testutils.CreateRootCertAndKey(fmt.Sprintf("root %d", i+2))
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done1, done2 := make(chan struct{}), make(chan struct{})
		rootCA := secConfig.RootCA()
		go func() {
			defer close(done1)
			s := ca.LocalSigner{}
			if signer, err := rootCA.Signer(); err == nil {
				s = *signer
			}
			updatedRootCA, err := ca.NewRootCA(append(rootCA.Certs, cert...), s.Cert, s.Key, ca.DefaultNodeCertExpiration, nil)
			require.NoError(t, err)
			externalPool.AppendCertsFromPEM(cert)
			require.NoError(t, secConfig.UpdateRootCA(&updatedRootCA, externalPool))
		}()

		go func() {
			defer close(done2)
			require.NoError(t, ca.RenewTLSConfigNow(ctx, secConfig, tc.ConnBroker))
		}()

		<-done1
		<-done2

		newCert, err := ioutil.ReadFile(paths.Node.Cert)
		require.NoError(t, err)

		require.NotEqual(t, newCert, leafCert)
		leafCert = newCert

		// at the start of this loop had i+1 certs, afterward should have added one more
		require.Len(t, secConfig.ClientTLSCreds.Config().RootCAs.Subjects(), i+2)
		require.Len(t, secConfig.ServerTLSCreds.Config().RootCAs.Subjects(), i+2)
		// no matter what, the external CA still has the extra external CA root cert
		_, err = secConfig.ExternalCA().Sign(context.Background(), signReq)
		require.NoError(t, err)
	}
}

func writeAlmostExpiringCertToDisk(t *testing.T, tc *testutils.TestCA, cn, ou, org string) {
	s, err := tc.RootCA.Signer()
	require.NoError(t, err)

	// Create a new RootCA, and change the policy to issue 6 minute certificates
	// Because of the default backdate of 5 minutes, this issues certificates
	// valid for 1 minute.
	newRootCA, err := ca.NewRootCA(tc.RootCA.Certs, s.Cert, s.Key, ca.DefaultNodeCertExpiration, nil)
	assert.NoError(t, err)
	newSigner, err := newRootCA.Signer()
	require.NoError(t, err)
	newSigner.SetPolicy(&cfconfig.Signing{
		Default: &cfconfig.SigningProfile{
			Usage:  []string{"signing", "key encipherment", "server auth", "client auth"},
			Expiry: 6 * time.Minute,
		},
	})

	// Issue a new certificate with the same details as the current config, but with 1 min expiration time, and
	// overwrite the existing cert on disk
	_, _, err = newRootCA.IssueAndSaveNewCertificates(ca.NewKeyReadWriter(tc.Paths.Node, nil, nil), cn, ou, org)
	assert.NoError(t, err)
}

func TestRenewTLSConfigWorker(t *testing.T) {
	t.Parallel()

	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get a new nodeConfig with a TLS cert that has the default Cert duration, but overwrite
	// the cert on disk with one that expires in 1 minute
	nodeConfig, err := tc.WriteNewNodeConfig(ca.WorkerRole)
	assert.NoError(t, err)
	c := nodeConfig.ClientTLSCreds
	writeAlmostExpiringCertToDisk(t, tc, c.NodeID(), c.Role(), c.Organization())

	renew := make(chan struct{})
	updates := ca.RenewTLSConfig(ctx, nodeConfig, tc.ConnBroker, renew)
	select {
	case <-time.After(10 * time.Second):
		assert.Fail(t, "TestRenewTLSConfig timed-out")
	case certUpdate := <-updates:
		assert.NoError(t, certUpdate.Err)
		assert.NotNil(t, certUpdate)
		assert.Equal(t, ca.WorkerRole, certUpdate.Role)
	}

	root, err := helpers.ParseCertificatePEM(tc.RootCA.Certs)
	assert.NoError(t, err)

	issuerInfo := nodeConfig.IssuerInfo()
	assert.NotNil(t, issuerInfo)
	assert.Equal(t, root.RawSubjectPublicKeyInfo, issuerInfo.PublicKey)
	assert.Equal(t, root.RawSubject, issuerInfo.Subject)
}

func TestRenewTLSConfigManager(t *testing.T) {
	t.Parallel()

	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get a new nodeConfig with a TLS cert that has the default Cert duration, but overwrite
	// the cert on disk with one that expires in 1 minute
	nodeConfig, err := tc.WriteNewNodeConfig(ca.WorkerRole)
	assert.NoError(t, err)
	c := nodeConfig.ClientTLSCreds
	writeAlmostExpiringCertToDisk(t, tc, c.NodeID(), c.Role(), c.Organization())

	renew := make(chan struct{})
	updates := ca.RenewTLSConfig(ctx, nodeConfig, tc.ConnBroker, renew)
	select {
	case <-time.After(10 * time.Second):
		assert.Fail(t, "TestRenewTLSConfig timed-out")
	case certUpdate := <-updates:
		assert.NoError(t, certUpdate.Err)
		assert.NotNil(t, certUpdate)
		assert.Equal(t, ca.WorkerRole, certUpdate.Role)
	}

	root, err := helpers.ParseCertificatePEM(tc.RootCA.Certs)
	assert.NoError(t, err)

	issuerInfo := nodeConfig.IssuerInfo()
	assert.NotNil(t, issuerInfo)
	assert.Equal(t, root.RawSubjectPublicKeyInfo, issuerInfo.PublicKey)
	assert.Equal(t, root.RawSubject, issuerInfo.Subject)
}

func TestRenewTLSConfigWithNoNode(t *testing.T) {
	t.Parallel()

	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get a new nodeConfig with a TLS cert that has the default Cert duration, but overwrite
	// the cert on disk with one that expires in 1 minute
	nodeConfig, err := tc.WriteNewNodeConfig(ca.WorkerRole)
	assert.NoError(t, err)
	c := nodeConfig.ClientTLSCreds
	writeAlmostExpiringCertToDisk(t, tc, c.NodeID(), c.Role(), c.Organization())

	// Delete the node from the backend store
	err = tc.MemoryStore.Update(func(tx store.Tx) error {
		node := store.GetNode(tx, nodeConfig.ClientTLSCreds.NodeID())
		assert.NotNil(t, node)
		return store.DeleteNode(tx, nodeConfig.ClientTLSCreds.NodeID())
	})
	assert.NoError(t, err)

	renew := make(chan struct{})
	updates := ca.RenewTLSConfig(ctx, nodeConfig, tc.ConnBroker, renew)
	select {
	case <-time.After(10 * time.Second):
		assert.Fail(t, "TestRenewTLSConfig timed-out")
	case certUpdate := <-updates:
		assert.Error(t, certUpdate.Err)
		assert.Contains(t, certUpdate.Err.Error(), "not found when attempting to renew certificate")
	}
}

func TestForceRenewTLSConfig(t *testing.T) {
	t.Parallel()

	tc := testutils.NewTestCA(t)
	defer tc.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get a new managerConfig with a TLS cert that has 15 minutes to live
	nodeConfig, err := tc.WriteNewNodeConfig(ca.ManagerRole)
	assert.NoError(t, err)

	renew := make(chan struct{}, 1)
	updates := ca.RenewTLSConfig(ctx, nodeConfig, tc.ConnBroker, renew)
	renew <- struct{}{}
	select {
	case <-time.After(10 * time.Second):
		assert.Fail(t, "TestForceRenewTLSConfig timed-out")
	case certUpdate := <-updates:
		assert.NoError(t, certUpdate.Err)
		assert.NotNil(t, certUpdate)
		assert.Equal(t, certUpdate.Role, ca.ManagerRole)
	}
}
