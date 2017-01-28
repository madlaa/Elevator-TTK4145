//Authors: Mads Laastad & Tommy Berntzen

package network

import(
	"fmt"
	."net"
	"encoding/json"
	"../global"
	"time"
)

func NetworkInit() (*UDPConn, *UDPConn, *UDPConn, *UDPConn, *UDPConn, *UDPConn, []byte){
	OrderReceiveAddr, _ :=  ResolveUDPAddr("udp", global.BROADCAST+global.ORDER_RECEIVE_PORT)
	LocalOrderReceiveListen,_ := ListenUDP("udp", OrderReceiveAddr)
	OrderReceiveConn, _ := DialUDP("udp", nil, OrderReceiveAddr)
	
	CostAddr, _ := ResolveUDPAddr("udp", global.BROADCAST+global.COST_PORT)
	LocalCostListen,_ := ListenUDP("udp", CostAddr)
	CostConn, _ := DialUDP("udp", nil, CostAddr)
	
	OrderRemoveAddr, _ := ResolveUDPAddr("udp", global.BROADCAST+global.ORDER_REMOVAL_PORT)
	LocalOrderRemoveListen,_ := ListenUDP("udp", OrderRemoveAddr)
	OrderRemoveConn, _ := DialUDP("udp", nil, OrderRemoveAddr)

	Buffer := make([]byte,1024)

	return LocalOrderReceiveListen, OrderReceiveConn, LocalCostListen, CostConn, LocalOrderRemoveListen, OrderRemoveConn, Buffer
}

func GlobalOrderSender(orderReceiveConn *UDPConn, order global.Button){
	stopSending := 0
	for{
		b, err := json.Marshal(order)
		if err == nil{
			_,writeErr := orderReceiveConn.Write([]byte(b))
			if writeErr != nil{
				fmt.Println("Could not write order, error:",writeErr)
			}
		}else{
			fmt.Println("Could not Marshal order: ", err)
		}
		stopSending ++
		if stopSending > global.NUMBER_OF_NETWORK_PACKAGES {
			break
		}
		time.Sleep(global.SENDING_SLEEP_TIME)
	}
}

func GlobalOrderReceiver(globalOrderChan chan global.Button, localOrderReceiveListen *UDPConn, buffer []byte){
	var lastKnownOrder global.Button
	for{
		n,_,readErr := localOrderReceiveListen.ReadFromUDP(buffer) //Implement error-handling?
		if readErr != nil{
			fmt.Println("Could not read order ", readErr)
		}
		var temp global.Button
		err:= json.Unmarshal(buffer[:n], &temp)
		if err == nil{
			if lastKnownOrder != temp || (lastKnownOrder == temp && global.IsResend){
				globalOrderChan <- temp
				global.IsResend = false
			}
		} else {
			fmt.Println("Could not receive orderdata, check type ", err)
		}
		lastKnownOrder = temp
		time.Sleep(global.RECIVING_SLEEP_TIME)
	}
}

func GlobalOrderRemoveSender(orderRemoveConn *UDPConn, orderRemove global.Button){
	stopSending := 0
	for{
		b, err := json.Marshal(orderRemove)
		if err == nil{
			_,writeErr := orderRemoveConn.Write([]byte(b))
			if writeErr != nil{
				fmt.Println("Could not write orderRemove, error:", writeErr)
			}
		}else{
			fmt.Println("Could not send Marshal orderRemove: ", err)
		}
		stopSending ++
		if stopSending > global.NUMBER_OF_NETWORK_PACKAGES {
			break
		}
		time.Sleep(global.SENDING_SLEEP_TIME)
	}
}

func GlobalOrderRemoveReceiver(globalOrderRemoveChan chan global.Button, localOrderRemoveListen *UDPConn, buffer []byte){
	for{
		n,_,readErr := localOrderRemoveListen.ReadFromUDP(buffer)
		if readErr != nil{
			fmt.Println("Could not read orderRemove", readErr)
			panic(readErr)
		}
		var temp global.Button
		err:= json.Unmarshal(buffer[:n], &temp)
		if err == nil{
			globalOrderRemoveChan <- temp
		} else {
			fmt.Println("Could not receive orderRemove. Check type and error:", err)
		}
		time.Sleep(global.SENDING_SLEEP_TIME)
	}
}

func GlobalCostSender(costConn *UDPConn, cost global.Order){
	stopSending := 0
	for{
		b, err := json.Marshal(cost)
		if err == nil{
			_,writeErr := costConn.Write([]byte(b))
			if writeErr != nil{
				fmt.Println("Could not write cost, error:", writeErr)
			}
		}else{
			fmt.Println("Could not Marshal cost: ", err)
		}
		stopSending ++
		if stopSending > global.NUMBER_OF_NETWORK_PACKAGES {
			break
		}
		time.Sleep(global.SENDING_SLEEP_TIME)
	}
}

func GlobalCostReceiver(globalCostChan chan global.Order, localCostListen *UDPConn, buffer []byte){
	var lastKnownCost global.Order
	for{
		n,_,readErr := localCostListen.ReadFromUDP(buffer)
		if readErr != nil{
			fmt.Println("Could not read order", readErr)
			panic(readErr)
		}
		var temp global.Order
		err:= json.Unmarshal(buffer[:n], &temp)
		if err == nil{
			if lastKnownCost != temp{
				globalCostChan <- temp
			}
		} else {
			fmt.Println("Could not Receive data, check type", err)
		}
		lastKnownCost = temp
		time.Sleep(global.RECIVING_SLEEP_TIME)
	}
}
