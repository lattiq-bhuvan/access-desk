import {
  Badge,
  Button,
  Card,
  CardContent,
  DataTable,
  MetricCard,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  toast,
} from '@lattiq/design-system';
import type { ColumnDef } from '@tanstack/react-table';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { ConfirmDialog } from '@/components/ConfirmDialog';
import { getErrorMessage } from '@/lib/api/errors';
import { decideRequest, listRequests } from '@/lib/api/requests';
import { useAuthStore } from '@/stores/useAuthStore';
import type { AccessRequest, RequestStatus } from '@/types/domain';

const PAGE_SIZE = 10;

const STATUS_BADGE: Record<RequestStatus, 'secondary' | 'default' | 'destructive'> = {
  PENDING: 'secondary',
  APPROVED: 'default',
  REJECTED: 'destructive',
};

type StatusFilter = RequestStatus | 'ALL';

interface PendingDecision {
  request: AccessRequest;
  decision: 'APPROVE' | 'REJECT';
}

export function Requests() {
  const auth = useAuthStore(state => state.auth);
  const isApprover = !!auth?.user?.roles?.includes('approver');

  const [rows, setRows] = useState<AccessRequest[]>([]);
  const [total, setTotal] = useState(0);
  const [pageIndex, setPageIndex] = useState(0);
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('ALL');
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [pending, setPending] = useState<PendingDecision | null>(null);
  const [deciding, setDeciding] = useState(false);

  // Stretch: pending/approved/rejected counts for the metric-card row.
  // The backend has no aggregate endpoint, so read each status's `total`
  // from a page_size=1 request — cheap, and only fetched for approvers.
  const [counts, setCounts] = useState<Record<RequestStatus, number> | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setLoadError(null);
    try {
      const res = await listRequests({
        status: statusFilter === 'ALL' ? undefined : statusFilter,
        page: pageIndex + 1,
        page_size: PAGE_SIZE,
      });
      setRows(res.requests);
      setTotal(res.total);
    } catch (err) {
      setLoadError(getErrorMessage(err, 'Failed to load requests'));
    } finally {
      setLoading(false);
    }
  }, [statusFilter, pageIndex]);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    if (!isApprover) return;
    let cancelled = false;
    Promise.all(
      (['PENDING', 'APPROVED', 'REJECTED'] as RequestStatus[]).map(status =>
        listRequests({ status, page: 1, page_size: 1 }).then(res => [status, res.total] as const)
      )
    ).then(entries => {
      if (cancelled) return;
      setCounts(Object.fromEntries(entries) as Record<RequestStatus, number>);
    });
    return () => {
      cancelled = true;
    };
  }, [isApprover, rows]);

  async function confirmDecision() {
    if (!pending) return;
    setDeciding(true);
    try {
      await decideRequest(pending.request.id, { decision: pending.decision });
      toast.success(
        pending.decision === 'APPROVE' ? 'Request approved' : 'Request rejected'
      );
      setPending(null);
      await load();
    } catch (err) {
      toast.error(getErrorMessage(err, 'Could not record the decision'));
    } finally {
      setDeciding(false);
    }
  }

  const columns = useMemo<ColumnDef<AccessRequest, unknown>[]>(() => {
    const base: ColumnDef<AccessRequest, unknown>[] = [
      {
        accessorKey: 'dataset',
        header: 'Dataset',
        cell: ({ row }) => row.original.dataset?.name ?? `#${row.original.dataset_id}`,
      },
      { accessorKey: 'requester', header: 'Requester' },
      { accessorKey: 'reason', header: 'Reason' },
      {
        accessorKey: 'status',
        header: 'Status',
        cell: ({ row }) => (
          <Badge variant={STATUS_BADGE[row.original.status]}>{row.original.status}</Badge>
        ),
      },
      {
        accessorKey: 'created_at',
        header: 'Requested',
        cell: ({ row }) => new Date(row.original.created_at).toLocaleString(),
      },
    ];

    if (isApprover) {
      base.push({
        id: 'actions',
        header: 'Actions',
        cell: ({ row }) => {
          const req = row.original;
          if (req.status !== 'PENDING') {
            return <span className="text-muted-foreground text-sm">—</span>;
          }
          return (
            <div className="flex gap-2">
              <Button
                size="sm"
                onClick={() => setPending({ request: req, decision: 'APPROVE' })}
              >
                Approve
              </Button>
              <Button
                size="sm"
                variant="destructive"
                onClick={() => setPending({ request: req, decision: 'REJECT' })}
              >
                Reject
              </Button>
            </div>
          );
        },
      });
    }

    return base;
  }, [isApprover]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Requests</h1>
          <p className="text-muted-foreground">
            {isApprover
              ? 'All access requests. Approve or reject anything pending.'
              : 'Access requests you have submitted.'}
          </p>
        </div>

        <Select
          value={statusFilter}
          onValueChange={value => {
            setStatusFilter(value as StatusFilter);
            setPageIndex(0);
          }}
        >
          <SelectTrigger className="w-40">
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="ALL">All statuses</SelectItem>
            <SelectItem value="PENDING">Pending</SelectItem>
            <SelectItem value="APPROVED">Approved</SelectItem>
            <SelectItem value="REJECTED">Rejected</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {isApprover && counts && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <MetricCard title="Pending" value={counts.PENDING} />
          <MetricCard title="Approved" value={counts.APPROVED} />
          <MetricCard title="Rejected" value={counts.REJECTED} />
        </div>
      )}

      {loadError && (
        <Card>
          <CardContent className="pt-6 text-destructive">{loadError}</CardContent>
        </Card>
      )}

      {!loadError && rows.length === 0 && !loading && (
        <Card>
          <CardContent className="pt-6 text-muted-foreground">
            No requests to show.
          </CardContent>
        </Card>
      )}

      {!loadError && (rows.length > 0 || loading) && (
        <DataTable
          columns={columns}
          data={rows}
          loading={loading}
          serverSide
          totalRows={total}
          pageIndex={pageIndex}
          pageSize={PAGE_SIZE}
          pageCount={Math.max(1, Math.ceil(total / PAGE_SIZE))}
          onPageChange={setPageIndex}
        />
      )}

      <ConfirmDialog
        open={pending !== null}
        title={pending?.decision === 'APPROVE' ? 'Approve this request?' : 'Reject this request?'}
        description={
          pending
            ? `${pending.request.requester} requested "${pending.request.dataset?.name}". This cannot be undone.`
            : ''
        }
        confirmText={pending?.decision === 'APPROVE' ? 'Approve' : 'Reject'}
        destructive={pending?.decision === 'REJECT'}
        loading={deciding}
        onConfirm={confirmDecision}
        onCancel={() => setPending(null)}
      />
    </div>
  );
}
