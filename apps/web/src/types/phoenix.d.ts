declare module "phoenix" {
  export class Socket {
    constructor(
      endPoint: string,
      opts?: {
        params?: Record<string, string> | (() => Promise<Record<string, string>>);
        timeout?: number;
        heartbeatIntervalMs?: number;
        reconnectAfterMs?: (tries: number) => number;
        logger?: (kind: string, msg: string, data: unknown) => void;
        longpollerTimeout?: number;
        encode?: (payload: object, callback: (result: string) => void) => void;
        decode?: (payload: string, callback: (result: object) => void) => void;
        transport?: typeof WebSocket;
      },
    );
    connect(params?: Record<string, string>): void;
    disconnect(callback?: () => void, code?: number, reason?: string): void;
    channel(topic: string, chanParams?: Record<string, unknown>): Channel;
    onOpen(callback: () => void): void;
    onClose(callback: () => void): void;
    onError(callback: (error: Event) => void): void;
    onMessage(callback: (msg: unknown) => void): void;
  }

  type PushStatus = "ok" | "error" | "timeout";

  export class Push {
    receive(status: PushStatus, callback: (response: unknown) => void): Push;
  }

  export class Channel {
    join(timeout?: number): Push;
    leave(timeout?: number): Push;
    push(event: string, payload: Record<string, unknown>, timeout?: number): Push;
    on(event: string, callback: (response: unknown) => void): number;
    off(event: string, ref?: number): void;
  }
}
