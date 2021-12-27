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
		message, err := domain.NewMessage("test_id", "", []string{"#hashtag", "#other_hashtag"})
		require.NoError(t, err)
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_1", ChannelID: uuid.New().String()})
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_2", ChannelID: uuid.New().String()})
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
	})

	t.Run("Ok: Message without forwards", func(t *testing.T) {
		t.Parallel()
		message, err := domain.NewMessage(uuid.New().String(), "", []string{"#hashtag", "#other_hashtag"})
		require.NoError(t, err)
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
	})

	t.Run("Ok: Rewrite message", func(t *testing.T) {
		t.Parallel()
		id := uuid.New().String()
		message, err := domain.NewMessage(id, "", []string{"#hashtag", "#other_hashtag"})
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
		message, err := domain.NewMessage("test_id", "", []string{"#hashtag", "#other_hashtag"})
		require.NoError(t, err)
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_1", ChannelID: uuid.New().String()})
		message.AddChannelMessage(domain.ChannelMessage{ID: "test_id_for_channel_2", ChannelID: uuid.New().String()})
		err = repo.Save(context.Background(), message)
		require.NoError(t, err)
		requestedMsg, err := repo.GetByID(context.Background(), message.ID)
		require.NoError(t, err)
		assert.Equal(t, message, requestedMsg)
	})
}

func newPostgresMessageRepository(t *testing.T) *adapter.PostgresMessageRepository {
	t.Helper()
	conn, err := adapter.NewPostgresConnPool()
	require.NoError(t, err)
	return adapter.NewPostgresMessageRepository(conn)
}
