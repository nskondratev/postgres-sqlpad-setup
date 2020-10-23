package users

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nskondratev/postgres-sqlpad-setup/webhooks/sqlpad"
)

const editorRole = "editor"

// Cleanuper получает список пользователей и удаляет аккаунты с ролью editor, которые слишом старые
type Cleanuper struct {
	cli     *sqlpad.Client
	userTTL time.Duration
}

// NewCleanuper ...
func NewCleanuper(cli *sqlpad.Client, userTTL time.Duration) *Cleanuper {
	return &Cleanuper{
		cli:     cli,
		userTTL: userTTL,
	}
}

// Cleanup ...
func (c *Cleanuper) Cleanup(ctx context.Context) error {
	users, err := c.getUsersToCleanUp(ctx)
	if err != nil {
		return err
	}

	return c.cleanUpUsers(ctx, users)
}

func (c *Cleanuper) getUsersToCleanUp(ctx context.Context) ([]sqlpad.User, error) {
	users, err := c.cli.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("[users_cleanuper] failed to get users: %w", err)
	}

	res := make([]sqlpad.User, 0, len(users))
	for _, u := range users {
		if u.Role == editorRole && time.Now().After(u.CreatedAt.Add(c.userTTL)) {
			log.Printf(
				"Delete user with email <%s>, he was created at: %s\n",
				u.Email,
				u.CreatedAt.Format(time.RFC3339),
			)
			res = append(res, u)
		}
	}

	return res, nil
}

func (c *Cleanuper) cleanUpUsers(ctx context.Context, users []sqlpad.User) error {
	if len(users) == 0 {
		return nil
	}

	for _, u := range users {
		err := c.cli.DeleteUser(ctx, u.ID)
		if err != nil {
			return fmt.Errorf("[users_cleanuper] failed to delete user(%s): %w", u.ID, err)
		}
	}

	return nil
}
