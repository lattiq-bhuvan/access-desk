import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  FormWrapper,
  Skeleton,
  TextareaField,
  toast,
  useFormValidation,
} from '@lattiq/design-system';
import { useEffect, useState } from 'react';
import { z } from 'zod';

import { getErrorMessage } from '@/lib/api/errors';
import { createRequest } from '@/lib/api/requests';
import { listDatasets } from '@/lib/api/datasets';
import type { Dataset } from '@/types/domain';

const requestAccessSchema = z.object({
  reason: z
    .string()
    .min(10, 'Give at least 10 characters of context')
    .max(500, 'Keep it under 500 characters'),
});

export function Catalog() {
  const [datasets, setDatasets] = useState<Dataset[] | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [activeDataset, setActiveDataset] = useState<Dataset | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const form = useFormValidation({ schema: requestAccessSchema, defaultValues: { reason: '' } });

  useEffect(() => {
    let cancelled = false;
    listDatasets()
      .then(data => {
        if (!cancelled) setDatasets(data);
      })
      .catch(err => {
        if (!cancelled) setLoadError(getErrorMessage(err, 'Failed to load datasets'));
      });
    return () => {
      cancelled = true;
    };
  }, []);

  function openRequestDialog(dataset: Dataset) {
    form.reset({ reason: '' });
    setActiveDataset(dataset);
  }

  async function handleSubmit(data: { reason: string }) {
    if (!activeDataset) return;
    setSubmitting(true);
    try {
      await createRequest({ dataset_id: activeDataset.id, reason: data.reason });
      toast.success(`Access request sent for ${activeDataset.name}`);
      setActiveDataset(null);
    } catch (err) {
      toast.error(getErrorMessage(err, 'Could not submit the request'));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold">Dataset catalog</h1>
        <p className="text-muted-foreground">
          Request access to a dataset. An approver will review it on the Requests page.
        </p>
      </div>

      {loadError && (
        <Card>
          <CardContent className="pt-6 text-destructive">{loadError}</CardContent>
        </Card>
      )}

      {!datasets && !loadError && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3].map(i => (
            <Skeleton key={i} className="h-40 w-full" />
          ))}
        </div>
      )}

      {datasets && datasets.length === 0 && (
        <Card>
          <CardContent className="pt-6 text-muted-foreground">
            No datasets are available yet.
          </CardContent>
        </Card>
      )}

      {datasets && datasets.length > 0 && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {datasets.map(dataset => (
            <Card key={dataset.id}>
              <CardHeader>
                <CardTitle>{dataset.name}</CardTitle>
                <CardDescription>{dataset.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <Button onClick={() => openRequestDialog(dataset)}>Request access</Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog
        open={activeDataset !== null}
        onOpenChange={open => !open && setActiveDataset(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Request access</DialogTitle>
            <DialogDescription>
              {activeDataset ? `Tell the approver why you need "${activeDataset.name}".` : ''}
            </DialogDescription>
          </DialogHeader>

          {activeDataset && (
            <FormWrapper
              form={form}
              showCard={false}
              submitText={submitting ? 'Submitting...' : 'Submit request'}
              isLoading={submitting}
              onSubmit={handleSubmit}
            >
              <TextareaField
                form={form}
                name="reason"
                label="Reason"
                placeholder="e.g. Need this to validate the onboarding risk model"
                rows={4}
              />
            </FormWrapper>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
