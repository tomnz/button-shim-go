package buttonshim

// SetColor sets the RGB pixel to the given color.
func (d *Driver) SetColor(r, g, b byte) {
	d.ledColor <- color{r, g, b}
}

// SetBrightness sets the brightness of the pixel.
func (d *Driver) SetBrightness(brightness byte) {
	// TODO: Use the actual brightness range of the LED IC
	// https://github.com/pimoroni/button-shim/issues/3
	d.brightness = brightness
}

// updateLed continually waits for new color signals, building and sending a
// corresponding command to the LED for each.
// Designed to be called in a long-running goroutine.
func (d *Driver) updateLed() {
	for {
		color := <-d.ledColor

		pinQueue := newQueue()
		pinQueue.writeByte(0)
		pinQueue.writeByte(0)
		pinQueue.writeByte(0xef)
		pinQueue.writeByte(d.options.gamma[d.scaleVal(color.b)])
		pinQueue.writeByte(d.options.gamma[d.scaleVal(color.g)])
		pinQueue.writeByte(d.options.gamma[d.scaleVal(color.r)])
		pinQueue.writeByte(0)
		pinQueue.writeByte(0)

		bytes := pinQueue.bytes

		msg := []byte{regOutput}
		msg = append(msg, bytes...)
		if err := d.i2c.Tx(msg, nil); err != nil {
			panic(err)
		}
	}
}

func (d *Driver) scaleVal(val byte) byte {
	return byte(uint16(val) * uint16(d.brightness) / 255)
}

// TODO: Revisit all the below... Unnecessarily messy

// queue provides helper functions for constructing a message to transmit to the LED.
type queue struct {
	bytes []byte
}

func newQueue() *queue {
	return &queue{bytes: []byte{0}}
}

func (q *queue) next() {
	if len(q.bytes) == 0 {
		q.bytes = []byte{0}
	} else {
		q.bytes = append(q.bytes, q.bytes[len(q.bytes)-1])
	}
}

func (q *queue) setBit(pin, value byte) {
	if value != 0 {
		q.bytes[len(q.bytes)-1] |= (1 << pin)
	} else {
		q.bytes[len(q.bytes)-1] &= ^(1 << pin)
	}
}

func (q *queue) writeByte(b byte) {
	for i := 0; i < 8; i++ {
		q.next()
		q.setBit(pinLedClock, 0)
		q.setBit(pinLedData, b&0x80) // 0b10000000
		q.next()
		q.setBit(pinLedClock, 1)
		b <<= 1
	}
}

const chunkSize = 32
