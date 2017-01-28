//Authors: Mads Laastad & Tommy Berntzen

package main

import(
	"fmt"
	"time"
	"./global"
	"./network"
	"./state"
	."net"
	"container/list"
	"os"
	"os/signal"
)
var globalOrderList list.List

func sendingOrderUpdatesToNetwork(orderReceiveConn *UDPConn, orderRemoveConn *UDPConn, externalOrderChan chan global.Button, localOrderRemoveChan chan global.Button){
	for{
		select{
			case newGlobalOrder := <- externalOrderChan:
				network.GlobalOrderSender(orderReceiveConn, newGlobalOrder)
			case newGlobalRemoveOrder := <- localOrderRemoveChan:
				network.GlobalOrderRemoveSender(orderRemoveConn, newGlobalRemoveOrder)
		}
		time.Sleep(global.SENDING_SLEEP_TIME)
	}
}

func networkPolling(costConn *UDPConn, globalOrderChan chan global.Button, globalCostChan chan global.Order, localOrderChan chan global.Button, setButtonLightChan chan global.Button, globalOrderRemoveChan chan global.Button, externalOrderChan chan global.Button, clearButtonLightChan chan global.Button){
	var newGlobalOrder global.Button
	for{	
		select{
			case newGlobalOrder = <- globalOrderChan:
				globalOrderInQueue := false
				for e := globalOrderList.Front(); e != nil; e = e.Next(){
					if newGlobalOrder == e.Value.(global.Button){
						globalOrderInQueue = true
					}
				}
				if !globalOrderInQueue{
					globalOrderList.PushBack(newGlobalOrder)
					go orderAuction(costConn, newGlobalOrder, globalCostChan, localOrderChan, setButtonLightChan, globalOrderRemoveChan, externalOrderChan, clearButtonLightChan)
					//We run the auctions parallell by go'ing them. Once one auction is compleate the go-function will die.
					fmt.Println("Staring a new goroutine for costOrderHandler")
				}
		}
		time.Sleep(global.RECIVING_SLEEP_TIME)
	}
}

func orderAuction(costConn *UDPConn, newGlobalOrder global.Button, globalCostChan chan global.Order, localOrderChan chan global.Button, setButtonLightChan chan global.Button, globalOrderRemoveChan chan global.Button, externalOrderChan chan global.Button, clearButtonLightChan chan global.Button){
	auctionComplete := false
	localCostOrder := state.CostFunction(newGlobalOrder)
	go network.GlobalCostSender(costConn, localCostOrder)
	auctionTimer := time.NewTimer(time.Duration(global.AUCTION_TIME_DURATION))
	bestElevator := localCostOrder //Type Order{ID, Cost}
	for !auctionComplete{
		select{
			case newGlobalCost := <- globalCostChan:
				if localCostOrder.ID == newGlobalCost.ID{ //&& localCostOrder.ID.Dir == newGlobalCost.ID.Dir
					if newGlobalCost.Cost >= global.COST_MULTIPLIER && newGlobalCost.Cost < bestElevator.Cost{
						fmt.Println("We updated the bestElevator from ", bestElevator.Cost, "to: ", newGlobalCost.Cost)
						bestElevator = newGlobalCost
					}
				}
			case <- auctionTimer.C:
				auctionComplete = true
		}
	}
	go waitingForOrderConfirmation(bestElevator, globalOrderRemoveChan, externalOrderChan, clearButtonLightChan)
	if bestElevator == localCostOrder{
		localOrderChan <- bestElevator.ID //If we win the auction, we send the order to the local elevator
		fmt.Println("Auction complete! WE WON with the cost,", localCostOrder.Cost)
	}
	if bestElevator != localCostOrder{
		fmt.Println("Lost the auction with the cost,", localCostOrder.Cost)
		setButtonLightChan <- bestElevator.ID //When one elevator has taken the order, all the lights should be set.
	}
}

func waitingForOrderConfirmation(bestElevator global.Order, globalOrderRemoveChan chan global.Button, externalOrderChan chan global.Button, clearButtonLightChan chan global.Button){
	confirmationTimer := time.NewTimer(global.WORSTCASE_HANDLETIME)
	waitingForConfirmation := true
	for waitingForConfirmation{
		select{
			case orderConfirmed := <- globalOrderRemoveChan:
				if orderConfirmed == bestElevator.ID{
					for e := globalOrderList.Front(); e != nil; e = e.Next(){
						if orderConfirmed == e.Value.(global.Button){
							globalOrderList.Remove(e)
							fmt.Println("Removing order")
						}
					}
					clearButtonLightChan <- bestElevator.ID
					waitingForConfirmation = false
				}
			case <- confirmationTimer.C:
				for e := globalOrderList.Front(); e != nil; e = e.Next(){
					if bestElevator.ID == e.Value.(global.Button){
						globalOrderList.Remove(e)
					}
				}
				fmt.Println("An elevator did not handle the order,", bestElevator, ". Re-sending the order.", bestElevator.ID)
				for i := 0; i < global.NUMBER_OF_NETWORK_PACKAGES; i++{ 
					externalOrderChan <- bestElevator.ID //Sends the order back to the system as a new order.
					time.Sleep(global.RESEND_SLEEP_TIME)
				}
				global.IsResend = true
				waitingForConfirmation = false
		}
	}
}



func main(){
	buttonDetectChan := make(chan global.Button)
	floorDetectChan := make(chan int)
	setDirectionChan := make(chan global.Direction)
	setButtonLightChan := make(chan global.Button)
	clearButtonLightChan := make(chan global.Button)
	openDoorChan := make(chan bool)
	doorClosedChan := make(chan bool)
	lastKnownStateChan := make(chan global.Button)
	externalOrderChan := make(chan global.Button)
	localOrderChan := make(chan global.Button)
	localOrderRemoveChan := make(chan global.Button)
	globalOrderChan := make(chan global.Button)
	globalOrderRemoveChan := make(chan global.Button)
	globalCostChan := make(chan global.Order)

	keyboardInterrupt := make(chan os.Signal)                                       
	signal.Notify(keyboardInterrupt, os.Interrupt)  
	
	state.ElevatorInit(buttonDetectChan, floorDetectChan, setDirectionChan, setButtonLightChan, clearButtonLightChan, openDoorChan, localOrderChan, localOrderRemoveChan, externalOrderChan, globalOrderRemoveChan, lastKnownStateChan, doorClosedChan)

	LocalOrderReceiveListen, OrderReceiveConn, LocalCostListen, CostConn, LocalOrderRemoveListen, OrderRemoveConn, Buffer := network.NetworkInit()

	fmt.Println(global.IsResend)
	go state.StatePolling(lastKnownStateChan)
	go networkPolling(CostConn, globalOrderChan, globalCostChan, localOrderChan, setButtonLightChan, globalOrderRemoveChan, externalOrderChan, clearButtonLightChan)
	go sendingOrderUpdatesToNetwork(OrderReceiveConn, OrderRemoveConn, externalOrderChan, localOrderRemoveChan)
	go network.GlobalOrderReceiver(globalOrderChan, LocalOrderReceiveListen, Buffer)
	go network.GlobalCostReceiver(globalCostChan, LocalCostListen, Buffer)
	go network.GlobalOrderRemoveReceiver(globalOrderRemoveChan, LocalOrderRemoveListen, Buffer)
	go state.HandleKeyboardInterrupt(keyboardInterrupt)
	select{}
}
