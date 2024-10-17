package events

import "sync"

type eventListener[T any] struct {
	eventName string
	callback  func(T)
}

type EventHandler[T any] struct {
	listenersLock sync.RWMutex
	listeners     map[string][]*eventListener[T]
	wg            sync.WaitGroup
}

func NewEventHandler[T any]() *EventHandler[T] {
	return &EventHandler[T]{
		listeners: make(map[string][]*eventListener[T]),
	}
}

func (eh *EventHandler[T]) On(eventName string, callback func(T)) (removeListener func()) {
	eh.listenersLock.Lock()
	defer eh.listenersLock.Unlock()

	el := &eventListener[T]{
		eventName: eventName,
		callback:  callback,
	}

	eh.listeners[eventName] = append(eh.listeners[eventName], el)

	return func() {
		eh.listenersLock.Lock()
		defer eh.listenersLock.Unlock()

		filter := make([]*eventListener[T], 0)
		for _, l := range eh.listeners[eventName] {
			if l != el {
				filter = append(filter, l)
			}
		}

		eh.listeners[eventName] = filter

		if len(filter) == 0 {
			delete(eh.listeners, eventName)
		}
	}
}

func (eh *EventHandler[T]) Emit(eventName string, data T) {
	eh.listenersLock.Lock()
	defer eh.listenersLock.Unlock()

	listeners := eh.listeners[eventName]
	eh.wg.Add(len(listeners))

	for _, l := range listeners {
		// Run this inside goroutine
		go func() {
			defer eh.wg.Done()
			l.callback(data)
		}()
	}
}

func (eh *EventHandler[T]) Wait() {
	eh.wg.Wait()
}
