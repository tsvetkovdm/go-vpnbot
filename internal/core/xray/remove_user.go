package core_xray

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/serial"
)

func (c *Client) RemoveUser(ctx context.Context, inboundTag string, email string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.handler.AlterInbound(ctx, &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: email,
		}),
	})
	if err != nil {
		return fmt.Errorf("remove user email=%s tag=%s: %w", email, inboundTag, err)
	}

	return nil
}
