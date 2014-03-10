package types


type ActionType int    

const (
	InvalidMsg ActionType  = iota  //Only used to check if the message recieved is of type OrderMsg.
	NewOrder		 
	DeleteOrder
	Tender
	AddOrder
)
type OrderType struct{
	Button 	int			// Holds the button on the floor
	Floor 	int			// Holds the floor
}

type OrderMsg struct {
	Action    	ActionType   	// Holds what the information of what to do with the message
	Order 		OrderType 		// Holds the floor and button of the order
	TenderVal 	int				// If the action is a Tender, this will hold the cost from the tender, that is, the value from the cost function for this order
}
