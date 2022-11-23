package yaml

import (
	"bytes"
)

// Flush the buffer if needed.
func flush(emitter *yaml_emitter_t) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) {
		return yaml_emitter_flush(emitter)
	}
	return true
}

// Put a character to the output buffer.
func put(emitter *yaml_emitter_t, value byte) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) {
		return false
	}
	emitter.buffer[emitter.buffer_pos] = value
	emitter.buffer_pos++
	emitter.column++
	return true
}

// Put a line break to the output buffer.
func put_break(emitter *yaml_emitter_t) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) {
		return false
	}
	switch emitter.line_break {
	case yaml_CR_BREAK:
		emitter.buffer[emitter.buffer_pos] = '\r'
		emitter.buffer_pos += 1
	case yaml_LN_BREAK:
		emitter.buffer[emitter.buffer_pos] = '\n'
		emitter.buffer_pos += 1
	case yaml_CRLN_BREAK:
		emitter.buffer[emitter.buffer_pos+0] = '\r'
		emitter.buffer[emitter.buffer_pos+1] = '\n'
		emitter.buffer_pos += 2
	default:
		panic("unknown line break setting")
	}
	emitter.column = 0
	emitter.line++
	return true
}

// Copy a character from a string into buffer.
func write(emitter *yaml_emitter_t, s []byte, i *int) bool {
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) {
		return false
	}
	p := emitter.buffer_pos
	w := width(s[*i])
	switch w {
	case 4:
		emitter.buffer[p+3] = s[*i+3]
		fallthrough
	case 3:
		emitter.buffer[p+2] = s[*i+2]
		fallthrough
	case 2:
		emitter.buffer[p+1] = s[*i+1]
		fallthrough
	case 1:
		emitter.buffer[p+0] = s[*i+0]
	default:
		panic("unknown character width")
	}
	emitter.column++
	emitter.buffer_pos += w
	*i += w
	return true
}

// Write a whole string into buffer.
func write_all(emitter *yaml_emitter_t, s []byte) bool {
	for i := 0; i < len(s); {
		if !write(emitter, s, &i) {
			return false
		}
	}
	return true
}

// Copy a line break character from a string into buffer.
func write_break(emitter *yaml_emitter_t, s []byte, i *int) bool {
	if s[*i] == '\n' {
		if !put_break(emitter) {
			return false
		}
		*i++
	} else {
		if !write(emitter, s, i) {
			return false
		}
		emitter.column = 0
		emitter.line++
	}
	return true
}

// Set an emitter error and return false.
func yaml_emitter_set_emitter_error(emitter *yaml_emitter_t, problem string) bool {
	emitter.error = yaml_EMITTER_ERROR
	emitter.problem = problem
	return false
}

// Emit an event.
func yaml_emitter_emit(emitter *yaml_emitter_t, event *yaml_event_t) bool {
	emitter.events = append(emitter.events, *event)
	for !yaml_emitter_need_more_events(emitter) {
		event := &emitter.events[emitter.events_head]
		if !yaml_emitter_analyze_event(emitter, event) {
			return false
		}
		if !yaml_emitter_state_machine(emitter, event) {
			return false
		}
		yaml_event_delete(event)
		emitter.events_head++
	}
	return true
}

// Check if we need to accumulate more events before emitting.
//
// We accumulate extra
//  - 1 event for DOCUMENT-START
//  - 2 events for SEQUENCE-START
//  - 3 events for MAPPING-START
//
func yaml_emitter_need_more_events(emitter *yaml_emitter_t) bool {
	if emitter.events_head == len(emitter.events) {
		return true
	}
	var accumulate int
	switch emitter.events[emitter.events_head].typ {
	case yaml_DOCUMENT_START_EVENT:
		accumulate = 1
		break
	case yaml_SEQUENCE_START_EVENT:
		accumulate = 2
		break
	case yaml_MAPPING_START_EVENT:
		accumulate = 3
		break
	default:
		return false
	}
	if len(emitter.events)-emitter.events_head > accumulate {
		return false
	}
	var level int
	for i := emitter.events_head; i < len(emitter.events); i++ {
		switch emitter.events[i].typ {
		case yaml_STREAM_START_EVENT, yaml_DOCUMENT_START_EVENT, yaml_SEQUENCE_START_EVENT, yaml_MAPPING_START_EVENT:
			level++
		case yaml_STREAM_END_EVENT, yaml_DOCUMEN