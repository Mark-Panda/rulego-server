import { request } from '@src/utils/request';

// 获取业务文档列表
export function getDocs(params = {}) {
  const query = {
    page: params.page || 1,
    size: params.size || 10,
    keyword: params.keywords || params.keyword || undefined,
  };
  return request
    .get('/doc/list', { params: query })
    .then((res) => {
      // 统一返回结构为 { items, total } 以兼容业务文档页
      const rawItems = Array.isArray(res?.items)
        ? res.items
        : Array.isArray(res?.list)
        ? res.list
        : [];
      const total = res?.total ?? (Array.isArray(rawItems) ? rawItems.length : 0);
      const items = rawItems.map((it) => ({
        id: it.id ?? it.docId ?? it._id ?? undefined,
        name: it.name ?? it.title ?? '-',
        description: it.description ?? it.desc ?? '-',
        content: it.content ?? '',
        chainId: it.chainId ?? it.ruleChainId ?? it.chain_id ?? undefined,
        chainName: it.chainName ?? it.ruleChainName ?? it.chain_name ?? '-',
        createTime: it.createTime ?? it.createdAt ?? '-',
        updateTime: it.updateTime ?? it.updatedAt ?? '-',
        _original: it,
      }));
      return { items, total };
    });
}

// 获取业务文档详情
export function getDoc(id) {
  return request({
    url: `/api/docs/${id}`,
    method: 'get'
  });
}

// 创建业务文档
export function createDoc(data) {
  // 适配后端接口字段：title、content、desc
  const payload = {
    title: data?.name ?? data?.title,
    content: data?.content,
    desc: data?.description ?? data?.desc,
  };
  return request.post('/doc/create', payload);
}

// 更新业务文档
export function updateDoc(id, data) {
  const payload = {
    id: id ?? data?.id,
    title: data?.name ?? data?.title,
    content: data?.content,
    desc: data?.description ?? data?.desc,
  };
  return request.post('/doc/edit', payload);
}

// 删除业务文档
export function deleteDoc(id) {
  return request({
    url: `/api/docs/${id}`,
    method: 'delete'
  });
}