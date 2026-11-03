package core_xray

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vless"
)

func (c *Client) AddUser(ctx context.Context, inboundTag string, user User) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.handler.AlterInbound(ctx, &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Email: user.Email,
				Account: serial.ToTypedMessage(&vless.Account{
					Id: user.ID,
				}),
			},
		}),
	})
	if err != nil {
		return fmt.Errorf("add user email=%s tag=%s: %w", user.Email, inboundTag, err)
	}

	return nil
}
