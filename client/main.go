package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/farinas09/go-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cc.Close()

	client := testpb.NewTestServiceClient(cc)

	DoUnary(client)
	//DoClientStreaming(client)
	DoServerStreaming(client)
	DoBidirectionalStreaming(client)
}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "test1",
	}
	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get test: %v", err)
	}
	log.Printf("Response from server: %v", res)
}

func DoClientStreaming(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "10",
			Question: "What is the capital of France?",
			Answer:   "Paris",
			TestId:   "test1",
		},
		{
			Id:       "11",
			Question: "What is the capital of Germany?",
			Answer:   "Berlin",
			TestId:   "test1",
		},
	}
	stream, err := c.SetQuestions(context.Background())
	if err != nil {
		log.Fatalf("Failed to set questions: %v", err)
	}
	for _, question := range questions {
		log.Println("Sending question: ", question.Id)
		err := stream.Send(question)
		time.Sleep(1 * time.Second)
		if err != nil {
			log.Fatalf("Failed to send question: %v", err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to close and receive: %v", err)
	}
	log.Printf("Response from server: %v", res)
}

func DoServerStreaming(c testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "test1",
	}
	res, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get students per test: %v", err)
	}
	for {
		student, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive student: %v", err)
		}
		log.Println("Student: ", student.Id, student.Name, student.Age)
	}
}

func DoBidirectionalStreaming(c testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer: "42",
	}

	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("Failed to take test: %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			stream.Send(&answer)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive question: %v", err)
				break
			}
			log.Println("Question: ", res.Id, res.Question)
		}
		close(waitChannel)
	}()
	<-waitChannel
}
