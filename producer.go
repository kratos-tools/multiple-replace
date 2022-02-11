package multiple

type Producer interface {
	Yield() interface{}
	CanShard() bool
	Shard() []Producer
	IsDone() bool
}
