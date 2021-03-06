//Authors: Mads Laastad & Tommy Berntzen

package state

import(
	"fmt"
	"time"
	"../driver"
	"../global"
	"os"
	"log"
)

func localOrderPolling(buttonDetectChan chan global.Button, externalOrderChan chan global.Button, localOrderChan chan global.Button){ //Filtrates internal orders
	for{
		select{
			case newButtonOrder := <- buttonDetectChan:
				if newButtonOrder.Dir == global.NONE{
					localOrderChan <- newButtonOrder
				} else if newButtonOrder.Dir != global.NONE{
					externalOrderChan <- newButtonOrder
				}
			}
		time.Sleep(global.RECIVING_SLEEP_TIME)
	}
}

func elevatorPollingStartup(buttonDetectChan chan global.Button, floorDetectChan chan int, setDirectionChan chan global.Direction, setButtonLightChan chan global.Button, clearButtonLightChan chan global.Button, openDoorChan chan bool, doorClosedChan chan bool, localOrderChan chan global.Button, externalOrderChan chan global.Button, lastKnownStateChan chan global.Button){
	go driver.DetectButton(buttonDetectChan)
	go driver.DetectFloor(floorDetectChan)
	go driver.SetMotorDirection(setDirectionChan)
	go driver.ButtonLightHandler(setButtonLightChan, clearButtonLightChan)
	go driver.DoorHandler(openDoorChan, doorClosedChan)
	go localOrderPolling(buttonDetectChan, externalOrderChan, localOrderChan)
	
}

func HandleKeyboardInterrupt(keyboardInterrupt chan os.Signal){
  	for sig := range keyboardInterrupt {
  		fmt.Println(" ")                                             
    		log.Printf("\nCaptured %v, stopping elevator and exiting.", sig)
    		driver.IoWriteAnalog((0x100+0), global.MOTOR_STOP)
    		os.Exit(1)                                                     
  	}                                                                

}
