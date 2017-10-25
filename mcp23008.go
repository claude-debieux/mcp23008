package mcp23008

import "golang.org/x/exp/io/i2c"

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

func relayInit(d *i2c.Device) {
	// SetAllDirection
	err := d.WriteReg(MCP23008_REG_IODIR, []byte{0})
	if err != nil {
		panic(err)
	}

	// SetAllPullUp
	err = d.WriteReg(MCP23008_REG_GPPU, []byte{0})
	if err != nil {
		panic(err)
	}
}

func relaySwitchOn(d *i2c.Device, swi byte) {
	buff := []byte{0}

	stateOn := byte(swi + 1)

	d.ReadReg(MCP23008_REG_GPIO, buff)
	d.WriteReg(MCP23008_REG_GPIO,[]byte{stateOn | buff[0]})
}

func relaySwitchOff(d *i2c.Device, swi byte) {
	buff := []byte{0}

	stateOff := byte(swi + 1) ^ 0xf

	d.ReadReg(MCP23008_REG_GPIO, buff)
	d.WriteReg(MCP23008_REG_GPIO,[]byte{stateOff & buff[0]})
}

func relaySwitchRead(d *i2c.Device, swi byte) byte {
	buff := []byte{0}

	bitPos := swi + 1

	d.ReadReg(MCP23008_REG_GPIO, buff)
	state := (buff[0] & bitPos) >> swi

	return state
}

