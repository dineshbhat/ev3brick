package ev3brick

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ev3go/ev3dev"
)

const (
	//OutPortA port A
	OutPortA = "ev3-ports:outA"
	//OutPortB port B
	OutPortB = "ev3-ports:outB"
	//OutPortC port B
	OutPortC = "ev3-ports:outC"
	//OutPortD port B
	OutPortD = "ev3-ports:outD"
	//InPort1 port 1
	InPort1 = "ev3-ports:in1"
	//InPort2  port 2
	InPort2 = "ev3-ports:in2"
	//InPort3 port 3
	InPort3 = "ev3-ports:in3"
	//InPort4 port4
	InPort4 = "ev3-ports:in4"
	//I2C1 - I2C port1
	I2C1 = "i2c1"
	//I2C2 - I2C port2
	I2C2 = "i2c2"
	//I2C3 - I2C port3
	I2C3 = "i2c3"
	//I2C4 - I2C port4
	I2C4 = "i2c4"
)

//EV3Brick - brick details. Models the ev3 brick for easier coding
type EV3Brick struct {
	Motors [4]*ev3dev.TachoMotor
	Sensor [4]*ev3dev.Sensor

	err error
}

//Brick Instance of Brick
var Brick = &EV3Brick{}

//MotorDrivers - Types of motordrivers
var MotorDrivers = []string{"lego-ev3-m-motor", "lego-ev3-l-motor"}

//SensorDrivers - types of sensor drivers
var SensorDrivers = []string{
	"ht-nxt-color-v2",
	"lego-ev3-touch",
	"lego-nxt-touch",
	"lego-ev3-color",
	"lego-ev3-us",
	"lego-nxt-us",
	"lego-ev3-gyro",
}

//OutputOutPorts - array of output ports
var OutputOutPorts = []string{
	OutPortA,
	OutPortB,
	OutPortC,
	OutPortD,
}

//InputOutPorts - array of input ports
var InputOutPorts = []string{
	InPort1,
	InPort2,
	InPort3,
	InPort4,
}

//I2CPorts i2c ports
var I2CPorts = []string{
	"",
	I2C1,
	I2C2,
	I2C3,
	I2C4,
}

//Init - Initialize EV3 EV3Brick
//Read input and output ports and load the motor and sensors
func (e *EV3Brick) Init() error {
	// Get the motors
	for i, p := range OutputOutPorts {
		for _, d := range MotorDrivers {
			m, err := ev3dev.TachoMotorFor(p, d)
			if err != nil {
				log.Printf("failed to find %s on %s", d, p)
			} else {
				e.Motors[i] = m
				break
			}
		}
	}
	for i, p := range InputOutPorts {
		for _, d := range SensorDrivers {
			for _, i2c := range I2CPorts {
				np := p
				if i2c != "" {
					np = fmt.Sprintf("%s:%s", p, i2c)
				}
				s, err := ev3dev.SensorFor(np, d)
				if err != nil {
					log.Printf("failed to find %s on %s", d, np)
				} else {
					e.Sensor[i] = s
					break
				}
			}
		}
	}
	return nil
}

func (e EV3Brick) String() string {
	var buf bytes.Buffer
	buf.WriteString("EV3Brick:\n")
	for i, m := range e.Motors {
		if m != nil {
			buf.WriteString(fmt.Sprintf("	%d]%+v\n", i, *m))
		}
	}
	for i, s := range e.Sensor {
		if s != nil {
			buf.WriteString(fmt.Sprintf("	%d]%+v\n", i, *s))
		}
	}
	return buf.String()
}

//RunMotorForever - run motor for ever
// func (e *EV3EV3Brick) RunMotorForever(port )
