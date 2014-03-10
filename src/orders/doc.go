package orders
/* 
Manages orders locally and over the network, checks for events, and sets the order lights.
All variables and communication internally in the module is handled by OrderHandler. Orderhandler contains all used variables in the module to eliminate the 
use of package variables to have better control over where the variables can be changed. 

Functions:
func OrderHandler(orderReachedEvent chan<- bool, newOrderEvent chan<- bool, switchDirEvent chan<- Direction, noOrdersEvent chan<- bool)
- The module boss. Everything in the module will go through the orderhandler. It directs everything in the module and launches events on the event channels.
func msgHandler(msg types.OrderMsg, locOrdMat *[Floors][Buttons] int, aTen *map[types.OrderType] TenderType, lTen *map[types.OrderType] time.Time, prevFloor int, dir Direction, netAlive bool)(newOrder bool)
- The messages handler does, as the name implies, handle messages. Almost everything order related is sent as a OrderMsg struct to MsgHandler who acts accoringly
to the message type. Takes care of both the local order matrix and the
func checkTenderMaps(aTenders map[types.OrderType] TenderType, lTenders map[types.OrderType] time.Time)(tenderAction bool, tenderOrders []types.OrderMsg)
- Checks if the an order in either the lost- or active tender map have run out.
func atOrder(locOrdMat[Floors][Buttons] int, prevDir Direction, prevFloor *int) (orderReached bool, del bool, delOrders []types.OrderMsg)
- Checks if the elevator is at an order.
func checkMsg(msg types.OrderMsg) bool 
-Checks that a message is valid
func isLocOrdMatEmpty(locOrdMat [Floors][Buttons] int) bool
- Checks if the the local order list is empty
func checkForOrders(locOrdMat[Floors][Buttons] int, prevFloor int)(ordersAtCurFloor[Buttons] bool, ordersInDir[2] bool)
- Checks if there are an order at the given floor or orders further in that direction
func getOrders(locOrdMat *[Floors][Buttons] int, aTenders map[types.OrderType] TenderType, lTenders map[types.OrderType] time.Time )(newOrders bool, orders []types.OrderMsg ) {
- Gets new orders from the buttons
func getDir(direction Direction, prevFloor int, locOrdMat[Floors][Buttons] int) Direction
- Gets the next direction of travel
func cost(orderFloor int, orderType int, locOrdMat [Floors][Buttons] int, prevFloor int, direction Direction) (cost int)
- Gets the cost for the elevator to do a given order. Sub-functions: getTravelCost(), getWaitCost(), getDirectionCost().
*/
