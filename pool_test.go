package pool

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func metaFactory() Factory {
    initial := -1
    return func(context.Context) interface{} {
        initial++
        return initial
    }
}

func TestPool_Basic(t *testing.T) {
    g := metaFactory()
    pool := New(g)

    assert.Equal(t, 0, pool.Len())
    elm := pool.Get()
    assert.Equal(t, 0, elm.(int))
    pool.Put(elm)
    assert.Equal(t, 1, pool.Len())
    assert.Equal(t, 0, pool.Get().(int))
    assert.Equal(t, 1, pool.Get().(int))
}

func TestPool_Capacity(t *testing.T) {
    g := metaFactory()
    size := 5
    pool := New(g, OptCapacity(size))

    assert.Equal(t, size, pool.Cap())

    items := []interface{}{}

    for i := 0; i < size; i++ {
        items = append(items, pool.Get())
    }

    extra := pool.Get()
    assert.Equal(t, size, extra.(int))

    for _, item := range items {
        pool.Put(item)
    }

    pool.Put(extra)

    for _, item := range items {
        assert.Equal(t, item.(int), pool.Get().(int))
    }
}

func TestPool_Lease(t *testing.T) {
    g := metaFactory()
    pool := New(g, OptLeaseMS(20))

    pool.Put(pool.Get())

    elm := pool.Get()
    assert.Equal(t, 0, elm.(int))
    pool.Put(elm)

    time.Sleep(time.Millisecond * 22)
    assert.Equal(t, 1, pool.Get().(int))
}

