package adapter

import (
	"context"
	"database/sql"

	"github.com/stdevHsequeda/SendToChannelsBot/adapter/model"
	"github.com/stdevHsequeda/SendToChannelsBot/domain"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ChannelMessageModel struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
}

type PostgresMessageRepository struct {
	conn *sql.DB
}

func NewPostgresMessageRepository(conn *sql.DB) *PostgresMessageRepository {
	return &PostgresMessageRepository{conn}
}

func (p *PostgresMessageRepository) Save(ctx context.Context, msg domain.Message) error {
	messageModel, err := unmarshalMessage(msg)
	if err != nil {
		return err
	}

	return messageModel.Upsert(ctx, p.conn, true, nil, boil.Infer(), boil.Infer())
}

func (p *PostgresMessageRepository) GetByID(ctx context.Context, messageID string) (domain.Message, error) {
	message, err := model.FindMessage(ctx, p.conn, messageID)
	if err != nil {
		return domain.Message{}, err
	}

	return marshalMessage(message)
}

func marshalMessage(model *model.Message) (domain.Message, error) {
	var channelMessageModels []ChannelMessageModel
	if err := model.ChannelMessages.Unmarshal(&channelMessageModels); err != nil {
		return domain.Message{}, err
	}

	var channelMessages []domain.ChannelMessage
	if len(channelMessageModels) != 0 {
		channelMessages = make([]domain.ChannelMessage, len(channelMessageModels), cap(channelMessageModels))
		for i, chMsgModel := range channelMessageModels {
			channelMessages[i] = domain.ChannelMessage{
				ID:        chMsgModel.ID,
				ChannelID: chMsgModel.ChannelID,
			}
		}
	}

	messageID, err := domain.NewMessageIDFromStr(model.ID)
	if err != nil {
		return domain.Message{}, err
	}

	return domain.Message{
		ID:              messageID,
		Hashtags:        model.Hashtags,
		ChannelMessages: channelMessages,
	}, nil
}

func unmarshalMessage(message domain.Message) (*model.Message, error) {
	var channelMessageModels = make([]ChannelMessageModel, len(message.ChannelMessages), cap(message.ChannelMessages))
	for i, chMsg := range message.ChannelMessages {
		channelMessageModels[i] = ChannelMessageModel{
			ID:        chMsg.ID,
			ChannelID: chMsg.ChannelID,
		}
	}

	model := &model.Message{
		ID:       message.ID.String(),
		Hashtags: message.Hashtags,
	}
	if err := model.ChannelMessages.Marshal(channelMessageModels); err != nil {
		return nil, err
	}

	return model, nil
}
