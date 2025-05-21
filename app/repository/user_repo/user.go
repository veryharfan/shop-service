package userrepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"shop-service/app/domain"
	"shop-service/pkg"
)

type userRepository struct {
	httpClient         *http.Client
	baseURL            string
	internalAuthHeader string
}

func NewUserRepository(httpClient *http.Client, baseURL string, internalAuthHeader string) domain.UserRepository {
	return &userRepository{
		httpClient:         httpClient,
		baseURL:            baseURL,
		internalAuthHeader: internalAuthHeader,
	}
}

func (r *userRepository) PatchUserShop(ctx context.Context, userID int64, req domain.UserShopUpdateRequest) error {
	url := fmt.Sprintf("%s/internal/user-service/users/%d/shop", r.baseURL, userID)
	reqBody, err := json.Marshal(req)
	if err != nil {
		slog.ErrorContext(ctx, "[userRepository] PatchUserShop", "json.Marshal", err)
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(reqBody))
	if err != nil {
		slog.ErrorContext(ctx, "[userRepository] PatchUserShop", "http.NewRequestWithContext", err)
		return err
	}

	pkg.AddRequestHeader(ctx, r.internalAuthHeader, httpReq)

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		slog.ErrorContext(ctx, "[userRepository] PatchUserShop", "httpClient.Do", err)
		return err
	}
	defer resp.Body.Close()

	var response any
	if err := pkg.DecodeResponseBody(resp, &response); err != nil {
		slog.ErrorContext(ctx, "[userRepository] PatchUserShop", "DecodeResponseBody", err)
		return err
	}

	return nil
}
