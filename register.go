package main

import (
	"login/mq"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Fullname    string `json:"fullname,omitempty"`
	PhoneNumner string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Birthdate   string `json:"birthdate,omitempty"`
	LastLogin   string `json:"last_login,omitempty"`
}

type RegisterResponse struct {
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
}

type RegisterRequestMsg struct {
	Id      string          `json:"id"`
	Request RegisterRequest `json:"register_request"`
}

type RegisterResponseMsg struct {
	Id       string           `json:"id"`
	Response RegisterResponse `json:"register_response"`
}

type registerHandler struct {
	mq    mq.MessageQueue
	idGen IdGenerator
}

// Build message and send to message queue for processing
// Get response from message queue and return to client

// TODO: Add log, retry, monitor mechanism
func (rh *registerHandler) handle(c *gin.Context) {
	// Parse request
	var req RegisterRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	// Generate id
	id := rh.idGen.Generate()

	// Build message
	msg := RegisterRequestMsg{
		Id:      id,
		Request: req,
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

func (rh *registerHandler) getResponse(id string) RegisterResponse {
	return RegisterResponse{
		Code:   200,
		Status: "OK",
	}
}

// TODO: Implement selector
func (rh *registerHandler) initSelector() {

}

func newRegisterHandler(mq mq.MessageQueue, idGen IdGenerator) *registerHandler {
	reg := &registerHandler{
		mq:    mq,
		idGen: idGen,
	}

	reg.initSelector()

	return reg
}
