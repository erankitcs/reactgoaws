import React, { useState } from 'react';
import ChatMessage from "./chat/ChatMessage";
import useCustomWebSocket from '../hooks/useWebSocket';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircle } from '@fortawesome/free-solid-svg-icons';

 // Load Chat History Function
 const fetchHistoricalMessages = async(id) => {
  console.log("Loading chat history")
  const headers = new Headers();
  headers.append("Content-Type","application/json");
  const requestOptions = {
      method: "GET",
      headers: headers
  }
    
  const response = await fetch(`/movies/${id}/chats`,requestOptions)
  const data = await response.json();
  return data;
};

const MovieChat = (props) => {
  const movieID = props.movieID;

  const { messages, sendJsonMessage, connectionStatus } = useCustomWebSocket(`http://172.21.246.236:8080/movies/${movieID}/chatws`, fetchHistoricalMessages, movieID );
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
      from: 'Ankit', // Replace with your actual username
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
            userid={1}
            username={message.payload.from}
            text={message.payload.message}
            date={message.payload.sent}
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
