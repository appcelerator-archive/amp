package etcd

import (
	"github.com/appcelerator/amp/data/storage"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"strings"
	"sync"
)

const (
	// We have set a buffer in order to reduce times of context switches.
	incomingBufSize = 100
	outgoingBufSize = 100
)

// watchChan implements watch.Interface.
type watchChan struct {
	etcd              *etcd
	key               string
	initialRev        int64
	recursive         bool
	filter            storage.Filter
	ctx               context.Context
	cancel            context.CancelFunc
	incomingEventChan chan *event
	resultChan        chan storage.Event
	errChan           chan error
}

type event struct {
	key       string
	value     []byte
	rev       int64
	isDeleted bool
	isCreated bool
}

func (s *etcd) watch(ctx context.Context, key string, rev int64, filter storage.Filter, recursive bool) (storage.WatchInterface, error) {
	if recursive && !strings.HasSuffix(key, "/") {
		key += "/"
	}
	wc := s.createWatchChan(ctx, key, rev, recursive, filter)
	go wc.run()
	return wc, nil
}

func (s *etcd) createWatchChan(ctx context.Context, key string, rev int64, recursive bool, filter storage.Filter) *watchChan {
	wc := &watchChan{
		etcd:              s,
		key:               key,
		initialRev:        rev,
		recursive:         recursive,
		filter:            filter,
		incomingEventChan: make(chan *event, incomingBufSize),
		resultChan:        make(chan storage.Event, outgoingBufSize),
		errChan:           make(chan error, 1),
	}
	wc.ctx, wc.cancel = context.WithCancel(ctx)
	return wc
}

func (wc *watchChan) run() {
	go wc.startWatching()

	var resultChanWG sync.WaitGroup
	resultChanWG.Add(1)
	go wc.processEvent(&resultChanWG)

	select {
	case err := <-wc.errChan:
		errResult := parseError(err)
		if errResult != nil {
			// error result is guaranteed to be received by user before closing ResultChan.
			select {
			case wc.resultChan <- *errResult:
			case <-wc.ctx.Done(): // user has given up all results
			}
		}
		wc.cancel()
	case <-wc.ctx.Done():
	}
	// we need to wait until resultChan wouldn't be sent to anymore
	resultChanWG.Wait()
	close(wc.resultChan)
}

func (wc *watchChan) Stop() {
	wc.cancel()
}

func (wc *watchChan) ResultChan() <-chan storage.Event {
	return wc.resultChan
}

// sync tries to retrieve existing data and send them to process.
// The revision to watch will be set to the revision in response.
func (wc *watchChan) sync() error {
	opts := []clientv3.OpOption{}
	if wc.recursive {
		opts = append(opts, clientv3.WithPrefix())
	}
	getResp, err := wc.etcd.client.Get(wc.ctx, wc.key, opts...)
	if err != nil {
		return err
	}
	wc.initialRev = getResp.Header.Revision

	for _, kv := range getResp.Kvs {
		wc.sendEvent(parseKV(kv))
	}
	return nil
}

// startWatching does:
// - get current objects if initialRev=0; set initialRev to current rev
// - watch on given key and send events to process.
func (wc *watchChan) startWatching() {
	if wc.initialRev == 0 {
		if err := wc.sync(); err != nil {
			wc.sendError(err)
			return
		}
	}
	opts := []clientv3.OpOption{clientv3.WithRev(wc.initialRev + 1)}
	if wc.recursive {
		opts = append(opts, clientv3.WithPrefix())
	}
	wch := wc.etcd.client.Watch(wc.ctx, wc.key, opts...)
	for wres := range wch {
		if wres.Err() != nil {
			// If there is an error on server (e.g. compaction), the channel will return it before closed.
			wc.sendError(wres.Err())
			return
		}
		for _, e := range wres.Events {
			wc.sendEvent(parseEvent(e))
		}
	}
}

// processEvent processes events from etcd watcher and sends results to resultChan.
func (wc *watchChan) processEvent(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case e := <-wc.incomingEventChan:
			res := wc.transform(e)
			if res == nil {
				continue
			}
			// If user couldn't receive results fast enough, we also block incoming events from watcher.
			// Because storing events in local will cause more memory usage.
			// The worst case would be closing the fast watcher.
			select {
			case wc.resultChan <- *res:
			case <-wc.ctx.Done():
				return
			}
		case <-wc.ctx.Done():
			return
		}
	}
}

// transform transforms an event into a result for user if not filtered.
func (wc *watchChan) transform(e *event) *storage.Event {
	event := &storage.Event{
		Key:       e.key,
		Value:     e.value,
		Revision:  e.rev,
		IsCreated: e.isCreated,
		IsDeleted: e.isDeleted,
	}
	return event
}

func (wc *watchChan) sendError(err error) {
	// Context.canceled is an expected behavior.
	// We should just stop all goroutines in watchChan without returning error.
	// TODO: etcd client should return context.Canceled instead of grpc specific error.
	if grpc.Code(err) == codes.Canceled || err == context.Canceled {
		return
	}
	select {
	case wc.errChan <- err:
	case <-wc.ctx.Done():
	}
}

func (wc *watchChan) sendEvent(e *event) {
	if len(wc.incomingEventChan) == incomingBufSize {
		log.Printf("Fast watcher, slow processing. Number of buffered events: %d."+
			"Probably caused by slow decoding, user not receiving fast, or other processing logic",
			incomingBufSize)
	}
	select {
	case wc.incomingEventChan <- e:
	case <-wc.ctx.Done():
	}
}

func parseKV(kv *mvccpb.KeyValue) *event {
	return &event{
		key:       string(kv.Key),
		value:     kv.Value,
		rev:       kv.ModRevision,
		isCreated: kv.ModRevision == kv.CreateRevision,
		isDeleted: false,
	}
}

func parseEvent(e *clientv3.Event) *event {
	return &event{
		key:       string(e.Kv.Key),
		value:     e.Kv.Value,
		rev:       e.Kv.ModRevision,
		isDeleted: e.Type == clientv3.EventTypeDelete,
		isCreated: e.IsCreate(),
	}
}

func parseError(err error) *storage.Event {
	return &storage.Event{
		IsError: true,
		Error:   err,
	}
}
