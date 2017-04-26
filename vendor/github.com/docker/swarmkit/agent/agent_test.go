package agent

import (
	"errors"
	"fmt"
	"testing"
	"time"

	events "github.com/docker/go-events"
	agentutils "github.com/docker/swarmkit/agent/testutils"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	cautils "github.com/docker/swarmkit/ca/testutils"
	"github.com/docker/swarmkit/connectionbroker"
	"github.com/docker/swarmkit/remotes"
	"github.com/docker/swarmkit/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestAgent(t *testing.T) {
	// TODO(stevvooe): The current agent is fairly monolithic, making it hard
	// to test without implementing or mocking an entire master. We'd like to
	// avoid this, as these kinds of tests are expensive to maintain.
	//
	// To support a proper testing program, the plan is to decouple the agent
	// into the following components:
	//
	// 	Connection: Manages the RPC connection and the available managers. Must
	// 	follow lazy grpc style but also expose primitives to force reset, which
	// 	is currently exposed through remotes.
	//
	//	Session: Manages the lifecycle of an agent from Register to a failure.
	//	Currently, this is implemented as Agent.session but we'd prefer to
	//	encapsulate it to keep the agent simple.
	//
	// 	Agent: With the above scaffolding, the agent reduces to Agent.Assign
	// 	and Agent.Watch. Testing becomes as simple as assigning tasks sets and
	// 	checking that the appropriate events come up on the watch queue.
	//
	// We may also move the Assign/Watch to a Driver type and have the agent
	// oversee everything.
}

func TestAgentStartStop(t *testing.T) {
	tc := cautils.NewTestCA(t)
	defer tc.Stop()

	agentSecurityConfig, err := tc.NewNodeConfig(ca.WorkerRole)
	require.NoError(t, err)

	addr := "localhost:4949"
	remotes := remotes.NewRemotes(api.Peer{Addr: addr})

	db, cleanup := storageTestEnv(t)
	defer cleanup()

	agent, err := New(&Config{
		Executor:    &agentutils.TestExecutor{},
		ConnBroker:  connectionbroker.New(remotes),
		Credentials: agentSecurityConfig.ClientTLSCreds,
		DB:          db,
		NodeTLSInfo: &api.NodeTLSInfo{},
	})
	require.NoError(t, err)
	assert.NotNil(t, agent)

	ctx, _ := context.WithTimeout(context.Background(), 5000*time.Millisecond)

	assert.Equal(t, errAgentNotStarted, agent.Stop(ctx))
	assert.NoError(t, agent.Start(ctx))

	if err := agent.Start(ctx); err != errAgentStarted {
		t.Fatalf("expected agent started error: %v", err)
	}

	assert.NoError(t, agent.Stop(ctx))
}

func TestHandleSessionMessageNetworkManagerChanges(t *testing.T) {
	nodeChangeCh := make(chan *NodeChanges, 1)
	defer close(nodeChangeCh)
	tester := agentTestEnv(t, nodeChangeCh, nil)
	defer tester.cleanup()

	currSession, closedSessions := tester.dispatcher.GetSessions()
	require.NotNil(t, currSession)
	require.NotNil(t, currSession.Description)
	require.Empty(t, closedSessions)

	var messages = []*api.SessionMessage{
		{
			Managers: []*api.WeightedPeer{
				{&api.Peer{NodeID: "node1", Addr: "10.0.0.1"}, 1.0}},
			NetworkBootstrapKeys: []*api.EncryptionKey{{}},
		},
		{
			Managers: []*api.WeightedPeer{
				{&api.Peer{NodeID: "node1", Addr: ""}, 1.0}},
			NetworkBootstrapKeys: []*api.EncryptionKey{{}},
		},
		{
			Managers: []*api.WeightedPeer{
				{&api.Peer{NodeID: "node1", Addr: "10.0.0.1"}, 1.0}},
			NetworkBootstrapKeys: nil,
		},
		{
			Managers: []*api.WeightedPeer{
				{&api.Peer{NodeID: "", Addr: "10.0.0.1"}, 1.0}},
			NetworkBootstrapKeys: []*api.EncryptionKey{{}},
		},
		{
			Managers: []*api.WeightedPeer{
				{&api.Peer{NodeID: "node1", Addr: "10.0.0.1"}, 0.0}},
			NetworkBootstrapKeys: []*api.EncryptionKey{{}},
		},
	}

	for _, m := range messages {
		m.SessionID = currSession.SessionID
		tester.dispatcher.SessionMessageChannel() <- m
		select {
		case nodeChange := <-nodeChangeCh:
			require.FailNow(t, "there should be no node changes with these messages: %v", nodeChange)
		case <-time.After(100 * time.Millisecond):
		}
	}

	currSession, closedSessions = tester.dispatcher.GetSessions()
	require.NotEmpty(t, currSession)
	require.Empty(t, closedSessions)
}

