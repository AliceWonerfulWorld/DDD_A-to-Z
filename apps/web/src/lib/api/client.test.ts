import { describe, it, expect, vi, beforeEach } from "vitest";
import { apiFetch, ApiError } from "./client";

function mockFetch(status: number, body: unknown, ok?: boolean) {
  const isOk = ok ?? (status >= 200 && status < 300);
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      ok: isOk,
      status,
      json: () => Promise.resolve(body),
    }),
  );
}

beforeEach(() => {
  vi.unstubAllGlobals();
});

describe("ApiError", () => {
  it("name が ApiError になっている", () => {
    const err = new ApiError(404, null);
    expect(err.name).toBe("ApiError");
    expect(err.status).toBe(404);
  });

  it("instanceof Error が true", () => {
    expect(new ApiError(500, null) instanceof Error).toBe(true);
  });
});

describe("apiFetch", () => {
  it("成功時にレスポンスのJSONを返す", async () => {
    mockFetch(200, { user: "alice" });
    const result = await apiFetch<{ user: string }>("/me");
    expect(result).toEqual({ user: "alice" });
  });

  it("204 の場合は undefined を返す", async () => {
    mockFetch(204, null);
    const result = await apiFetch("/noop");
    expect(result).toBeUndefined();
  });

  it("!res.ok の場合は ApiError を投げる", async () => {
    mockFetch(401, { message: "Unauthorized" });
    await expect(apiFetch("/me")).rejects.toBeInstanceOf(ApiError);
  });

  it("エラー時の status が正しく伝わる", async () => {
    mockFetch(403, { message: "Forbidden" });
    const err = await apiFetch("/secret").catch((e) => e);
    expect(err).toBeInstanceOf(ApiError);
    expect((err as ApiError).status).toBe(403);
  });

  it("credentials: include が設定されている", async () => {
    mockFetch(200, {});
    const spy = vi
      .fn()
      .mockResolvedValue({ ok: true, status: 200, json: () => Promise.resolve({}) });
    vi.stubGlobal("fetch", spy);
    await apiFetch("/me");
    expect(spy).toHaveBeenCalledWith(
      "/api/me",
      expect.objectContaining({ credentials: "include" }),
    );
  });

  it("GET リクエストには Content-Type を付けない", async () => {
    const spy = vi
      .fn()
      .mockResolvedValue({ ok: true, status: 200, json: () => Promise.resolve({}) });
    vi.stubGlobal("fetch", spy);
    await apiFetch("/me");
    const headers = spy.mock.calls[0][1].headers as Headers;
    expect(headers.has("Content-Type")).toBe(false);
  });

  it("POST リクエストには Content-Type: application/json を付ける", async () => {
    const spy = vi
      .fn()
      .mockResolvedValue({ ok: true, status: 200, json: () => Promise.resolve({}) });
    vi.stubGlobal("fetch", spy);
    await apiFetch("/items", { method: "POST", body: JSON.stringify({ x: 1 }) });
    const headers = spy.mock.calls[0][1].headers as Headers;
    expect(headers.get("Content-Type")).toBe("application/json");
  });

  it("呼び出し元が Content-Type を上書きできる", async () => {
    const spy = vi
      .fn()
      .mockResolvedValue({ ok: true, status: 200, json: () => Promise.resolve({}) });
    vi.stubGlobal("fetch", spy);
    await apiFetch("/upload", {
      method: "POST",
      headers: { "Content-Type": "multipart/form-data" },
    });
    const headers = spy.mock.calls[0][1].headers as Headers;
    expect(headers.get("Content-Type")).toBe("multipart/form-data");
  });

  it("tuple array 形式の headers を渡せる", async () => {
    const spy = vi
      .fn()
      .mockResolvedValue({ ok: true, status: 200, json: () => Promise.resolve({}) });
    vi.stubGlobal("fetch", spy);
    await apiFetch("/items", {
      headers: [["Accept", "application/json"]],
    });
    const headers = spy.mock.calls[0][1].headers as Headers;
    expect(headers.get("Accept")).toBe("application/json");
  });
});
