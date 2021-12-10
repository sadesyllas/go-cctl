/// <reference types="@sveltejs/kit" />
interface ImportMeta {
  env: {
    VITE_API_PORT: string;
  };
}

type Fetch = (info: RequestInfo, init?: RequestInit) => Promise<Response>;
