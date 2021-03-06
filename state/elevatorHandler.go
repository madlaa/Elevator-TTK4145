//Authors: Mads Laastad & Tommy Berntzen

package state

import (
	"time"
	"../global"
	"container/list"
)

func elevatorHandler(initFloor int, floorDetectChan chan int, openDoorChan chan bool, setDirectionChan chan global.Direction, setButtonLightChan chan global.Button, clearButtonLightChan chan global.Button, localOrderChan chan global.Button, localOrderRemoveChan chan global.Button, lastKnownStateChan chan global.Button, doorClosedChan chan bool){
	lastKnownState := global.Button{initFloor, global.NONE}
	var doorsClosed bool
	var currentFloor int
	var localPriorityQueue list.List
	lastKnownStateChan <- lastKnownState
	for{
		select{
			case newOrder := <- localOrderChan:
				if localPriorityQueue.Front() == nil && newOrder.Floor != lastKnownState.Floor && doorsClosed{
					localPriorityQueue.PushBack(newOrder)
					setButtonLightChan <- newOrder
					if lastKnownState.Floor > localPriorityQueue.Front().Value.(global.Button).Floor{
						setDirectionChan <- global.DOWN
						lastKnownState.Dir = global.DOWN
						lastKnownStateChan <- lastKnownState
					} else if lastKnownState.Floor < localPriorityQueue.Front().Value.(global.Button).Floor{
						setDirectionChan <- global.UP
						lastKnownState.Dir = global.UP
						lastKnownStateChan <- lastKnownState
					}
				} else if localPriorityQueue.Front() == nil && lastKnownState.Dir == global.NONE && newOrder.Floor == lastKnownState.Floor && doorsClosed{
					setDirectionChan <- global.NONE
					doorsClosed = false
					openDoorChan <- true
					//send confirmed to network
					if newOrder.Dir != global.NONE{
						localOrderRemoveChan <- newOrder
					}
				}
				orderInQueue := false
				for e := localPriorityQueue.Front(); e != nil; e = e.Next(){
					if newOrder == e.Value.(global.Button){
						orderInQueue = true
					}
				}
				if !orderInQueue && newOrder.Floor != lastKnownState.Floor{
					setButtonLightChan <- newOrder
					localPriorityQueue.PushBack(newOrder)
				}
			case doorsClosed = <- doorClosedChan:
				if localPriorityQueue.Front() != nil{
					if lastKnownState.Floor > localPriorityQueue.Front().Value.(global.Button).Floor{
						setDirectionChan <- global.DOWN
						lastKnownState.Dir = global.DOWN
						lastKnownStateChan <- lastKnownState
					} else if lastKnownState.Floor < localPriorityQueue.Front().Value.(global.Button).Floor{
						setDirectionChan <- global.UP
						lastKnownState.Dir = global.UP
						lastKnownStateChan <- lastKnownState
					}
				} else if localPriorityQueue.Front() == nil && currentFloor != -1{
					setDirectionChan <- global.NONE
					lastKnownState.Dir = global.NONE
					lastKnownStateChan <- lastKnownState
				}
			case currentFloor = <- floorDetectChan:
				beenHereBefore := false
				if currentFloor != -1 && currentFloor != lastKnownState.Floor{
					lastKnownState.Floor = currentFloor
					lastKnownStateChan <- lastKnownState
				}
				if localPriorityQueue.Front() != nil{
					for e := localPriorityQueue.Front(); e != nil; e = e.Next() {
						if lastKnownState.Floor == e.Value.(global.Button).Floor && (e.Value.(global.Button).Dir == lastKnownState.Dir || e.Value.(global.Button).Dir == global.NONE){
							if !beenHereBefore && currentFloor != -1 {
								setDirectionChan <- global.NONE
								doorsClosed = false
								openDoorChan <- true
							}
							clearButtonLightChan <- e.Value.(global.Button)
							if e.Value.(global.Button).Dir != global.NONE{
								localOrderRemoveChan <- e.Value.(global.Button)
							}
							localPriorityQueue.Remove(e)
							beenHereBefore = true
						}
					}
				}
				if (localPriorityQueue.Front() != nil && lastKnownState.Floor == localPriorityQueue.Front().Value.(global.Button).Floor){
					if !beenHereBefore && currentFloor != -1 {
						setDirectionChan <- global.NONE
						doorsClosed = false
						openDoorChan <- true
						
					}
					clearButtonLightChan <- localPriorityQueue.Front().Value.(global.Button)
					if localPriorityQueue.Front().Value.(global.Button).Dir != global.NONE{
						localOrderRemoveChan <- localPriorityQueue.Front().Value.(global.Button)
					}
					localPriorityQueue.Remove(localPriorityQueue.Front())			
					beenHereBefore = true
				}
			}
		}
		time.Sleep(global.RECIVING_SLEEP_TIME)
}
