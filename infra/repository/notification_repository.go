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
		SetIsRead(false).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	// Create() 직후의 결과에는 Edges.Receiver(관계)가 채워지지 않아서
	// toDomainNotification에서 n.Edges.Receiver.ID 접근 시 panic 났었음.
	// 저장 후 다시 조회하면서 Receiver 관계를 명시적으로 불러옴(WithReceiver).
	created, err = r.client.Notification.
		Query().
		Where(notification.IDEQ(created.ID)).
		WithReceiver().
		Only(ctx)
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
		WithReceiver().
		Order(ent.Desc(notification.FieldCreatedAt)).
		All(c)

	if err != nil {
		return nil, err
	}

	return collections.Map(notifications, toDomainNotification), nil
}

func (r *notificationRepository) FindByID(c context.Context, id *domain.ID) (*domain.Notification, error) {
	n, err := r.client.Notification.
		Query().
		Where(notification.IDEQ(*id)).
		WithReceiver().
		Only(c)
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
func (r *notificationRepository) Delete(c context.Context, id *domain.ID) error {
	return r.client.Notification.
		DeleteOneID(*id).
		Exec(c)
}
