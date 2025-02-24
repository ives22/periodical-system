package impl

import (
	"context"
	"periodical/internal/apps/user"
	"testing"
)

var (
	userSvc = NewUserImpl()
	ctx     = context.Background()
)

func TestCreateUser(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.Username = "admin"
	req.Password = "123456"

	u, err := userSvc.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(u)
}
