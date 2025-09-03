package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJWTAuth(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "fastproxy-test"
	auth := NewJWTAuth(secretKey, issuer)

	Convey("Given a JWTAuth instance", t, func() {
		Convey("When generating a token", func() {
			serviceName := "test-service"
			permissions := []string{"read", "write"}
			expiresAt := time.Now().Add(time.Hour)

			token, err := auth.GenerateToken(serviceName, permissions, expiresAt)

			Convey("Should generate token without error", func() {
				So(err, ShouldBeNil)
				So(token, ShouldNotBeEmpty)
			})

			Convey("When validating the token", func() {
				claims, err := auth.ValidateToken(token)

				Convey("Should validate token successfully", func() {
					So(err, ShouldBeNil)
					So(claims, ShouldNotBeNil)
					So(claims.ServiceName, ShouldEqual, serviceName)
					So(claims.Permissions, ShouldResemble, permissions)
					So(claims.Issuer, ShouldEqual, issuer)
				})
			})

			Convey("When validating an invalid token", func() {
				_, err := auth.ValidateToken("invalid-token")

				Convey("Should return error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("When validating an expired token", func() {
				expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
					ServiceName: serviceName,
					Permissions: permissions,
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    issuer,
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				})
				tokenString, _ := expiredToken.SignedString(secretKey)

				_, err := auth.ValidateToken(tokenString)

				Convey("Should return error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestAuthMiddleware(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "fastproxy-test"
	auth := NewJWTAuth(secretKey, issuer)

	Convey("Given a JWTAuth middleware", t, func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		middleware := auth.AuthMiddleware(handler)

		Convey("When request has no Authorization header", func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			Convey("Should return 401 Unauthorized", func() {
				So(w.Code, ShouldEqual, http.StatusUnauthorized)
			})
		})

		Convey("When request has invalid Authorization header format", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "InvalidFormat")
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			Convey("Should return 401 Unauthorized", func() {
				So(w.Code, ShouldEqual, http.StatusUnauthorized)
			})
		})

		Convey("When request has invalid token", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer invalid-token")
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			Convey("Should return 401 Unauthorized", func() {
				So(w.Code, ShouldEqual, http.StatusUnauthorized)
			})
		})

		Convey("When request has valid token", func() {
			serviceName := "test-service"
			permissions := []string{"read", "write"}
			expiresAt := time.Now().Add(time.Hour)

			token, _ := auth.GenerateToken(serviceName, permissions, expiresAt)
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			Convey("Should call next handler", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, "success")
			})
		})
	})
}