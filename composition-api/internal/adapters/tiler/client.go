package tiler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Client interface {
	GetDZI(ctx context.Context, filePath string) (string, error)
	GetTile(ctx context.Context, filePath string, level, col, row int, format string) ([]byte, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) Client {
	// Добавляем протокол, если его нет
	url := strings.TrimSpace(baseURL)
	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	return &client{
		baseURL:    strings.TrimSuffix(url, "/"),
		httpClient: &http.Client{},
	}
}

func (c *client) GetDZI(ctx context.Context, filePath string) (string, error) {
	// Формируем URL для DZI
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join("/dzi", filePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *client) GetTile(ctx context.Context, filePath string, level, col, row int, format string) ([]byte, error) {
	// Формируем URL для тайла
	// Формат: /dzi/{file_path}/files/{level}/{col}_{row}.{format}
	tilePath := filePath + "/files/" + fmt.Sprintf("%d/%d_%d.%s", level, col, row, format)

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join("/dzi", tilePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
