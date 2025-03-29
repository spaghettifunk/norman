package eventmanager

import "sync"

type EventType string

const (
	SegmentInitialized EventType = "SEGMENT_INITIALIZED"
	SegmentCreated     EventType = "SEGMENT_CREATED"
	SegmentIndexed     EventType = "SEGMENT_INDEXED"
	SegmentErrored     EventType = "SEGMENT_ERRORED"
)

var lock = &sync.Mutex{}

// Event represents the data associated with an event.
type Event struct {
	Type EventType
	Data any
}

// Observer defines the interface for event observers.
type Observer interface {
	OnNotify(Event)
}

// eventManager manages event subscriptions and notifications.
type eventManager struct {
	observers map[EventType][]Observer
}

var singleInstance *eventManager

func GetEventManager() *eventManager {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &eventManager{
				observers: make(map[EventType][]Observer),
			}
		}
	}
	return singleInstance
}

// Subscribe adds an observer to listen for a specific event.
func (em *eventManager) Subscribe(eventType EventType, observer Observer) {
	em.observers[eventType] = append(em.observers[eventType], observer)
}

// Unsubscribe removes an observer from listening to a specific event.
func (em *eventManager) Unsubscribe(eventType EventType, observer Observer) {
	observers, ok := em.observers[eventType]
	if !ok {
		return
	}

	for i, obs := range observers {
		if obs == observer {
			em.observers[eventType] = append(observers[:i], observers[i+1:]...)
			break
		}
	}
}

// Notify sends an event to all registered observers for a specific event.
func (em *eventManager) Notify(event Event) {
	observers, ok := em.observers[event.Type]
	if !ok {
		return
	}
	for _, observer := range observers {
		observer.OnNotify(event)
	}
}
