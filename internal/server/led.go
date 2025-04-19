package server

import (
	"github.com/Binozo/EchoGo/v2/internal/bindings/led"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) ledsHandler(c *gin.Context) {
	if s.ledController == nil {
		c.String(http.StatusInternalServerError, "Waiting for LED setup")
		return
	}
	var leds []led.Led

	if err := c.ShouldBind(&leds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.ledController.SetLEDs(leds...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
