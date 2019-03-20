package ev3brick

import (
	"bytes"
	"fmt"
	"log"
	"sync"

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

//mutext to lock while updating err
var mut sync.Mutex

//EV3Brick - brick details. Models the ev3 brick for easier coding
type EV3Brick struct {
	MotorsMap   map[string]*Motor
	SensorsMap  map[string]*ev3dev.Sensor
	Drive       *Drive
	LeftSensor  *EV3ColorSensor
	RightSensor *EV3ColorSensor
	err         error
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
	e.MotorsMap = make(map[string]*Motor)
	// Get the motors
	for _, p := range OutputPorts {
		for _, d := range MotorDrivers {
			m, err := MotorFor(p, d)
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

//SetDrive sets Drive to the 2 specified motors
func (e *EV3Brick) SetDrive(leftMotorPort, rightMotorPort, leftColorPort, rightColorPort string) {
	lc := EV3ColorSensorFor(e.SensorsMap[leftColorPort])
	rc := EV3ColorSensorFor(e.SensorsMap[rightColorPort])
	e.Drive = DriveFor(e.MotorsMap[leftMotorPort], e.MotorsMap[rightMotorPort], lc, rc)
}
