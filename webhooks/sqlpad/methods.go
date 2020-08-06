package sqlpad

import (
	"context"
	"fmt"
	"net/http"
)

const (
	queriesPath = "/api/queries"
	usersPath   = "/api/users"
)

// GetUsers возвращает список пользователей
func (c *Client) GetUsers(ctx context.Context) (users []User, err error) {
	err = c.doRequest(ctx, http.MethodGet, usersPath, nil, &users)
	if err != nil {
		return nil, fmt.Errorf("[sqlpad_client] failed to get users: %w", err)
	}

	return users, nil
}

// DeleteUser удаляет пользователя
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", usersPath, userID), nil, nil)
	if err != nil {
		return fmt.Errorf("[sqlpad_client] failed to delete user: %w", err)
	}

	return nil
}

// DeleteQuery удаляет пользователя
func (c *Client) DeleteQuery(ctx context.Context, queryID string) error {
	err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", queriesPath, queryID), nil, nil)
	if err != nil {
		return fmt.Errorf("[sqlpad_client] failed to delete query: %w", err)
	}

	return nil
}

// GetQuery возвращает запрос по id
func (c *Client) GetQuery(ctx context.Context, queryID string) (query Query, err error) {
	err = c.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", queriesPath, queryID), nil, &query)
	if err != nil {
		return query, fmt.Errorf("[sqlpad_client] failed to get query: %w", err)
	}

	return query, nil
}

// ClearQueryACL очищает ACL для запроса
func (c *Client) ClearQueryACL(ctx context.Context, queryID string) (err error) {
	err = c.doRequest(ctx, http.MethodPut, fmt.Sprintf("%s/%s", queriesPath, queryID), ACLRequest{ACL: []ACL{}}, nil)
	if err != nil {
		return fmt.Errorf("[sqlpad_client] failed to put query: %w", err)
	}

	return nil
}
