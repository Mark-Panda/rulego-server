package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	pb "kratos-server/api/rulego/v1"
	"kratos-server/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// LogService 日志服务
type LogService struct {
	pb.UnimplementedLogServiceServer

	runLogUc *biz.RunLogUsecase
	log      *log.Helper
}

// NewLogService 创建日志服务
func NewLogService(runLogUc *biz.RunLogUsecase, logger log.Logger) *LogService {
	return &LogService{
		runLogUc: runLogUc,
		log:      log.NewHelper(logger),
	}
}

// GetDebugLogs 获取调试日志
func (s *LogService) GetDebugLogs(ctx context.Context, req *pb.GetDebugLogsRequest) (*pb.GetDebugLogsResponse, error) {
	// 这里应该从 RuleGo 的调试数据中获取实际的日志
	// 暂时返回模拟数据
	logs := []*pb.DebugLogEntry{}

	// 模拟生成一些调试日志
	for i := 0; i < 10; i++ {
		logs = append(logs, &pb.DebugLogEntry{
			ChainId:   req.ChainId,
			NodeId:    fmt.Sprintf("%s-node-%d", req.NodeId, i),
			Message:   fmt.Sprintf("Debug message %d for chain %s", i, req.ChainId),
			Timestamp: time.Now().Add(time.Duration(-i) * time.Minute).Unix(),
			Level:     []string{"DEBUG", "INFO", "WARN", "ERROR"}[i%4],
			Metadata: map[string]string{
				"source":   "node",
				"flowType": "input",
				"msgType":  "JSON",
			},
		})
	}

	// 分页处理
	total := len(logs)
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (int(req.Page) - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return &pb.GetDebugLogsResponse{
			Logs:     []*pb.DebugLogEntry{},
			Total:    int32(total),
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}

	end := offset + pageSize
	if end > total {
		end = total
	}

	return &pb.GetDebugLogsResponse{
		Logs:     logs[offset:end],
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ListRunLogs 获取运行日志列表
func (s *LogService) ListRunLogs(ctx context.Context, req *pb.ListRunLogsRequest) (*pb.ListRunLogsResponse, error) {
	runLogs, total, err := s.runLogUc.ListRunLogs(ctx, req.Username, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	pbLogs := make([]*pb.RunLogEntry, len(runLogs))
	for i, runLog := range runLogs {
		pbLogs[i] = &pb.RunLogEntry{
			Id:           int64(runLog.ID),
			ChainId:      runLog.ChainID,
			MsgType:      runLog.MsgType,
			Data:         runLog.Data,
			Result:       runLog.Result,
			StartTime:    runLog.StartTime.Unix(),
			EndTime:      runLog.EndTime.Unix(),
			Status:       runLog.Status,
			ErrorMessage: runLog.ErrorMessage,
			Username:     runLog.Username,
		}
	}

	return &pb.ListRunLogsResponse{
		Logs:     pbLogs,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteRunLogs 删除运行日志
func (s *LogService) DeleteRunLogs(ctx context.Context, req *pb.DeleteRunLogsRequest) (*pb.DeleteRunLogsResponse, error) {
	ids := make([]uint, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = uint(id)
	}

	deletedCount, err := s.runLogUc.DeleteRunLogs(ctx, ids, req.Username)
	if err != nil {
		return &pb.DeleteRunLogsResponse{
			Success:      false,
			Message:      err.Error(),
			DeletedCount: 0,
		}, nil
	}

	return &pb.DeleteRunLogsResponse{
		Success:      true,
		Message:      "Run logs deleted successfully",
		DeletedCount: int32(deletedCount),
	}, nil
}

// LocaleService 本地化服务
type LocaleService struct {
	pb.UnimplementedLocaleServiceServer

	log *log.Helper
}

// NewLocaleService 创建本地化服务
func NewLocaleService(logger log.Logger) *LocaleService {
	return &LocaleService{
		log: log.NewHelper(logger),
	}
}

// GetLocales 获取本地化配置
func (s *LocaleService) GetLocales(ctx context.Context, req *pb.GetLocalesRequest) (*pb.GetLocalesResponse, error) {
	// 从文件系统读取本地化配置
	localesDir := "locales"
	localeFile := filepath.Join(localesDir, req.Lang+".json")

	// 默认本地化内容
	defaultLocales := map[string]string{
		"app.title":            "RuleGo Server",
		"menu.dashboard":       "仪表盘",
		"menu.rules":           "规则链",
		"menu.components":      "组件管理",
		"menu.logs":            "日志管理",
		"menu.settings":        "系统设置",
		"button.save":          "保存",
		"button.cancel":        "取消",
		"button.delete":        "删除",
		"button.edit":          "编辑",
		"button.add":           "新增",
		"button.refresh":       "刷新",
		"button.search":        "搜索",
		"button.export":        "导出",
		"button.import":        "导入",
		"message.success":      "操作成功",
		"message.error":        "操作失败",
		"message.confirm":      "确认操作",
		"message.loading":      "加载中...",
		"form.required":        "该字段为必填项",
		"form.invalid":         "输入格式不正确",
		"table.actions":        "操作",
		"table.no_data":        "暂无数据",
		"pagination.total":     "共 {total} 条",
		"rule.name":            "规则链名称",
		"rule.description":     "规则链描述",
		"rule.status.active":   "已启用",
		"rule.status.inactive": "已禁用",
		"component.name":       "组件名称",
		"component.type":       "组件类型",
		"component.version":    "版本号",
		"log.level.debug":      "调试",
		"log.level.info":       "信息",
		"log.level.warn":       "警告",
		"log.level.error":      "错误",
	}

	// 根据语言返回不同的本地化内容
	if req.Lang == "en_us" {
		defaultLocales = map[string]string{
			"app.title":            "RuleGo Server",
			"menu.dashboard":       "Dashboard",
			"menu.rules":           "Rule Chains",
			"menu.components":      "Components",
			"menu.logs":            "Logs",
			"menu.settings":        "Settings",
			"button.save":          "Save",
			"button.cancel":        "Cancel",
			"button.delete":        "Delete",
			"button.edit":          "Edit",
			"button.add":           "Add",
			"button.refresh":       "Refresh",
			"button.search":        "Search",
			"button.export":        "Export",
			"button.import":        "Import",
			"message.success":      "Success",
			"message.error":        "Error",
			"message.confirm":      "Confirm",
			"message.loading":      "Loading...",
			"form.required":        "This field is required",
			"form.invalid":         "Invalid format",
			"table.actions":        "Actions",
			"table.no_data":        "No data",
			"pagination.total":     "Total {total} items",
			"rule.name":            "Rule Name",
			"rule.description":     "Description",
			"rule.status.active":   "Active",
			"rule.status.inactive": "Inactive",
			"component.name":       "Component Name",
			"component.type":       "Component Type",
			"component.version":    "Version",
			"log.level.debug":      "Debug",
			"log.level.info":       "Info",
			"log.level.warn":       "Warning",
			"log.level.error":      "Error",
		}
	}

	// 尝试从文件读取本地化配置
	if data, err := os.ReadFile(localeFile); err == nil {
		var fileLocales map[string]string
		if err := json.Unmarshal(data, &fileLocales); err == nil {
			// 合并文件中的配置和默认配置
			for k, v := range fileLocales {
				defaultLocales[k] = v
			}
		}
	}

	return &pb.GetLocalesResponse{
		Locales: defaultLocales,
	}, nil
}

// SaveLocales 保存本地化配置
func (s *LocaleService) SaveLocales(ctx context.Context, req *pb.SaveLocalesRequest) (*pb.SaveLocalesResponse, error) {
	// 保存本地化配置到文件系统
	localesDir := "locales"

	// 创建目录如果不存在
	if err := os.MkdirAll(localesDir, 0755); err != nil {
		return &pb.SaveLocalesResponse{
			Success: false,
			Message: "Failed to create locales directory: " + err.Error(),
		}, nil
	}

	localeFile := filepath.Join(localesDir, req.Lang+".json")

	// 将本地化数据序列化为JSON
	data, err := json.MarshalIndent(req.Locales, "", "  ")
	if err != nil {
		return &pb.SaveLocalesResponse{
			Success: false,
			Message: "Failed to marshal locales: " + err.Error(),
		}, nil
	}

	// 写入文件
	if err := os.WriteFile(localeFile, data, 0644); err != nil {
		return &pb.SaveLocalesResponse{
			Success: false,
			Message: "Failed to save locales file: " + err.Error(),
		}, nil
	}

	s.log.WithContext(ctx).Infof("Saved locales for language %s to %s", req.Lang, localeFile)

	return &pb.SaveLocalesResponse{
		Success: true,
		Message: "Locales saved successfully",
	}, nil
}
