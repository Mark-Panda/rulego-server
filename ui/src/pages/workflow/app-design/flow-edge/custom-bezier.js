import { BezierEdge, BezierEdgeModel } from '@logicflow/core';

class CustomEdge extends BezierEdge {}

class CustomEdgeModel extends BezierEdgeModel {
  constructor(data, graphModel) {
    super(data, graphModel);
    // 节流参数，用于控制updatePathByAnchor的调用频率
    this._updateTimer = null;
    this._lastUpdateTime = 0;
    this._updateThrottle = 16; // 约 60fps
  }
  // 设置边样式
  getEdgeStyle() {
    const style = super.getEdgeStyle();
    if (this.isSelected) {
      style.stroke = '#466dfd';
    } else {
      style.stroke = '#d1d5dd';
    }
    return style;
  }

  // 设置 hover 轮廓样式
  getOutlineStyle() {
    const style = super.getOutlineStyle();
    style.stroke = 'none';
    style.hover.stroke = 'none';
    return style;
  }

  /**
   * 计算两点之间的距离
   */
  calculateDistance(point1, point2) {
    const dx = point2.x - point1.x;
    const dy = point2.y - point1.y;
    return Math.sqrt(dx * dx + dy * dy);
  }

  /**
   * 判断两点是否接近水平或垂直直线
   */
  isNearStraightLine(point1, point2) {
    const dx = Math.abs(point2.x - point1.x);
    const dy = Math.abs(point2.y - point1.y);
    // 降低阈值，使直线检测更加严格
    const threshold = 10; // 允许的偏差阈值
    
    // 接近水平线
    if (dy < threshold) return true;
    // 接近垂直线  
    if (dx < threshold) return true;
    
    return false;
  }

  /**
   * 判断是否需要避免节点遮挡和交叉
   */
  needsAvoidObstacle(startPoint, endPoint) {
    const dx = endPoint.x - startPoint.x;
    const dy = endPoint.y - startPoint.y;
    
    // 如果目标节点在源节点的左上方或左下方，需要避免直线连接
    if (dx < 0) return true;
    
    // 如果水平距离小于垂直距离且水平距离较小，可能造成遮挡
    if (Math.abs(dx) < Math.abs(dy) && Math.abs(dx) < 80) return true;
    
    return false;
  }

  /**
   * 智能选择连线类型和offset值
   */
  calculateSmartOffset(startPoint, endPoint) {
    const distance = this.calculateDistance(startPoint, endPoint);
    const isNearStraight = this.isNearStraightLine(startPoint, endPoint);
    const needsAvoid = this.needsAvoidObstacle(startPoint, endPoint);
    
    // 距离阈值设置
    const shortDistance = 100;  // 短距离
    const mediumDistance = 200; // 中等距离
    
    if (distance < shortDistance) {
      // 短距离处理
      if (isNearStraight && !needsAvoid) {
        return 0.5; // 直线
      } else {
        return 15; // 平缓弧线
      }
    } else if (distance < mediumDistance) {
      // 中等距离处理
      if (needsAvoid) {
        return 40; // 中等S线
      } else {
        return 25; // 平缓弧线
      }
    } else {
      // 长距离处理
      return 60; // 标准S线，但不过于弯曲
    }
  }

  /**
   * 重写initEdgeData方法，在边初始化时调整参数
   */
  initEdgeData(data) {
    super.initEdgeData(data);
    
    // 使用智能算法计算offset
    this.offset = this.calculateSmartOffset(this.startPoint, this.endPoint);
  }

  /**
   * 给边自定义方案，使其支持基于锚点的位置更新边的路径
   */
  updatePathByAnchor() {
    // 节流控制，避免拖动过程中频繁更新
    const now = Date.now();
    if (now - this._lastUpdateTime < this._updateThrottle) {
      if (this._updateTimer) {
        clearTimeout(this._updateTimer);
      }
      this._updateTimer = setTimeout(() => {
        this._doUpdatePathByAnchor();
      }, this._updateThrottle);
      return;
    }
    
    this._lastUpdateTime = now;
    this._doUpdatePathByAnchor();
  }

  /**
   * 实际执行路径更新的方法
   */
  _doUpdatePathByAnchor() {
    const sourceNodeModel = this.graphModel.getNodeModelById(this.sourceNodeId);
    const sourceAnchor = sourceNodeModel
      ?.getDefaultAnchor()
      .find((anchor) => anchor.id === this.sourceAnchorId);
    const targetNodeModel = this.graphModel.getNodeModelById(this.targetNodeId);
    const targetAnchor = targetNodeModel
      ?.getDefaultAnchor()
      .find((anchor) => anchor.id === this.targetAnchorId);

    if (sourceAnchor) {
      const startPoint = {
        x: sourceAnchor?.x,
        y: sourceAnchor?.y,
      };
      this.updateStartPoint(startPoint);
    }
    if (targetAnchor) {
      const endPoint = {
        x: targetAnchor?.x,
        y: targetAnchor?.y,
      };
      this.updateEndPoint(endPoint);
    }
    
    // 使用智能算法重新计算offset
    this.offset = this.calculateSmartOffset(this.startPoint, this.endPoint);
    
    // 这里需要将原有的pointsList设置为空，才能触发重新计算控制点
    this.pointsList = [];
    this.initPoints();
  }
}

export default {
  type: 'custom-bezier',
  view: CustomEdge,
  model: CustomEdgeModel,
};
