package spotify

import (
	"context"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/config"
)

const (
	taskPlaylistUID = "task-playlist"
	taskAlbumUID    = "task-album"
	taskShowUID     = "task-show"
	taskTrackUID    = "task-track"
	taskUserUID     = "task-user"
	taskHistoryUID  = "task-history"
)

func (c *client) taskRegister() error {
	if err := task.Manager.Add(context.Background(), task.NewTask(
		taskPlaylistUID,
		"Playlist",
		config.GetDefaultDuration("task.playlist_s", 6*60*60),
		c.taskWrap(c.taskPlaylist),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(context.Background(), task.NewTask(
		taskAlbumUID,
		"Album",
		config.GetDefaultDuration("task.album_s", 6*60*60),
		c.taskWrap(c.taskAlbum),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(context.Background(), task.NewTask(
		taskShowUID,
		"Show",
		config.GetDefaultDuration("task.show_s", 6*60*60),
		c.taskWrap(c.taskShow),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(context.Background(), task.NewTask(
		taskTrackUID,
		"Track",
		config.GetDefaultDuration("task.track_s", 60*60),
		c.taskWrap(c.taskTrack),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(context.Background(), task.NewTask(
		taskUserUID,
		"User",
		config.GetDefaultDuration("task.user_s", 6*60*60),
		c.taskWrap(c.taskUser),
	)); err != nil {
		return err
	}

	if err := task.Manager.Add(context.Background(), task.NewTask(
		taskHistoryUID,
		"History",
		config.GetDefaultDuration("task.history_s", 10*60),
		c.taskWrap(c.taskHistory),
	)); err != nil {
		return err
	}

	return nil
}

func (c *client) taskWrap(fn func(context.Context, model.User) (string, error)) func(context.Context, []model.User) []task.TaskResult {
	return func(ctx context.Context, users []model.User) []task.TaskResult {
		results := make([]task.TaskResult, 0, len(users))

		for _, user := range users {
			msg, err := fn(ctx, user)
			results = append(results, task.TaskResult{
				User:    user,
				Message: msg,
				Error:   err,
			})
		}

		return results
	}
}

func (c *client) taskPlaylist(ctx context.Context, user model.User) (string, error) {
	if err := c.playlistSync(ctx, user); err != nil {
		return "", fmt.Errorf("synchronize playlists %w", err)
	}

	if err := c.playlistUpdate(ctx, user); err != nil {
		return "", fmt.Errorf("update playlists %w", err)
	}

	if err := c.playlistCoverSync(ctx, user); err != nil {
		return "", fmt.Errorf("synchronize playlist covers %w", err)
	}

	return "", nil
}

func (c *client) taskAlbum(ctx context.Context, user model.User) (string, error) {
	if err := c.albumSync(ctx, user); err != nil {
		return "", fmt.Errorf("synchronize albums %w", err)
	}

	if err := c.albumUpdate(ctx, user); err != nil {
		return "", fmt.Errorf("update albums %w", err)
	}

	if err := c.albumCoverSync(ctx, user); err != nil {
		return "", fmt.Errorf("synchronize album covers %w", err)
	}

	return "", nil
}

func (c *client) taskShow(ctx context.Context, user model.User) (string, error) {
	if err := c.showSync(ctx, user); err != nil {
		return "", fmt.Errorf("synchronize shows %w", err)
	}

	if err := c.showUpdate(ctx, user); err != nil {
		return "", fmt.Errorf("update shows %w", err)
	}

	if err := c.showCoverSync(ctx, user); err != nil {
		return "", fmt.Errorf("synchronize show covers %w", err)
	}

	return "", nil
}

func (c *client) taskTrack(ctx context.Context, user model.User) (string, error) {
	msg1, err := c.playlistTrackSync(ctx, user)
	if err != nil {
		return "", fmt.Errorf("synchronize tracks %w", err)
	}

	msg2, err := c.tracksSync(ctx, user)
	if err != nil {
		return "", fmt.Errorf("update playlist tracks based on links %w", err)
	}

	return fmt.Sprintf("%s | %s", msg1, msg2), nil
}

func (c *client) taskUser(ctx context.Context, user model.User) (string, error) {
	if err := c.syncUser(ctx, user); err != nil {
		return "", fmt.Errorf("synchronize users %w", err)
	}

	return "", nil
}

func (c *client) taskHistory(ctx context.Context, user model.User) (string, error) {
	if _, err := c.historySync(ctx, user); err != nil {
		return "", fmt.Errorf("get history %w", err)
	}

	return "", nil
}
