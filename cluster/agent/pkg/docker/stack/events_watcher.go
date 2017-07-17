package stack

import (
	"context"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type EventsWatcherOptions struct {
	Options types.EventsOptions
}

func NewEventsWatcherOptions(eventTypes ...string) *EventsWatcherOptions {
	o := &EventsWatcherOptions{
		Options: types.EventsOptions{
			Filters: filters.NewArgs(),
		},
	}
	o.AddTypeFilters(eventTypes...)
	return o
}

func NewEventsWatcherOptionsSince(since string, eventTypes ...string) *EventsWatcherOptions {
	o := NewEventsWatcherOptions(eventTypes...)
	o.Since(since)
	return o
}

func NewEventsWatcherOptionsUntil(until string, eventTypes ...string) *EventsWatcherOptions {
	o := NewEventsWatcherOptions(eventTypes...)
	o.Until(until)
	return o
}

func NewEventsWatcherOptionsBetween(since, until string, eventTypes ...string) *EventsWatcherOptions {
	o := NewEventsWatcherOptions(eventTypes...)
	o.Since(since)
	o.Until(until)
	return o
}

// SetEventsOptions is used if to provide an explicit set of options;
// otherwise, use a combination of the other convenience option helpers
// (Since, Until, Add*) before calling Subscribe
func (o *EventsWatcherOptions) SetEventsOptions(opts types.EventsOptions) {
	o.Options = opts
}

// Since is a duration (relative to now) or a RFC3339 timestamp or a unix timestamp
// See: githuo.com/docker/docker/api/types/time/timestamp.go
func (o *EventsWatcherOptions) Since(value string) {
	o.Options.Since = value
}

// Until is a duration (relative to now) or a RFC3339 timestamp or a unix timestamp
// See: githuo.com/docker/docker/api/types/time/timestamp.go
func (o *EventsWatcherOptions) Until(value string) {
	o.Options.Until = value
}

// AddScopeFilter adds a filter to match an event.Scope to the provided scope
func (o *EventsWatcherOptions) AddScopeFilter(scope string) {
	o.Options.Filters.Add("scope", scope)
}

// AddAttributesFilter adds a filter to match the supplied key-value pair
func (o *EventsWatcherOptions) AddAttributesFilter(key, value string) {
	o.Options.Filters.Add("label", key + "=" + value)
}

// AddAttributesFilters adds a filter to match the supplied attributes map
func (o *EventsWatcherOptions) AddAttributesFilters(attrs map[string]string) {
	for k, v := range attrs {
		o.AddAttributesFilter(k, v)
	}
}

// AddImageFilter adds a filter that matches "image" to an image name or ID
func (o *EventsWatcherOptions) AddImageFilter(id string) {
	o.Options.Filters.Add("image", id)
}

// AddTypeFilter adds a filter arg for each specified type (or use the explicit
// convenience helper methods (AddXXXEventTypeFilter))
// See: githuo.com/docker/docker/api/types/events/events.go
// - "container"
// - "daemon"
// - "image"
// - "network"
// - "plugin"
// - "volume"
// - "service"
// - "node"
// - "secret"
// AddTypeFilters adds a list of type filters at once
func (o *EventsWatcherOptions) AddTypeFilters(eventTypes ...string) {
	for _, e := range eventTypes {
		o.Options.Filters.Add("type", e)
	}
}

// AddDaemonEventTypeFilter adds a filter that matches event "type" == "daemon"
func (o *EventsWatcherOptions) AddDaemonEventTypeFilter() {
	o.AddTypeFilters(events.DaemonEventType)
}

// AddContainerEventTypeFilter adds a filter that matches event "type" == "container"
func (o *EventsWatcherOptions) AddContainerEventTypeFilter() {
	o.AddTypeFilters(events.ContainerEventType)
}

// AddPluginEventTypeFilter adds a filter that matches event "type" == "plugin"
func (o *EventsWatcherOptions) AddPluginEventTypeFilter() {
	o.AddTypeFilters(events.PluginEventType)
}

// AddVolumeEventTypeFilter adds a filter that matches event "type" == "volume"
func (o *EventsWatcherOptions) AddVolumeEventTypeFilter() {
	o.AddTypeFilters(events.VolumeEventType)
}

// AddNetworkEventTypeFilter adds a filter that matches event "type" == "network"
func (o *EventsWatcherOptions) AddNetworkEventTypeFilter() {
	o.AddTypeFilters(events.NetworkEventType)
}

// AddServiceEventTypeFilter adds a filter that matches event "type" == "service"
func (o *EventsWatcherOptions) AddServiceEventTypeFilter() {
	o.AddTypeFilters(events.ServiceEventType)
}

// AddNodeEventTypeFilter adds a filter that matches event "type" == "node"
func (o *EventsWatcherOptions) AddNodeEventTypeFilter() {
	o.AddTypeFilters(events.NodeEventType)
}

// AddSecretEventTypeFilter adds a filter that matches event "type" == "secret"
func (o *EventsWatcherOptions) AddSecretEventTypeFilter() {
	o.AddTypeFilters(events.SecretEventType)
}

// EventsWatcher is the interface for watching and reacting to events
// EventsWatcher inspired by github.com/docker/cli/cli/command/event_utils.go:EventHandler
type EventsWatcher interface {
	On(action string, handler func(events.Message))
	OnError(handler func(error))
	Watch()
	Cancel()
}

// NewEventsWatcher returns a new watcher instance
func NewEventsWatcher(ctx context.Context, apiClient client.APIClient, opts *EventsWatcherOptions) EventsWatcher {
	if opts == nil {
		opts = NewEventsWatcherOptions()
	}
	w := &eventsWatcher{
		ctx:       ctx,
		apiClient: apiClient,
		opts:      opts,
		handlers:  make(map[string]func(message events.Message)),
	}
	return w
}

// NewEventsWatcherWithCancel returns a new cancelable watcher instance
func NewEventsWatcherWithCancel(ctx context.Context, apiClient client.APIClient, opts *EventsWatcherOptions) EventsWatcher {
	ctx, cancel := context.WithCancel(ctx)
	w := NewEventsWatcher(ctx, apiClient, opts)
	// downcast to store the cancel func in the eventsWatcher instance
	// should always be ok, so if not let it panic
	ew := w.(*eventsWatcher)
	ew.cancel = cancel
	return w
}

type eventsWatcher struct {
	mu           sync.Mutex
	ctx          context.Context
	apiClient    client.APIClient
	opts         *EventsWatcherOptions
	events       <-chan events.Message
	handlers     map[string]func(events.Message)
	errors       <-chan error
	errorHandler func(error)
	cancel       context.CancelFunc
	subscribed   bool
}

func (w *eventsWatcher) subscribe() {
	w.mu.Lock()
	if w.events == nil {
		w.events, w.errors = w.apiClient.Events(w.ctx, w.opts.Options)
		w.subscribed = true
	}
	w.mu.Unlock()
}

// Watch ranges over events and dispactches event messages to any
// subscribed handlers for a given action using a goroutine
func (w *eventsWatcher) Watch() {
	// subscribe to events if not already subscribed
	w.mu.Lock()
	evts := w.events
	w.mu.Unlock()
	if evts == nil {
		w.subscribe()
	}

	// on errors, call error handler if still subscribed
	go func() {
		for e := range w.errors {
			w.mu.Lock()
			subscribed := w.subscribed
			w.mu.Unlock()
			if !subscribed {
				break
			}
			w.mu.Lock()
			h := w.errorHandler
			w.mu.Unlock()
			if h == nil {
				continue
			}
			go h(e)
		}
	}()

	// on events, call event handlers if still subscribed
	go func() {
		for e := range w.events {
			w.mu.Lock()
			subscribed := w.subscribed
			w.mu.Unlock()
			if !subscribed {
				break
			}
			w.mu.Lock()
			h, exists := w.handlers[e.Action]
			w.mu.Unlock()
			if !exists {
				// if no specific handler found, then and only then check if there is a "*' catch-all handler
				w.mu.Lock()
				h, exists = w.handlers["*"]
				w.mu.Unlock()
				if !exists {
					continue
				}
			}
			go h(e)
		}
	}()
}

// Cancel is used to stop watching events
func (w *eventsWatcher) Cancel() {
	if w.cancel == nil {
		return
	}

	w.mu.Lock()
	w.cancel()
	w.cancel = nil
	w.events = nil
	w.errors = nil
	w.subscribed = false
	w.mu.Unlock()
}

// On is used by watchers to register a handler that will be invoked for the subscribed action
// The handler will replace a previous handler
func (w *eventsWatcher) On(action string, handler func(message events.Message)) {
	w.mu.Lock()
	w.handlers[action] = handler
	w.mu.Unlock()
}

// OnError is used by watchers to register a handler for errors
// The handler will replace a previous handler
func (w *eventsWatcher) OnError(handler func(err error)) {
	w.mu.Lock()
	w.errorHandler = handler
	w.mu.Unlock()
}



