package synology

import (
	"fmt"
	"io"
)

const (
	SYNOFileStationList            = "SYNO.FileStation.List"
	SYNOFileStationInfo            = "SYNO.FileStation.Info"
	SYNOFileStationSearch          = "SYNO.FileStation.Search"
	SYNOFileStationVirtualFolder   = "SYNO.FileStation.VirtualFolder"
	SYNOFileStationFavorite        = "SYNO.FileStation.Favorite"
	SYNOFileStationThumb           = "SYNO.FileStation.Thumb"
	SYNOFileStationDirSize         = "SYNO.FileStation.DirSize"
	SYNOFileStationMD5             = "SYNO.FileStation.MD5"
	SYNOFileStationCheckPermission = "SYNO.FileStation.CheckPermission"
	SYNOFileStationUpload          = "SYNO.FileStation.Upload"
	SYNOFileStationDownload        = "SYNO.FileStation.Download"
	SYNOFileStationSharing         = "SYNO.FileStation.Sharing"
	SYNOFileStationCreateFolder    = "SYNO.FileStation.CreateFolder"
	SYNOFileStationRename          = "SYNO.FileStation.Rename"
	SYNOFileStationCopyMove        = "SYNO.FileStation.CopyMove"
	SYNOFileStationDelete          = "SYNO.FileStation.Delete"
	SYNOFileStationExtract         = "SYNO.FileStation.Extract"
	SYNOFileStationCompress        = "SYNO.FileStation.Compress"
	SYNOFileStationBackgroundTask  = "SYNO.FileStation.BackgroundTask"
)

type FileStationService interface {
	ListShares() ([]FileInfo, error)
	List(path string) ([]FileInfo, error)
	Stat(path string) ([]FileInfo, error)
	Download(path string, w io.Writer) error
}

type FileInfo struct {
	Isdir bool   `json:"isdir"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Stat  struct {
		Owner struct {
			Gid   int    `json:"gid"`
			Group string `json:"group"`
			UID   int    `json:"uid"`
			User  string `json:"user"`
		} `json:"owner"`
		Size uint64 `json:"size"`
		Time struct {
			Atime  int `json:"atime"`
			Crtime int `json:"crtime"`
			Ctime  int `json:"ctime"`
			Mtime  int `json:"mtime"`
		} `json:"time"`
	} `json:"additional"`
}

type FileStationServiceOp struct {
	c *Client
}

type listResponse struct {
	Data struct {
		Offset int `json:"offset"`
		// Note tis is really a union
		Shares []FileInfo `json:"shares"`
		Files  []FileInfo `json:"files"`
		Total  int        `json:"total"`
	} `json:"data"`
	Success bool `json:"success"`
}

func (s *FileStationServiceOp) ListShares() ([]FileInfo, error) {
	api := s.c.GetApi(SYNOFileStationList)

	params := map[string]string{
		"api":     SYNOFileStationList,
		"version": fmt.Sprintf("%d", api.MaxVersion),
		"method":  "list_share",
	}

	resp := &listResponse{}

	err := s.c.do("GET", api.Path, params, resp)
	return resp.Data.Shares, err
}

func (s *FileStationServiceOp) List(path string) ([]FileInfo, error) {
	api := s.c.GetApi(SYNOFileStationList)

	params := map[string]string{
		"api":         SYNOFileStationList,
		"version":     fmt.Sprintf("%d", api.MaxVersion),
		"method":      "list",
		"folder_path": path,
		"additional":  "[\"size\",\"owner\",\"time\",\"perm\"]",
	}

	resp := &listResponse{}

	err := s.c.do("GET", api.Path, params, resp)
	return resp.Data.Files, err
}

func (s *FileStationServiceOp) Stat(path string) ([]FileInfo, error) {
	api := s.c.GetApi(SYNOFileStationList)

	params := map[string]string{
		"api":        SYNOFileStationList,
		"version":    fmt.Sprintf("%d", api.MaxVersion),
		"method":     "getinfo",
		"path":       path,
		"additional": "[\"size\",\"owner\",\"time\",\"perm\"]",
	}

	//resp := &listResponse{}

	err := s.c.do("GET", api.Path, params, nil)
	return nil, err
}

func (s *FileStationServiceOp) Download(path string, w io.Writer) error {
	api := s.c.GetApi(SYNOFileStationDownload)

	params := map[string]string{
		"api":     api.Name,
		"version": fmt.Sprintf("%d", api.MaxVersion),
		"method":  "download",
		"path":    path,
		"mode":    "download",
	}

	err := s.c.download("GET", api.Path, params, w)
	return err
}
