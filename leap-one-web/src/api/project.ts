/** 项目服务API */

import { apiClient } from './client';
import type {
  Project,
  ProjectListParams,
  ProjectMember,
  Iteration,
} from '@/types/project';

const BASE_URL = '/project';

/** 获取项目列表 */
export function getProjectListApi(params?: ProjectListParams) {
  return apiClient.getPage<Project>(`${BASE_URL}/list`, params as Record<string, unknown>);
}

/** 获取项目详情 */
export function getProjectDetailApi(id: number): Promise<Project> {
  return apiClient.get<Project>(`${BASE_URL}/${id}`).then((res) => res.data);
}

/** 创建项目 */
export function createProjectApi(data: Partial<Project>): Promise<Project> {
  return apiClient.post<Project>(BASE_URL, data).then((res) => res.data);
}

/** 更新项目 */
export function updateProjectApi(id: number, data: Partial<Project>): Promise<Project> {
  return apiClient.put<Project>(`${BASE_URL}/${id}`, data).then((res) => res.data);
}

/** 删除项目 */
export function deleteProjectApi(id: number): Promise<void> {
  return apiClient.delete(`${BASE_URL}/${id}`).then(() => undefined);
}

/** 获取项目成员列表 */
export function getProjectMembersApi(projectId: number): Promise<ProjectMember[]> {
  return apiClient
    .get<ProjectMember[]>(`${BASE_URL}/${projectId}/members`)
    .then((res) => res.data);
}

/** 更新项目成员 */
export function updateProjectMembersApi(
  projectId: number,
  members: Omit<ProjectMember, 'joinedAt'>[]
): Promise<void> {
  return apiClient.put(`${BASE_URL}/${projectId}/members`, { members }).then(() => undefined);
}

/** 获取迭代列表 */
export function getIterationListApi(projectId: number) {
  return apiClient.getPage<Iteration>(`${BASE_URL}/${projectId}/iteration`);
}

/** 创建迭代 */
export function createIterationApi(
  projectId: number,
  data: Omit<Iteration, 'id' | 'projectId' | 'taskCount' | 'completedTaskCount' | 'progress'>
): Promise<Iteration> {
  return apiClient.post<Iteration>(`${BASE_URL}/${projectId}/iteration`, data).then((res) => res.data);
}
