package access_transport_worker

import (
	"context"
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	"go.uber.org/zap"
)

type ReconcileWorker struct {
	accessService AccessService
	config        Config
	log           *core_logger.Logger
}

type AccessService interface {
	Reconcile(
		ctx context.Context,
		dryRun bool,
	) (domain.XraySyncReport, error)
}

func NewReconcileWorker(
	accessService AccessService,
	config Config,
	log *core_logger.Logger,
) *ReconcileWorker {
	return &ReconcileWorker{
		accessService: accessService,
		config:        config,
		log:           log,
	}
}

func (w *ReconcileWorker) Run(ctx context.Context) error {
	w.log.Warn("start xray reconciler",
		zap.Duration("interval", w.config.Interval),
		zap.Bool("dry_run", w.config.DryRun),
	)

	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()

	w.reconcile(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Warn("xray reconciler stopped")
			return nil
		case <-ticker.C:
			w.reconcile(ctx)
		}
	}
}

func (w *ReconcileWorker) reconcile(ctx context.Context) {
	report, err := w.accessService.Reconcile(ctx, w.config.DryRun)
	if ctx.Err() != nil {
		return
	}
	if err != nil {
		w.log.Error("reconcile xray", zap.Error(err))
		return
	}

	if report.IsEmpty() {
		w.log.Debug("xray in sync")
		return
	}

	prefix := ""
	if report.DryRun {
		prefix = "[dry-run] would be "
	}

	for _, c := range report.Removed {
		w.log.Info(prefix+"removed from xray", zap.String("inbound", c.InboundTag), zap.String("email", c.Email))
	}
	for _, c := range report.Added {
		w.log.Info(prefix+"added to xray", zap.String("inbound", c.InboundTag), zap.String("email", c.Email))
	}
	for _, e := range report.Errors {
		w.log.Error("reconcile xray", zap.Error(e))
	}
}
