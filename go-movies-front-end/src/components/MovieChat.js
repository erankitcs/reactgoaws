import React, { useState } from 'react';
import ChatMessage from "./chat/ChatMessage";
import useCustomWebSocket from '../hooks/useWebSocket';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircle } from '@fortawesome/free-solid-svg-icons';
import { useOutletContext } from "react-router-dom";
 // Load Chat History Function
 const fetchHistoricalMessages = async(id, jwtToken) => {
  console.log("Loading chat history")
  const headers = new Headers();
  headers.append("Content-Type","application/json");
  headers.append("Authorization", "Bearer "+ jwtToken);
  const requestOptions = {
      method: "GET",
      headers: headers,
      credentials: "include"
  }
    
  const response = await fetch(`/protected/movies/${id}/chats`,requestOptions)
  const data = await response.json();
  return data;
};

const MovieChat = (props) => {
  const movieID = props.movieID;
  const { jwtToken, userName } = useOutletContext();
  const [token] = useState(jwtToken);
  console.log(token);
  const { messages, sendJsonMessage, connectionStatus  } = useCustomWebSocket(`http://172.21.246.236:8080/protected/movies/${movieID}/chatws`, fetchHistoricalMessages, movieID, token  );

  // Custom Hook for WebSocket Connection
 
  const [message, setMessage] = useState();

  // useEffect(() => {
  //   fetchChatHistory(movieID);
  // }, [movieID]);


 // Send Msg to socket
 const sendMessageEvent = (e) => {
  e.preventDefault();
  if (message) {
  console.log("Sending message to socket");
  console.log(message);
  const newMessage = {
    type: 'send_message',
    payload: {
      message: message,
    }
  };
  sendJsonMessage(newMessage)
  setMessage('');
}
};

return (
    <div>
      <h2>Chat </h2>
      { connectionStatus === 'Open' ? (<p>Connection Status: <FontAwesomeIcon icon={faCircle} beatFade size="xs" style={{color: "#1cca96",}} /></p>):
      (<p>Connection Status: <FontAwesomeIcon icon={faCircle} beatFade size="xs" style={{color: "#ff0000", }} /></p>)}  
      <div className="chat-box border p-3" style={{ maxHeight: '400px', overflowY: 'auto' }}>
      {messages.map((message, index) => (
           
          <ChatMessage 
            key={index}
            type = {message.type}
            username={message.payload.from}
            text={message.payload.message}
            date={message.payload.sent}
            msgByMe={ message.payload.from === userName}
          />
        ))}
      </div>
      <form onSubmit={sendMessageEvent} className="mt-3">
        <div className="input-group">
          <input
            type="text"
            className="form-control"
            placeholder="Type your message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
          />
          <button type="submit" className="btn btn-primary">
            Send
          </button>
        </div>
      </form>
    </div>
  );
};

export default MovieChat;
