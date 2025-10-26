import { useEffect, useState, useCallback, useRef } from 'react';
import { websocketService } from '../services/websocket';

interface UseWebSocketOptions {
  autoConnect?: boolean;
  reconnectOnClose?: boolean;
}

interface UseWebSocketResult {
  isConnected: boolean;
  connect: () => Promise<void>;
  disconnect: () => void;
  send: (message: unknown) => void;
  on: (event: string, handler: (data: unknown) => void) => () => void;
  off: (event: string, handler?: (data: unknown) => void) => void;
}

export const useWebSocket = (options: UseWebSocketOptions = {}): UseWebSocketResult => {
  const { autoConnect = true } = options;
  const [isConnected, setIsConnected] = useState(false);
  const connectPromiseRef = useRef<Promise<void> | null>(null);

  const connect = useCallback(async () => {
    if (connectPromiseRef.current) {
      return connectPromiseRef.current;
    }

    try {
      connectPromiseRef.current = websocketService.connect();
      await connectPromiseRef.current;
      setIsConnected(true);
    } catch (error) {
      console.error('WebSocket connection failed:', error);
      setIsConnected(false);
    } finally {
      connectPromiseRef.current = null;
    }
  }, []);

  const disconnect = useCallback(() => {
    websocketService.disconnect();
    setIsConnected(false);
  }, []);

  const send = useCallback((message: unknown) => {
    websocketService.send(message);
  }, []);

  const on = useCallback((event: string, handler: (data: unknown) => void) => {
    return websocketService.on(event, handler);
  }, []);

  const off = useCallback((event: string, handler?: (data: unknown) => void) => {
    websocketService.off(event, handler);
  }, []);

  useEffect(() => {
    if (autoConnect) {
      connect();
    }

    const checkInterval = setInterval(() => {
      const connected = websocketService.isConnected();
      setIsConnected(connected);
    }, 1000);

    return () => {
      clearInterval(checkInterval);
      if (autoConnect) {
        disconnect();
      }
    };
  }, [autoConnect, connect, disconnect]);

  return {
    isConnected,
    connect,
    disconnect,
    send,
    on,
    off,
  };
};

export const useWebSocketEvent = <T = unknown>(
  event: string,
  handler: (data: T) => void,
  deps: React.DependencyList = []
): void => {
  const handlerRef = useRef(handler);

  useEffect(() => {
    handlerRef.current = handler;
  }, [handler]);

  useEffect(() => {
    const wrappedHandler = (data: unknown) => {
      handlerRef.current(data as T);
    };

    const unsubscribe = websocketService.on(event, wrappedHandler);

    return () => {
      unsubscribe();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [event, ...deps]);
};
