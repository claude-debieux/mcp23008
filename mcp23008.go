// This is a library to manage Microchip MCP23008 This chip is used on onion.io Omega2 relay expension
package mcp23008

import (
	"golang.org/x/exp/io/i2c"
	"math"
)

const (
	MCP23008_REG_IODIR = 0x00
	MCP23008_REG_IPOL = 0x01
	MCP23008_REG_GPINTEN = 0x02
	MCP23008_REG_DEFVAL = 0x03
	MCP23008_REG_INTCON = 0x04
	MCP23008_REG_IOCON = 0x05
	MCP23008_REG_GPPU = 0x06
	MCP23008_REG_INTF = 0x07
	MCP23008_REG_INTCAP = 0x08
	MCP23008_REG_GPIO = 0x09
	MCP23008_REG_OLAT = 0x0A
)

// McpInit function initialize MCP28003 after boot or restart of device
func McpInit(d *i2c.Device) error {
	// SetAllDirection
	err := d.WriteReg(MCP23008_REG_IODIR, []byte{0})
	if err != nil {
		return err
	}

	// SetAllPullUp
	err = d.WriteReg(MCP23008_REG_GPPU, []byte{0})
	return err
}

// McpGpioToggle change state of selected GPIO other one are unchanged
func McpGpioToggle(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	d.ReadReg(MCP23008_REG_GPIO, regValue)

	// Write ON to selected GPIO other one keep unchanged
	d.WriteReg(MCP23008_REG_GPIO,[]byte{regValue[0] ^ mask})
}


// McpGpioOn set GPIO to ON/High state other one are unchanged
func McpGpioOn(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	d.ReadReg(MCP23008_REG_GPIO, regValue)

	// Write ON to selected GPIO other one keep unchanged
	d.WriteReg(MCP23008_REG_GPIO,[]byte{mask | regValue[0]})
}

// Set all GPIO to ON/High state
func McpGpioAllOn(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(MCP23008_REG_GPIO,[]byte{0xf})
}

// McpGpioOff set GPIO to OFF/Low state other one are unchanged
func McpGpioOff(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 0 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio))) ^ 0xf

	// Read current state of all GPIO's
	d.ReadReg(MCP23008_REG_GPIO, regValue)

	// Write OFF to selected GPIO other one keep unchanged
	d.WriteReg(MCP23008_REG_GPIO,[]byte{mask & regValue[0]})
}

// Set all GPIO to OFF/Low state
func McpGpioAllOff(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(MCP23008_REG_GPIO,[]byte{0x0})
}

// This function return state of selected GPIO 1 for ON/High or 0 for OFF/Low state
func McpReadGpio(d *i2c.Device, gpio byte) byte {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	d.ReadReg(MCP23008_REG_GPIO, regValue)
	return (regValue[0] & mask) >> gpio
}
