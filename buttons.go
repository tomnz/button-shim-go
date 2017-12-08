package buttonshim

import "time"

// Button represents a button.
type Button byte

const (
	// ButtonA represents button A.
	ButtonA Button = iota
	// ButtonB represents button B.
	ButtonB
	// ButtonC represents button C.
	ButtonC
	// ButtonD represents button D.
	ButtonD
	// ButtonE represents button E.
	ButtonE
)

var (
	// Buttons lists all of the available buttons in order. If you don't want to hardcode the number
	// of buttons in your application, you can use len(buttonshim.Buttons) instead.
	Buttons = []Button{
		ButtonA,
		ButtonB,
		ButtonC,
		ButtonD,
		ButtonE,
	}
)

// String returns the letter corresponding to the button.
func (b Button) String() string {
	switch b {
	case ButtonA:
		return "A"
	case ButtonB:
		return "B"
	case ButtonC:
		return "C"
	case ButtonD:
		return "D"
	case ButtonE:
		return "E"
	}
	panic("unknown button!")
}

// pollButtons drives a polling loop to check the pressed state of each button.
// Will emit press/release events to the respective handlers. Poll rate is configurable
// as an Option.
// Designed to be called in a long-running goroutine.
func (d *Driver) pollButtons() {
	for {
		stateBytes := make([]byte, 1)
		if err := d.i2c.Tx([]byte{regInput}, stateBytes); err != nil {
			panic(err)
		}
		state := stateBytes[0]

		now := time.Now()
		for _, btn := range Buttons {
			currentPress := (state>>btn)&1 == 1
			start, wasPressed := d.buttonState[btn]
			if currentPress && !wasPressed {
				// This is a new press - mark the time and signal
				d.buttonState[btn] = now
				if handlers, ok := d.pressHandlers[btn]; ok {
					for _, handler := range handlers {
						// Non-blocking send so we don't hang if a handler stops receiving
						// These additional messages will just be dropped
						select {
						case handler <- struct{}{}:
						default:
						}
					}
				}
			} else if !currentPress && wasPressed {
				// This is a button release
				delete(d.buttonState, btn)
				elapsed := now.Sub(start)
				if handlers, ok := d.releaseHandlers[btn]; ok {
					for _, handler := range handlers {
						select {
						case handler <- elapsed:
						default:
						}
					}
				}
			}
			// In other scenarios, nothing needs to be done!
		}
		time.Sleep(d.options.buttonPoll)
	}
}

// ButtonPressChan returns a new receive channel that can be used to handle press events
// for a specific button.
// The actual message carries no payload - use it as a signalling mechanism only.
// For example:
//
//	pressA := shim.ButtonPressChan(buttonshim.ButtonA)
//	pressB := shim.ButtonPressChan(buttonshim.ButtonB)
//	go func() {
//		for {
//			select {
//			case <-pressA:
//				fmt.Println("Button A pressed!")
//			case <-pressB:
//				fmt.Println("Button B pressed!")
//			}
//		}
//	}()
//
func (d *Driver) ButtonPressChan(button Button) <-chan struct{} {
	// Keep a small buffer in case the client doesn't handle the event before our next poll
	handler := make(chan struct{}, handlerBufLen)
	d.pressHandlers[button] = append(d.pressHandlers[button], handler)
	return handler
}

// ButtonReleaseChan returns a new receive channel that can be used to handle release events
// for a specific button.
// The actual message indicates the final duration that the button was observed as being
// pressed for. The accuracy of this value is limited by the poll rate.
func (d *Driver) ButtonReleaseChan(button Button) <-chan time.Duration {
	handler := make(chan time.Duration, handlerBufLen)
	d.releaseHandlers[button] = append(d.releaseHandlers[button], handler)
	return handler
}
