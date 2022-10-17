package chat_bot_connector

import (
	"context"
	"fmt"

	"github.com/wanderer69/go-chat-bot-connector/src/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func GrpcInit(address string, port int) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	ss := fmt.Sprintf("%v:%v", address, port) // 5300
	conn, err := grpc.Dial(ss, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	return conn, nil
}

// Check(context.Context, *CheckRequest) (*CheckResponse, error)
func GrpcCheck(conn *grpc.ClientConn, query string) (string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.CheckRequest{
		Query: query,
	}
	response, err := client.Check(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "error", err
	}
	//fmt.Println(response)
	return response.Result, nil
}

// ParsePhrase(context.Context, *ParsePhraseRequest) (*ParsePhraseResponse, error)
func GrpcParsePhrase(conn *grpc.ClientConn, userID string, sessionID string, sentence string, sequenseNum int) (string, string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.ParsePhraseRequest{
		UserId:      userID,
		SessionId:   sessionID,
		Phrase:      sentence,
		SequenseNum: int32(sequenseNum),
	}
	//	fmt.Printf("request %#v\r\n", request)
	response, err := client.ParsePhrase(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", "", "", err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}
	return response.Result, response.QueryId, e, nil
}

// CheckParsePhrase(context.Context, *CheckParsePhraseRequest) (*CheckParsePhraseResponse, error)
func GrpcCheckParsePhrase(conn *grpc.ClientConn, id string) (string, string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.CheckParsePhraseRequest{
		QueryId: id,
	}
	response, err := client.CheckParsePhrase(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", "", "", err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}

	return response.Result, response.Phrase, e, nil
}

// SetSettings(context.Context, *SetSettingsRequest) (*SetSettingsResponse, error)
func GrpcSetSettings(conn *grpc.ClientConn, state string, rl string, ptoc string, gl string) (string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.SetSettingsRequest{
		State: state,
		Settings: &proto.Settings{
			RelationList:    rl,
			PathToOpCorpora: ptoc,
			GrammaticsList:  gl,
		},
	}
	response, err := client.SetSettings(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", "", err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}

	return response.Result, e, nil
}

// Stat(context.Context, *StatRequest) (*StatResponse, error)
func GrpcStat(conn *grpc.ClientConn, mode string) (string, string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.StatRequest{
		Mode: mode,
	}
	response, err := client.Stat(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", "", "", err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}
	return response.Result, response.Info, e, nil
}
