package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecurityMiddleware(t *testing.T) {
	Convey("Given a SecurityMiddleware instance", t, func() {
		jwtSecret := []byte("test-secret-key")
		issuer := "fastproxy-test"
		sm := NewSecurityMiddleware(jwtSecret, issuer)

		Convey("When creating a new SecurityMiddleware", func() {
			Convey("Should create instance successfully", func() {
				So(sm, ShouldNotBeNil)
				So(sm.jwtAuth, ShouldNotBeNil)
				So(sm.keyManager, ShouldNotBeNil)
			})
		})

		Convey("When generating a service token", func() {
			serviceName := "test-service"
			permissions := []string{"read", "write"}
			duration := time.Hour

			token, err := sm.GenerateServiceToken(serviceName, permissions, duration)

			Convey("Should generate token without error", func() {
				So(err, ShouldBeNil)
				So(token, ShouldNotBeEmpty)
			})
		})

		Convey("When testing security headers middleware", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middleware := sm.SecurityHeadersMiddleware(handler)
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			Convey("Should add security headers", func() {
				So(w.Header().Get("X-Content-Type-Options"), ShouldEqual, "nosniff")
				So(w.Header().Get("X-Frame-Options"), ShouldEqual, "DENY")
				So(w.Header().Get("X-XSS-Protection"), ShouldEqual, "1; mode=block")
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When testing CORS middleware", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			middleware := sm.CORSMiddleware(handler)

			Convey("For regular requests", func() {
				req := httptest.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()

				middleware.ServeHTTP(w, req)

				Convey("Should add CORS headers", func() {
					So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, "*")
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})

			Convey("For OPTIONS requests", func() {
				req := httptest.NewRequest("OPTIONS", "/", nil)
				w := httptest.NewRecorder()

				middleware.ServeHTTP(w, req)

				Convey("Should return 200 OK", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
				})
			})
		})

		Convey("When testing key manager integration", func() {
			err := sm.RotateKey("test-key", "Test key", 32)

			Convey("Should rotate key successfully", func() {
				So(err, ShouldBeNil)

				key, exists := sm.GetActiveKey("test-key")
				So(exists, ShouldBeTrue)
				So(key, ShouldNotBeNil)
			})
		})
	})
}