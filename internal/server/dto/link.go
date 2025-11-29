package dto

import "github.com/topvennie/spotify_organizer/internal/database/model"

type Link struct {
	ID                int `json:"id"`
	SourceDirectoryID int `json:"source_directory_id,omitzero"`
	SourcePlaylistID  int `json:"source_playlist_id,omitzero"`
	TargetDirectoryID int `json:"target_directory_id,omitzero"`
	TargetPlaylistID  int `json:"target_playlist_id,omitzero"`
}

func LinkDTO(l *model.Link) Link {
	return Link{
		ID:                l.ID,
		SourceDirectoryID: l.SourceDirectoryID,
		SourcePlaylistID:  l.SourcePlaylistID,
		TargetDirectoryID: l.TargetDirectoryID,
		TargetPlaylistID:  l.TargetPlaylistID,
	}
}

func (l *Link) ToModel() *model.Link {
	return &model.Link{
		ID:                l.ID,
		SourceDirectoryID: l.SourceDirectoryID,
		SourcePlaylistID:  l.SourcePlaylistID,
		TargetDirectoryID: l.TargetDirectoryID,
		TargetPlaylistID:  l.TargetPlaylistID,
	}
}
