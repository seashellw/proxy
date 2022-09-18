export class Fetcher {
  static spliceURLSearchParams = (url: string, params?: any) => {
    if (!params) return url;
    let query = new URLSearchParams(url.split("?")[1] || undefined);

    if (params instanceof URLSearchParams) {
      for (const [key, value] of params) {
        query.set(key, value);
      }
    } else {
      for (const key in params) {
        query.set(key, `${params[key]}`);
      }
    }

    return `${url.split("?")[0]}?${query.toString()}`;
  };

  constructor(option?: {
    onRequest?: (request: {
      input: string | URL;
      init?: RequestInit;
    }) => typeof request | undefined;
    onResponse?: (response: Response) => typeof response | undefined;
    base?: string;
  }) {
    let { onRequest, onResponse, base } = option || {};

    this.fetch = async (input, init) => {
      input = new URL(input, base || undefined);
      let req: Parameters<NonNullable<typeof onRequest>>[0] | undefined = {
        input,
        init,
      };
      if (onRequest) {
        req = onRequest({ input, init });
      }
      if (!req) return;
      ({ input, init } = req);
      const res: Response | undefined = await fetch(input, init);
      if (onResponse) {
        return onResponse(res);
      }
      if (!res.ok) {
        console.error("fetch not ok", res);
      }
      return res;
    };

    this.get = async (url, option) => {
      url = Fetcher.spliceURLSearchParams(url, option?.query);
      const res = await this.fetch(url, {
        method: "GET",
        headers: option?.headers,
      });
      if (!res) return;
      if (res.headers.get("content-type")?.includes("application/json")) {
        return res.json();
      }
      return res.text();
    };

    this.post = async (url, body, option) => {
      let { headers, query } = option || {};
      url = Fetcher.spliceURLSearchParams(url, query);
      let _body: string | FormData | undefined;
      if (!body) {
        _body = undefined;
      } else if (body instanceof FormData) {
        _body = body;
        headers = { "Content-Type": "multipart/form-data", ...headers };
      } else if (body instanceof Object) {
        _body = JSON.stringify(body);
        headers = { "Content-Type": "application/json", ...headers };
      } else {
        _body = undefined;
        console.error("Body must be an Object or FormData");
      }
      const res = await this.fetch(url, {
        method: "POST",
        headers: option?.headers,
        body: _body,
      });
      if (!res) return;
      if (res.headers.get("content-type")?.includes("application/json")) {
        return res.json();
      }
      return res.text();
    };
  }

  fetch: (
    input: string | URL,
    init?: RequestInit
  ) => Promise<Response | undefined>;

  get: <Response = any, Query = Record<string, string | number>>(
    url: string,
    option?: { headers?: Record<string, string>; query?: Query }
  ) => Promise<Response>;

  post: <Response = any, Body = any>(
    url: string,
    body?: Body,
    option?: {
      headers?: Record<string, string>;
      query?: Record<string, string | number>;
    }
  ) => Promise<Response>;
}

export const fetcher = new Fetcher();

const { get, post } = fetcher;

export { get, post };
