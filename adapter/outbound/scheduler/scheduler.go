package scheduler

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/robertobff/nexpos/adapter/outbound/auth"
	"github.com/robertobff/nexpos/utils"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type Scheduler struct {
	cron   *cron.Cron
	fb     *auth.Firebase
	logger *zap.SugaredLogger
}

var Module = fx.Module(
	"scheduler",
	fx.Provide(NewScheduler),
)

func NewScheduler(logger *zap.SugaredLogger, fb *auth.Firebase) *Scheduler {
	s := &Scheduler{
		cron:   cron.New(),
		fb:     fb,
		logger: logger,
	}
	s.reloadPendingDeletions(context.Background())
	return s
}

func (s *Scheduler) reloadPendingDeletions(ctx context.Context) {
	s.logger.Info("reloading pending deletions")
	iter := s.fb.Iter(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		s.logger.Error("failed to reload deletion requests: ", err)
		return
	}

	for _, doc := range docs {
		var req auth.DeletionRequest
		if err := doc.DataTo(&req); err != nil {
			s.logger.Error("failed to decode deletion request ", doc.Ref.ID, ": ", err)
			continue
		}

		if time.Now().After(*req.ExpiresAt) {
			err := s.fb.DisableUser(ctx, req.UID)
			if err != nil {
				s.logger.Error("failed to disable user ", *req.UID, ": ", err)
				continue
			}
			s.logger.Info("user", *req.UID, "  disabled (expired)")

			err = s.fb.UpdateFirestore(req.UID, []firestore.Update{
				{Path: "is_completed", Value: true},
			})
			if err != nil {
				s.logger.Error("failed to mark request as completed for UID ", *req.UID, ": ", err)
			}
			continue
		}

		jobID, err := s.cron.AddFunc("0 0 * * *", func() {
			innerReq, err := s.fb.GetDeletionRequest(ctx, req.UID)
			if err != nil {
				s.logger.Error("failed to fetch deletion request for UID ", *req.UID, ": ", err)
				return
			}

			if innerReq == nil || *innerReq.Cancelled {
				s.logger.Info("deletion for UID ", *req.UID, "  was canceled, skipping")
				return
			}

			if time.Now().After(*innerReq.ExpiresAt) || time.Now().Equal(*innerReq.ExpiresAt) {
				err = s.fb.DisableUser(ctx, req.UID)
				if err != nil {
					s.logger.Error("failed to disable user ", *req.UID, ": ", err)
					return
				}
				s.logger.Info("user ", *req.UID, " disabled after 30 days")

				err = s.fb.UpdateFirestore(req.UID, []firestore.Update{
					{Path: "is_completed", Value: true},
				})
				if err != nil {
					s.logger.Error("failed to mark deletion request as completed for UID ", *req.UID, ": ", err)
				}
			}
		})
		if err != nil {
			s.logger.Error("failed to reschedule deletion for UID ", *req.UID, ": ", err)
			continue
		}

		err = s.fb.UpdateFirestore(req.UID, []firestore.Update{
			{Path: "cron_job_id", Value: int(jobID)},
		})
		if err != nil {
			s.logger.Error("failed to update CronJobID for UID ", *req.UID, ": ", err)
			s.cron.Remove(jobID)
			continue
		}

		s.logger.Info("rescheduled deletion for UID ", *req.UID, " at ", req.ExpiresAt.String(), " (Job ID: ", jobID, ")")
	}
}

func (s *Scheduler) ScheduleUserDeletion(ctx context.Context, id, externalId *string) error {
	if id == nil || externalId == nil {
		return errors.New("id is required")
	}

	expiresAt := time.Now().AddDate(0, 0, 30)

	jobID, err := s.cron.AddFunc("0 0 * * *", func() {
		req, err := s.fb.GetDeletionRequest(ctx, id)
		if err != nil {
			s.logger.Error("error verifying delete request for UID ", *externalId, " : ", err)
			return
		}

		if req == nil || *req.Cancelled {
			s.logger.Info("deletion for UID:  ", *externalId, " was canceled, ignoring")
			return
		}

		err = s.fb.DisableUser(ctx, externalId)
		if err != nil {
			s.logger.Error("error disabling user %s: %v", *externalId, err)
			return
		}
		s.logger.Info("user: ", *externalId, " disabled after 30 days")

		err = s.fb.UpdateFirestore(id, []firestore.Update{
			{Path: "is_completed", Value: true},
		})
		if err != nil {
			s.logger.Error("error marking delete request as completed for UID ", *externalId, " : ", err)
		}
	})
	if err != nil {
		s.logger.Error("error scheduling deletion for UID ", *externalId, " : ", err)
		return err
	}

	err = s.fb.SetFirestore(id, auth.DeletionRequest{
		UID:         externalId,
		CreatedAt:   utils.PTime(time.Now()),
		ExpiresAt:   utils.PTime(expiresAt),
		CronJobID:   utils.PInt(int(jobID)),
		IsCompleted: utils.PBool(false),
		Cancelled:   utils.PBool(false),
	})
	if err != nil {
		s.logger.Error("error scheduling deletion for UID ", *externalId, " : ", err)
		s.cron.Remove(jobID)
		return err
	}

	s.logger.Info("scheduled deletion for UID ", *externalId, " on ", expiresAt.String(), " (Job ID: ", jobID, ")")
	return nil
}

func (s *Scheduler) CancelUserDeletion(ctx context.Context, id *string) error {
	req, err := s.fb.GetDeletionRequest(ctx, id)
	if err != nil {
		s.logger.Error("error verifying delete request for UID ", *id, " : ", err)
		return err
	}

	if req == nil || *req.Cancelled {
		s.logger.Info("deletion for UID:  ", *id, " was canceled, ignoring")
		return nil
	}

	s.cron.Remove(cron.EntryID(*req.CronJobID))

	err = s.fb.UpdateFirestore(id, []firestore.Update{
		{Path: "cancelled", Value: true},
		{Path: "is_completed", Value: true},
	})
	if err != nil {
		s.logger.Error("error marking delete request as completed for UID ", *id, " : ", err)
		return err
	}

	s.logger.Info("scheduled deletion canceled for UID ", id)
	return nil
}
