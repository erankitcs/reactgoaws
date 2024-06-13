// src/hooks/useWebSocket.js
import { useEffect, useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

const useCustomWebSocket = (url, fetchHistoricalMessages, movieID) => {
  const [messages, setMessages] = useState([]);

  const {
    sendJsonMessage,
    lastJsonMessage,
    readyState,
  } = useWebSocket(url, {
    onOpen: () => console.log('WebSocket connection established'),
    onClose: () => console.log('WebSocket connection closed'),
    shouldReconnect: (closeEvent) => true,
  });

  useEffect(() => {
    if (fetchHistoricalMessages) {
      fetchHistoricalMessages(movieID).then((historicalMessages) => {
        if (historicalMessages) {
          setMessages(historicalMessages);
        }
        
      });
      
    }
  }, [fetchHistoricalMessages, movieID ]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      console.log(lastJsonMessage);
      setMessages((prevMessages) => [...prevMessages, lastJsonMessage]);
    }
  }, [lastJsonMessage]);

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Connecting',
    [ReadyState.OPEN]: 'Open',
    [ReadyState.CLOSING]: 'Closing',
    [ReadyState.CLOSED]: 'Closed',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  return { messages, sendJsonMessage, connectionStatus };
};

export default useCustomWebSocket;
