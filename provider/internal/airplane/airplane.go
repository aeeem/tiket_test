package airplane

type AirplaneRequest struct {
	MessageID  string `json:"-"`
	From       string `json:"from"`
	To         string `json:"to"`
	Date       string `json:"date"`
	Passengers string `json:"passengers"`
}

type Airplane struct {
	ID                      string `json:"id"`
	Airline                 string `json:"airline"`
	FlightNumber            string `json:"flight_number"`
	From                    string `json:"from"`
	To                      string `json:"to"`
	DepartureTime           string `json:"departure_time"`
	ArrivalTime             string `json:"arrival_time"`
	Price                   int    `json:"price"`
	Currency                string `json:"currency"`
	Available               bool   `json:"available"`
	TotalAvailablePassanger int    `json:"total_available_passanger"` //asuming availibility based on total passenger
}

type AirplaneRepository interface {
	Fetch(search AirplaneRequest) (res []Airplane, err error)
	Subscribe() (search []AirplaneRequest, msgID []string, err error)
	ACK(msgID []string) (err error)
	Publish(id string, airplane Airplane) (msgID string, err error)
}

type AirplaneUsecase interface {
	Subscribe() (search []AirplaneRequest, err error)
	Publish(search AirplaneRequest) (err error)
}
