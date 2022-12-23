import { serveDir } from "https://deno.land/std@0.170.0/http/file_server.ts";
import {
  Handler,
  serve,
  serveTls,
} from "https://deno.land/std@0.170.0/http/server.ts";
import { createRouter } from "https://esm.sh/radix3@1.0.0";

const config = JSON.parse(await Deno.readTextFile("./config.json"));

const router = createRouter();

const serviceConfig: Record<string, string> = config.service;
for (const key in serviceConfig) {
  const proxy = new URL(serviceConfig[key]);
  const data: { handler: Handler } = {
    handler: (req) => {
      const url = new URL(req.url);
      url.pathname = url.pathname.replace(key, proxy.pathname);
      url.host = proxy.host;
      url.port = proxy.port;
      url.protocol = proxy.protocol;
      return fetch(url, req);
    },
  };
  router.insert(key, data);
  if (key.endsWith("/")) router.insert(key + "**", data);
  else router.insert(key + "/**", data);
}

const staticConfig: Record<string, string> = config.static;
for (const key in staticConfig) {
  const dir = staticConfig[key];
  const root = key.replace("/", "");
  const opt = { fsRoot: dir, urlRoot: root, showIndex: true };
  const data: { handler: Handler } = {
    handler: async (req) => {
      const url = new URL(req.url);
      let response = await serveDir(req, opt);
      if (response.status === 404) {
        url.pathname = key;
        response = await serveDir(new Request(url), opt);
      }
      return response;
    },
  };
  router.insert(key, data);
  if (key.endsWith("/")) router.insert(key + "**", data);
  else router.insert(key + "/**", data);
}

const handler: Handler = (req) => {
  const match = router.lookup(new URL(req.url).pathname);
  const res = match?.handler(req);
  if (!res) return new Response("Not Found", { status: 404 });
  return res;
};

if (config.HTTPS) {
  const CertFile = config.HTTPS.CertFile;
  const KeyFile = config.HTTPS.KeyFile;
  serveTls(handler, {
    port: 443,
    certFile: CertFile,
    keyFile: KeyFile,
  });
} else {
  serve(handler, { port: 80 });
}

const cdnList: URL[] = config.cdn.map((url: string) => new URL(url));
const cdnHandler: Handler = async (req) => {
  for (const item of cdnList) {
    const url = new URL(req.url);
    url.host = item.host;
    url.port = item.port;
    url.protocol = item.protocol;
    console.log(url.toString());
    const res = await fetch(url);
    if (!res.ok) continue;
    return res;
  }
  return new Response("Not Found", { status: 404 });
};

if (config.HTTPS) {
  const CertFile = config.HTTPS.CertFile;
  const KeyFile = config.HTTPS.KeyFile;
  serveTls(handler, {
    port: 9002,
    certFile: CertFile,
    keyFile: KeyFile,
  });
} else {
  serve(cdnHandler, { port: 9002 });
}
