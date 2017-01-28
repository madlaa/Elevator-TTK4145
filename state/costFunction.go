//Authors: Mads Laastad & Tommy Berntzen

package state

import(
	"../global"
	"time"
	"math/rand"
)
var currentState global.Button

func StatePolling(lastKnownStateChan chan global.Button){
	for{
		select{
			case currentState = <- lastKnownStateChan:
		}
		time.Sleep(global.RECIVING_SLEEP_TIME)
	}
}

func CostFunction(newGlobalOrder global.Button) global.Order{

	var numOfFloors int

	if currentState.Dir == global.UP && newGlobalOrder.Dir == global.UP{
		if newGlobalOrder.Floor > currentState.Floor{
			numOfFloors = newGlobalOrder.Floor - currentState.Floor
		}else if newGlobalOrder.Floor <= currentState.Floor{
			numOfFloors = (2*global.N_FLOORS -2) - (currentState.Floor + newGlobalOrder.Floor)
		}
		}else if currentState.Dir == global.DOWN && newGlobalOrder.Dir == global.DOWN{
			if newGlobalOrder.Floor < currentState.Floor{
				numOfFloors = currentState.Floor - newGlobalOrder.Floor
			}else if newGlobalOrder.Floor >= currentState.Floor{
				numOfFloors = currentState.Floor + newGlobalOrder.Floor
			}
		}else if currentState.Dir == global.UP && newGlobalOrder.Dir == global.DOWN{
			if currentState.Floor > newGlobalOrder.Floor{
				numOfFloors = (2*global.N_FLOORS -2) - (currentState.Floor + newGlobalOrder.Floor)
			}else if currentState.Floor < newGlobalOrder.Floor{
				numOfFloors = 2*global.N_FLOORS - 2 - newGlobalOrder.Floor
			}
		}else if currentState.Dir == global.DOWN && newGlobalOrder.Dir == global.UP{
			if currentState.Floor > newGlobalOrder.Floor{
				numOfFloors = currentState.Floor + newGlobalOrder.Floor
			}else if currentState.Floor <= newGlobalOrder.Floor{
				numOfFloors = currentState.Floor + newGlobalOrder.Floor
			}
		}else if currentState.Dir == global.NONE{
			if currentState.Floor >= newGlobalOrder.Floor{
				numOfFloors = currentState.Floor - newGlobalOrder.Floor
			}else if currentState.Floor < newGlobalOrder.Floor{
				numOfFloors = newGlobalOrder.Floor - currentState.Floor
			}
		}
	rand.Seed(time.Now().UTC().UnixNano())
	return global.Order{newGlobalOrder, numOfFloors*global.COST_MULTIPLIER+rand.Intn(global.COST_MULTIPLIER-1)}
}
