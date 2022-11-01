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
	fmt.Printf("'%v'\r\n", ss)
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

// SetVersion(context.Context, *SetVersionRequest) (*SetVersionResponse, error)
func GrpcSetVersion(conn *grpc.ClientConn, state string, id string, date string, rl string, gl string) (string, string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.SetVersionRequest{
		State: state,
		Version: &proto.Version{
			VersionId: &proto.VersionId{
				Id:    id,
				Date:  date,
			},
			RelationList:    rl,
			GrammaticsList:  gl,
		},
	}
	response, err := client.SetVersion(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", "", "", err
	}
	e := ""
	idn := ""
	if response.Error != nil {
		e = *response.Error
	} else {
		if response.VersionId != nil {
			idn = response.VersionId.Id
		}
	}
	return response.Result, idn, e, nil
}

// GetVersion(context.Context, *GetVersionRequest) (*GetVersionResponse, error)
func GrpcGetVersion(conn *grpc.ClientConn, state string, id string, date string) (string, string, string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.GetVersionRequest{
		State: state,
		VersionId: &proto.VersionId{
			Id:    id,
			Date:  date,
		},
	}
	response, err := client.GetVersion(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", "", "", "", err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}
	rl := ""
	gl := ""
	if response.Version != nil {
		rl = response.Version.RelationList
		gl = response.Version.GrammaticsList
	}

	return response.Result, e, rl, gl, nil
}

type FilterField struct {
	Id string
	Values []string
	Strict bool
}

type SortField struct {
	Id string
	Type string
}

type Pagination struct {
	Limit int32
	Offset int32
}

type GetLogsIn struct {
        FilterFields []FilterField
        SortFields SortField
        TableFields []string
        Pagination Pagination
}

type LogRow struct {
    Id string
    Data string
    UserRequest string
    BotResponse string
    Version string
    UserId string
}

type GetLogsOut struct {
    Count int32
    Rows []LogRow 
}

// GetLogs(context.Context, *GetLogsRequest) (*GetLogsResponse, error)
func GrpcGetLogs(conn *grpc.ClientConn, gl GetLogsIn) (string, string, GetLogsOut, error) {
	client := proto.NewChatBotClient(conn)
	ffa :=  []*proto.FilterField{}
	for i, _ := range gl.FilterFields {
		ff :=  proto.FilterField{
			Id: gl.FilterFields[i].Id,
                        Values: gl.FilterFields[i].Values,
                        Strict: gl.FilterFields[i].Strict,
		}
		ffa = append(ffa, &ff)
	}
	request := &proto.GetLogsRequest{
		FilterFields: ffa,
		SortFields: &proto.SortField{Id: gl.SortFields.Id, Type: gl.SortFields.Type},
		TableFields: gl.TableFields,
		Pagination: &proto.Pagination{Limit: gl.Pagination.Limit, Offset: gl.Pagination.Offset},
	}
	response, err := client.GetLogs(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", "", GetLogsOut{}, err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}
	glo := GetLogsOut{}
	if response.LogResult != nil {
	        glo.Count = response.LogResult.Count
	        for i, _ := range response.LogResult.Rows {
	        	r := LogRow{
	        		response.LogResult.Rows[i].Id,
	        		response.LogResult.Rows[i].Data,
	        		response.LogResult.Rows[i].UserRequest,
	        		response.LogResult.Rows[i].BotResponse,
	        		response.LogResult.Rows[i].Version,
	        		response.LogResult.Rows[i].UserId,
			}
			glo.Rows = append(glo.Rows, r)
	        }
	}

	return response.Result, e, glo, nil
}

type ListVersionsItem struct {
	Id string
	Date string
}

// ListVersions(context.Context, *ListVersionsRequest) (*ListVersionsResponse, error)
func GrpcListVersions(conn *grpc.ClientConn, state string) (string, []ListVersionsItem, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.ListVersionsRequest{
		State: state,
	}
	response, err := client.ListVersions(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", nil, "", err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}
        lvia := []ListVersionsItem{}
	if response.ListVersions != nil {
	     for i, _ := range response.ListVersions.VersionId {
	     	lvi := ListVersionsItem{response.ListVersions.VersionId[i].Id, response.ListVersions.VersionId[i].Date}
	     	lvia = append(lvia, lvi)
	     }
	}

	return response.Result, lvia, e, nil
}

// TestVersion(context.Context, *TestVersionRequest) (*TestVersionResponse, error)
func GrpcTestVersion(conn *grpc.ClientConn, userID string, phrase string, id string, date string, rl string, gl string, sequenseNum int) (string, string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.TestVersionRequest{
		UserId:      userID,
		Phrase:      phrase,
		SequenseNum: int32(sequenseNum),
		Version: &proto.Version{
			VersionId: &proto.VersionId{
				Id:    id,
				Date:  date,
			},
			RelationList:    rl,
			GrammaticsList:  gl,
		},
	}
	//	fmt.Printf("request %#v\r\n", request)
	response, err := client.TestVersion(context.Background(), request)
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

/*
// CheckTestVersion(context.Context, *CheckTestVersionRequest) (*CheckTestVersionResponse, error)
func GrpcCheckTestVersion(conn *grpc.ClientConn, id string) (string, string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.CheckTestVersionRequest{
		QueryId: id,
	}
	response, err := client.CheckTestVersion(context.Background(), request)
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
*/
// GetWord(context.Context, *GetWordRequest) (*GetWordResponse, error)
func GrpcGetWord(conn *grpc.ClientConn, word string) (string, [][]string, string, error) {
	client := proto.NewChatBotClient(conn)
	request := &proto.GetWordRequest{
		Word: word,
	}
	response, err := client.GetWord(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", nil, "", err
	}
	e := ""
	if response.Error != nil {
		e = *response.Error
	}
	wpl := [][]string{}
	if response.WordProps != nil {
                for i, _ := range response.WordProps.WordProps {
                       wp := []string{response.WordProps.WordProps[i].BaseWord, response.WordProps.WordProps[i].PartOfSpeach}
                       wpl = append(wpl, wp)
                }
	}

	return response.Result, wpl, e, nil
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
