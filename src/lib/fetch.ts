/**
 * fetch封装类，提供get、post方法，提供请求拦截器和响应拦截器，
 * 支持定义公共基础路径，支持解构使用
 */
export class Fetcher {
  /**
   * 将请求参数拼接至url
   */
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
    // 请求前拦截器
    onRequest?: (request: {
      input: string | URL;
      init?: RequestInit;
    }) => typeof request | undefined;
    // 响应前拦截器
    onResponse?: (response: Response) => typeof response | undefined;
    // 请求基础路径
    base?: string;
  }) {
    let { onRequest, onResponse, base } = option || {};

    this.fetch = async (input, init) => {
      if (typeof input === "string" && base) {
        input = base + input;
      }
      let req: Parameters<NonNullable<typeof onRequest>>[0] | undefined = {
        input,
        init,
      };
      if (onRequest) {
        req = onRequest({ input, init });
      }
      // 若请求拦截器返回空值，则不予发送
      if (!req) return;
      ({ input, init } = req);
      const res: Response | undefined = await fetch(input, init);
      if (onResponse) {
        // 若响应拦截器返回空值，则返回空值
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
        // 如果是FormData类型
        _body = body;
        headers = { "Content-Type": "multipart/form-data", ...headers };
      } else if (body instanceof Object) {
        // 如果是对象，设置json格式
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

  /**
   * 原生fetch，可用于自定义请求，已加载拦截器和基础路径
   */
  fetch: (
    input: string | URL,
    init?: RequestInit
  ) => Promise<Response | undefined>;

  /**
   * 发起get请求
   */
  get: <Response = any, Query = Record<string, string | number>>(
    url: string,
    option?: { headers?: Record<string, string>; query?: Query }
  ) => Promise<Response>;

  /**
   * 发起post请求
   */
  post: <Response = any, Body = any>(
    url: string,
    body?: Body,
    option?: {
      headers?: Record<string, string>;
      query?: Record<string, string | number>;
    }
  ) => Promise<Response>;
}

export const fetcher = new Fetcher({
  onResponse: (res) => {
    if (res.status === 403) {
      // 密码不正确
      localStorage.removeItem("password");
      return;
    }
    return res;
  },
});

const { get, post } = fetcher;

export { get, post };
