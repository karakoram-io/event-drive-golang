// File explaining the basic flow of the system working on the basis of event driven approach
// @author Shahrukh
// @email skhan@salucro.com
package main

import (
	"fmt"
	"sync"
	"time"
)

// Event represents a generic event.
type Event struct {
	Name      string
	Timestamp time.Time
}

// EventEmitter represents an event emitter.
type EventEmitter struct {
	listeners map[string][]chan Event
	mutex     sync.Mutex
}

// NewEventEmitter creates a new event emitter.
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]chan Event),
	}
}

// AddListener adds a listener for a specific event.
func (emitter *EventEmitter) AddListener(eventName string, listener chan Event) {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()

	emitter.listeners[eventName] = append(emitter.listeners[eventName], listener)
}

// Emit sends an event to all registered listeners.
func (emitter *EventEmitter) Emit(event Event) {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()

	if listeners, ok := emitter.listeners[event.Name]; ok {
		for _, listener := range listeners {
			go func(ch chan Event) {
				ch <- event
			}(listener)
		}
	}
}

// Listener1 is a sample event listener.
func Listener1(events chan Event) {
	for {
		event := <-events
		fmt.Println("Listener 1 received event:", event)
	}
}

// Listener2 is another sample event listener.
func Listener2(events chan Event) {
	for {
		event := <-events
		fmt.Println("Listener 2 received event:", event)
	}
}

func main() {
	emitter := NewEventEmitter()

	// Create channels for listeners
	events1 := make(chan Event)
	events2 := make(chan Event)

	// Register listeners
	emitter.AddListener("event_type_1", events1)
	emitter.AddListener("event_type_2", events2)

	// Start listeners
	go Listener1(events1)
	go Listener2(events2)

	// Emit some events
	for i := 0; i < 5; i++ {
		event1 := Event{Name: "event_type_1", Timestamp: time.Now()}
		event2 := Event{Name: "event_type_2", Timestamp: time.Now()}

		emitter.Emit(event1)
		emitter.Emit(event2)

		time.Sleep(time.Second)
	}

	// Allow time for listeners to process events
	time.Sleep(5 * time.Second)
}
