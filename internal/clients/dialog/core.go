package dialog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	inpb "github.com/syth0le/dialog-service/proto/internalapi"

	"github.com/syth0le/social-network/internal/model"
)

type Client interface {
	CreateDialog(ctx context.Context, userId model.UserID, participants []model.UserID) (*model.Dialog, error)
	CreateMessage(ctx context.Context, dialogID model.DialogID, senderID model.UserID, text string) error
	GetDialogMessages(ctx context.Context, dialogID model.DialogID, userID model.UserID) ([]*model.Message, error)
}

type ClientImpl struct {
	logger *zap.Logger
	client inpb.DialogServiceClient
}

func NewClientImpl(logger *zap.Logger, conn *grpc.ClientConn) *ClientImpl {
	return &ClientImpl{
		logger: logger,
		client: inpb.NewDialogServiceClient(conn),
	}
}

func (c *ClientImpl) CreateDialog(ctx context.Context, userId model.UserID, participants []model.UserID) (*model.Dialog, error) {
	participantsIDs := make([]string, len(participants))
	for idx, item := range participants {
		participantsIDs[idx] = item.String()
	}

	dialog, err := c.client.CreateDialog(ctx, &inpb.CreateDialogRequest{
		UserId:          userId.String(),
		ParticipantsIds: participantsIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("create dialog: %w", err)
	}

	return &model.Dialog{
		ID:              model.DialogID(dialog.DialogId),
		ParticipantsIDs: participants,
	}, nil
}

func (c *ClientImpl) CreateMessage(ctx context.Context, dialogID model.DialogID, senderID model.UserID, text string) error {
	_, err := c.client.CreateMessage(ctx, &inpb.CreateMessageRequest{
		DialogId: dialogID.String(),
		SenderId: senderID.String(),
		Text:     text,
	})
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	return nil
}

func (c *ClientImpl) GetDialogMessages(ctx context.Context, dialogID model.DialogID, userID model.UserID) ([]*model.Message, error) {
	pbMessages, err := c.client.GetDialogMessages(ctx, &inpb.GetDialogMessagesRequest{
		DialogId: dialogID.String(),
		UserId:   userID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	messages := make([]*model.Message, len(pbMessages.Messages))
	for idx, item := range pbMessages.Messages {
		messages[idx] = &model.Message{
			ID:       model.MessageID(item.Id),
			DialogID: model.DialogID(item.DialogId),
			SenderID: model.UserID(item.SenderId),
			Text:     item.Text,
		}
	}

	return messages, nil
}
