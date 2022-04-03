package main

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"log"
)

func echoHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	// do something with the invocation here
	out = &common.Content{
		Data:        in.Data,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}

type Operands struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Result struct {
	Outcome float32 `json:"outcome"`
	Comment string  `json:"comment"`
}

func calculationHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	log.Printf("calculation - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	var operands Operands
	json.Unmarshal(in.Data, &operands)
	calculationResult := operands.X + 324*operands.Y + 19/operands.Z

	result := Result{Outcome: calculationResult, Comment: "Greetings"}
	responseBody, err := json.Marshal(result)
	out = &common.Content{
		Data:        responseBody,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}

func main() {
	log.Printf("Go initialize Dapr Service")
	s, err := daprd.NewService(":50001")
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
	if err := s.AddServiceInvocationHandler("echo", echoHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}
	if err := s.AddServiceInvocationHandler("calculate", calculationHandler); err != nil {
		log.Fatalf("error adding calculation invocation handler: %v", err)
	}
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
	log.Printf("SomeService is running = echo handler has been attached to service")
}
