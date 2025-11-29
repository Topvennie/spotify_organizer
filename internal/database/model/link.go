package model

import "github.com/topvennie/spotify_organizer/pkg/sqlc"

type Link struct {
	ID                int
	SourceDirectoryID int
	SourcePlaylistID  int
	TargetDirectoryID int
	TargetPlaylistID  int
}

func LinkModel(l sqlc.Link) *Link {
	sourceDirectoryID := 0
	if l.SourceDirectoryID.Valid {
		sourceDirectoryID = int(l.SourceDirectoryID.Int32)
	}
	sourcePlaylistID := 0
	if l.SourcePlaylistID.Valid {
		sourcePlaylistID = int(l.SourcePlaylistID.Int32)
	}
	targetDirectoryID := 0
	if l.TargetDirectoryID.Valid {
		targetDirectoryID = int(l.TargetDirectoryID.Int32)
	}
	targetPlaylistID := 0
	if l.TargetPlaylistID.Valid {
		targetPlaylistID = int(l.TargetPlaylistID.Int32)
	}

	return &Link{
		ID:                int(l.ID),
		SourceDirectoryID: sourceDirectoryID,
		SourcePlaylistID:  sourcePlaylistID,
		TargetDirectoryID: targetDirectoryID,
		TargetPlaylistID:  targetPlaylistID,
	}
}

func (l *Link) Equal(l2 Link) bool {
	return l.SourceDirectoryID == l2.SourceDirectoryID && l.SourcePlaylistID == l2.SourcePlaylistID && l.TargetDirectoryID == l2.TargetDirectoryID && l.TargetPlaylistID == l2.TargetDirectoryID
}
