package access_service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_xray "github.com/tsvetkovdm/go-vpnbot/internal/core/xray"
)

func (s *AccessService) Reconcile(
	ctx context.Context,
	dryRun bool,
) (domain.XraySyncReport, error) {
	tags, err := s.accessRepository.ListInboundTags(ctx)
	if err != nil {
		return domain.XraySyncReport{}, fmt.Errorf("reconcile: %w", err)
	}

	entries, err := s.accessRepository.ListDesiredAccess(ctx)
	if err != nil {
		return domain.XraySyncReport{}, fmt.Errorf("reconcile: %w", err)
	}

	desired := make(map[string]map[string]string, len(tags))
	for _, e := range entries {
		if desired[e.InboundTag] == nil {
			desired[e.InboundTag] = make(map[string]string)
		}
		desired[e.InboundTag][e.Email] = e.XrayUUID.String()
	}

	report := domain.XraySyncReport{DryRun: dryRun}

	for _, tag := range tags {
		actual, err := s.xrayClient.ListUsers(ctx, tag)
		if err != nil {
			report.Errors = append(report.Errors, err)
			continue
		}

		toAdd, toRemove := diffInbound(desired[tag], actual)

		for _, email := range toRemove {
			if !dryRun {
				if err := s.xrayClient.RemoveUser(ctx, tag, email); err != nil {
					report.Errors = append(report.Errors, err)
					continue
				}
			}
			report.Removed = append(report.Removed, domain.XraySyncChange{InboundTag: tag, Email: email})
		}

		for _, user := range toAdd {
			if !dryRun {
				if err := s.xrayClient.AddUser(ctx, tag, user); err != nil {
					report.Errors = append(report.Errors, err)
					continue
				}
			}
			report.Added = append(report.Added, domain.XraySyncChange{InboundTag: tag, Email: user.Email})
		}
	}

	return report, nil
}

func diffInbound(
	desired map[string]string,
	actual []core_xray.User,
) (toAdd []core_xray.User, toRemove []string) {
	present := make(map[string]string, len(actual))

	for _, u := range actual {
		if !domain.IsBotUserEmail(u.Email) {
			continue // ручной клиент из config.json — не наш
		}
		present[u.Email] = u.ID

		wantID, ok := desired[u.Email]
		if !ok || !strings.EqualFold(wantID, u.ID) {
			toRemove = append(toRemove, u.Email)
		}
	}

	for email, id := range desired {
		haveID, ok := present[email]
		if !ok || !strings.EqualFold(haveID, id) {
			toAdd = append(toAdd, core_xray.User{Email: email, ID: id})
		}
	}

	sort.Strings(toRemove)
	sort.Slice(toAdd, func(i, j int) bool { return toAdd[i].Email < toAdd[j].Email })

	return toAdd, toRemove
}
