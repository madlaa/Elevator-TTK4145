//Authors: Mads Laastad & Tommy Berntzen

package global

import(
	"time"
)

const (
	N_FLOORS = 4 //Set number of floors.
	MOTOR_RUN = 2800
	MOTOR_STOP = 0
	WORSTCASE_HANDLETIME = DOORS_OPENTIME*N_FLOORS+(2*N_FLOORS-2)*2*time.Second //approx 24 sec
	RECIVING_SLEEP_TIME = 15*time.Millisecond
	SENDING_SLEEP_TIME = 30*time.Millisecond
	RESEND_SLEEP_TIME = 60*time.Millisecond
	DOORS_OPENTIME = 3000*time.Millisecond
	STOP_MOTOR_DELAY = 100*time.Millisecond
	AUCTION_TIME_DURATION = NUMBER_OF_NETWORK_PACKAGES*SENDING_SLEEP_TIME*2 //100*time.Millisecond
	COST_MULTIPLIER = 100000
	NUMBER_OF_NETWORK_PACKAGES = 2*N_FLOORS
	ORDER_RECEIVE_PORT string = ":58742"
	ORDER_REMOVAL_PORT string = ":58744"
	COST_PORT string = ":58746"
	BROADCAST string = "129.241.187.255"
)

var IsResend bool

type Direction int

const (
	DOWN Direction = iota // = 0
	UP		        // = 1
	NONE			// = 2
)

type Button struct{
	Floor int
	Dir Direction
}

type Order struct{
	ID Button
	Cost int
}
