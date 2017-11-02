// This is a library to manage Microchip MCP23008 This chip is used on onion.io Omega2 relay expension
package mcp23008

import (
	"golang.org/x/exp/io/i2c"
	"math"
)

type Mcp23008 struct {
	Device      *i2c.Device
	Name        string `json:"name"`
	Address     int    `json:"address"`
	Count       byte   `json:"count"`
	Description string `json:"description"`
	Gpios       []int8 `json:"gpios"` //Using int8 array instead byte for enconding with json.Marshall....
}

const (
	iodirReg   = 0x00
	ipolReg    = 0x01
	gpintenReg = 0x02
	defvalReg  = 0x03
	intconReg  = 0x04
	ioconReg   = 0x05
	gppuReg    = 0x06
	intfReg    = 0x07
	intcapReg  = 0x08
	gpioReg    = 0x09
	olatReg    = 0x0A
)

func New(device string, name string, address int, count byte, description string) (Mcp23008, error) {
	var err error
	module := Mcp23008{nil, name, address, count, description, nil}
	if count < 1 && count > 8 {
		count = 8
	}
	if device != "" {
		err = Init(device, module.Address, &module)
	}
	return module, err
}

// Init function initialize MCP28003 after boot or restart of device
func Init(device string, add int, module *Mcp23008) error {

	var err error

	module.Device, err = i2c.Open(&i2c.Devfs{Dev: device}, add)
	if err != nil {
		module.Device = nil
		return err
	}

	if module.Count > 0 && module.Count <= 8 {
		module.Gpios = make([]int8, module.Count)
		for g := range module.Gpios {
			module.Gpios[g] = int8(McpReadGpio(module.Device, byte(g)))
		}
	}

	// SetAllDirection
	err = module.Device.WriteReg(iodirReg, []byte{0})
	if err != nil {
		return err
	}

	// SetAllPullUp
	err = module.Device.WriteReg(gppuReg, []byte{0})
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
	d.WriteReg(gpioReg, []byte{regValue[0] ^ mask})
}

// McpGpioOn set GPIO to ON/High state other one are unchanged
func McpGpioOn(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	d.ReadReg(gpioReg, regValue)

	// Write ON to selected GPIO other one keep unchanged
	d.WriteReg(gpioReg, []byte{mask | regValue[0]})
}

// Set all GPIO to ON/High state
func McpGpioAllOn(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(gpioReg, []byte{0xf})
}

// McpGpioOff set GPIO to OFF/Low state other one are unchanged
func McpGpioOff(d *i2c.Device, gpio byte) {
	regValue := []byte{0}

	// Set 0 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio))) ^ 0xf

	// Read current state of all GPIO's
	d.ReadReg(gpioReg, regValue)

	// Write OFF to selected GPIO other one keep unchanged
	d.WriteReg(gpioReg, []byte{mask & regValue[0]})
}

// Set all GPIO to OFF/Low state
func McpGpioAllOff(d *i2c.Device) {
	// Write ON to all GPIO
	d.WriteReg(gpioReg, []byte{0x0})
}

// This function return state of selected GPIO 1 for ON/High or 0 for OFF/Low state
func McpReadGpio(d *i2c.Device, gpio byte) byte {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	d.ReadReg(gpioReg, regValue)
	return (regValue[0] & mask) >> gpio
}
