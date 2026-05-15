import { useCallback, useRef } from "react";

type Navigate = (path: string) => void | Promise<void>;

export function useGuardedNavigation(onNavigate: Navigate) {
  const navigationInProgressRef = useRef(false);

  return useCallback(
    async (path: string) => {
      if (navigationInProgressRef.current) {
        return;
      }

      navigationInProgressRef.current = true;
      try {
        await onNavigate(path);
      } finally {
        navigationInProgressRef.current = false;
      }
    },
    [onNavigate],
  );
}
