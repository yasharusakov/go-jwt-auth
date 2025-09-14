package nats

import (
	"context"
	"encoding/json"
	"log"
	"user-service/internal/service"

	"github.com/nats-io/nats.go"
)

type NATSHandler struct {
	service service.UserService
	nc      *nats.Conn
}

func NewNATSHandler(service service.UserService, nc *nats.Conn) *NATSHandler {
	return &NATSHandler{service, nc}
}

func (h *NATSHandler) Register() {
	_, err := h.nc.Subscribe("user.getByEmail", h.GetUserByEmail)
	if err != nil {
		log.Fatal("failed to subscribe to user.getByEmail:", err)
	}

	_, err = h.nc.Subscribe("user.getById", h.GetUserByID)
	if err != nil {
		log.Fatal("failed to subscribe to user.getById:", err)
	}

	_, err = h.nc.Subscribe("user.checkByEmail", h.CheckUserExistsByEmail)
	if err != nil {
		log.Fatal("failed to subscribe to user.checkByEmail:", err)
	}

	_, err = h.nc.Subscribe("user.register", h.RegisterUser)
	if err != nil {
		log.Fatal("failed to subscribe to user.register:", err)
	}

	log.Println("NATS subscriptions registered successfully")
}

func (h *NATSHandler) GetUserByEmail(msg *nats.Msg) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		msg.Respond([]byte(`{"error":"invalid request"}`))
		return
	}

	user, err := h.service.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		msg.Respond([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	if user == nil {
		msg.Respond([]byte(`{"error":"user not found"}`))
		return
	}

	resp, _ := json.Marshal(user)
	msg.Respond(resp)
}

func (h *NATSHandler) GetUserByID(msg *nats.Msg) {
	var req struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		msg.Respond([]byte(`{"error":"invalid request"}`))
		return
	}

	user, err := h.service.GetUserByID(context.Background(), req.ID)
	if err != nil {
		msg.Respond([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	if user == nil {
		msg.Respond([]byte(`{"error":"user not found"}`))
		return
	}

	resp, _ := json.Marshal(user)
	log.Println("GetUserByID response:", string(resp))
	msg.Respond(resp)
}

func (h *NATSHandler) CheckUserExistsByEmail(msg *nats.Msg) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		msg.Respond([]byte(`{"error":"invalid request"}`))
		return
	}

	exists, err := h.service.CheckUserExistsByEmail(context.Background(), req.Email)
	if err != nil {
		msg.Respond([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	resp, _ := json.Marshal(map[string]bool{"exists": exists})
	msg.Respond(resp)
}

func (h *NATSHandler) RegisterUser(msg *nats.Msg) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		msg.Respond([]byte(`{"error":"invalid request"}`))
		return
	}

	userID, err := h.service.RegisterUser(context.Background(), req.Email, []byte(req.Password))
	if err != nil {
		msg.Respond([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}

	resp, _ := json.Marshal(map[string]string{"id": userID})
	msg.Respond(resp)
}
