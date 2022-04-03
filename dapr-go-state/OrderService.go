package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

func main() {
	client, err := dapr.NewClient() // initializes Dapr Client using environment variable DAPR_GRPC_PORT
	STATE_STORE_NAME := "statestore"
	if err != nil {
		panic(err)
	}
	defer client.Close()

	for i := 0; i < 10; i++ {
		time.Sleep(5000)
		orderId := rand.Intn(1000-1) + 1
		ctx := context.Background()
		//Save the order under key order_# with # the random orderId
		if err := client.SaveState(ctx, STATE_STORE_NAME, fmt.Sprintf("order_%d", orderId), []byte(strconv.Itoa(orderId))); err != nil {
			panic(err)
		}
		//Get the data that was just saved
		result, err := client.GetState(ctx, STATE_STORE_NAME, fmt.Sprintf("order_%d", orderId))
		if err != nil {
			panic(err)
		}
		log.Println("Result after get: ")
		log.Printf("Value retrieved from state store %s (stored under key %s)", result.Value, result.Key)
	}
}
