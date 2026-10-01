package ratelimit

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const sweepEvery = 5 * time.Minute

type visitor struct {
	limiter *rate.Limiter
	seen    time.Time
}

type Limiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	key      func(*http.Request) string
	rate     rate.Limit
	burst    int
	ttl      time.Duration
	swept    time.Time
}

// New a token bucket per caller. The bucket refills at rate and tolerates burst requests at once.
func New(r rate.Limit, burst int, key func(*http.Request) string) *Limiter {
	return &Limiter{
		visitors: make(map[string]*visitor),
		key:      key,
		rate:     r,
		burst:    burst,
		ttl:      sweepEvery * 2,
		swept:    time.Now(),
	}
}

// Middleware refuses with 429 once the caller runs out of tokens.
func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l.allow(l.key(r)) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", retryAfter(l.rate))
		w.WriteHeader(http.StatusTooManyRequests)

		json.NewEncoder(w).Encode(map[string]string{
			"code":    "too_many_requests",
			"message": "too many requests, try again later",
		})
	})
}

func (l *Limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if now.Sub(l.swept) > sweepEvery {
		for k, v := range l.visitors {
			if now.Sub(v.seen) > l.ttl {
				delete(l.visitors, k)
			}
		}

		l.swept = now
	}

	v, ok := l.visitors[key]

	if !ok {
		v = &visitor{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.visitors[key] = v
	}

	v.seen = now

	return v.limiter.Allow()
}

// ByIP the caller as the edge sees it. Behind the ingress gateway the socket is the proxy, so the forwarded chain is what identifies the client.
func ByIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if first, _, found := strings.Cut(fwd, ","); found {
			return strings.TrimSpace(first)
		}

		return strings.TrimSpace(fwd)
	}

	if real := r.Header.Get("X-Real-Ip"); real != "" {
		return strings.TrimSpace(real)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return r.RemoteAddr
	}

	return host
}

// ByToken the bearer token when there is one. An office shares an address but not a session, and this endpoint costs money per call.
func ByToken(r *http.Request) string {
	if token := r.Header.Get("Authorization"); token != "" {
		return token
	}

	return ByIP(r)
}

func retryAfter(r rate.Limit) string {
	if r <= 0 {
		return "60"
	}

	seconds := int(1 / float64(r))

	if seconds < 1 {
		seconds = 1
	}

	return strconv.Itoa(seconds)
}
