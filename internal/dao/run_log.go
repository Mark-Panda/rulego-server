package dao

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego-server/config/logger"
	"github.com/rulego/rulego-server/internal/constants"
	"github.com/rulego/rulego-server/internal/model"
	"github.com/rulego/rulego-server/internal/utils/file"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/fs"
	"github.com/rulego/rulego/utils/json"
)

type EventDao struct {
	*FileStorage
	config config.Config
}

func NewEventDao(config config.Config) (*EventDao, error) {
	return &EventDao{
		config: config,
	}, nil
}

// SaveRunLog 保存工作流运行日志快照
func (s *EventDao) SaveRunLog(username string, ctx types.RuleContext, snapshot types.RuleChainRunSnapshot) error {
	var paths = []string{s.config.DataDir, constants.DirWorkflows}
	chainId := ctx.RuleChain().GetNodeId().Id
	paths = append(paths, username, constants.DirWorkflowsRun, chainId)
	pathStr := path.Join(paths...)
	//创建文件夹
	_ = fs.CreateDirs(pathStr)
	snapshot.Id = time.Now().Format("20060102150405000") + "_" + snapshot.Id
	//保存到文件
	if byteV, err := json.Marshal(snapshot); err != nil {
		logger.Logger.Printf("dao/EventDao:SaveRunLog marshal error %v", err)
		return err
	} else {
		//v, _ := json.Format(byteV)
		//保存规则链到文件
		if err = fs.SaveFile(filepath.Join(pathStr, snapshot.Id), byteV); err != nil {
			logger.Logger.Printf("dao/EventDao:SaveRunLog save file error %v", err)
			return err
		}
	}
	return nil
}

func (s *EventDao) SaveRunLogToDataBase(username string, ctx types.RuleContext, snapshot types.RuleChainRunSnapshot) error {
	nodelogs, _ := json.Marshal(snapshot.Logs)
	additionalInfo, _ := json.Marshal(snapshot.AdditionalInfo)
	ruleChainInfo, _ := json.Marshal(snapshot.RuleChain)
	metadata, _ := json.Marshal(snapshot.Metadata)
	t := time.Now()
	runLog := model.RunLog{
		RunId:          snapshot.Id,
		ChainId:        snapshot.RuleChain.RuleChain.ID,
		ChainName:      snapshot.RuleChain.RuleChain.Name,
		NodeLog:        string(nodelogs),
		AdditionalInfo: string(additionalInfo),
		RuleChainInfo:  string(ruleChainInfo),
		Metadata:       string(metadata),
		StartTs:        snapshot.StartTs,
		EndTs:          snapshot.EndTs,
		CreatedAt:      &t,
		UpdatedAt:      &t,
	}
	return model.DBClient.Client.Create(&runLog).Error
}

func (s *EventDao) Delete(username string, chainId, id string) error {
	var paths = []string{s.config.DataDir, constants.DirWorkflows}
	if id != "" {
		paths = append(paths, username, constants.DirWorkflowsRun, chainId, id)
	} else {
		paths = append(paths, username, constants.DirWorkflowsRun, chainId)
	}

	pathStr := path.Join(paths...)
	return os.RemoveAll(pathStr)
}

func (s *EventDao) DeleteDataBaseByRunId(username string, chainId, id string) error {
	return model.DBClient.Client.Where("run_id = ?", id).Delete(&model.RunLog{}).Error
}

func (s *EventDao) DeleteByChainId(username string, chainId string) error {
	var paths = []string{s.config.DataDir, constants.DirWorkflows}
	paths = append(paths, username, constants.DirWorkflowsRun, chainId)
	pathStr := path.Join(paths...)
	return os.RemoveAll(pathStr)
}

func (s *EventDao) DeleteDataBaseByChainId(username string, chainId string) error {
	return model.DBClient.Client.Where("chain_id = ?", chainId).Delete(&model.RunLog{}).Error
}

func (s *EventDao) visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

