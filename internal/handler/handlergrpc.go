package handler

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/proto"
	"github.com/Piktet/azopkov.git/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p *StorageServer) ShortenURL(ctx context.Context, x *proto.URLShortenRequest) (*proto.URLShortenResponse, error) {
	full := strings.TrimSpace(x.GetUrl())
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Error("error parsing request", zap.Error(err))
		return nil, err
	}

	user, err := getUserFromMd(ctx)
	if err != nil {
		return nil, err
	}
	short, shorterr := p.GetShort(ctx, full, user)
	if shorterr != nil && !errors.Is(shorterr, utils.ErrConflict) {
		if errors.Is(shorterr, utils.ErrorDeleted) {
			logger.Log().Error("error getting short", zap.Error(shorterr))
			return nil, status.Error(codes.NotFound, shorterr.Error())
		}
		logger.Log().Error("error getting short", zap.Error(shorterr))
		return nil, status.Error(codes.Unknown, shorterr.Error())
	}

	p.sendAudit(ctx, model.ActionShorten, user, full)
	return proto.URLShortenResponse_builder{Result: &short}.Build(), nil
}

func (p *StorageServer) ExpandURL(ctx context.Context, x *proto.URLExpandRequest) (*proto.URLExpandResponse, error) {
	if x.GetId() == "" {
		return nil, errors.New("not allowed")
	}

	user, err := getUserFromMd(ctx)
	if err != nil {
		return nil, err
	}
	full, err := p.GetFull(ctx, x.GetId())
	if err != nil {
		if errors.Is(err, utils.ErrorDeleted) {
			logger.Log().Error("error getting full (is deleted)", zap.Error(err))
			return nil, err
		}
		logger.Log().Error("error getting full", zap.Error(err))
		return nil, err
	}

	p.sendAudit(ctx, model.ActionFollow, user, full)

	return proto.URLExpandResponse_builder{Result: &full}.Build(), nil
}

func (p *StorageServer) ListUserURLs(ctx context.Context, x *emptypb.Empty) (*proto.UserURLsResponse, error) {

	user, err := getUserFromMd(ctx)
	if err != nil {
		return nil, err
	}
	response, err := p.GetUserList(ctx, user)
	if err != nil {
		logger.Log().Error("error getting short", zap.Error(err))
		return nil, err
	}

	if len(response) == 0 {
		logger.Log().Error("response is empty", zap.Error(status.Error(codes.NotFound, "response is empty")))
		return nil, status.Error(codes.NotFound, "response is empty")
	}

	logger.Log().Info("response", zap.Int("count", len(response)))

	data := make([]*proto.URLData, 0, len(response))
	for _, v := range response {
		short := p.format(v.Short)
		data = append(data, proto.URLData_builder{OriginalUrl: &v.Full, ShortUrl: &short}.Build())
	}

	logger.Log().Info("response", zap.Any("data", data))

	return proto.UserURLsResponse_builder{Url: data}.Build(), nil
}

func getUserFromMd(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", utils.ErrUserNotFound
	}
	token := md.Get(model.ContextValueAuth)
	if len(token) > 0 {
		return token[0], nil
	}
	return "", utils.ErrUserNotFound

}
