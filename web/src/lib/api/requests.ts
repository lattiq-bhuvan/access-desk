import { api } from '@/lib/api/client';
import type {
  AccessRequest,
  CreateRequestBody,
  DecisionBody,
  ListRequestsParams,
  ListRequestsResponse,
} from '@/types/domain';

export async function createRequest(body: CreateRequestBody): Promise<AccessRequest> {
  return api.post<AccessRequest>('/v1/requests', body);
}

export async function listRequests(params: ListRequestsParams = {}): Promise<ListRequestsResponse> {
  const query = new URLSearchParams();
  if (params.status) query.set('status', params.status);
  if (params.page) query.set('page', String(params.page));
  if (params.page_size) query.set('page_size', String(params.page_size));

  const qs = query.toString();
  return api.get<ListRequestsResponse>(`/v1/requests${qs ? `?${qs}` : ''}`);
}

export async function decideRequest(id: number, body: DecisionBody): Promise<AccessRequest> {
  return api.post<AccessRequest>(`/v1/requests/${id}/decision`, body);
}
