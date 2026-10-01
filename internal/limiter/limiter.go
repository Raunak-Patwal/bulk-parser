package limiter

type RequestLimiter struct {
	sem chan struct{}
}

func NewRequestLimiter(maxRequests int) *RequestLimiter {
	return &RequestLimiter{
		sem: make(chan struct{}, maxRequests),
	}
}

func (l *RequestLimiter) Acquire() {
	l.sem <- struct{}{}
}

func (l *RequestLimiter) Release() {
	<-l.sem
}
