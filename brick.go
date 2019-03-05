package ev3brick

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ev3go/ev3dev"
	"github.com/pkg/errors"
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

//mutext to lock while updating err
var mut sync.Mutex

//EV3Brick - brick details. Models the ev3 brick for easier coding
type EV3Brick struct {
	MotorsMap  map[string]*ev3dev.TachoMotor
	SensorsMap map[string]*ev3dev.Sensor

	err error
}

//MotorParams motor parameters
type MotorParams struct {
	Speed     int
	TimeInMs  int
	Rotations int
	Degrees   int
	StopMode  string
	Mode      string
	RampUp    int
	RampDown  int
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

//OutputPorts - array of output ports
var OutputPorts = []string{
	OutPortA,
	OutPortB,
	OutPortC,
	OutPortD,
}

//InputPorts - array of input ports
var InputPorts = []string{
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
	e.MotorsMap = make(map[string]*ev3dev.TachoMotor)
	// Get the motors
	for _, p := range OutputPorts {
		for _, d := range MotorDrivers {
			m, err := ev3dev.TachoMotorFor(p, d)
			if err != nil {
				log.Printf("failed to find %s on %s", d, p)
			} else {
				e.MotorsMap[p] = m
				break
			}
		}
	}
	e.SensorsMap = make(map[string]*ev3dev.Sensor)
	var wg sync.WaitGroup
	wg.Add(len(InputPorts))
	for i, p := range InputPorts {
		go func(port string, index int) {
			for _, d := range SensorDrivers {
				for _, i2c := range I2CPorts {
					np := port
					if i2c != "" {
						np = fmt.Sprintf("%s:%s", port, i2c)
					}
					s, err := ev3dev.SensorFor(np, d)
					if err != nil {
						log.Printf("failed to find %s on %s", d, np)
					} else {
						e.SensorsMap[port] = s
						break
					}
				}
			}
			wg.Done()
		}(p, i)
	}
	wg.Wait()
	return nil
}

func (e EV3Brick) String() string {
	var buf bytes.Buffer
	buf.WriteString("EV3Brick:\n")
	for k, v := range e.MotorsMap {
		buf.WriteString(fmt.Sprintf("	%s]%+v\n", k, *v))
	}
	for k, v := range e.SensorsMap {
		buf.WriteString(fmt.Sprintf("	%s]%+v\n", k, *v))
	}
	return buf.String()
}

//RunForTime - Run Motor for time
func (e *EV3Brick) RunForTime(port string, speed int, timeInSecs int) error {
	m := e.MotorsMap[port]
	if m == nil {
		return errors.New(fmt.Sprintf("Cannot find motor at port %s", port))
	}

	m.SetSpeedSetpoint(speed)
	m.SetTimeSetpoint(time.Duration(timeInSecs) * time.Millisecond)
	m.Command("run-timed")
	return m.Err()
}

//RunForRotation - Run Motor for rotations
func (e *EV3Brick) RunForRotation(port string, speed int, rot int) error {
	m := e.MotorsMap[port]
	if m == nil {
		return errors.New(fmt.Sprintf("Cannot find motor at port %s", port))
	}
	cpr := m.CountPerRot()
	m.SetSpeedSetpoint(speed)
	m.SetPositionSetpoint(cpr * rot)
	m.Command("run-to-rel-pos")
	return m.Err()
}

//MoveTankRotation - Move left and right motors for rotation
func (e *EV3Brick) MoveTankRotation(portLeft string, portRight string, speedLeft int, speedRight int, rot int) error {
	go func() {
		err := e.RunForRotation(portLeft, speedLeft, rot)
		if err != nil {
			mut.Lock()
			defer mut.Unlock()
			e.err = errors.Wrap(e.err, err.Error())
		}
	}()
	go func() {
		err := e.RunForRotation(portRight, speedRight, rot)
		if err != nil {
			mut.Lock()
			defer mut.Unlock()
			e.err = errors.Wrap(e.err, err.Error())
		}
	}()

	return e.err
}
