package server

import (
	"github.com/gin-gonic/gin"
	"log"
)

func (s *Server) speakerHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	buffer := make([]byte, 4096)
	for {
		if c.Request.Context().Err() != nil {
			return
		}
		n, err := c.Request.Body.Read(buffer)
		if err != nil {
			break
		}
		if err = s.speaker.Pump(buffer[:n]); err != nil {
			log.Printf("Failed to pump audio to speaker: %v\n", err)
			break
		}
	}
}
