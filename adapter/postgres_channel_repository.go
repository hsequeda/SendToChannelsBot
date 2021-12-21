package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/stdevHsequeda/SendToChannelsBot/adapter/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/stdevHsequeda/SendToChannelsBot/domain"
)

type PostgresChannelRepository struct {
	conn *sql.DB
}

func NewPostgresChannelRepository(conn *sql.DB) *PostgresChannelRepository {
	return &PostgresChannelRepository{conn}
}

func (r *PostgresChannelRepository) GetChannelsByHashtags(ctx context.Context, hashtags []string) ([]domain.Channel, error) {
	var hashtagsWithQuotes = make([]string, len(hashtags), cap(hashtags))
	for e := range hashtags {
		hashtagsWithQuotes[e] = fmt.Sprintf("'%s'", hashtags[e])
	}

	channelModels, err := model.Channels(qm.Where(fmt.Sprintf("hashtags && ARRAY[%s]::varchar[]", strings.Join(hashtagsWithQuotes, ",")))).All(ctx, r.conn)
	if err != nil {
		return nil, err
	}
	channels := make([]domain.Channel, len(channelModels))
	for index := range channelModels {
		channels[index] = unmarshalChannel(channelModels[index])
	}

	return channels, nil
}

func (r *PostgresChannelRepository) Save(ctx context.Context, channel domain.Channel) error {
	model := marshalChannel(channel)
	return model.Insert(ctx, r.conn, boil.Infer())
}

func (r *PostgresChannelRepository) DeleteAll(ctx context.Context) error {
	_, err := model.Channels().DeleteAll(ctx, r.conn)
	return err
}

func unmarshalChannel(channelModel *model.Channel) domain.Channel {
	return domain.Channel{
		Id:       channelModel.ID,
		Hashtags: channelModel.Hashtags,
	}
}

func marshalChannel(channel domain.Channel) *model.Channel {
	return &model.Channel{
		ID:       channel.Id,
		Hashtags: channel.Hashtags,
	}
}
