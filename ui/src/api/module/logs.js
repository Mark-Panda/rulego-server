import { request } from '@src/utils/request';

export function getLogsDebug(params) {
  return request.get('/logs/debug', {
    params,
  });
}

export function getRuntimeLogs(params) {
  return request.get('/logs/runs', {
    params,
  });
}

