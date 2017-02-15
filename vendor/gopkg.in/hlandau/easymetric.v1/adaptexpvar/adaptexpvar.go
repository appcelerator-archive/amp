package adaptexpvar

import "gopkg.in/hlandau/measurable.v1"
import "sync"
import "expvar"
import "fmt"

var UnregisteredValue = "null"

type adaptor struct {
	mutex      sync.RWMutex
	measurable measurable.Measurable
}

var measurablesMutex sync.RWMutex
var measurables = map[string]*adaptor{}

func (a *adaptor) getMeasurable() measurable.Measurable {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.measurable
}

func (a *adaptor) String() string {
	m := a.getMeasurable()
	if m == nil {
		return UnregisteredValue
	}

	mi, ok := m.(interface {
		MsInt64() int64
	})
	if !ok {
		return UnregisteredValue
	}

	return fmt.Sprintf("%v", mi.MsInt64())
}

func hook(m measurable.Measurable, event measurable.HookEvent) {
	switch m.MsType() {
	case measurable.CounterType, measurable.GaugeType:
		// ok
	default:
		return // not supported
	}

	name := m.MsName()

	switch event {
	case measurable.RegisterEvent, measurable.RegisterCatchupEvent:
		a := &adaptor{
			measurable: m,
		}

		measurablesMutex.Lock()
		defer measurablesMutex.Unlock()

		measurables[name] = a
		expvar.Publish(name, a)

	case measurable.UnregisterEvent:
		a := measurables[name]
		if a != nil { // this should always be the case, but whatever
			a.mutex.Lock()
			defer a.mutex.Unlock()

			a.measurable = nil
		}
	}
}

var once sync.Once
var hookKey int

func Register() {
	once.Do(func() {
		measurable.RegisterHook(&hookKey, hook)
	})
}
