import type { ApiError } from '@lattiq/webtk';

function isApiError(err: unknown): err is ApiError {
  return typeof err === 'object' && err !== null && 'status' in err && 'code' in err;
}

/**
 * Extracts a human-readable message from whatever createApiClient's request()
 * rejected with. Backend errors already carry {code, message, details} (see
 * foundry's errors.APIError) — webtk's client re-shapes those into ApiError.
 */
export function getErrorMessage(err: unknown, fallback = 'Something went wrong'): string {
  if (isApiError(err)) return err.message || fallback;
  if (err instanceof Error) return err.message;
  return fallback;
}

export function isUnauthorized(err: unknown): boolean {
  return isApiError(err) && err.status === 401;
}
