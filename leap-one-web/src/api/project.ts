/**
 * 项目服务 API
 *
 * 提供项目 CRUD、成员管理、迭代管理等接口
 */

import { apiClient } from './client';
import type {
  Project,
  ProjectListParams,
  ProjectMember,
  Iteration,
} from '@/types/project';

const BASE = '/project';

// ── 项目查询 ─────────────────────────────────────────────────

/** 获取项目列表（分页） */
export const getProjectListApi = (params?: ProjectListParams) =>
  apiClient.getPage<Project>(`${BASE}/list`, params as Record<string, unknown>);

/** 获取项目详情 */
export const getProjectDetailApi = (id: number): Promise<Project> =>
  apiClient.get<Project>(`${BASE}/${id}`).then((res) => res.data);

// ── 项目写入 ─────────────────────────────────────────────────

/** 创建项目 */
export const createProjectApi = (data: Partial<Project>): Promise<Project> =>
  apiClient.post<Project>(BASE, data).then((res) => res.data);

/** 更新项目 */
export const updateProjectApi = (id: number, data: Partial<Project>): Promise<Project> =>
  apiClient.put<Project>(`${BASE}/${id}`, data).then((res) => res.data);

/** 删除项目 */
export const deleteProjectApi = (id: number): Promise<void> =>
  apiClient.delete(`${BASE}/${id}`).then(() => undefined);

// ── 成员管理 ─────────────────────────────────────────────────

/** 获取项目成员列表 */
export const getProjectMembersApi = (projectId: number): Promise<ProjectMember[]> =>
  apiClient
    .get<ProjectMember[]>(`${BASE}/${projectId}/members`)
    .then((res) => res.data);

/** 更新项目成员 */
export const updateProjectMembersApi = (
  projectId: number,
  members: Omit<ProjectMember, 'joinedAt'>[]
): Promise<void> =>
  apiClient.put(`${BASE}/${projectId}/members`, { members }).then(() => undefined);

// ── 迭代管理 ─────────────────────────────────────────────────

/** 获取迭代列表 */
export const getIterationListApi = (projectId: number) =>
  apiClient.getPage<Iteration>(`${BASE}/${projectId}/iteration`);

/** 创建迭代 */
export const createIterationApi = (
  projectId: number,
  data: Omit<Iteration, 'id' | 'projectId' | 'taskCount' | 'completedTaskCount' | 'progress'>
): Promise<Iteration> =>
  apiClient.post<Iteration>(`${BASE}/${projectId}/iteration`, data).then((res) => res.data);