func TestHandleSessionMessageNodeChanges(t *testing.T) {
	nodeChangeCh := make(chan *NodeChanges, 1)
	defer close(nodeChangeCh)
	tester := agentTestEnv(t, nodeChangeCh, nil)
	defer tester.cleanup()

	currSession, closedSessions := tester.dispatcher.GetSessions()
	require.NotNil(t, currSession)
	require.NotNil(t, currSession.Description)
	require.Empty(t, closedSessions)

	var testcases = []struct {
		msg      *api.SessionMessage
		change   *NodeChanges
		errorMsg string
	}{
		{
			msg: &api.SessionMessage{
				Node: &api.Node{},
			},
			change:   &NodeChanges{Node: &api.Node{}},
			errorMsg: "the node changed, but no notification of node change",
		},
		{
			msg: &api.SessionMessage{
				RootCA: []byte("new root CA"),
			},
			change:   &NodeChanges{RootCert: []byte("new root CA")},
			errorMsg: "the root cert changed, but no notification of node change",
		},
		{
			msg: &api.SessionMessage{
				Node:   &api.Node{ID: "something"},
				RootCA: []byte("new root CA"),
			},
			change: &NodeChanges{
				Node:     &api.Node{ID: "something"},
				RootCert: []byte("new root CA"),
			},
			errorMsg: "the root cert and node both changed, but no notification of node change",
		},
		{
			msg: &api.SessionMessage{
				Node:   &api.Node{ID: "something"},
				RootCA: tester.testCA.RootCA.Certs,
			},
			errorMsg: "while a node and root cert were provided, nothing has changed so no node changed",
		},
	}

	for _, tc := range testcases {
		tc.msg.SessionID = currSession.SessionID
		tester.dispatcher.SessionMessageChannel() <- tc.msg
		if tc.change != nil {
			select {
			case nodeChange := <-nodeChangeCh:
				require.Equal(t, tc.change, nodeChange, tc.errorMsg)
			case <-time.After(100 * time.Millisecond):
				require.FailNow(t, tc.errorMsg)
			}
		} else {
			select {
			case nodeChange := <-nodeChangeCh:
				require.FailNow(t, "%s: but got change: %v", tc.errorMsg, nodeChange)
			case <-time.After(100 * time.Millisecond):
			}
		}
	}

	currSession, closedSessions = tester.dispatcher.GetSessions()
	require.NotEmpty(t, currSession)
	require.Empty(t, closedSessions)
}

