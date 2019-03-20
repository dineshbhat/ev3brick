package ev3brick

import (
	"strconv"

	"github.com/ev3go/ev3dev"
)

//EV3ColorSensor - lego color sensor details
type EV3ColorSensor struct {
	*ev3dev.Sensor
}

//Color values
const (
	None = iota
	Black
	Blue
	Green
	Yellow
	Red
	White
	Brown
)

//EV3ColorSensorFor - returns a EV3 Color Sensor attached to the specified port
func EV3ColorSensorFor(s *ev3dev.Sensor) *EV3ColorSensor {
	return &EV3ColorSensor{
		Sensor: s,
	}
}

//SetColorMode sets the Sensor to color mode
func (l *EV3ColorSensor) SetColorMode() error {
	l.SetMode("COL-COLOR")
	return l.Err()
}

//SetReflectMode sets the Sensor to color mode
func (l *EV3ColorSensor) SetReflectMode() error {
	l.SetMode("COL-REFLECT")
	return l.Err()
}

//SetAmbientMode sets the Sensor to color mode
func (l *EV3ColorSensor) SetAmbientMode() error {
	l.SetMode("COL-AMBIENT")
	return l.Err()
}

//GetValue returns value0 for basic modes
func (l *EV3ColorSensor) GetValue() (int, error) {
	v, err := l.Value(0)
	if err != nil {
		return -1, err
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return -1, err
	}
	return i, nil
}
