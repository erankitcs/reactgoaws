import React, { useState } from 'react';

const ChatMessage = ({ author, text, date }) => {
  return (
    <div className={`d-flex flex-row ${author === 'You' ? 'justify-content-end' : 'justify-content-start'} mb-3`}>
      <div className={`chat-bubble ${author === 'You' ? 'bg-primary text-white' : 'bg-light'} p-3 rounded`}>
        <div className="d-flex justify-content-between">
          <small className="fw-bold">{author}</small>
          <small className="text-muted">{date}</small>
        </div>
        <p className="mb-0">{text}</p>
      </div>
    </div>
  );
};

const ChatBox = () => {
  const [messages, setMessages] = useState([
    {
      author: 'You',
      text: 'Hello, how can I help you today?',
      date: '2023-05-01 10:00 AM',
    },
    {
      author: 'John Doe',
      text: 'Hi, I have a question about your product.',
      date: '2023-05-01 10:02 AM',
    },
    {
      author: 'You',
      text: 'Sure, please go ahead and ask your question.',
      date: '2023-05-01 10:03 AM',
    },
    {
      author: 'John Doe',
      text: 'Can you explain the pricing for your premium plan?',
      date: '2023-05-01 10:05 AM',
    },
  ]);
  const [newMessage, setNewMessage] = useState('');

  const handleMessageChange = (e) => {
    setNewMessage(e.target.value);
  };

  const handleMessageSubmit = (e) => {
    e.preventDefault();
    if (newMessage.trim()) {
      const currentDate = new Date().toLocaleString();
      const newMessageObj = {
        author: 'You',
        text: newMessage,
        date: currentDate,
      };
      setMessages([...messages, newMessageObj]);
      setNewMessage('');

      try {
        fetch('/api/chat/messages', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(newMessageObj),
        })
        .then((response) => {response.json()})
        .then((msg) => {
          console.log(msg);
          setMessages([...messages, newMessageObj, botResponse]);
        });
        const response = await axios.post('/api/chat/messages', newMessageObj);
        const botResponse = response.data;
        
      } catch (error) {
        console.error('Error sending message:', error);
      }
    }
  };

  return (
    <div>
      <h2>Chat</h2>
      <div className="chat-box border p-3" style={{ maxHeight: '400px', overflowY: 'auto' }}>
        {messages.map((message, index) => (
          <ChatMessage
            key={index}
            author={message.author}
            text={message.text}
            date={message.date}
          />
        ))}
      </div>
      <form onSubmit={handleMessageSubmit} className="mt-3">
        <div className="input-group">
          <input
            type="text"
            className="form-control"
            placeholder="Type your message..."
            value={newMessage}
            onChange={handleMessageChange}
          />
          <button type="submit" className="btn btn-primary">
            Send
          </button>
        </div>
      </form>
    </div>
  );
};

export default ChatBox;
