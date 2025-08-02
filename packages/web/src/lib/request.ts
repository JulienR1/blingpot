import z from "zod";

type RequestOptions = Omit<RequestInit, "method" | "credentials">;
type ResponseType = "json" | "text" | "none";

const methods = ["get", "post", "put", "delete"] as const;
type Method = (typeof methods)[number];

function make<F>(factory: (method: string) => F): Record<Method, F> {
  return methods.reduce(
    (acc, method) => ({ ...acc, [method]: factory(method) }),
    {} as Record<Method, F>,
  );
}

const fcts = (url: string) => ({
  json: make(
    (method) =>
      <S extends z.ZodType>(schema: S, opts?: RequestOptions) =>
        execute(url, "json", method, schema, opts),
  ),
  text: make(
    (method) => (opts?: RequestOptions) =>
      execute(url, "text", method, z.string(), opts),
  ),
  none: make(
    (method) => (opts?: RequestOptions) =>
      execute(url, "none", method, z.null(), opts),
  ),
});

export function request<T extends ResponseType = "json">(
  url: string,
  type?: T,
): ReturnType<typeof fcts>[T extends "json" ? "json" : T] {
  // eslint-disable-next-line
  return fcts(url)[type ?? "json"] as any;
}

async function execute<S extends z.ZodType>(
  url: string,
  type: ResponseType,
  method: string,
  schema: S,
  opts: RequestInit = {},
) {
  try {
    const response = await fetch(url, {
      method,
      credentials: "include",
      ...opts,
    });

    if (type === "none") {
      return null;
    }

    const promise = type === "json" ? response.json : response.text;
    const data = await promise();
    return schema.parse(data);
  } catch {
    return null;
  }
}
