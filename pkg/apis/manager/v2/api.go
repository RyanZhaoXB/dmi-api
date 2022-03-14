package v2


type DeviceModel struct {
	ID				string
	Name			string
	Description		string
	Manufacturer	string
	Protocol 		string
	States 			[]DeviceStateType
	Commands 		[]DeviceCommandType
	Properties		[]DevicePropertyType
	Data			[]DeviceDataType
}

type Device struct {
	ID					string
	Name				string
	Description			string
	Annotations			interface{}
	Model				string
	LastOnline  		string
	State				string
	Address 			interface{}
	DataAccess			interface{}
	Properties			map[string]DeviceProperty
	PropertyVisitors	[]DevicePropertyVisitor
}

// BaseMessage the base struct of event message
type BaseMessage struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

type Protocol struct {
	Name			string
	Description 	string
	ProtocolConfig	interface{}
}