package jobs

import (
	"context"
	"github.com/google/uuid"
	"go-server/pkg/logger"
	"go.uber.org/zap"
	"time"
)

type runLogKey struct{}

func jobLogger(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(runLogKey{}).(*zap.Logger); ok {
		return log
	}
	return logger.GetLogger()
}

// executeRun owns lifecycle emission. A killed process intentionally leaves an
// unmatched start; in-memory telemetry cannot claim durable completion.
func executeRun(ctx context.Context, name, runID string, run func(context.Context) (RunResult, error)) {
	if runID == "" {
		runID = uuid.NewString()
	}
	log := logger.With(zap.String("job_name", name), zap.String("run_id", runID))
	ctx = context.WithValue(ctx, runLogKey{}, log)
	started := time.Now()
	complete := BeginRun(name)
	log.Info("Job started", zap.String("event", "job_started"))
	result, err := run(ctx)
	outcome := OutcomeSuccess
	switch {
	case isContextCancellation(err):
		outcome = OutcomeInterrupted
	case err != nil:
		outcome = OutcomeFailure
	case !result.LockAcquired:
		outcome = OutcomeSkipped
	}
	complete(outcome)
	fields := []zap.Field{zap.String("event", "job_completed"), zap.String("outcome", outcome), zap.Float64("duration_ms", float64(time.Since(started))/float64(time.Millisecond))}
	if err != nil {
		fields = append(fields, zap.String("internal_error", logger.Sanitize(err.Error())))
	}
	if outcome == OutcomeFailure {
		log.Error("Job completed", fields...)
	} else {
		log.Info("Job completed", fields...)
	}
}
