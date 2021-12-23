package adapter_test

import (
	"context"
	"testing"

	"github.com/stdevHsequeda/SendToChannelsBot/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresMessageRepository_Save(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name    string
		Message domain.Message
	}{
		{
			Name: "Ok",
			Message: domain.Message{
				ID: "test_id",
				Hashtags: []string{
					"#hashtag",
					"#other_hashtag",
				},
				ChannelMessages: []domain.ChannelMessage{
					{
						ID:        "test_id_for_channel_1",
						ChannelID: "1",
					},
					{
						ID:        "test_id_for_channel_2",
						ChannelID: "2",
					},
				},
			},
		},
		{
			Name: "Without Channel Messages",
			Message: domain.Message{
				ID: "test_id2",
				Hashtags: []string{
					"#hashtag",
					"#other_hashtag",
				},
				ChannelMessages: nil,
			},
		},
	}

	repo := newPostgresMessageRepository(t)
	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			message := c.Message
			err := repo.Save(ctx, message)
			require.NoError(t, err)
			requestMessage, err := repo.GetByID(ctx, message.ID)
			require.NoError(t, err)
			assert.Equal(t, message, requestMessage)
		})
	}
}

func newPostgresMessageRepository(t *testing.T) *adapter.PostgresMessageRepository {
	t.Helper()
	conn, err := adapter.NewPostgresConnPool()
	require.NoError(t, err)
	return adapter.NewPostgresMessageRepository(conn)
}
