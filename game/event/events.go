package event

type eventType uint8

const (
	GameMenuEvent = iota
	GameInfoEvent
	GameRoundEvent
	GameQuitEvent
	PlayerGameHitEvent
	PlayerGameStandEvent
)

type eventHandler func()

type Event struct {
	EventType eventType
}

var eventHandlers map[eventType]eventHandler

var events []Event

func init() {
	eventHandlers = make(map[eventType]eventHandler)
}

func RegisterEventHandler(eventType eventType, eventHandler eventHandler) {
	eventHandlers[eventType] = eventHandler
}

func AddEvent(event *Event) {
	events = append(events, *event)
}

func HandleEvents() {
	for len(events) > 0 {
		event := &events[0]
		events = append(events[:0], events[1:]...)
		eventHandlers[event.EventType]()
	}
}
