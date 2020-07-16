package pool

import (
    "context"
    "reflect"
    "testing"
    "time"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
    if a == b {
        return
    }
    t.Fatalf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

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

    assertEqual(t, 0, pool.Len())
    elm := pool.Get()
    assertEqual(t, 0, elm.(int))
    pool.Put(elm)
    assertEqual(t, 1, pool.Len())
    assertEqual(t, 0, pool.Get().(int))
    assertEqual(t, 1, pool.Get().(int))
}

func TestPool_Capacity(t *testing.T) {
    g := metaFactory()
    size := 5
    pool := New(g, OptCapacity(size))

    assertEqual(t, size, pool.Cap())

    items := []interface{}{}

    for i := 0; i < size; i++ {
        items = append(items, pool.Get())
    }

    extra := pool.Get()
    assertEqual(t, size, extra.(int))

    for _, item := range items {
        pool.Put(item)
    }

    pool.Put(extra)

    for _, item := range items {
        assertEqual(t, item.(int), pool.Get().(int))
    }
}

func TestPool_Lease(t *testing.T) {
    g := metaFactory()
    pool := New(g, OptLeaseMS(20))

    pool.Put(pool.Get())

    elm := pool.Get()
    assertEqual(t, 0, elm.(int))
    pool.Put(elm)

    time.Sleep(time.Millisecond * 22)
    assertEqual(t, 1, pool.Get().(int))
}

