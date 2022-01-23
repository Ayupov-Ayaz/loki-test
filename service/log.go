package service

type Label struct {
	App string `json:"app"`
}

type Stream struct {
	Stream Label      `json:"stream"`
	Values [][]string `json:"values"`
}

type Log struct {
	Streams []Stream `json:"streams"`
}

func NewLog(app string, values [][]string) *Log {
	return &Log{
		Streams: []Stream{
			{
				Stream: Label{
					App: app,
				},
				Values: values,
			},
		},
	}
}

type QueryData struct {
	Result []Stream `json:"result"`
}

type Query struct {
	Data QueryData `json:"data"`
}
