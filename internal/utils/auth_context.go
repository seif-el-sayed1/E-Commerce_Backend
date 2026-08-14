package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/**
	this helper functions is used to get data from the context
**/

func GetUserID(c *gin.Context) uuid.UUID {
	val, exists := c.Get("userId")
	if !exists {
		panic("GetUserID: userId not found in context — is the route protected?")
	}

	id, ok := val.(uuid.UUID)
	if !ok {
		panic("GetUserID: userId in context is not a uuid.UUID")
	}

	return id
}

func GetUserRole(c *gin.Context) string {
	val, exists := c.Get("role")
	if !exists {
		panic("GetUserRole: role not found in context — is the route protected?")
	}

	role, ok := val.(string)
	if !ok {
		panic("GetUserRole: role in context is not a string")
	}

	return role
}

func GetUserData(c *gin.Context) interface{} {
	val, exists := c.Get("user")
	if !exists {
		panic("GetUserData: user not found in context — is the route protected?")
	}

	return val
}
