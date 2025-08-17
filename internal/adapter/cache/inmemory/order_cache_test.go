package inmemory_test

import (
	"sync"
	"testing"

	"github.com/folivorra/get_order/internal/adapter/cache/inmemory"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

func newTestCache(cap int) *inmemory.InMemOrderCache {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	return inmemory.NewInMemOrderCache(logger, cap)
}

func newOrder(id string) *domain.Order {
	uid := uuid.MustParse(id)
	return &domain.Order{
		OrderUID:    uid,
		TrackNumber: "TRACK-" + id,
	}
}

func TestCacheSetGet(t *testing.T) {
	cache := newTestCache(5)

	order := newOrder("11111111-1111-1111-1111-111111111111")
	cache.Set(order)

	got, err := cache.Get(order.OrderUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OrderUID != order.OrderUID {
		t.Errorf("expected uid %v, got %v", order.OrderUID, got.OrderUID)
	}
}

func TestCacheEviction(t *testing.T) {
	cache := newTestCache(2)

	o1 := newOrder("11111111-1111-1111-1111-111111111111")
	o2 := newOrder("22222222-2222-2222-2222-222222222222")
	o3 := newOrder("33333333-3333-3333-3333-333333333333")

	cache.Set(o1)
	cache.Set(o2)
	cache.Set(o3) // должно вытеснить o1

	if _, err := cache.Get(o1.OrderUID); err == nil {
		t.Errorf("expected o1 to be evicted")
	}
	if _, err := cache.Get(o2.OrderUID); err != nil {
		t.Errorf("expected o2 to exist")
	}
	if _, err := cache.Get(o3.OrderUID); err != nil {
		t.Errorf("expected o3 to exist")
	}
}

func TestCacheUpdateMovesToFront(t *testing.T) {
	cache := newTestCache(3)

	o1 := newOrder("11111111-1111-1111-1111-111111111111")
	o2 := newOrder("22222222-2222-2222-2222-222222222222")
	o3 := newOrder("33333333-3333-3333-3333-333333333333")

	cache.Set(o1)
	cache.Set(o2)
	cache.Set(o3)

	// обновляем o1, он должен стать самым новым
	o1.TrackNumber = "UPDATED"
	cache.Set(o1)

	// добавляем новый, должен вытеснить o2, так как o1 был недавно обновлен
	o4 := newOrder("44444444-4444-4444-4444-444444444444")
	cache.Set(o4)

	// проверяем, что o2 был удален
	if _, err := cache.Get(o2.OrderUID); err == nil {
		t.Errorf("expected o2 to be evicted")
	}

	// проверяем, что o1 обновлен
	got, err := cache.Get(o1.OrderUID)
	if err != nil {
		t.Fatalf("expected o1 to exist, got error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected o1 to be non-nil")
	}
	if got.TrackNumber != "UPDATED" {
		t.Errorf("expected o1 TrackNumber to be 'UPDATED', got %v", got.TrackNumber)
	}
}

func TestCacheGetNotFound(t *testing.T) {
	cache := newTestCache(2)
	uid := uuid.New()

	_, err := cache.Get(uid)
	if err != inmemory.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestCacheConcurrency(t *testing.T) {
	cache := newTestCache(100)
	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			o := newOrder(uuid.New().String())
			cache.Set(o)
			_, _ = cache.Get(o.OrderUID)
		}(i)
	}
	wg.Wait()
}
