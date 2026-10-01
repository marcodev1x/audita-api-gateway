package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func ok(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }

func call(h http.Handler, ip, token string) int {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)

	if ip != "" {
		req.Header.Set("X-Forwarded-For", ip)
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	return rec.Code
}

// TestMiddlewareRefusesAfterBurst the bucket holds burst requests, and the one after it is refused.
func TestMiddlewareRefusesAfterBurst(t *testing.T) {
	h := New(rate.Every(time.Hour), 3, ByIP).Middleware(http.HandlerFunc(ok))

	for i := 1; i <= 3; i++ {
		if code := call(h, "203.0.113.10", ""); code != http.StatusOK {
			t.Fatalf("requisição %d = %d, want 200", i, code)
		}
	}

	if code := call(h, "203.0.113.10", ""); code != http.StatusTooManyRequests {
		t.Errorf("requisição 4 = %d, want 429", code)
	}
}

// TestMiddlewareIsolatesCallers one caller running out of tokens must not refuse anybody else.
func TestMiddlewareIsolatesCallers(t *testing.T) {
	h := New(rate.Every(time.Hour), 1, ByIP).Middleware(http.HandlerFunc(ok))

	if code := call(h, "203.0.113.10", ""); code != http.StatusOK {
		t.Fatalf("primeiro chamador = %d, want 200", code)
	}

	if code := call(h, "203.0.113.10", ""); code != http.StatusTooManyRequests {
		t.Fatalf("primeiro chamador repetido = %d, want 429", code)
	}

	if code := call(h, "198.51.100.7", ""); code != http.StatusOK {
		t.Errorf("segundo chamador = %d, want 200", code)
	}
}

// TestMiddlewareRefills waiting for the bucket to refill lets the caller back in.
func TestMiddlewareRefills(t *testing.T) {
	h := New(rate.Every(20*time.Millisecond), 1, ByIP).Middleware(http.HandlerFunc(ok))

	if code := call(h, "203.0.113.10", ""); code != http.StatusOK {
		t.Fatalf("primeira = %d, want 200", code)
	}

	if code := call(h, "203.0.113.10", ""); code != http.StatusTooManyRequests {
		t.Fatalf("segunda = %d, want 429", code)
	}

	time.Sleep(30 * time.Millisecond)

	if code := call(h, "203.0.113.10", ""); code != http.StatusOK {
		t.Errorf("depois da recarga = %d, want 200", code)
	}
}

// TestByTokenSeparatesSessions the AI endpoint costs money per call, so an office sharing one address must not share one bucket.
func TestByTokenSeparatesSessions(t *testing.T) {
	h := New(rate.Every(time.Hour), 1, ByToken).Middleware(http.HandlerFunc(ok))

	if code := call(h, "203.0.113.10", "Bearer aaa"); code != http.StatusOK {
		t.Fatalf("sessão a = %d, want 200", code)
	}

	if code := call(h, "203.0.113.10", "Bearer aaa"); code != http.StatusTooManyRequests {
		t.Fatalf("sessão a repetida = %d, want 429", code)
	}

	if code := call(h, "203.0.113.10", "Bearer bbb"); code != http.StatusOK {
		t.Errorf("sessão b no mesmo IP = %d, want 200", code)
	}
}

// TestByIPReadsTheForwardedChain behind the ingress gateway the socket is the proxy, not the client.
func TestByIPReadsTheForwardedChain(t *testing.T) {
	cases := []struct {
		name       string
		forwarded  string
		realIP     string
		remoteAddr string
		want       string
	}{
		{name: "cadeia encaminhada", forwarded: "203.0.113.10, 10.0.0.1", want: "203.0.113.10"},
		{name: "encaminhado único", forwarded: "203.0.113.10", want: "203.0.113.10"},
		{name: "x-real-ip", realIP: "198.51.100.7", want: "198.51.100.7"},
		{name: "socket", remoteAddr: "192.0.2.5:54321", want: "192.0.2.5"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)

			if c.forwarded != "" {
				req.Header.Set("X-Forwarded-For", c.forwarded)
			}

			if c.realIP != "" {
				req.Header.Set("X-Real-Ip", c.realIP)
			}

			if c.remoteAddr != "" {
				req.RemoteAddr = c.remoteAddr
			}

			if got := ByIP(req); got != c.want {
				t.Errorf("ByIP() = %q, want %q", got, c.want)
			}
		})
	}
}

// TestMiddlewareAnswersWithRetryAfter a refused caller is told when to come back.
func TestMiddlewareAnswersWithRetryAfter(t *testing.T) {
	h := New(rate.Every(12*time.Second), 1, ByIP).Middleware(http.HandlerFunc(ok))

	call(h, "203.0.113.10", "")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}

	if rec.Header().Get("Retry-After") != "12" {
		t.Errorf("Retry-After = %q, want 12", rec.Header().Get("Retry-After"))
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
}
