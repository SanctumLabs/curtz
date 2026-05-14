package inmemqueue

import "carduka/bidsvc/pkg/infra/queue"

// publishers
type InMemQueueMessagePublisher interface {
	queue.MessagePublisher
}

type InMemQueueEventPublisher interface {
	queue.EventPublisher
}

//  consumers

type InMemQueueEventConsumer interface {
	queue.EventConsumer
}

type InMemQueueMessageConsumer interface {
	queue.MessageConsumer
}
