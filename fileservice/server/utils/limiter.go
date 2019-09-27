package utils

// 限制请求数
type Limiter struct {
	ch    chan struct{}
	count int
}

func NewLimiter(count int) *Limiter {
	return &Limiter{ch: make(chan struct{}, count), count: count}
}

func (p *Limiter) Get() bool {
	if len(p.ch) >= p.count {
		return false
	}
	p.ch <- struct{}{}
	return true
}

func (p *Limiter) Put() {
	<-p.ch
}
