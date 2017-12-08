# Button SHIM Go library

Provides a Go implementation for interfacing with Pimoroni's [Button SHIM](https://shop.pimoroni.com/products/button-shim). The top-level library provides a lot of the same functionality as the reference [Python library](http://docs.pimoroni.com/buttonshim/), including:

* Handle press/release events for the buttons
* Update the color of the RGB LED.

## Overview

The library depends on the [periph.io](https://periph.io) framework for low level device communication.

The library implements button press and release handlers using channels. For example:

    aPress := shim.ButtonPressChan(buttonshim.ButtonA)
    bPress := shim.ButtonPressChan(buttonshim.ButtonB)
    go func() {
        for {
            select {
            case <-aPress:
                fmt.Println("Button A pressed!)
            case <-bPress:
                fmt.Println("Button B pressed!)
            }
        }
    }()

Please refer to the [godocs](https://godoc.org/github.com/tomnz/button-shim-go) for full API reference.
