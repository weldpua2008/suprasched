package handlers

import (
	"github.com/mustafaturan/bus"
)

func Trace(e *bus.Event) {
	log.Tracef("Event for %s: %+v", e.Topic, e)
}
