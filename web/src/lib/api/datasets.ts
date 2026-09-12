import { api } from '@/lib/api/client';
import type { Dataset } from '@/types/domain';

interface ListDatasetsResponse {
  datasets: Dataset[];
}

export async function listDatasets(): Promise<Dataset[]> {
  const res = await api.get<ListDatasetsResponse>('/v1/datasets');
  return res.datasets;
}
