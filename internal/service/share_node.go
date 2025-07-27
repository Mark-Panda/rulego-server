// Package service provides the business logic for the share node service.
package service

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego-server/config/logger"
	"github.com/rulego/rulego-server/internal/constants"
	"github.com/rulego/rulego-server/internal/dao"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/fs"
	"github.com/rulego/rulego/utils/json"
)

// ShareNodeService represents the service for managing share nodes.
type ShareNodeService struct {
	// userSettingDao is the DAO for managing user settings.
	userSettingDao *dao.UserSettingDao
	// ShareNodeDao is the DAO for managing share nodes.
	shareNodeDao *dao.ShareNodeDao
	// config is the configuration of the server.
	config config.Config
	// username is the username of the user who owns the share nodes.
	username string
}

// NewShareNodeService creates a new ShareNodeService.
func NewShareNodeService(c config.Config, username string) (*ShareNodeService, error) {
	shareNodeService := &ShareNodeService{
		config:   c,
		username: username,
	}

	userSettingDao, err := dao.NewUserSettingDao(c, fmt.Sprintf(
		"%s/%s/%s", c.DataDir, constants.DirWorkflows, username))
	if err != nil {
		return &ShareNodeService{}, err
	}

	shareNodeService.userSettingDao = userSettingDao
	if err = shareNodeService.Init(); err != nil {
		return &ShareNodeService{}, err
	}

	shareNodeDao, err := dao.NewShareNodeDao(c, username)
	if err != nil {
		return nil, err
	}
	shareNodeService.shareNodeDao = shareNodeDao

	return shareNodeService, nil
}

// Init initializes all share nodes from the data directory.
func (s *ShareNodeService) Init() error {
	dirs := constants.GetShareNodesDir()
	for idx := range dirs {
		oneShareNode := dirs[idx]
		folderPath := path.Join(s.config.DataDir,
			constants.DirWorkflows, s.username, constants.DirWorkflowsShareNodes, oneShareNode)
		if err := fs.CreateDirs(folderPath); err != nil {
			logger.Logger.Printf("failed to create share node directory: %s\n", err)
			return err
		}
	}

	return nil
}

// Load loads user's all share nodes from the data directory.
func (s *ShareNodeService) Load(ruleConfig *types.Config) error {
	dirs := constants.GetShareNodesDir()
	for idx := range dirs {
		oneShareNode := dirs[idx]
		folderPath := path.Join(s.config.DataDir,
			constants.DirWorkflows, s.username, constants.DirWorkflowsShareNodes, oneShareNode)
		entries, err := os.ReadDir(folderPath)
		if err != nil {
			return err
		}

		var names []string
		for index := range entries {
			entry := entries[index]
			if entry.IsDir() {
				continue
			}

			fileName := entry.Name()
			// 或使用 filepath.Ext(fileName) == ".json"
			if strings.HasSuffix(fileName, ".json") {
				names = append(names, fileName)
			}
		}

		if len(names) == 0 {
			continue
		}

		for index := range names {
			filePath := path.Join(folderPath, names[index])
			content := fs.LoadFile(filePath)
			if content == nil {
				continue
			}

			switch oneShareNode {
			case constants.TypeShareNode:
				shareNode := types.RuleNode{}
				if err = json.Unmarshal(content, &shareNode); err != nil {
					return err
				}
				if _, err = ruleConfig.NodePool.NewFromRuleNode(shareNode); err != nil {
					return err
				}
			case constants.TypeShareEndpoint:
				endpoint := types.EndpointDsl{}
				if err = json.Unmarshal(content, &endpoint); err != nil {
					return err
				}

				if _, err = ruleConfig.NodePool.NewFromEndpoint(endpoint); err != nil {
					return err
				}
			default:
			}
		}
	}
	return nil
}

// Upsert updates or inserts a share node with the given node type and node info.
func (s *ShareNodeService) Upsert(nodeType string, nodeInfo []byte) error {
	if nodeType == constants.TypeShareNode {
		return s.shareNodeDao.UpsertNode(nodeInfo)
	}
	return s.shareNodeDao.UpsertEndpoint(nodeInfo)
}

// GetByID returns the share node with the given ID and node type.
func (s *ShareNodeService) GetByID(nodeID string, nodeType string) ([]byte, error) {
	return s.shareNodeDao.GetByID(nodeID, nodeType)
}

// Del deletes the share node with the given ID and node type.
func (s *ShareNodeService) Del(nodeID string, nodeType string) error {
	if nodeType == constants.TypeShareNode {
		return s.shareNodeDao.DelNode(nodeID)
	}
	return s.shareNodeDao.DelEndpoint(nodeID)
}

// ListNode returns a list of share nodes with the given keywords and pagination information.
func (s *ShareNodeService) ListNode(keywords string, page int, size int) ([]types.RuleNode, int, error) {
	return s.shareNodeDao.ListNode(keywords, page, size, constants.TypeShareNode)
}

// ListEndpoint returns a list of share nodes with the given keywords and pagination information.
func (s *ShareNodeService) ListEndpoint(keywords string, page int, size int) ([]types.EndpointDsl, int, error) {
	return s.shareNodeDao.ListEndpoint(keywords, page, size, constants.TypeShareEndpoint)
}
