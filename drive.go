package ev3brick

import (
	"time"

	"github.com/ev3go/ev3dev/motorutil"
)

//drive.go manages 2 motor drive

//Drive - 2 motor drive
type Drive struct {
	*motorutil.Steering
	leftMotor        *Motor
	rightMotor       *Motor
	leftColorSensor  *EV3ColorSensor
	rightColorSensor *EV3ColorSensor
}

//DriveFor returns drive
func DriveFor(leftMotor, rightMotor *Motor, leftColorSensor, rightColorSensor *EV3ColorSensor) *Drive {
	return &Drive{
		Steering: &motorutil.Steering{
			Left:  leftMotor.TachoMotor,
			Right: rightMotor.TachoMotor,
		},
		leftMotor:        leftMotor,
		rightMotor:       rightMotor,
		leftColorSensor:  leftColorSensor,
		rightColorSensor: rightColorSensor,
	}
}

//RunForTime - Run Motor for time
func (d *Drive) RunForTime(speed int, timeInSecs int) error {
	go d.leftMotor.RunForTime(speed, timeInSecs)
	go d.rightMotor.RunForTime(speed, timeInSecs)
	//need to figure out how to return error
	return nil
}

//RunForRotation - Run Motor for rotations
func (d *Drive) RunForRotation(speed int, rot int) error {
	go d.leftMotor.RunForRotation(speed, rot)
	go d.rightMotor.RunForRotation(speed, rot)
	//need to figure out how to return error
	return nil
}

//RunForever - runs for ever until stopped
func (d *Drive) RunForever(speed int) error {
	go d.leftMotor.RunForever(speed)
	go d.rightMotor.RunForever(speed)
	//need to figure out how to return error
	return nil
}

//Stop - stop drive
func (d *Drive) Stop() error {
	go d.leftMotor.Stop()
	go d.rightMotor.Stop()
	//need to figure out how to return error
	return nil
}

//MoveUntilBlack - moves forward until black
func (d *Drive) MoveUntilBlack(speed int) error {
	d.leftColorSensor.SetColorMode()
	go d.RunForever(speed)
	for {
		v, err := d.leftColorSensor.GetValue()
		if err != nil {
			return err
		}
		if v == Black {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	d.Stop()
	return nil
}

//LineFollow -follows a line until line ends
func (d *Drive) LineFollow(speed int) {
	//set color sensors o reflective mode
	d.leftColorSensor.SetReflectMode()

}
