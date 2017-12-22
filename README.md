# Button SHIM Go library

[![build](https://travis-ci.org/tomnz/button-shim-go.svg?branch=master)](https://travis-ci.org/tomnz/button-shim-go)
[![godocs](https://godoc.org/github.com/tomnz/button-shim-go?status.svg)](https://godoc.org/github.com/tomnz/button-shim-go)

Provides a Go implementation for interfacing with Pimoroni's [Button SHIM](https://shop.pimoroni.com/products/button-shim). The top-level library provides a lot of the same functionality as the reference [Python library](http://docs.pimoroni.com/buttonshim/), including:

* Handle press/release events for the buttons
* Update the color of the RGB LED.

## Overview

## Installation

First, clone the project into your Go path:

```bash
go get github.com/tomnz/button-shim-go
```

The library depends on the [periph.io](https://periph.io) framework for low level device communication. You can install this manually with `go get`, or (preferred) use `dep`:

```bash
go get -u github.com/golang/dep/cmd/dep
cd $GOPATH/src/github.com/tomnz/button-shim-go
dep ensure
```

## Usage

First, initialize a periph.io I2C bus, and instantiate the device with it:

```go
package main

import (
    "github.com/tomnz/button-shim-go"
    "periph.io/x/periph/conn/i2c/i2creg"
    "periph.io/x/periph/host"
)

func main() {
    // TODO: Handle errors
    _, _ := host.Init()
    bus, _ := i2creg.Open("1")
    shim, _ := buttonshim.New(bus)
}
```

The library implements button press and release handlers using channels. For example:

```go
aPress := shim.ButtonPressChan(buttonshim.ButtonA)
bPress := shim.ButtonPressChan(buttonshim.ButtonB)
aRelease := shim.ButtonReleaseChan(buttonshim.ButtonA)
go func() {
    for {
        select {
        case <-aPress:
            fmt.Println("Button A pressed!")
        case <-bPress:
            fmt.Println("Button B pressed!")
        case holdDuration := <-aRelease:
            fmt.Printf("Button A held for %s!", holdDuration)
        }
    }
}()
```

The color and brightness of the pixel can also be changed:

```go
shim.SetColor(255, 0, 0)
shim.SetBrightness(127)
```

Please refer to the [godocs](https://godoc.org/github.com/tomnz/button-shim-go) for full API reference.

## Contributing

Contributions welcome! Please refer to the [contributing guide](https://github.com/tomnz/button-shim-go/blob/master/CONTRIBUTING.md).
