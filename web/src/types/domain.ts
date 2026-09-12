// Mirrors accessdesk's Go models (internal/model) and DTOs (internal/dto) —
// JSON field names must match the `json:"..."` tags on those Go structs.

export type RequestStatus = 'PENDING' | 'APPROVED' | 'REJECTED';

export interface Dataset {
  id: number;
  slug: string;
  name: string;
  description: string;
  created_at: string;
}

export interface AccessRequest {
  id: number;
  dataset_id: number;
  dataset: Dataset;
  requester: string;
  reason: string;
  status: RequestStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateRequestBody {
  dataset_id: number;
  reason: string;
}

export interface DecisionBody {
  decision: 'APPROVE' | 'REJECT';
  note?: string;
}

export interface ListRequestsResponse {
  requests: AccessRequest[];
  page: number;
  page_size: number;
  total: number;
}

export interface ListRequestsParams {
  status?: RequestStatus;
  page?: number;
  page_size?: number;
}
