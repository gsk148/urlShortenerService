package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gsk148/urlShorteningService/internal/app/api"
	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/hashutil"
	"github.com/gsk148/urlShorteningService/internal/app/logger"
	pb "github.com/gsk148/urlShorteningService/internal/app/proto"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

type ShortenerService struct {
	pb.UnimplementedShortenerServiceServer
	strg storage.Storage
	log  zap.SugaredLogger
}

const (
	HeaderUserID = "x-user-id"
)

var (
	ErrMissingMetadata = errors.New("failed to get metadata from context")
)

func (s *ShortenerService) BatchShortenAPI(ctx context.Context, in *pb.BatchShortenAPIRequest) (*pb.BatchShortenAPIResponse, error) {
	var resp pb.BatchShortenAPIResponse
	urls := in.GetEntities()
	shortenerData := protoURLInfoToModel(urls)
	_, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	result := shortenerData
	resp.Entities = modelURLInfoToProto(result)
	return &resp, nil
}

func (s *ShortenerService) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	var resp pb.DeleteURLsResponse
	urls := in.GetShortUrl()
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	for _, v := range urls {
		err := s.strg.DeleteByUserIDAndShort(userID, v)
		if err != nil {
			return nil, err
		}
	}
	return &resp, nil
}

func (s *ShortenerService) FindByShortLink(ctx context.Context, in *pb.FindByShortLinkRequest) (*pb.URLInfo, error) {
	var resp pb.URLInfo
	url := in.GetShortUrl()
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "no url in request")
	}
	res, err := s.strg.Get(url)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while get short url in storage")
	}
	resp.OriginalUrl = res.OriginalURL
	return &resp, nil
}

func (s *ShortenerService) FindUserURLS(ctx context.Context, in *pb.FindUserURLSRequest) (*pb.BatchShortenAPIResponse, error) {
	var resp pb.BatchShortenAPIResponse
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}
	results, err := s.strg.GetBatchByUserID(userID)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while get urls in storage")
	}
	urls := modelShortenedDataToProto(results)
	resp.Entities = urls
	return &resp, nil
}

func (s *ShortenerService) GetStats(context.Context, *pb.GetStatisticRequest) (*pb.GetStatisticResponse, error) {
	var resp pb.GetStatisticResponse
	stat := s.strg.GetStatistic()
	resp.Urls = int32(stat.URLs)
	resp.Users = int32(stat.Users)
	return &resp, nil
}

func (s *ShortenerService) Ping(context.Context, *pb.PingRequest) (*pb.PingResponse, error) {
	err := s.strg.Ping()
	return &pb.PingResponse{}, err
}

func (s *ShortenerService) ShortenAPI(ctx context.Context, in *pb.ShortenAPIRequest) (*pb.ShortenAPIResponse, error) {
	var resp pb.ShortenAPIResponse
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	originUrl := in.GetUrl()
	shortURL := hashutil.Encode([]byte(originUrl))

	shortenedData := api.ShortenedData{
		UserID:      userID,
		UUID:        uuid.New().String(),
		ShortURL:    shortURL,
		OriginalURL: originUrl,
		IsDeleted:   false,
	}

	res, err := s.strg.Store(shortenedData)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while post long url in storage")
	}
	resp.Result = res.ShortURL
	return &resp, nil
}

func (s *ShortenerService) Shorten(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	resp := pb.ShortenResponse{}
	url := in.GetOriginalUrl()
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no userID in metadata")
	}

	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "no url in request")
	}

	shortURL := hashutil.Encode([]byte(url))

	shortenedData := api.ShortenedData{
		UserID:      userID,
		UUID:        uuid.New().String(),
		ShortURL:    shortURL,
		OriginalURL: url,
		IsDeleted:   false,
	}

	short, err := s.strg.Store(shortenedData)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "error while post long url in storage")
	}
	resp.ShortUrl = short.ShortURL
	return &resp, nil
}

func main() {
	cfg := config.Load()
	log := logger.NewLogger()
	myLog := logger.NewLogger()
	store, err := storage.NewStorage(*cfg, *myLog)
	if err != nil {
		log.Fatal(err)
	}
	st := store
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterShortenerServiceServer(
		s, &ShortenerService{
			strg:                                st,
			UnimplementedShortenerServiceServer: pb.UnimplementedShortenerServiceServer{},
		})
	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func protoURLInfoToModel(urls []*pb.URLInfo) []api.URLInfo {
	var convertedURLS []api.URLInfo
	for _, v := range urls {
		newURL := api.URLInfo{
			UUID:          uuid.NewString(),
			UserID:        v.UserID,
			CorrelationID: v.CorrelationId,
			OriginalURL:   v.OriginalUrl,
			ShortURL:      v.ShortUrl,
			IsDeleted:     v.GetIsDeleted(),
		}
		convertedURLS = append(convertedURLS, newURL)
	}
	return convertedURLS
}

func modelURLInfoToProto(urls []api.URLInfo) []*pb.URLInfo {
	var convertedURLS []*pb.URLInfo
	for _, v := range urls {
		newURL := pb.URLInfo{
			Uuid:          v.UUID,
			UserID:        v.UserID,
			CorrelationId: v.CorrelationID,
			OriginalUrl:   v.OriginalURL,
			ShortUrl:      v.ShortURL,
			IsDeleted:     v.IsDeleted,
		}
		convertedURLS = append(convertedURLS, &newURL)
	}
	return convertedURLS
}

func modelShortenedDataToProto(urls []api.ShortenedData) []*pb.URLInfo {
	var convertedURLS []*pb.URLInfo
	for _, v := range urls {
		newURL := pb.URLInfo{
			Uuid:          v.UUID,
			UserID:        v.UserID,
			CorrelationId: v.OriginalURL,
			OriginalUrl:   v.OriginalURL,
			ShortUrl:      v.ShortURL,
			IsDeleted:     v.IsDeleted,
		}
		convertedURLS = append(convertedURLS, &newURL)
	}
	return convertedURLS
}

func getUserIDFromMD(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingMetadata
	}
	value, ok := GetMetadataValue(md, HeaderUserID)
	if !ok {
		return "", fmt.Errorf("failed to get %s header from metadata: no values", HeaderUserID)
	}
	return value, nil
}

func GetMetadataValue(md metadata.MD, name string) (string, bool) {
	values := md.Get(name)
	if len(values) == 0 {
		return "", false
	}

	for _, v := range values {
		if v != "" {
			return v, true
		}
	}

	return "", false
}
