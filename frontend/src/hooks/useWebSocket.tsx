import { useRef, useState } from 'react';

type WebSocketMessage = {
  type: string;
  [key: string]: any;
};

type WebSocketHookOptions = {
  url: string;
  onMessage: (message?: WebSocketMessage) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (error: Event) => void;
  reconnectAttempts?: number;
  reconnectInterval?: number;
};

export const useWebSocket = (options: WebSocketHookOptions) => {
  const {
    url,
    onMessage,
    onOpen,
    onClose,
    onError,
  } = options;

  // const [playerId, setPlayerId] = useState<string | null>(() => {
  //   // Initialize from localStorage if available
  //   if (typeof window !== 'undefined') {
  //     // return localStorage.getItem('wsPlayerId');
  //   }
  //   return null;
  // });

  const wsRef = useRef<WebSocket | null>(null);

  const connect = () => {
    // Clear any existing connection
    if (wsRef.current) {
      // wsRef.current.close();
    }


    wsRef.current = new WebSocket(url);
    console.log(url)

    wsRef.current.onopen = () => {

      // If we have a playerId, send reconnect request
      // if (playerId) {
      //   send({ type: 'RECONNECT', playerId });
      // }

      if (onOpen) onOpen();
    };

    wsRef.current.onmessage = (event) => {
      let message: WebSocketMessage | undefined;
      if (event.data) {
        message = JSON.parse(event.data) as WebSocketMessage;
      }

      // Handle connection messages
      // if (message.type === 'CONNECTED' && message.playerId) {
      //   setPlayerId(message.playerId);
      //   localStorage.setItem('wsPlayerId', message.playerId);
      // }
      if (message) {
        onMessage(message);
      }
    };

    wsRef.current.onclose = () => {
      console.log('WebSocket disconnected');
      if (onClose) onClose();

      // Attempt reconnection if we haven't exceeded max attempts
      // if (reconnectCountRef.current < reconnectAttempts) {
      //   reconnectCountRef.current += 1;
      //   reconnectTimeoutRef.current = setTimeout(
      //     connect,
      //     reconnectInterval
      //   );
      //   console.log(`Reconnecting attempt ${reconnectCountRef.current}`);
      // }
    };

    wsRef.current.onerror = (error) => {
      console.error('WebSocket error:', error);
      if (onError) onError(error);
    };
  };

  const send = (data: WebSocketMessage) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data));
    } else {
      console.warn('WebSocket not connected. Message not sent:', data);
    }
  };

  const disconnect = () => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
  };

  return [connect, disconnect, send];
}
