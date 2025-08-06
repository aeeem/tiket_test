package airplane_stream

type AirplaneStreamRequest struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Date       string `json:"date"`
	Passengers int    `json:"passengers"`
}
