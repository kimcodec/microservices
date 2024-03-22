package main

import (
	"context"
	"log"
	"time"

	"github.com/kimcodec/microservices/lesson_1/grpc/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to server: ", err.Error())
	}
	defer conn.Close()

	c := note_v1.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &note_v1.GetRequest{Id: 1})
	if err != nil {
		log.Println("Failed to get note: ", err.Error())
	}

	log.Printf("Note info:\n %+v", r.Note)
}
