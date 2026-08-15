package repository

import (
	"context"

	"github.com/janghanul090801/pico-backend/domain"
	"github.com/janghanul090801/pico-backend/ent"
	"github.com/janghanul090801/pico-backend/ent/notification"
	"github.com/janghanul090801/pico-backend/ent/user"
	"github.com/janghanul090801/pico-backend/internal/collections"
)

type notificationRepository struct {
	client *ent.Client
}

func NewNotificationRepository(client *ent.Client) domain.NotificationRepository {
	return &notificationRepository{
		client: client,
	}
}

func (r *notificationRepository) Create(
	ctx context.Context,
	n *domain.Notification,
) (*domain.Notification, error) {
	created, err := r.client.Notification.
		Create().
		SetReceiverID(n.ReceiverID).
		SetType(int(n.Type)).
		SetMessage(n.Message).
		SetMetadata(n.Metadata).
		SetIsRead(n.IsRead).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return toDomainNotification(created), nil
}

func (r *notificationRepository) FindAllByReceiverID(
	c context.Context,
	receiverID *domain.ID,
) ([]*domain.Notification, error) {
	notifications, err := r.client.Notification.
		Query().
		Where(
			notification.HasReceiverWith(
				user.IDEQ(*receiverID),
			),
		).
		Order(ent.Desc(notification.FieldCreatedAt)).
		All(c)

	if err != nil {
		return nil, err
	}

	return collections.Map(notifications, toDomainNotification), nil
}

func (r *notificationRepository) FindByID(c context.Context, id *domain.ID) (*domain.Notification, error) {
	n, err := r.client.Notification.Get(c, *id)
	if err != nil {
		return nil, err
	}

	return toDomainNotification(n), nil
}

func (r *notificationRepository) UpdateIsReadAsTrue(
	c context.Context,
	id *domain.ID,
) error {
	_, err := r.client.Notification.
		Update().
		Where(
			notification.IDEQ(*id),
		).
		SetIsRead(true).
		Save(c)

	if err != nil {
		return err
	}

	return nil
}

func (r *notificationRepository) CountUnreadByReceiverID(
	c context.Context,
	receiverID *domain.ID,
) (int, error) {
	return r.client.Notification.
		Query().
		Where(
			notification.HasReceiverWith(
				user.IDEQ(*receiverID),
			),
			notification.IsReadEQ(false),
		).
		Count(c)
}
