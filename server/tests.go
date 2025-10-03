package server

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/farinas09/go-grpc/models"
	"github.com/farinas09/go-grpc/repository"
	"github.com/farinas09/go-grpc/studentpb"
	"github.com/farinas09/go-grpc/testpb"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, errors.New("test not found")
	}
	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	err := s.repo.SetTest(ctx, &models.Test{
		Id:   req.GetId(),
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}
	return &testpb.SetTestResponse{Id: req.GetId(), Name: req.GetName()}, nil
}

func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {
	for {
		question, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: true})
		}
		if err != nil {
			return err
		}
		err = s.repo.SetQuestion(context.Background(), &models.Question{
			Id:       question.GetId(),
			Question: question.GetQuestion(),
			Answer:   question.GetAnswer(),
			TestId:   question.GetTestId(),
		})
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}

}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: true})
		}
		if err != nil {
			return err
		}
		err = s.repo.SetEnrollment(context.Background(), &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		})
		if err != nil {
			log.Println("Error setting enrollment:", err)
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())
	if err != nil {
		return err
	}
	for _, student := range students {
		time.Sleep(1 * time.Second)
		err = stream.Send(&studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
