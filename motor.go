package ev3brick

import (
	"time"

	"github.com/ev3go/ev3dev"
)

//Motor management

//Motor - enhanced TachoMotor
type Motor struct {
	*ev3dev.TachoMotor
}

//MotorFor returns a motor for a given port and driver
func MotorFor(port string, driver string) (*Motor, error) {
	tm, err := ev3dev.TachoMotorFor(port, driver)
	if err != nil {
		return nil, err
	}
	m := &Motor{
		TachoMotor: tm,
	}
	return m, nil
}

//RunForTime - Run Motor for time
func (m *Motor) RunForTime(speed int, timeInSecs int) error {
	m.SetSpeedSetpoint(speed)
	m.SetTimeSetpoint(time.Duration(timeInSecs) * time.Millisecond)
	m.Command("run-timed")
	return m.Err()
}

//RunForRotation - Run Motor for rotations
func (m *Motor) RunForRotation(speed int, rot int) error {
	cpr := m.CountPerRot()
	m.SetSpeedSetpoint(speed)
	m.SetPositionSetpoint(cpr * rot)
	m.Command("run-to-rel-pos")
	return m.Err()
}

//RunForever - runs forever until stop command is issued
func (m *Motor) RunForever(speed int) error {
	m.SetSpeedSetpoint(speed)
	m.Command("run-forever")
	return m.Err()
}

//Stop - Stops the motor
func (m *Motor) Stop() error {
	m.Command("stop")
	return m.Err()
}
