package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/configs"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/middleware"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Configuração do Redis e do Middleware para Testes
var redisTestStore = storage.NewRedisStorage("localhost:6379")
var limitRequest = 100

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RateLimiter(&configs.Config{
		RateLimitIP:    limitRequest,
		RateLimitToken: limitRequest,
		BlockTime:      3,
	}, redisTestStore))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Request Allowed"})
	})
	return r
}

// Limpa as chaves do Redis após cada teste
func clearRedis() {
	redisTestStore.FlushAll()
}

// Testa se múltiplas requisições dentro do limite são permitidas
func TestRateLimiter_WithinLimit(t *testing.T) {
	clearRedis()

	router := setupTestRouter()
	for i := 0; i < limitRequest; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("API_KEY", "test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected request to be allowed")
	}
}

// Testa se o middleware bloqueia corretamente após exceder o limite
func TestRateLimiter_ExceedLimit(t *testing.T) {
	clearRedis()

	router := setupTestRouter()
	for i := 0; i < limitRequest; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("API_KEY", "test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected request to be allowed")
	}

	// A 6ª requisição deve ser bloqueada
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("API_KEY", "test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Expected request to be blocked")
}

// Testa se o bloqueio ocorre apenas uma vez e se expira corretamente
func TestRateLimiter_BlockingAndExpiration(t *testing.T) {
	clearRedis()

	router := setupTestRouter()

	// Enviar 5 requisições para ativar o bloqueio
	for i := 0; i < limitRequest; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("API_KEY", "test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Última requisição deve ser bloqueada
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("API_KEY", "test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Expected request to be blocked")

	// Esperar o tempo de bloqueio acabar
	time.Sleep(4 * time.Second)

	// Testar novamente após o bloqueio expirar
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("API_KEY", "test-token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected request to be allowed after block expires")
}

// Testa se o bloqueio ocorre apenas uma vez, mesmo com múltiplas requisições simultâneas
func TestRateLimiter_SingleBlocking(t *testing.T) {
	clearRedis()

	router := setupTestRouter()

	// Enviar 5 requisições rápidas para ativar o bloqueio
	for i := 0; i < limitRequest; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("API_KEY", "test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Várias requisições simultâneas
	// Todas devem ser bloqueadas, mas o bloqueio só deve ser registrado uma vez
	for i := 0; i < limitRequest; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("API_KEY", "test-token")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Expected request to be blocked")
		}()

	}

	// Esperar o tempo de bloqueio acabar
	time.Sleep(5 * time.Second)

	// Testar novamente após o bloqueio expirar
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("API_KEY", "test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected request to be allowed after block expires")
}
