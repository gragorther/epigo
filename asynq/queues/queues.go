package queues

const (
	QueueCritical = "critical"
	QueueVeryHigh = "veryHigh"
	QueueDefault  = "default"
	QueueLow      = "low"
	QueueVeryLow  = "veryLow"
	QueueHigh     = "high"
)

var Queues = map[string]int{
	QueueCritical: 6,
	QueueVeryHigh: 5,
	QueueHigh:     4,
	QueueDefault:  3,
	QueueLow:      2,
	QueueVeryLow:  1,
}
