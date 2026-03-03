package cache

type Response struct {
	Body    string   `json:"body"`
	Status  int      `json:"status"`
	Headers []Header `json:"headers"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
