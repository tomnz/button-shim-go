package buttonshim

import (
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/devices"
)

// New returns a new Button SHIM hardware driver. The driver encapsulates a set of buttons
// and a single RGB LED, handling interaction with the actual device.
// Connects to the device on the given I2C bus at its standard address.
func New(bus i2c.Bus, opts ...Option) (*Driver, error) {
	return NewWithConn(&i2c.Dev{Bus: bus, Addr: addr}, opts...)
}

// NewWithConn returns a new Button SHIM hardware driver, using the given periph.io
// conn.Conn object. Typically New should be used instead, but this may be useful for
// testing using mocks, or a custom I2C connection with a different address than the
// default.
func NewWithConn(periphConn conn.Conn, opts ...Option) (*Driver, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	d := &Driver{
		options:         options,
		pressHandlers:   map[Button][]chan<- struct{}{},
		releaseHandlers: map[Button][]chan<- time.Duration{},
		i2c:             periphConn,
		buttonState:     map[Button]time.Time{},
		ledColor:        make(chan color),
		brightness:      255,
	}
	if err := d.setup(); err != nil {
		return nil, err
	}
	return d, nil
}

// Driver is the primary struct for interacting with the Button SHIM device.
type Driver struct {
	options options
	// Device handle for I2C bus
	i2c             conn.Conn
	pressHandlers   map[Button][]chan<- struct{}
	releaseHandlers map[Button][]chan<- time.Duration
	frame           byte
	// Indicates the pressed time if the button is currently pressed
	buttonState map[Button]time.Time
	// Communicates new colors to the LED color goroutine
	ledColor   chan color
	brightness byte
}

type color struct {
	r, g, b byte
}

func (d *Driver) setup() error {
	// Initialize the hardware
	if err := d.i2c.Tx([]byte{regConfig, 0x1f}, nil); err != nil {
		return err
	}
	if err := d.i2c.Tx([]byte{regPolarity, 0x00}, nil); err != nil {
		return err
	}
	if err := d.i2c.Tx([]byte{regOutput, 0x00}, nil); err != nil {
		return err
	}

	go d.updateLed()
	go d.pollButtons()
	d.ledColor <- color{0, 0, 0}

	return nil
}

// Halt implements devices.Device.
func (d *Driver) Halt() error {
	// TODO: Clear the LED
	return nil
}

var _ devices.Device = &Driver{}
