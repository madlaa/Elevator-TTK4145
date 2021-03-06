//Authors: Mads Laastad & Tommy Berntzen

package driver

import ("fmt"
	"time"
	"../global"
)

var lampMatrix = [][]int{{LIGHT_DOWN1, LIGHT_UP1, LIGHT_COMMAND1},
			{LIGHT_DOWN2, LIGHT_UP2, LIGHT_COMMAND2},
			{LIGHT_DOWN3, LIGHT_UP3, LIGHT_COMMAND3},
			{LIGHT_DOWN4, LIGHT_UP4, LIGHT_COMMAND4}}
			
var buttonMatrix = [][]int{{BUTTON_DOWN1,BUTTON_UP1,BUTTON_COMMAND1},
			        {BUTTON_DOWN2,BUTTON_UP2,BUTTON_COMMAND2},
				{BUTTON_DOWN3,BUTTON_UP3,BUTTON_COMMAND3},
				{BUTTON_DOWN4,BUTTON_UP4,BUTTON_COMMAND4}}

func DriverInit(){
	if !IoInit(){
		fmt.Println("Local initialization failed")
	} else {
		for i, value := range lampMatrix {
			for j := range value {	
				IoClearBit(lampMatrix[i][j])
				}
			}
		}
	
}

func DetectButton(buttonDetectChan chan global.Button){
	var tempButtonMatrix [global.N_FLOORS][3]bool
	for{
                for i, value := range buttonMatrix{
                        for j := range value{
                                if IoReadBit(buttonMatrix[i][j]) && !tempButtonMatrix[i][j]{
		                	buttonDetectChan <- global.Button{i,global.Direction(j)}
		                        tempButtonMatrix[i][j] = true
		                } else if !IoReadBit(buttonMatrix[i][j]) && tempButtonMatrix[i][j]{
		                	tempButtonMatrix[i][j] = false
		                	}
                	}
        	}
        	time.Sleep(global.SENDING_SLEEP_TIME)
	}
}

func DetectFloor(floorDetectChan chan int) {
        temp := make([]bool, global.N_FLOORS+1)
        for{
	        switch{
		        case IoReadBit(SENSOR_FLOOR1) && !temp[0]:
		                for i := 0; i < global.N_FLOORS+1; i++{
		        	        temp[i] = false
		        	}
		        	temp[0] = true
		                select{
		                        case floorDetectChan <- 0:
		                }
		        	IoClearBit(LIGHT_FLOOR_IND1)
		        	IoClearBit(LIGHT_FLOOR_IND2)
		        case IoReadBit(SENSOR_FLOOR2) && !temp[1]:
		        	select{
		        	        case floorDetectChan <- 1:
		        	}
		        	IoClearBit(LIGHT_FLOOR_IND1)
		        	IoSetBit(LIGHT_FLOOR_IND2)
		        	for i := 0; i < global.N_FLOORS+1; i++{
		        	        temp[i] = false
		        	}
		        	temp[1] = true
		        case IoReadBit(SENSOR_FLOOR3) && !temp[2]:
		        	select{
		        	        case floorDetectChan <- 2:
		        	}
		        	IoSetBit(LIGHT_FLOOR_IND1)
		        	IoClearBit(LIGHT_FLOOR_IND2)
		        	for i := 0; i < global.N_FLOORS+1; i++{
		        	        temp[i] = false
		        	}
		        	temp[2] = true
		        case IoReadBit(SENSOR_FLOOR4) && !temp[3]:
		        	select{
		        	        case floorDetectChan <- 3:
		        	}
		        	IoSetBit(LIGHT_FLOOR_IND1)
		        	IoSetBit(LIGHT_FLOOR_IND2)
		        	for i := 0; i < global.N_FLOORS+1; i++{
		        	        temp[i] = false
		        	}
		        	temp[3] = true
		        case (!IoReadBit(SENSOR_FLOOR1) && !IoReadBit(SENSOR_FLOOR2) && !IoReadBit(SENSOR_FLOOR3) && !IoReadBit(SENSOR_FLOOR4)) && !temp[4]:
		        	floorDetectChan <- -1
		        	for i := 0; i < global.N_FLOORS+1; i++{
		        	        temp[i] = false
		        	}
		        	temp[4] = true
	        }
	        time.Sleep(global.SENDING_SLEEP_TIME)
	}
}

func SetMotorDirection(setDirectionChan chan global.Direction){
	for{
        	select{
        		case tempDir := <- setDirectionChan:
        			switch{
                			case tempDir == 0:
					        IoSetBit(MOTORDIR)
					        IoWriteAnalog(MOTOR, global.MOTOR_RUN)
					case tempDir == 1:
					        IoClearBit(MOTORDIR)
					        IoWriteAnalog(MOTOR, global.MOTOR_RUN)
					case tempDir == 2:
						time.Sleep(global.STOP_MOTOR_DELAY)
					        IoWriteAnalog(MOTOR, global.MOTOR_STOP)
				}
        	}
        	time.Sleep(global.RECIVING_SLEEP_TIME)
	}
}

func ButtonLightHandler(setButtonLightChan chan global.Button, clearButtonLightChan chan global.Button){
	for{
		select {
		        case tempSetLight := <-setButtonLightChan:
			        for i, value := range lampMatrix{
			                for j := range value{
			                        if tempSetLight.Floor == i && tempSetLight.Dir == global.Direction(j){
			                                IoSetBit(lampMatrix[i][j])
			                        }
			                }
			        }
		        case tempClearLight := <-clearButtonLightChan:
			        for i, value := range lampMatrix{
			                for j := range value{
			                        if tempClearLight.Floor == i && tempClearLight.Dir == global.Direction(j){
			                                IoClearBit(lampMatrix[i][j])
			                        }
			                }
			        }
		      }
	time.Sleep(global.RECIVING_SLEEP_TIME)
	}
}

func DoorHandler(openDoorChan chan bool, doorClosedChan chan bool){
	for{
		select{
			case tempDoor := <- openDoorChan:
				if tempDoor{
					IoSetBit(LIGHT_DOOR_OPEN)
					time.Sleep(global.DOORS_OPENTIME)
					IoClearBit(LIGHT_DOOR_OPEN)
					doorClosedChan <- true
				} else{fmt.Println("Unauthorized door command.")}

						
		}
	}
	time.Sleep(global.RECIVING_SLEEP_TIME)
}
