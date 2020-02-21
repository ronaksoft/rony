package grpcGateway

import (
	"context"
	"git.ronaksoftware.com/ronak/rony/gateway"
	"testing"
	"time"

	log "git.ronaksoftware.com/ronak/rony/logger"
	"git.ronaksoftware.com/ronak/rony/msg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type mockHandlers struct {
	ConnectHandlerFunc func(gateway.Conn) error
	MessageHandlerFunc func(gateway.Conn, *msg.MessageEnvelope)
	FlushFuncFunc      func(gateway.Conn) *msg.MessageEnvelope
}

func (m *mockHandlers) ConnectHandler(c gateway.Conn) error {
	return m.ConnectHandlerFunc(c)
}

func (m *mockHandlers) MessageHandler(c gateway.Conn, envelope *msg.MessageEnvelope) {
	m.MessageHandlerFunc(c, envelope)
}

func (m *mockHandlers) FlushFunc(c gateway.Conn) *msg.MessageEnvelope {
	return m.FlushFuncFunc(c)
}

func TestGateway(t *testing.T) {
	handlers := &mockHandlers{}
	handlers.MessageHandlerFunc = func(c gateway.Conn, envelope *msg.MessageEnvelope) {
		log.Info("message handler called")
	}
	handlers.ConnectHandlerFunc = func(c gateway.Conn) error {
		log.Info("connect handler called")
		return nil
	}

	gateway := New(Config{
		ConnectHandler: handlers.ConnectHandler,
		MessageHandler: handlers.MessageHandler,
		FlushFunc:      handlers.FlushFunc,
		ListenAddress:  ":0",
	})

	gateway.Run()

	// create a client stream to server
	addr := gateway.listener.Addr().String()
	_, err := clientStream(addr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func clientStream(addr string) (msg.BotProxy_ServeBotRequestClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := msg.NewBotProxyClient(conn)

	// send auth key in context
	md := metadata.Pairs("auth-key", "testauthkey")
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	stream, err := client.ServeBotRequest(ctx)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
