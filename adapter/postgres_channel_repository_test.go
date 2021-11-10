package adapter_test

import (
	"context"
	"testing"

	"github.com/stdevHsequeda/SendToChannelsBot/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresChannelRepository_Save(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name    string
		Channel domain.Channel
	}{
		{
			Name: "Channel1",
			Channel: domain.Channel{
				Id:       "1",
				Hashtags: []string{"cerveza", "refresco", "cocina"},
			},
		},
		{
			Name: "Channel2",
			Channel: domain.Channel{
				Id:       "2",
				Hashtags: []string{"cerveza", "cocina"},
			},
		},
	}

	repo := newPostgresChannelRepository(t)
	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			channel := c.Channel
			err := repo.Save(ctx, channel)
			require.NoError(t, err)
		})
	}
}

func TestPostgresChannelRepository_GetChannelByHashtags(t *testing.T) {
	t.Parallel()
	repo := newPostgresChannelRepository(t)

	err := repo.DeleteAll(context.Background())
	require.NoError(t, err)

	ctx := context.Background()

	channel1 := domain.Channel{
		Id:       "5",
		Hashtags: []string{"lavadora", "cafe"},
	}

	channel2 := domain.Channel{
		Id:       "6",
		Hashtags: []string{"reina", "cafe"},
	}

	channelsToAdd := []domain.Channel{
		channel1,
		channel2,
	}

	for _, ch := range channelsToAdd {
		err := repo.Save(ctx, ch)
		require.NoError(t, err)
	}

	channels, err := repo.GetChannelsByHashtags(context.Background(), []string{"cafe"})
	require.NoError(t, err)
	assert.Len(t, channels, 2)

	channels, err = repo.GetChannelsByHashtags(context.Background(), []string{"lavadora"})
	require.NoError(t, err)
	assert.Len(t, channels, 1)

	channels, err = repo.GetChannelsByHashtags(context.Background(), []string{"levadura"})
	require.NoError(t, err)
	assert.Len(t, channels, 0)
}

func newPostgresChannelRepository(t *testing.T) adapter.PostgresChannelRepository {
	t.Helper()
	conn, err := adapter.NewPostgresConnPool()
	require.NoError(t, err)
	return adapter.NewPostgresChannelRepository(conn)
}
