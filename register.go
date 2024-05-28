package main

import (
	"login/entity"
	"login/internal/utility"
	"login/mq"
	"sync"

	"github.com/gin-gonic/gin"
)

type registerHandler struct {
	mq    mq.MessageQueue
	idGen utility.IdGenerator
	resp  map[string]entity.RegisterResponse
	mu    sync.Mutex
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
	for retry := 0; retry < 3; retry++ {
		resp, ok := rh.resp[id]
		if ok {
			rh.mu.Lock()
			delete(rh.resp, id)
			rh.mu.Unlock()

			return resp
		}
	}
	return entity.RegisterResponse{}
}

// Basic implementation of message queue
func (rh *registerHandler) initSelector() {
	rh.resp = make(map[string]entity.RegisterResponse)
	rh.mu = sync.Mutex{}
	rh.mq.Subscribe("register_response", func(msg interface{}) error {
		resp, ok := msg.(entity.RegisterResponseMsg)
		if !ok {
			return nil
		}

		rh.mu.Lock()
		rh.resp[resp.Id] = resp.Response
		rh.mu.Unlock()
		return nil
	})
}

func newRegisterHandler(mq mq.MessageQueue, idGen utility.IdGenerator) *registerHandler {
	reg := &registerHandler{
		mq:    mq,
		idGen: idGen,
	}

	reg.initSelector()

	return reg
}
