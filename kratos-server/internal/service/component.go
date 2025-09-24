package service

import (
	"context"
	"time"

	pb "kratos-server/api/rulego/v1"
	"kratos-server/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// ComponentService 组件服务
type ComponentService struct {
	pb.UnimplementedComponentServiceServer

	uc  *biz.ComponentUsecase
	log *log.Helper
}

// NewComponentService 创建组件服务
func NewComponentService(uc *biz.ComponentUsecase, logger log.Logger) *ComponentService {
	return &ComponentService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// ListComponents 获取所有组件列表
func (s *ComponentService) ListComponents(ctx context.Context, req *pb.ListComponentsRequest) (*pb.ListComponentsResponse, error) {
	// 获取RuleGo的内置组件
	endpoints := []*pb.ComponentInfo{
		{
			Type:        "http",
			Name:        "HTTP Endpoint",
			Category:    "endpoint",
			Description: "HTTP接收端点",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
		{
			Type:        "mqtt",
			Name:        "MQTT Endpoint",
			Category:    "endpoint",
			Description: "MQTT接收端点",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
		{
			Type:        "websocket",
			Name:        "WebSocket Endpoint",
			Category:    "endpoint",
			Description: "WebSocket接收端点",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
	}

	nodes := []*pb.ComponentInfo{
		{
			Type:        "jsFilter",
			Name:        "JS Filter",
			Category:    "filter",
			Description: "JavaScript过滤器",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
		{
			Type:        "jsTransform",
			Name:        "JS Transform",
			Category:    "transform",
			Description: "JavaScript转换器",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
		{
			Type:        "log",
			Name:        "Log",
			Category:    "action",
			Description: "日志输出组件",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
		{
			Type:        "dbClient",
			Name:        "Database Client",
			Category:    "external",
			Description: "数据库客户端",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
		{
			Type:        "restApiCall",
			Name:        "REST API Call",
			Category:    "external",
			Description: "REST API调用",
			Version:     "1.0",
			Author:      "rulego",
			IsCustom:    false,
			IsInstalled: true,
		},
	}

	// 获取用户自定义组件
	customComponents, err := s.uc.ListComponents(ctx, req.Username)
	if err == nil {
		for _, comp := range customComponents {
			nodes = append(nodes, &pb.ComponentInfo{
				Type:        comp.Type,
				Name:        comp.Name,
				Category:    "custom",
				Description: comp.Description,
				Version:     comp.Version,
				Author:      comp.Author,
				IsCustom:    true,
				IsInstalled: comp.IsInstalled,
			})
		}
	}

	return &pb.ListComponentsResponse{
		Endpoints: endpoints,
		Nodes:     nodes,
		Builtins: map[string]string{
			"msgTypes":  "POST,GET,PUT,DELETE",
			"functions": "sin,cos,abs,sqrt,max,min",
		},
		SharedNodes: []*pb.SharedComponentInfo{
			{
				Id:           "shared-http",
				Type:         "http",
				Name:         "共享HTTP服务",
				Description:  "共享的HTTP接收端点",
				EndpointType: "http",
				Configuration: map[string]string{
					"port": "9090",
					"path": "/api/v1",
				},
			},
		},
	}, nil
}

// ListCustomComponents 获取自定义动态组件列表
func (s *ComponentService) ListCustomComponents(ctx context.Context, req *pb.ListCustomComponentsRequest) (*pb.ListCustomComponentsResponse, error) {
	components, err := s.uc.ListComponents(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	pbComponents := make([]*pb.CustomComponentDetail, len(components))
	for i, comp := range components {
		pbComponents[i] = &pb.CustomComponentDetail{
			Id:          comp.ID,
			Type:        comp.Type,
			Name:        comp.Name,
			Description: comp.Description,
			Version:     comp.Version,
			Author:      comp.Author,
			Dsl:         comp.DSL,
			CreatedTime: comp.CreatedAt.Unix(),
			UpdatedTime: comp.UpdatedAt.Unix(),
			IsInstalled: comp.IsInstalled,
		}
	}

	return &pb.ListCustomComponentsResponse{
		Components: pbComponents,
	}, nil
}

// GetCustomComponentDSL 获取自定义动态组件DSL
func (s *ComponentService) GetCustomComponentDSL(ctx context.Context, req *pb.GetCustomComponentDSLRequest) (*pb.GetCustomComponentDSLResponse, error) {
	component, err := s.uc.GetComponent(ctx, req.Id, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.GetCustomComponentDSLResponse{
		Component: &pb.CustomComponentDetail{
			Id:          component.ID,
			Type:        component.Type,
			Name:        component.Name,
			Description: component.Description,
			Version:     component.Version,
			Author:      component.Author,
			Dsl:         component.DSL,
			CreatedTime: component.CreatedAt.Unix(),
			UpdatedTime: component.UpdatedAt.Unix(),
			IsInstalled: component.IsInstalled,
		},
	}, nil
}

// UpgradeCustomComponent 安装/升级自定义动态组件
func (s *ComponentService) UpgradeCustomComponent(ctx context.Context, req *pb.UpgradeCustomComponentRequest) (*pb.UpgradeCustomComponentResponse, error) {
	component := &biz.Component{
		ID:          req.Id,
		DSL:         req.Dsl,
		Username:    req.Username,
		IsInstalled: true,
	}

	_, err := s.uc.UpdateComponent(ctx, component)
	if err != nil {
		return &pb.UpgradeCustomComponentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.UpgradeCustomComponentResponse{
		Success: true,
		Message: "Component upgraded successfully",
	}, nil
}

// UninstallCustomComponent 卸载自定义动态组件
func (s *ComponentService) UninstallCustomComponent(ctx context.Context, req *pb.UninstallCustomComponentRequest) (*pb.UninstallCustomComponentResponse, error) {
	err := s.uc.DeleteComponent(ctx, req.Id, req.Username)
	if err != nil {
		return &pb.UninstallCustomComponentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.UninstallCustomComponentResponse{
		Success: true,
		Message: "Component uninstalled successfully",
	}, nil
}

// ListSharedComponents 获取所有共享组件
func (s *ComponentService) ListSharedComponents(ctx context.Context, req *pb.ListSharedComponentsRequest) (*pb.ListSharedComponentsResponse, error) {
	// 返回系统中的共享组件列表
	sharedComponents := []*pb.SharedComponentInfo{
		{
			Id:           "shared-http",
			Type:         "http",
			Name:         "共享HTTP服务",
			Description:  "共享的HTTP接收端点",
			EndpointType: "http",
			Configuration: map[string]string{
				"port": "9090",
				"path": "/api/v1",
			},
		},
		{
			Id:           "shared-mqtt",
			Type:         "mqtt",
			Name:         "共享MQTT服务",
			Description:  "共享的MQTT接收端点",
			EndpointType: "mqtt",
			Configuration: map[string]string{
				"server": "tcp://localhost:1883",
				"topic":  "rulego/+/msg",
			},
		},
		{
			Id:           "shared-websocket",
			Type:         "websocket",
			Name:         "共享WebSocket服务",
			Description:  "共享的WebSocket接收端点",
			EndpointType: "websocket",
			Configuration: map[string]string{
				"port": "9092",
				"path": "/ws",
			},
		},
	}

	return &pb.ListSharedComponentsResponse{
		Components: sharedComponents,
	}, nil
}

// ListMarketplaceComponents 获取组件市场组件列表
func (s *ComponentService) ListMarketplaceComponents(ctx context.Context, req *pb.ListMarketplaceComponentsRequest) (*pb.ListMarketplaceComponentsResponse, error) {
	// 返回市场中可用的组件列表
	marketplaceComponents := []*pb.CustomComponentDetail{
		{
			Id:          "market-js-advanced",
			Type:        "jsAdvancedFilter",
			Name:        "高级JS过滤器",
			Description: "带有更多功能的JavaScript过滤器",
			Version:     "2.0",
			Author:      "community",
			Dsl:         `{"type":"jsAdvancedFilter","name":"高级JS过滤器","configuration":{"jsScript":"return msg.temperature > 25;"}}`,
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
			IsInstalled: false,
		},
		{
			Id:          "market-email-sender",
			Type:        "emailSender",
			Name:        "邮件发送器",
			Description: "发送邮件通知的组件",
			Version:     "1.0",
			Author:      "community",
			Dsl:         `{"type":"emailSender","name":"邮件发送器","configuration":{"smtpServer":"smtp.gmail.com","port":587}}`,
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
			IsInstalled: false,
		},
		{
			Id:          "market-mqtt-publisher",
			Type:        "mqttPublisher",
			Name:        "MQTT发布器",
			Description: "向MQTT主题发布消息的组件",
			Version:     "1.5",
			Author:      "community",
			Dsl:         `{"type":"mqttPublisher","name":"MQTT发布器","configuration":{"server":"tcp://localhost:1883","topic":"sensors/data"}}`,
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
			IsInstalled: false,
		},
		{
			Id:          "market-data-aggregator",
			Type:        "dataAggregator",
			Name:        "数据聚合器",
			Description: "聚合多个数据源的组件",
			Version:     "1.2",
			Author:      "community",
			Dsl:         `{"type":"dataAggregator","name":"数据聚合器","configuration":{"windowSize":10,"aggregateType":"avg"}}`,
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
			IsInstalled: false,
		},
	}

	// 如果用户要求查看自己的组件，可以加入用户自己发布的组件
	if req.CheckMy {
		userComponents, err := s.uc.ListComponents(ctx, req.Username)
		if err == nil {
			for _, comp := range userComponents {
				// 添加用户自己的组件到市场列表
				marketplaceComponents = append(marketplaceComponents, &pb.CustomComponentDetail{
					Id:          comp.ID,
					Type:        comp.Type,
					Name:        comp.Name,
					Description: comp.Description,
					Version:     comp.Version,
					Author:      comp.Author,
					Dsl:         comp.DSL,
					CreatedTime: comp.CreatedAt.Unix(),
					UpdatedTime: comp.UpdatedAt.Unix(),
					IsInstalled: comp.IsInstalled,
				})
			}
		}
	}

	return &pb.ListMarketplaceComponentsResponse{
		Components: marketplaceComponents,
	}, nil
}
