Repository description
------
Student project developed during the spring of 2015 as part of the course "TTK4145 Sanntidsprogrammering" at the Norwegian University of Science and Technology (NTNU). The project was a collaboration between Mads Laastad and Tommy Berntzen. Both authors had no prior experience with Golang before starting work on this project. 

We had a lot of fun while working on this project and decided to produce an unconventional and original solution to the problem. It is not the most efficient solution, but it was an interesting learning experience and scored 20.1/25 in the factory stress test and subsequent code review. Our solution implemented a UDP based communication system between the elevators with the intention to never identify the other elevators in the system. The resulting real-time system proved to be robust in stress tests regarding a sudden loss of one or more elevators, but the inability to directly communicate could sometimes result in a long response time to redistribute uncompleated orders from a sudden elevator loss.

The project assignment is listed below.

Summary
-------
Create software for controlling `n` elevators working in parallel across `m` floors.


Main requirements
-----------------
Be reasonable: There may be semantic hoops that you can jump through to create something that is "technically correct". But do not hesitate to contact us if you feel that something is ambiguous or missing from these requirements.

No orders are lost
 - Once the light on an external button (calling an elevator to that floor; top 6 buttons on the control panel) is turned on, an elevator should arrive at that floor
 - Similarly for an internal button (telling the elevator what floor you want to exit at; front 4 buttons on the control panel), but only the elevator at that specific workspace should take the order
 - This means handling both losing network connection, losing power (both to the elevator and the machine that controls it), and software that crashes
   - For internal orders, handling loss of power/software crash implies that the orders are executed once service is restored
   - The time used to detect these failures should be reasonable, ie. on the order of magnitude of seconds (not minutes)
 - If the elevator is disconnected from the network, it should still serve whatever orders are currently "in the system" (ie. whatever lights are showing)
   - It should also keep serving internal orders, so that people can exit the elevator even if it is disconnected

Multiple elevators should be more efficient than one
 - The orders should be distributed across the elevators in a reasonable way
   - Ex: If all three elevators are idle and two of them are at the bottom floor, then a new order at the top floor should be handled by the closest elevator (ie. neither of the two at the bottom).
 - You are free to choose and design your own "cost function" of some sort: Minimal movement, minimal waiting time, etc.
 - The project is not about creating the "best" or "optimal" distribution of orders. It only has to be clear that the elevators are cooperating and communicating.
 
An individual elevator should behave sensibly and efficiently
 - No stopping at every floor "just to be safe"
 - The external "call upward" and "call downward" buttons should behave differently
   - Ex: If the elevator is moving from floor 1 up to floor 4 and there is a downward order at floor 3, then the elevator should not stop on its way upward, but should return back to floor 3 on its way down
 
The lights should function as expected
 - The lights on the external buttons should show the same thing on all `n` workspaces
 - The internal lights should not be shared between elevators
 - The "door open" lamp should be used as a substitute for an actual door, and as such should not be switched on while the elevator is moving
   - The duration for keeping the door open should be in the 1-5 second range

 
Start with `1 <= n <= 3` elevators, and `m == 4` floors. Try to avoid hard-coding these values: You should be able to add a fourth elevator with no extra configuration, or change the number of floors with minimal configuration. You do, however, not need to test for `n > 3` and `m != 4`.

Unspecified behaviour
---------------------
Some things are left intentionally unspecified. Their implementation will not be explicitly tested, and are therefore up to you.

Which orders are cleared when stopping at a floor
 - You can clear only the orders in the direction of travel, or assume that everyone enters/exits the elevator when the door opens
 
How the elevator behaves when it cannot connect to the network during initialization
 - You can either enter a "single-elevator" mode, or refuse to start
 
How the external (call up, call down) buttons work when the elevator is disconnected from the network
 - You can optionally refuse to take these new orders
 

Permitted simplifications and assumptions
-----------------------------------------
Try to create something that works at a base level first, before adding more advanced features. You are of course free to include any or all (or more) of these optional features from the start.

You can make these simplifications and still get full score:
 - At least one elevator is always working normally
 - Stop button & obstruction switch are disabled
   - Their functionality (if/when implemented) is up to you.
 - No multiple simultaneous errors: Only one error happens at a time, but the system must still return to a fail-safe state after this error
 - No network partitioning: Situations where there are multiple sets of two or more elevators with no connection between them can be ignored
   
   

Evaluation
----------

 - Completion  
   This is a test of the complete system, where the student assistants will use a checklist of various scenarios to test that the elevator system works in accordance to the main requirements listed above. If some behaviour/feature is not part of these requirements, it will not be tested. 
   
 - Design  
   This is a 10 minute presentation driven by you, followed by a 5 minute Q&A. The presentation should demonstrate that:
   - The system is robust: It will always return to a fail-safe state where no orders are lost, regardless of any foreseen or unforeseen events
   - You have made conscious and reasoned decisions about network topology, module responsibilities, and other design issues
   - You understand any weaknesses your design has, and how they are (or could be) addressed
   
   A small hand-in containing diagrams and other visual aids should also be submitted. Details about time and place is found on blackboard.
   
 - Code Review  
   You will review other groups code, other groups will review your code and teachers will review your code. The practical details are yet to be determined. 
