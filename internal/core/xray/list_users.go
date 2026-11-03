package core_xray

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/proxy/vless"
)

func (c *Client) ListUsers(ctx context.Context, inboundTag string) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.handler.GetInboundUsers(ctx, &command.GetInboundUserRequest{
		Tag: inboundTag,
	})
	if err != nil {
		return nil, fmt.Errorf("get inbound users tag=%s: %w", inboundTag, err)
	}

	users := make([]User, 0, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		instance, err := u.GetAccount().GetInstance()
		if err != nil {
			return nil, fmt.Errorf("decode account email=%s tag=%s: %w", u.GetEmail(), inboundTag, err)
		}

		account, ok := instance.(*vless.Account)
		if !ok {
			return nil, fmt.Errorf("account email=%s tag=%s: not VLESS (%T)", u.GetEmail(), inboundTag, instance)
		}

		users = append(users, User{
			Email: u.GetEmail(),
			ID:    account.GetId(),
		})
	}

	return users, nil
}
