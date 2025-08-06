package airplane

type AirplaneRequest struct {
	MessageID  string `json:"-"`
	From       string `json:"from"`
	To         string `json:"to"`
	Date       string `json:"date" example:"2025-07-10"`
	Passengers int    `json:"passengers"`
}

type Airplane struct {
	ID            string `json:"id"`
	Airline       string `json:"airline"`
	FlightNumber  string `json:"flight_number"`
	From          string `json:"from"`
	To            string `json:"to"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
	Price         int    `json:"price"`
	Currency      string `json:"currency"`
	Available     bool   `json:"available"`
}

type AirplaneRedisResponse struct {
	ID            string `json:"id"`
	Airline       string `json:"airline"`
	FlightNumber  string `json:"flight_number"`
	From          string `json:"from"`
	To            string `json:"to"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
	Price         string `json:"price"`
	Currency      string `json:"currency"`
	Available     string `json:"available"`
}

type AirplaneRepository interface {
	ACK(msgID string) (err error)
	GetAirplanes(search string) (res []Airplane, err error)
	AddSearchQueue(AirplaneRequest AirplaneRequest) (msgID string, err error)
}

type AirplaneUsecase interface {
	GetAirplanes(search string) (res []Airplane, err error)
	AddSearchQueue(AirplaneRequest AirplaneRequest) (msgID string, err error)
}
