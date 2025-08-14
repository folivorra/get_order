package inmemory

import (
	"container/list"
	"errors"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Node struct {
	Key   string
	Value domain.Order
}

type InMemOrderCache struct {
	logger   *slog.Logger
	capacity int
	nodes    map[string]*list.Element
	queue    *list.List
}

func NewInMemOrderCache(logger *slog.Logger, capacity int) *InMemOrderCache {
	return &InMemOrderCache{
		logger:   logger,
		capacity: capacity,
		nodes:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (c *InMemOrderCache) Set(order *domain.Order) {
	if element, exists := c.nodes[order.OrderUID.String()]; exists {
		c.queue.MoveToFront(element)
		element.Value = *order
		c.logger.Debug("item exists and moved to front",
			slog.String("key", order.OrderUID.String()),
		)
		return
	}

	if c.queue.Len() >= c.capacity {
		if element := c.queue.Back(); element != nil {
			item := c.queue.Remove(element).(*Node)
			delete(c.nodes, item.Key)
			c.logger.Debug("item removed from cache",
				slog.String("key", item.Key),
			)
		}
	}

	item := &Node{
		Key:   order.OrderUID.String(),
		Value: *order,
	}

	element := c.queue.PushFront(item)
	c.nodes[order.OrderUID.String()] = element
	c.logger.Debug("item added to cache",
		slog.String("key", order.OrderUID.String()),
	)
}

func (c *InMemOrderCache) Get(uid uuid.UUID) (*domain.Order, error) {
	element, exists := c.nodes[uid.String()]
	if !exists {
		c.logger.Debug("item does not exist",
			slog.String("key", uid.String()),
		)
		return nil, ErrKeyNotFound
	}
	c.queue.MoveToFront(element)
	c.logger.Debug("item moved to front",
		slog.String("key", uid.String()),
	)
	node := element.Value.(*Node)
	return &node.Value, nil
}
