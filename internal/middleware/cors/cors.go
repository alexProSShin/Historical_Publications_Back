package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSOptions настраивает параметры CORS
type CORSOptions struct {
	AllowedOrigins   []string // Список разрешенных Origin
	AllowedMethods   []string // Список разрешенных HTTP-методов
	AllowedHeaders   []string // Список разрешенных заголовков
	ExposedHeaders   []string // Список заголовков, которые могут быть доступны клиенту
	AllowCredentials bool     // Разрешить передачу куки и заголовков авторизации
	MaxAge           int      // Время жизни Preflight-запроса
}

// New создает новое middleware для обработки CORS в Gin
func New(opts CORSOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin == "" || !isOriginAllowed(origin, opts.AllowedOrigins) {
			c.Next()
			return
		}

		// Устанавливаем основные заголовки CORS
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Vary", "Origin")

		if opts.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Если это Preflight-запрос (OPTIONS)
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowedMethods, ", "))
			c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowedHeaders, ", "))
			if opts.MaxAge > 0 {
				c.Writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(opts.MaxAge))
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Устанавливаем Exposed Headers, если они заданы
		if len(opts.ExposedHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(opts.ExposedHeaders, ", "))
		}

		// Передаем управление следующему обработчику
		c.Next()
	}
}

// isOriginAllowed проверяет, разрешен ли запрашиваемый Origin
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, o := range allowedOrigins {
		if o == origin || o == "*" {
			return true
		}
	}
	return false
}
