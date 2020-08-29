// This is a library to manage Microchip MCP23008 This chip is used on onion.io Omega2 relay expension
package mcp23008

import (
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"

	//"golang.org/x/exp/io/i2c"
	"log"
	"math"
	"time"
)

type Mcp23008 struct {
	Device      *i2c.Dev
	Name        string `json:"name"`
	Address     int    `json:"address"`
	Count       byte   `json:"count"`
	Description string `json:"description"`
	Gpios       []int8 `json:"gpios"` //Using int8 array instead byte for enconding with json.Marshall....
}

const (
	iodirReg   = 0x00 // I/O Direction register
	ipolReg    = 0x01 // Input Polarity register
	gpintenReg = 0x02 // Interrupt on change Control register
	defvalReg  = 0x03 // Default compare register for interrupt on change
	intconReg  = 0x04 // Interrupt control register
	ioconReg   = 0x05 // I/O Expander configuration register
	gppuReg    = 0x06 // GPIO Pull-up resistor register
	intfReg    = 0x07 // Interrupt flag register
	intcapReg  = 0x08 // Interrupt captured register
	gpioReg    = 0x09 // Port register
	olatReg    = 0x0A // Output Latch Register
)

func New(device string, name string, address int, count byte, description string) (Mcp23008, error) {
	log.Printf("I2C Module %s initialization...\n", name)
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

	//var err error

	//module.Device, err = i2c.Open(&i2c.Devfs{Dev: device}, add)
	b, err := i2creg.Open("")
	if err != nil {
		module.Device = nil
		return err
	}
	module.Device = &i2c.Dev{Addr: uint16(add), Bus: b}

	// Set All pin direction to Output
	//if err := module.Device.WriteReg(iodirReg, []byte{0}); err != nil {
	//	return err
	//}
	if err := module.Device.Tx([]byte{iodirReg, 0}, nil); err != nil {
		return err
	}

	// Disable pullup resistor for all pin
	//if err := module.Device.WriteReg(gppuReg, []byte{0}); err != nil {
	//	return err
	//}
	if err := module.Device.Tx([]byte{gppuReg, 0}, nil); err != nil {
		return err
	}

	// Set INTCON to 0 for all (Interrupt comparison with previous pin value)
	//if err := module.Device.WriteReg(intconReg, []byte{0}); err != nil {
	//	return err
	//}
	if err := module.Device.Tx([]byte{intconReg, 0}, nil); err != nil {
		return err
	}

	// Reading state off all pin
	if module.Count > 0 && module.Count <= 8 {
		module.Gpios = make([]int8, module.Count)
		for g := range module.Gpios {
			module.Gpios[g] = int8(ReadGpio(module, byte(g)))
		}
	}

	return err
}

func GpioSetRead(module *Mcp23008, gpio byte) error {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	log.Printf("GpioSetRead %v\n", mask|regValue[0])

	// Set pin direction to read
	//module.Device.ReadReg(iodirReg, regValue)
	module.Device.Tx([]byte{iodirReg}, regValue)
	//if err := module.Device.WriteReg(iodirReg, []byte{mask | regValue[0]}); err != nil {
	//	return err
	//}
	if err := module.Device.Tx([]byte{iodirReg, mask | regValue[0]}, nil); err != nil {
		return err
	}

	// Enable internal 100 k Ohms pull up resistor
	//module.Device.ReadReg(gppuReg, regValue)
	module.Device.Tx([]byte{gppuReg}, regValue)
	//if err := module.Device.WriteReg(gppuReg, []byte{mask | regValue[0]}); err != nil {
	//	return err
	//}
	if _, err := module.Device.Write([]byte{gppuReg, mask | regValue[0]}); err != nil {
		return err
	}

	// Reverse value of register
	//module.Device.ReadReg(ipolReg, regValue)
	module.Device.Tx([]byte{ipolReg}, regValue)
	//if err := module.Device.WriteReg(ipolReg, []byte{mask | regValue[0]}); err != nil {
	//	return err
	//}
	if _, err := module.Device.Write([]byte{ipolReg, mask | regValue[0]}); err != nil {
		return err
	}

	// Enable GPIO interrupt on change event
	//module.Device.ReadReg(gpintenReg, regValue)
	module.Device.Tx([]byte{gpintenReg}, regValue)
	//if err := module.Device.WriteReg(gpintenReg, []byte{mask | regValue[0]}); err != nil {
	//	return err
	//}
	if _, err := module.Device.Write([]byte{gpintenReg, mask | regValue[0]}); err != nil {
		return err
	}

	return nil
}

