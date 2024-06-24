package dialog

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/syth0le/social-network/internal/model"
)

type ClientMock struct {
	logger *zap.Logger
}

func NewClientMock(logger *zap.Logger) *ClientMock {
	return &ClientMock{
		logger: logger,
	}
}

func (c ClientMock) CreateDialog(ctx context.Context, userId model.UserID, participants []model.UserID) (*model.Dialog, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c ClientMock) CreateMessage(ctx context.Context, dialogID model.DialogID, senderID model.UserID, text string) error {
	return fmt.Errorf("not implemented")
}

func (c ClientMock) GetDialogMessages(ctx context.Context, dialogID model.DialogID, userID model.UserID) ([]*model.Message, error) {
	return nil, fmt.Errorf("not implemented")
}
