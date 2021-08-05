package types

type Metric struct {
	Id        string
	TimeStamp string
	Server    string
	Instance  string
	Metrics   []ABLApplication
}
