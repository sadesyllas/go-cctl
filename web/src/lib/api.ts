import type { Writable } from 'svelte/store';

import { browser } from '$app/env';

const API_PORT = <string>import.meta.env.VITE_API_PORT;

let hostname = './';

if (browser) {
  hostname = `${window.location.protocol}//${window.location.hostname}:${API_PORT}`;
}

export type ApiError = { status: number; statusText: string; url: string; data?: { [key: string]: unknown } };

type ApiFetchTracker = Pick<Writable<boolean>, 'set' | 'subscribe'>;

type ApiOptions = RequestInit & {
  baseUrl?: string;
  fetch?: Fetch;
  params?: string | string[][] | Record<string, string> | URLSearchParams;
  tracker?: ApiFetchTracker;
  parseResponseAsText?: boolean;
};

export async function send<T>(path: string, _options?: ApiOptions): Promise<{ headers: Headers; body: T }> {
  const { baseUrl = hostname, fetch: f, params, tracker, ...options } = { ..._options };

  if (tracker) {
    tracker.set(true);
  }

  if (options.body && !(options.body instanceof File) && !(options.body instanceof FormData)) {
    options.headers = { ...options.headers, 'Content-Type': 'application/json' };
    options.body = typeof options.body === 'string' ? options.body : JSON.stringify(options.body);
  }

  let url = buildFullPath(baseUrl, path);
  if (params) {
    url += (url.includes('?') ? '&' : '?') + new URLSearchParams(params).toString();
  }

  return (f || fetch)(url, options).then(async (r) => {
    if (tracker) {
      tracker.set(false);
    }

    const { status, statusText, headers } = r;
    if (status < 200 || status > 299) {
      let data: unknown;
      try {
        data = await parseBody(r, options);
        // eslint-disable-next-line no-empty
      } catch { }
      throw <ApiError>{ status, statusText, url, data };
    }

    return { headers, body: await parseBody<T>(r, options) };
  });
}

export async function get<T>(path: string, options?: ApiOptions): Promise<T> {
  return await send<T>(path, { ...options, method: 'GET' }).then((r) => r.body);
}

export async function getPage<T>(path: string, options?: ApiOptions): Promise<{ data: T[]; total: number; headers }> {
  return await send<T[]>(path, { ...options, method: 'GET' }).then(async (response) => {
    const { headers, body } = response;
    const total = headers.has('content-range') ? +headers.get('content-range').split('/')[1] : undefined;

    return { data: body, total, headers };
  });
}

export async function del<T>(path: string, options?: ApiOptions): Promise<T> {
  return await send<T>(path, { ...options, method: 'DELETE' }).then((r) => r.body);
}

export async function post<T>(path: string, body: BodyInit | null, options?: ApiOptions): Promise<T> {
  return await send<T>(path, { ...options, method: 'POST', body }).then((r) => r.body);
}

export async function put<T>(path: string, body: BodyInit | null, options?: ApiOptions): Promise<T> {
  return await send<T>(path, { ...options, method: 'PUT', body }).then((r) => r.body);
}

export async function patch<T>(path: string, body: BodyInit | null, options?: ApiOptions): Promise<T> {
  return await send<T>(path, { ...options, method: 'PATCH', body }).then((r) => r.body);
}

async function parseBody<T>(response: Response, { parseResponseAsText }: ApiOptions): Promise<T> {
  const { headers } = response;
  if (headers.has('content-length') && +headers.get('content-length') === 0) {
    return;
  }

  if (parseResponseAsText || (headers.has('content-type') && headers.get('content-type').indexOf('text') === 0)) {
    return response.text() as never;
  }

  return await response.json();
}

function isAbsoluteURL(url: string): boolean {
  // A URL is considered absolute if it begins with "<scheme>://" or "//" (protocol-relative URL).
  // RFC 3986 defines scheme name as a sequence of characters beginning with a letter and followed
  // by any combination of letters, digits, plus, period, or hyphen.
  return /^([a-z][\d+.a-z-]*:)?\/\//i.test(url);
}

function combineURLs(baseURL: string, relativeURL: string) {
  return relativeURL ? baseURL.replace(/\/+$/, '') + '/' + relativeURL.replace(/^\/+/, '') : baseURL;
}

function buildFullPath(baseURL: string, requestedURL: string): string {
  return baseURL && !isAbsoluteURL(requestedURL) ? combineURLs(baseURL, requestedURL) : requestedURL;
}
