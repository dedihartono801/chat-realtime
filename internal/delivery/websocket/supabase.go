package websocket

import (
	"github.com/dedihartono801/chat-realtime/internal/app/usecase/supabase"
	"github.com/dedihartono801/chat-realtime/pkg/dto"
)

type SupabaseHandler interface {
	SaveMessages(req *dto.SaveMessageDto) error
}

type supabaseHandler struct {
	service supabase.Service
}

func NewSupabase(service supabase.Service) SupabaseHandler {
	return &supabaseHandler{service}
}

func (h *supabaseHandler) SaveMessages(req *dto.SaveMessageDto) error {
	_, err := h.service.SaveMessages(req)

	if err != nil {
		return err
	}
	return nil
}
