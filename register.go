package main

import (
	"login/entity"
	"login/internal/utility"
	"login/mq"

	"github.com/gin-gonic/gin"
)

type registerHandler struct {
	mq    mq.MessageQueue
	idGen utility.IdGenerator
}

// Build message and send to message queue for processing
// Get response from message queue and return to client

// TODO: Add log, retry, monitor mechanism
func (rh *registerHandler) handle(c *gin.Context) {
	// Parse request
	var req entity.User
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	// Generate id
	id := rh.idGen.Generate()

	// Build message
	msg := entity.RegisterRequestMsg{
		Id:   id,
		User: req,
	}

	// Send message to message queue
	// TODO: retry or circuit breaker here
	err = rh.mq.Public("register", msg)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "error",
		})
		return
	}

	resp := rh.getResponse(id)

	c.JSON(resp.Code, gin.H{
		"status": resp.Status,
	})
}

func (rh *registerHandler) getResponse(id string) entity.RegisterResponse {
	return entity.RegisterResponse{
		Code:   200,
		Status: "OK",
	}
}

// TODO: Implement selector
func (rh *registerHandler) initSelector() {

}

func newRegisterHandler(mq mq.MessageQueue, idGen utility.IdGenerator) *registerHandler {
	reg := &registerHandler{
		mq:    mq,
		idGen: idGen,
	}

	reg.initSelector()

	return reg
}
