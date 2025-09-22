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
    const threshold = 15; // 允许的偏差阈值
    
    // 接近水平线
    if (dy < threshold) return true;
    // 接近垂直线
    if (dx < threshold) return true;
    
    return false;
  }

  /**
   * 重写initEdgeData方法，在边初始化时调整参数
   */
  initEdgeData(data) {
    super.initEdgeData(data);
    
    // 计算距离并判断连线类型
    const distance = this.calculateDistance(this.startPoint, this.endPoint);
    const shortDistanceThreshold = 120;
    const isNearStraight = this.isNearStraightLine(this.startPoint, this.endPoint);
    
    if (distance < shortDistanceThreshold) {
      if (isNearStraight) {
        // 短距离且接近直线时，使用极小的偏移值保持直线
        this.offset = 2;
      } else {
        // 短距离但不是直线时，使用小偏移值创建平缓弧线
        this.offset = 20;
      }
    }
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
    
    // 重新计算距离并判断连线类型
    const distance = this.calculateDistance(this.startPoint, this.endPoint);
    const shortDistanceThreshold = 120;
    const isNearStraight = this.isNearStraightLine(this.startPoint, this.endPoint);
    
    if (distance < shortDistanceThreshold) {
      if (isNearStraight) {
        // 短距离且接近直线时，使用极小的偏移值保持直线
        this.offset = 2;
      } else {
        // 短距离但不是直线时，使用小偏移值创建平缓弧线
        this.offset = 20;
      }
    } else {
      // 长距离时使用默认偏移值
      this.offset = 100;
    }
    
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
