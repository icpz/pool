package pool

func OptDeleter(d Deleter) Option {
    return func (p *Pool) {
        p.deleter = d
    }
}

func OptLeaseMS(l int64) Option {
    return func (p *Pool) {
        p.leaseMS = l
    }
}

func OptMaxSize(maxSize int) Option {
    return func (p *Pool) {
        p.items = make(chan *poolItem, maxSize)
    }
}

