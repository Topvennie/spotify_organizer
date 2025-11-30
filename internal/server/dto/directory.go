package dto

import (
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/pkg/utils"
)

type Directory struct {
	ID        int         `json:"id"`
	Name      string      `json:"name" validate:"required"`
	Children  []Directory `json:"children,omitzero"`
	Playlists []Playlist  `json:"playlists" validate:"required"`
}

func DirectoryDTO(d *model.Directory, models []*model.Directory) Directory {
	childrenModels := utils.SliceFilter(models, func(m *model.Directory) bool { return m.ParentID == d.ID })
	children := make([]Directory, 0, len(childrenModels))

	for _, child := range childrenModels {
		children = append(children, DirectoryDTO(child, models))
	}

	return Directory{
		ID:        d.ID,
		Name:      d.Name,
		Children:  children,
		Playlists: utils.SliceMap(d.Playlists, func(p model.Playlist) Playlist { return PlaylistDTO(&p, &p.Owner) }),
	}
}

func (d Directory) ToModel(userID, parentID int) *model.Directory {
	return &model.Directory{
		ID:        d.ID,
		UserID:    userID,
		Name:      d.Name,
		ParentID:  parentID,
		Playlists: utils.SliceMap(d.Playlists, func(p Playlist) model.Playlist { return *p.ToModel(userID) }),
	}
}
