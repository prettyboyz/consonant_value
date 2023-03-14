package yaml

// Set the writer error and return false.
func yaml_emitter_set_writer_error(emitter *yaml_emitter_t, problem string) bool {
	emitter.error = yaml_WRITER_ERROR
	emitter.problem = problem
	return false
}

// Flush the output buffer.
func yaml_emitter_flush(emitter *yaml_emitter_t) bool {
	if emitter.write_handler == nil {
		panic("write handler not set")
	}

	// Check if the buffer is empty.
	if emitter.buffer_pos == 0 {
		return true
	}

	// If the output encoding is UTF-8, we don't need to recode the buffer.
	if emitter.encoding == yaml_UTF8_ENCODING {
		if err := emitter.write_handler(emitter, emitter.buffer[:emitter.buffer_pos]); err != nil {
			return yaml_emitter_set_writer_error(emitter, "write error: "+err.Error())
		}
		emitter.buffer_pos = 0
		return true
	}

	// Recode the buffer into the raw buffer.
	var low, high int
	if emitter.encoding == yaml_UTF16LE_ENCODING {
		low, high = 0, 1
	} else {
		high, low = 1, 0
	}

	pos := 0
	for pos < emitter.buffer_pos {
		// See the "reader.c" code for more details on UTF-8 encoding.  Note
		// that we assume that the buffer contains a valid UTF-8 sequence.

		// Read the next UTF-8 character.
		octet := emitter.buffer[pos]

		var w int
		var value rune
		switch {
		case octet&0x80 == 0x00:
			w, value = 1, rune(octet&0x7F)
		case octet&0xE0 == 0xC0:
			w, value = 2, rune(octet&0x1F)
		case octet&0xF0 == 0xE0:
			w, 