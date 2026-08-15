package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
)

type notificationUseCase struct {
	notificationRepository domain.NotificationRepository
	contextTimeout         time.Duration
}

func NewNotificationUseCase(
	notificationRepository domain.NotificationRepository,
	contextTimeout time.Duration,
) domain.NotificationUseCase {
	return &notificationUseCase{
		notificationRepository: notificationRepository,
		contextTimeout:         contextTimeout,
	}
}

func (u *notificationUseCase) Create(c context.Context, notification *domain.Notification) (*domain.Notification, error) {
	n, err := u.notificationRepository.Create(c, notification)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return n, nil
}

func (u *notificationUseCase) ListByReceiverID(
	c context.Context,
	userID *domain.ID,
) ([]*domain.Notification, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	notifications, err := u.notificationRepository.FindAllByReceiverID(
		ctx,
		userID,
	)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return notifications, nil
}

func (u *notificationUseCase) MarkAsRead(c context.Context, id *domain.ID, receiverID *domain.ID) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	n, err := u.notificationRepository.FindByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.NewNotFoundError(err)
		}
		return domain.NewBadRequestError(err)
	}
	if n.ReceiverID != *receiverID {
		return domain.NewForbiddenError(errors.New("forbidden"))
	}

	return u.notificationRepository.UpdateIsReadAsTrue(ctx, id)
}
func (u *notificationUseCase) CountUnreadByReceiverID(c context.Context, receiverID *domain.ID) (int, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	n, err := u.notificationRepository.CountUnreadByReceiverID(ctx, receiverID)
	if err != nil {
		return 0, domain.NewInternalServerError(err)
	}

	return n, nil
}