// when the node description changes, the session is restarted and propagated up to the dispatcher
func TestSessionRestartedOnNodeDescriptionChange(t *testing.T) {
	tlsCh := make(chan events.Event, 1)
	defer close(tlsCh)
	tester := agentTestEnv(t, nil, tlsCh)
	defer tester.cleanup()

	currSession, closedSessions := tester.dispatcher.GetSessions()
	require.NotNil(t, currSession)
	require.NotNil(t, currSession.Description)
	require.Empty(t, closedSessions)

	tester.executor.UpdateNodeDescription(&api.NodeDescription{
		Hostname: "testAgent",
	})
	var gotSession *api.SessionRequest
	require.NoError(t, testutils.PollFuncWithTimeout(nil, func() error {
		gotSession, closedSessions = tester.dispatcher.GetSessions()
		if gotSession == nil {
			return errors.New("no current session")
		}
		if len(closedSessions) != 1 {
			return fmt.Errorf("expecting 1 closed sessions, got %d", len(closedSessions))
		}
		return nil
	}, 2*time.Second))
	require.NotEqual(t, currSession, gotSession)
	require.NotNil(t, gotSession.Description)
	require.Equal(t, "testAgent", gotSession.Description.Hostname)
	currSession = gotSession

	newTLSInfo := &api.NodeTLSInfo{
		TrustRoot:           cautils.ECDSA256SHA256Cert,
		CertIssuerPublicKey: []byte("public key"),
		CertIssuerSubject:   []byte("subject"),
	}
	tlsCh <- newTLSInfo
	require.NoError(t, testutils.PollFuncWithTimeout(nil, func() error {
		gotSession, closedSessions = tester.dispatcher.GetSessions()
		if gotSession == nil {
			return errors.New("no current session")
		}
		if len(closedSessions) != 2 {
			return fmt.Errorf("expecting 2 closed sessions, got %d", len(closedSessions))
		}
		return nil
	}, 2*time.Second))
	require.NotEqual(t, currSession, gotSession)
	require.NotNil(t, gotSession.Description)
	require.Equal(t, "testAgent", gotSession.Description.Hostname)
	require.Equal(t, newTLSInfo, gotSession.Description.TLSInfo)
}

type agentTester struct {
	agent      *Agent
	dispatcher *agentutils.MockDispatcher
	executor   *agentutils.TestExecutor
	cleanup    func()
	testCA     *cautils.TestCA
}

func agentTestEnv(t *testing.T, nodeChangeCh chan *NodeChanges, tlsChangeCh chan events.Event) *agentTester {
	var cleanup []func()
	tc := cautils.NewTestCA(t)
	cleanup = append(cleanup, tc.Stop)

	agentSecurityConfig, err := tc.NewNodeConfig(ca.WorkerRole)
	require.NoError(t, err)
	managerSecurityConfig, err := tc.NewNodeConfig(ca.ManagerRole)
	require.NoError(t, err)

	mockDispatcher, mockDispatcherStop := agentutils.NewMockDispatcher(t, managerSecurityConfig)
	cleanup = append(cleanup, mockDispatcherStop)

	remotes := remotes.NewRemotes(api.Peer{Addr: mockDispatcher.Addr})

	db, cleanupStorage := storageTestEnv(t)
	cleanup = append(cleanup, func() { cleanupStorage() })

	executor := &agentutils.TestExecutor{}

	agent, err := New(&Config{
		Executor:         executor,
		ConnBroker:       connectionbroker.New(remotes),
		Credentials:      agentSecurityConfig.ClientTLSCreds,
		DB:               db,
		NotifyNodeChange: nodeChangeCh,
		NotifyTLSChange:  tlsChangeCh,
		NodeTLSInfo: &api.NodeTLSInfo{
			TrustRoot:           tc.RootCA.Certs,
			CertIssuerPublicKey: agentSecurityConfig.IssuerInfo().PublicKey,
			CertIssuerSubject:   agentSecurityConfig.IssuerInfo().Subject,
		},
	})
	require.NoError(t, err)
	agent.nodeUpdatePeriod = 200 * time.Millisecond

	go agent.Start(context.Background())
	cleanup = append(cleanup, func() {
		agent.Stop(context.Background())
	})

	getErr := make(chan error)
	go func() {
		getErr <- agent.Err(context.Background())
	}()
	select {
	case err := <-getErr:
		require.FailNow(t, "starting agent errored with: %v", err)
	case <-agent.Ready():
	case <-time.After(5 * time.Second):
		require.FailNow(t, "agent not ready within 5 seconds")
	}

	return &agentTester{
		agent:      agent,
		dispatcher: mockDispatcher,
		executor:   executor,
		testCA:     tc,
		cleanup: func() {
			// go in reverse order
			for i := len(cleanup) - 1; i >= 0; i-- {
				cleanup[i]()
			}
		},
	}
}
