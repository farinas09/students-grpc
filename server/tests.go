package server

import (
	"context"
	"errors"
	"io"

	"github.com/farinas09/go-grpc/models"
	"github.com/farinas09/go-grpc/repository"
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
