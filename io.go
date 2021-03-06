package gpio

import (
	"errors"
	"os"
)

// Pin represents a single pin, which can be used either for reading or writing
type Pin struct {
	Number    uint
	direction direction
	f         *os.File
}

// NewInput opens the given pin number for reading. The number provided should be the pin number known by the kernel
func NewInput(p uint) (Pin, error) {
	pin := Pin{
		Number: p,
	}
	err := exportGPIO(pin)
	if err != nil {
		return pin, err
	}
	pin.direction = inDirection
	err = setDirection(pin, inDirection, 0)
	if err != nil {
		return pin, err
	}
	return openPin(pin, false)
}

// NewOutput opens the given pin number for writing. The number provided should be the pin number known by the kernel
// NewOutput also needs to know whether the pin should be initialized high (true) or low (false)
func NewOutput(p uint, initHigh bool) (Pin, error) {
	pin := Pin{
		Number: p,
	}
	err := exportGPIO(pin)
	if err != nil {
		return pin, err
	}
	initVal := uint(0)
	if initHigh {
		initVal = uint(1)
	}
	pin.direction = outDirection
	err = setDirection(pin, outDirection, initVal)
	if err != nil {
		return pin, err
	}
	return openPin(pin, true)
}

// Close releases the resources related to Pin. This doen't unexport Pin, use Cleanup() instead
func (p Pin) Close() {
	if p.f != nil {
		p.f.Close()
		p.f = nil
	}
}

// Cleanup close Pin and unexport it
func (p Pin) Cleanup() error {
	p.Close()
	return unexportGPIO(p)
}

// Read returns the value read at the pin as reported by the kernel. This should only be used for input pins
func (p Pin) Read() (uint, error) {
	return readPin(p)
}

// SetLogicLevel sets the logic level for the Pin. This can be
// either "active high" or "active low"
func (p Pin) SetLogicLevel(logicLevel LogicLevel) error {
	return setLogicLevel(p, logicLevel)
}

// High sets the value of an output pin to logic high
func (p Pin) High() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return writePin(p, 1)
}

// Low sets the value of an output pin to logic low
func (p Pin) Low() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return writePin(p, 0)
}
