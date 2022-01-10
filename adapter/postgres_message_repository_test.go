package adapter_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stdevHsequeda/SendToChannelsBot/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresMessageRepository_Save(t *testing.T) {
	t.Parallel()

	repo := newPostgresMessageRepository(t)
	t.Run("Ok: Message with two forwards", func(t *testing.T) {
		t.Parallel()
		messageID, err := domain.NewMessageIDFromInt64(112345)
		require.NoError(t, err)
		message, err := domain.NewMessage(messageID, "", []string{"#hashtag", "#other_hashtag"})
		require.NoError(t, err)
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_1", ChannelID: uuid.New().String()})
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_2", ChannelID: uuid.New().String()})
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
	})

	t.Run("Ok: Message without forwards", func(t *testing.T) {
		t.Parallel()
		messageID, err := domain.NewMessageIDFromInt64(int64(uuid.New().ID()))
		require.NoError(t, err)
		message, err := domain.NewMessage(messageID, "", []string{"#hashtag", "#other_hashtag"})
		require.NoError(t, err)
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
	})

	t.Run("Ok: Rewrite message", func(t *testing.T) {
		t.Parallel()
		messageID, err := domain.NewMessageIDFromInt64(112345)
		require.NoError(t, err)
		message, err := domain.NewMessage(messageID, "", []string{"#hashtag", "#other_hashtag"})
		require.NoError(t, err)
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_1", ChannelID: uuid.New().String()})
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
	})
}

func TestPostgresMessageRepository_GetByID(t *testing.T) {
	t.Parallel()

	repo := newPostgresMessageRepository(t)
	t.Run("Ok: Message with two forwards", func(t *testing.T) {
		t.Parallel()
		messageID, err := domain.NewMessageIDFromInt64(int64(uuid.New().ID()))
		require.NoError(t, err)
		message, err := domain.NewMessage(messageID, "", []string{"#hashtag", "#other_hashtag"})
		require.NoError(t, err)
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_1", ChannelID: uuid.New().String()})
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_2", ChannelID: uuid.New().String()})
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
		requestedMsg, err := repo.GetByID(context.Background(), message.ID.String())
		require.NoError(t, err)
		assert.Equal(t, message.ID, requestedMsg.ID)
		assert.Equal(t, message.Hashtags, requestedMsg.Hashtags)
		assert.Equal(t, message.ChannelMessages, requestedMsg.ChannelMessages)
	})
}

func newPostgresMessageRepository(t *testing.T) *adapter.PostgresMessageRepository {
	t.Helper()
	conn, err := adapter.NewPostgresConnPool()
	require.NoError(t, err)
	return adapter.NewPostgresMessageRepository(conn)
}
