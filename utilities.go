package app

import (
	"github.com/elithrar/simple-scrypt"
	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

func generatePasswordHash(password string) ([]byte, error) {
	hash, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return hash, err
	}
	return hash, nil
}

func comparePassword(hash []byte, password []byte) bool {
	return scrypt.CompareHashAndPassword(hash, password) == nil
}
