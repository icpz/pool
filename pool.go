package pool

import (
    "context"
    "time"
)

const (
    DefaultCapacity = 10
)

type Factory = func(context.Context) interface{}
type Deleter = func(interface{})

type Option func(*Pool)

type poolItem struct {
    elem     interface{}
    lastTime time.Time
}

func (i *poolItem) Expired(leaseMS int64) bool {
    now := time.Now()
    return now.Sub(i.lastTime).Milliseconds() > leaseMS
}

type Pool struct {
    items    chan *poolItem
    factory  Factory
    deleter  Deleter
    leaseMS  int64
}

func (p *Pool) GetContext(ctx context.Context) interface{} {
    for {
        select {
        case item := <-p.items:
            elem := item.elem
            if p.leaseMS >= 0 && item.Expired(p.leaseMS) {
                if p.deleter != nil {
                    p.deleter(elem)
                }
                continue
            }
            return elem
        default:
            return p.factory(ctx)
        }
    }
}

func (p *Pool) Get() interface{} {
    return p.GetContext(context.Background())
}

func (p *Pool) Put(elem interface{}) {
    item := &poolItem{
        elem: elem,
        lastTime: time.Now(),
    }

    select {
    case p.items <- item:
        break
    default:
        if p.deleter != nil {
            p.deleter(elem)
        }
    }
    return
}

func (p *Pool) ReleaseAll() {
    close(p.items)
    for item := range p.items {
        if p.deleter != nil {
            p.deleter(item.elem)
        }
    }
}

func New(factory Factory, opts ...Option) *Pool {
    p := &Pool{
        items:   make(chan *poolItem, DefaultCapacity),
        factory: factory,
    }

    for _, opt := range opts {
        opt(p)
    }

    return p
}

