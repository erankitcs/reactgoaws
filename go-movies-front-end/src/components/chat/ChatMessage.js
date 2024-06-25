import DateFormater from "../../utils/DateFormater";
const ChatMessage = ({ type, username, text, date, msgByMe }) => {
    var infoType = '';
    if (type === 'user_joined') {
      text = `${username} has joined the chat.`;
      infoType = 'notification';
    } else if (type === 'user_left') {
      text = `${username} has left the chat.`;
      infoType = 'notification';
    }
    return (
      
      <div className={`d-flex flex-row ${msgByMe ? 'justify-content-end' : 'justify-content-start'} mb-3`}>
        <div className={`chat-bubble ${msgByMe ? 'border border-success text-white' : 'border border-primary'} p-3 rounded`}>
          <div className="d-flex justify-content-between">
            <small className="text-success me-3">{username}</small>
            <small className="text-muted">{DateFormater(date)}</small>
          </div>
          <div className={`text-break ${infoType === 'notification' ? 'text-info': 'text-dark'}`} style={{ maxWidth: '500px' }}>{text}</div>
        </div>
      </div>
    );
  };

export default ChatMessage;