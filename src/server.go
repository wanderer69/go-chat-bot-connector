package chat_bot_connector

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	proto "github.com/wanderer69/go-chat-bot-connector/src/proto_internal"
)

type server struct {
	find  func(context.Context, *proto.FindRequest) (*proto.FindResponse, error)
	stat  func(context.Context, *proto.StatRequest) (*proto.StatResponse, error)
	check func(context.Context, *proto.CheckRequest) (*proto.CheckResponse, error)
}

func (s *server) Check(ctx context.Context, in *proto.CheckRequest) (*proto.CheckResponse, error) {
	r, err := s.check(ctx, in)
	return r, err
}

func (s *server) Find(ctx context.Context, in *proto.FindRequest) (*proto.FindResponse, error) {
	r, err := s.find(ctx, in)
	return r, err
}

func (s *server) Stat(ctx context.Context, in *proto.StatRequest) (*proto.StatResponse, error) {
	r, err := s.stat(ctx, in)
	return r, err
}

type CheckFunc func(ctx context.Context, query string) (string, error)
type Find func(ctx context.Context, query string, desired string, essence string, fc []FindCondition) (string, string, []FindFound, error)
type Stat func(ctx context.Context, mode string) (string, string, string, error)

func G_RPC_server(s *Settings, tasker func(s *Settings, cmd_ch chan *Command, answer_ch chan *CommandAnswer)) error {
	ss := fmt.Sprintf(":%v", s.PortServer) // 5300
	listener, err := net.Listen("tcp", ss)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
		return err
	}

	cmd_ch := make(chan *Command)
	answer_ch := make(chan *CommandAnswer)
	go tasker(s, cmd_ch, answer_ch)

	exec_cmd := func(ca *CommandAnswer) (*CommandAnswer, error) {
		// надо дождаться ответа
		ca = nil
		flag := false
		for {
			select {
			case ca = <-answer_ch:
				flag = true
			}
			if flag {
				break
			}
		}
		return ca, nil
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	server_item := server{
		check: func(c context.Context, request *proto.CheckRequest) (*proto.CheckResponse, error) {
			// fmt.Printf("%v\r\n", request.Query)
			cmd := Command{}
			cmd.Cmd = "check"
			cmd.ID = request.Query
			cmd_ch <- &cmd
			var ca *CommandAnswer
			ca, err := exec_cmd(ca)
			result := ""
			if err != nil {
				result = fmt.Sprintf("Error %v", err)
			} else {
				if ca == nil {
					result = "error ca == nil!"
				} else {
					result = ca.Result + " " + ca.Error
				}
			}
			response := &proto.CheckResponse{
				Result: result,
			}
			return response, nil
		},
		find: func(ctx context.Context, request *proto.FindRequest) (*proto.FindResponse, error) {
			// fmt.Printf("%v\r\n", request.Query)
			//result := "Error"
			cmd := Command{}
			cmd.Cmd = "find"
			//cmd.Sentence = request.Sentence

			cmd_ch <- &cmd

			var ca *CommandAnswer
			ca, err := exec_cmd(ca)

			result := ""
			var error_ *string
			if err != nil {
				result = "Error"
				err_ := fmt.Sprintf("%v", err)
				error_ = &err_
			} else {
				result = ca.Result
				error_ = &ca.Error
			}
			response := &proto.FindResponse{
				Result: result,
				Error:  error_,
				//ReqId: ca.ID,
			}
			return response, nil
		},
		stat: func(ctx context.Context, request *proto.StatRequest) (*proto.StatResponse, error) {
			result := ""
			response := &proto.StatResponse{
				Result: result,
			}
			return response, nil
		},
	}
	proto.RegisterChatBotInternalServer(grpcServer, proto.ChatBotInternalServer(&server_item))
	grpcServer.Serve(listener)

	return nil
}
