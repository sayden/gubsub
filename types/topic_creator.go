package types

type TopicCreator interface {
	AddTopic(endpoint string) error
}
