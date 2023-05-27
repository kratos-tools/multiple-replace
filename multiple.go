package multiple

type Multiple struct {
	producer   Producer
	worker     Worker
	dispatcher Dispatcher
}
// fork file
func NewMultiple(producer Producer, worker Worker, dispatcher Dispatcher) *Multiple {
	return &Multiple{
		producer,
		worker,
		dispatcher,
	}
}

func (m *Multiple) Run() {
	m.dispatcher.Dispatch(m.producer, m.worker)
}
