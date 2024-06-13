// src/hooks/useWebSocket.js
import { useEffect, useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

const useCustomWebSocket = (url, fetchHistoricalMessages, movieID, jwtToken) => {
  const [messages, setMessages] = useState([]);
  const [error, setError] = useState(null);
  console.log(jwtToken)

  const {
    sendJsonMessage,
    lastJsonMessage,
    readyState,
  } = useWebSocket(url, {
    protocols: ["Authorization", jwtToken, "chat"],
    onOpen: () => console.log('WebSocket connection established'),
    onClose: () => console.log('WebSocket connection closed'),
    shouldReconnect: (closeEvent) => true,
    onError: (event) => setError(event),
  });

  useEffect(() => {
    if (fetchHistoricalMessages && jwtToken) {
      console.log("fetchHistoricalMessages")
      fetchHistoricalMessages(movieID, jwtToken).then((historicalMessages) => {
        if (historicalMessages) {
          setMessages(historicalMessages);
        }
        
      });
      
    }
  }, [fetchHistoricalMessages, movieID, jwtToken ]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      console.log(lastJsonMessage);
      setMessages((prevMessages) => [...prevMessages, lastJsonMessage]);
    }
  }, [lastJsonMessage]);

  useEffect(() => {
    if (error) {
      console.error('WebSocket error:', error);
      // You can add additional error handling logic here
    }
  }, [error]);

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Connecting',
    [ReadyState.OPEN]: 'Open',
    [ReadyState.CLOSING]: 'Closing',
    [ReadyState.CLOSED]: 'Closed',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  return { messages, sendJsonMessage, connectionStatus  };
};

export default useCustomWebSocket;
