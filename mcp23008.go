// This is a library to manage Microchip MCP23008 This chip is used on onion.io Omega2 relay expension
package mcp23008

import (
	"golang.org/x/exp/io/i2c"
	"math"
	"log"
)

const (
	iodir   = 0x00
	ipol    = 0x01
	gpinten = 0x02
	defval  = 0x03
	intcon  = 0x04
	iocon   = 0x05
	gppu    = 0x06
	intf    = 0x07
	intcap  = 0x08
	gpio    = 0x09
	olat    = 0x0A
)

// McpInit function initialize MCP28003 after boot or restart of device
func McpInit(d *i2c.Device) error {
	// SetAllDirection
	err := d.WriteReg(iodir, []byte{0})
	if err != nil {
		return err
	}

	// SetAllPullUp
	err = d.WriteReg(gppu, []byte{0})
	return err
}

// McpGpioToggle change state of selected GPIO other one are unchanged
func McpGpioToggle(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	d.ReadReg(gpio, regValue)

	// Write ON to selected GPIO other one keep unchanged
	d.WriteReg(gpio,[]byte{regValue[0] ^ mask})
}


// McpGpioOn set GPIO to ON/High state other one are unchanged
func McpGpioOn(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	d.ReadReg(gpio, regValue)

	// Write ON to selected GPIO other one keep unchanged
	d.WriteReg(gpio,[]byte{mask | regValue[0]})
}

// Set all GPIO to ON/High state
func McpGpioAllOn(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(gpio,[]byte{0xf})
}

// McpGpioOff set GPIO to OFF/Low state other one are unchanged
func McpGpioOff(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 0 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio))) ^ 0xf

	// Read current state of all GPIO's
	d.ReadReg(gpio, regValue)

	// Write OFF to selected GPIO other one keep unchanged
	d.WriteReg(gpio,[]byte{mask & regValue[0]})
}

// Set all GPIO to OFF/Low state
func McpGpioAllOff(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(gpio,[]byte{0x0})
}

// This function return state of selected GPIO 1 for ON/High or 0 for OFF/Low state
func McpReadGpio(d *i2c.Device, gpio byte) byte {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	d.ReadReg(gpio, regValue)
	log.Printf("McpReadGpio gpio <%08b> mask <%08b> value <%08b>", gpio, mask, regValue[0])
	return (regValue[0] & mask) >> gpio
}
