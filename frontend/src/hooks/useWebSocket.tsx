import { useCallback, useEffect, useRef, useState } from 'react';

type WebSocketMessage = {
  type: string;
  [key: string]: any;
};

type WebSocketHookOptions = {
  url: string;
  onMessage?: (message: WebSocketMessage) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (error: Event) => void;
  enabled?: boolean;
};

type WebSocketState = {
  isConnected: boolean;
  error?: Event;
};

type WebSocketHookReturn = {
  connect: () => void;
  disconnect: () => void;
  send: (data: WebSocketMessage) => void;
  state: WebSocketState;
};

export const useWebSocket = (options: WebSocketHookOptions): WebSocketHookReturn => {
  const {
    url,
    onMessage,
    onOpen,
    onClose,
    onError,
    enabled = true,
  } = options;

  const wsRef = useRef<WebSocket | null>(null);
  const [state, setState] = useState<WebSocketState>({
    isConnected: false,
    error: undefined,
  });

  const callbacksRef = useRef({
    onMessage,
    onOpen,
    onClose,
    onError,
  });

  useEffect(() => {
    callbacksRef.current = {
      onMessage,
      onOpen,
      onClose,
      onError,
    };
  }, [onMessage, onOpen, onClose, onError]);

  const connect = useCallback(() => {
    if (wsRef.current) {
      return;
    }

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      setState(prev => ({ ...prev, isConnected: true, error: undefined }));
      callbacksRef.current.onOpen?.();
    };

    ws.onmessage = (event) => {
      try {
        const message = event.data ? JSON.parse(event.data) as WebSocketMessage : undefined;
        if (message) {
          callbacksRef.current.onMessage?.(message);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    ws.onclose = () => {
      setState(prev => ({ ...prev, isConnected: false }));
      callbacksRef.current.onClose?.();
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      setState(prev => ({ ...prev, isConnected: false, error }));
      callbacksRef.current.onError?.(error);
    };
  }, [url]);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.onopen = null;
      wsRef.current.onmessage = null;
      wsRef.current.onclose = null;
      wsRef.current.onerror = null;

      wsRef.current.close();
      wsRef.current = null;
      setState(prev => ({ ...prev, isConnected: false }));
    }
  }, []);

  const send = useCallback((data: WebSocketMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data));
    } else {
      console.warn('WebSocket not connected. Message not sent:', data);
    }
  }, []);

  useEffect(() => {
    if (enabled) {
      connect();
    } else {
      disconnect();
    }

    return () => {
      disconnect();
    };
  }, [enabled, connect, disconnect]);

  return {
    connect,
    disconnect,
    send,
    state,
  };
};
