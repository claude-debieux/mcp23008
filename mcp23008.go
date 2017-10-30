// This is a library to manage Microchip MCP23008 This chip is used on onion.io Omega2 relay expension
package mcp23008

import (
	"golang.org/x/exp/io/i2c"
	"math"
	"log"
)

const (
	iodirReg 	= 0x00
	ipolReg     = 0x01
	gpintenReg  = 0x02
	defvalReg   = 0x03
	intconReg   = 0x04
	ioconReg    = 0x05
	gppuReg  	= 0x06
	intfReg     = 0x07
	intcapReg   = 0x08
	gpioReg     = 0x09
	olatReg     = 0x0A
)

// McpInit function initialize MCP28003 after boot or restart of device
func McpInit(d *i2c.Device) error {
	// SetAllDirection
	err := d.WriteReg(iodirReg, []byte{0})
	if err != nil {
		return err
	}

	// SetAllPullUp
	err = d.WriteReg(gppuReg, []byte{0})
	return err
}

// McpGpioToggle change state of selected GPIO other one are unchanged
func McpGpioToggle(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	d.ReadReg(gpioReg, regValue)

	// Write ON to selected GPIO other one keep unchanged
	d.WriteReg(gpioReg,[]byte{regValue[0] ^ mask})
}


// McpGpioOn set GPIO to ON/High state other one are unchanged
func McpGpioOn(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	d.ReadReg(gpioReg, regValue)

	// Write ON to selected GPIO other one keep unchanged
	d.WriteReg(gpioReg,[]byte{mask | regValue[0]})
}

// Set all GPIO to ON/High state
func McpGpioAllOn(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(gpioReg,[]byte{0xf})
}

// McpGpioOff set GPIO to OFF/Low state other one are unchanged
func McpGpioOff(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 0 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio))) ^ 0xf

	// Read current state of all GPIO's
	d.ReadReg(gpioReg, regValue)

	// Write OFF to selected GPIO other one keep unchanged
	d.WriteReg(gpioReg,[]byte{mask & regValue[0]})
}

// Set all GPIO to OFF/Low state
func McpGpioAllOff(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(gpioReg,[]byte{0x0})
}

// This function return state of selected GPIO 1 for ON/High or 0 for OFF/Low state
func McpReadGpio(d *i2c.Device, gpio byte) byte {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	err := d.ReadReg(gpioReg, regValue)
	if err != nil {
		panic(err)
	}
	log.Printf("McpReadGpio gpio <%08b> mask <%08b> value <%08b>", gpio, mask, regValue[0])
	return (regValue[0] & mask) >> gpio
}
