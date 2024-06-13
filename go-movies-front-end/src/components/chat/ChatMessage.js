const ChatMessage = ({ type, userid, username, text, date }) => {
    if (type === 'user_joined') {
      text = `${username} has joined the chat.`;
    } else if (type === 'user_left') {
      text = `${username} has left the chat.`;
    }
    return (
      
      <div className={`d-flex flex-row ${userid === 'You' ? 'justify-content-end' : 'justify-content-start'} mb-3`}>
        <div className={`chat-bubble ${userid === 'You' ? 'bg-primary text-white' : 'bg-light'} p-3 rounded`}>
          <div className="d-flex justify-content-between">
            <small className="fw-bold">{username}</small>
            <small className="text-muted">{date}</small>
          </div>
          <p className="mb-0">{text}</p>
        </div>
      </div>
    );
  };

export default ChatMessage;