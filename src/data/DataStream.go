package data

// DataStream structure
type DataStream struct {
	maxWorkers int
	// workerFunc TYPE
}

// Construct new data stream -- take in # workers & function they should call
func NewDataStream(workers int) *DataStream {
	return &DataStream {
		maxWorkers: workers
	}
}
