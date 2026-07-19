package helpers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// Render renders a templ component with the given HTTP status code.
func Render(c *gin.Context, status int, component templ.Component) {
	c.Status(status)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		// Fallback: log error, send plain text
		http.Error(c.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