func (s *EventDao) List(username string, chainId string, current, size int) ([]types.RuleChainRunSnapshot, int, error) {
	var snapshots []types.RuleChainRunSnapshot

	var paths = []string{s.config.DataDir, constants.DirWorkflows}
	if chainId == "" {
		//加载所有
		paths = append(paths, username, constants.DirWorkflowsRun)
	} else {
		paths = append(paths, username, constants.DirWorkflowsRun, chainId)
	}
	pathStr := path.Join(paths...)
	// 获取目录下所有运行日志文件
	var files []string
	if err := filepath.Walk(pathStr, s.visit(&files)); err != nil {
		return snapshots, 0, nil
	}
	// 按文件时间戳排序
	fileWithTimestamps := file.SortFilesByTimestamp(files)

	// 计算分页的起始索引
	start := (current - 1) * size
	end := start + size
	if end > len(files) {
		end = len(files)
	}

	// 遍历文件，每个文件对应一条 RuleChainRunSnapshot 记录
	for _, file := range fileWithTimestamps[start:end] {
		data, err := os.ReadFile(file.Path)
		if err != nil {
			return nil, 0, err
		}

		var snapshot types.RuleChainRunSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return nil, 0, err
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, len(files), nil
}

func (s *EventDao) ListByDataBase(username, chainId string, current, size int, startTime, endTime string) ([]types.RuleChainRunSnapshot, int, error) {
	var snapshots []types.RuleChainRunSnapshot
	var runLogs []model.RunLog
	var total int64
	where := ""
	if startTime != "" {
		where = fmt.Sprintf("created_at >= '%s'", startTime)
	}
	if endTime != "" {
		if where != "" {
			where = fmt.Sprintf("%s and created_at <= '%s'", where, endTime)
		} else {
			where = fmt.Sprintf("created_at <= '%s'", endTime)
		}
	}
	if chainId == "" {
		if where != "" {
			if err := model.DBClient.Client.Where(where).Order("created_at desc").Offset((current - 1) * size).Limit(size).Find(&runLogs).Error; err != nil {
				return snapshots, 0, err
			}
			if err := model.DBClient.Client.Model(&model.RunLog{}).Where(where).Count(&total).Error; err != nil {
				return snapshots, 0, err
			}
		} else {
			if err := model.DBClient.Client.Order("created_at desc").Offset((current - 1) * size).Limit(size).Find(&runLogs).Error; err != nil {
				return snapshots, 0, err
			}
			if err := model.DBClient.Client.Model(&model.RunLog{}).Count(&total).Error; err != nil {
				return snapshots, 0, err
			}
		}
	} else {
		if where != "" {
			where = fmt.Sprintf("%s and chain_id = '%s'", where, chainId)
			if err := model.DBClient.Client.Where(where).Order("created_at desc").Offset((current - 1) * size).Limit(size).Find(&runLogs).Error; err != nil {
				return snapshots, 0, err
			}
			if err := model.DBClient.Client.Model(&model.RunLog{}).Where(where).Count(&total).Error; err != nil {
				return snapshots, 0, err
			}
		} else {
			where = fmt.Sprintf("chain_id = '%s'", chainId)
			if err := model.DBClient.Client.Where(where).Order("created_at desc").Offset((current - 1) * size).Limit(size).Find(&runLogs).Error; err != nil {
				return snapshots, 0, err
			}
			if err := model.DBClient.Client.Model(&model.RunLog{}).Where(where).Count(&total).Error; err != nil {
				return snapshots, 0, err
			}
		}
		if err := model.DBClient.Client.Model(&model.RunLog{}).Where("chain_id = ?", chainId).Count(&total).Error; err != nil {
			return snapshots, 0, err
		}
	}
	for _, runLog := range runLogs {
		snapshot, err := s.RunLogToRuleChainRunSnapshot(runLog)
		if err != nil {
			return snapshots, 0, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, int(total), nil
}

func (s *EventDao) Get(username, chainId, snapshotId string) (types.RuleChainRunSnapshot, error) {
	var snapshot types.RuleChainRunSnapshot
	var paths = []string{s.config.DataDir, constants.DirWorkflows}
	paths = append(paths, username, constants.DirWorkflowsRun, chainId, snapshotId)
	file := path.Join(paths...)
	data, err := os.ReadFile(file)
	if err != nil {
		return snapshot, err
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return snapshot, err
	} else {
		return snapshot, nil
	}
}

func (s *EventDao) GetByDataBase(username, chainId, snapshotId string) (types.RuleChainRunSnapshot, error) {
	var snapshot types.RuleChainRunSnapshot
	var runLog model.RunLog
	if err := model.DBClient.Client.Where("run_id = ?", snapshotId).First(&runLog).Error; err != nil {
		return snapshot, err
	}
	snapshot, err := s.RunLogToRuleChainRunSnapshot(runLog)
	if err != nil {
		return snapshot, err
	}
	return snapshot, nil
}

func (s *EventDao) RunLogToRuleChainRunSnapshot(runLog model.RunLog) (types.RuleChainRunSnapshot, error) {
	var snapshot types.RuleChainRunSnapshot
	if err := json.Unmarshal([]byte(runLog.NodeLog), &snapshot.Logs); err != nil {
		return snapshot, err
	}
	if err := json.Unmarshal([]byte(runLog.AdditionalInfo), &snapshot.AdditionalInfo); err != nil {
		return snapshot, err
	}
	if err := json.Unmarshal([]byte(runLog.RuleChainInfo), &snapshot.RuleChain); err != nil {
		return snapshot, err
	}
	if err := json.Unmarshal([]byte(runLog.Metadata), &snapshot.Metadata); err != nil {
		return snapshot, err
	}
	snapshot.StartTs = runLog.StartTs
	snapshot.EndTs = runLog.EndTs
	snapshot.StartTs = runLog.StartTs
	snapshot.EndTs = runLog.EndTs
	return snapshot, nil
}