// GpioReverse change state of selected GPIO other one are unchanged
func GpioReverse(module *Mcp23008, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	//module.Device.ReadReg(gpioReg, regValue)
	module.Device.Tx([]byte{gpioReg}, regValue)

	// Write ON to selected GPIO other one keep unchanged
	//module.Device.WriteReg(gpioReg, []byte{regValue[0] ^ mask})
	module.Device.Write([]byte{gpioReg, regValue[0] ^ mask})
}

// GpioOn set GPIO to ON/High state other one are unchanged
func GpioOn(module *Mcp23008, gpio byte) {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	// Read current state of all GPIO's
	//module.Device.ReadReg(gpioReg, regValue)
	module.Device.Tx([]byte{gpioReg}, regValue)

	// Write ON to selected GPIO other one keep unchanged
	//module.Device.WriteReg(gpioReg, []byte{mask | regValue[0]})
	module.Device.Write([]byte{gpioReg, mask | regValue[0]})
}

// Set all GPIO to ON/High state
func GpioAllOn(module *Mcp23008) {
	// Write ON to all GPIO
	//module.Device.WriteReg(gpioReg, []byte{0xf})
	module.Device.Write([]byte{gpioReg, 0xf})
}

// GpioOff set GPIO to OFF/Low state other one are unchanged
func GpioOff(module *Mcp23008, gpio byte) {
	regValue := []byte{0}

	// Set 0 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio))) ^ 0xf

	// Read current state of all GPIO's
	//module.Device.ReadReg(gpioReg, regValue)
	module.Device.Tx([]byte{gpioReg}, regValue)

	// Write OFF to selected GPIO other one keep unchanged
	//module.Device.WriteReg(gpioReg, []byte{mask & regValue[0]})
	module.Device.Write([]byte{gpioReg, mask & regValue[0]})
}

// Set all GPIO to OFF/Low state
func GpioAllOff(module *Mcp23008) {
	// Write ON to all GPIO
	//module.Device.WriteReg(gpioReg, []byte{0x0})
	module.Device.Write([]byte{gpioReg, 0x0})
}

// This function return state of selected GPIO 1 for ON/High or 0 for OFF/Low state
func ReadGpio(module *Mcp23008, gpio byte) byte {
	regValue := []byte{0}

	// Set 1 to corresponding BIT of GPIO
	mask := byte(math.Pow(2, float64(gpio)))

	//module.Device.ReadReg(gpioReg, regValue)
	module.Device.Tx([]byte{gpioReg}, regValue)
	return (regValue[0] & mask) >> gpio
}

// This function handle event on GPIOs
func RegisterInterrupt(module *Mcp23008, interrupt chan byte) {
	regValue := []byte{0}

	for {
		//module.Device.ReadReg(intfReg, regValue)
		module.Device.Tx([]byte{intfReg}, regValue)
		if regValue[0] != 0 {
			//log.Printf("Interrupt occurs GPIO %d \n", binaryToGpio(regValue) )
			interrupt <- binaryToGpio(regValue)
			//module.Device.ReadReg(intcapReg, regValue)
			module.Device.Tx([]byte{intcapReg}, regValue)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func binaryToGpio(registry []byte) byte {
	return byte(math.Log10(float64(registry[0])) / math.Log10(float64(2)))
}
