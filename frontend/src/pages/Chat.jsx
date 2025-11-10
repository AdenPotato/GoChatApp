import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import websocketService from '../utils/websocket';
import api from '../utils/api';

function Chat() {
  const [messages, setMessages] = useState([]);
  const [users, setUsers] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [connected, setConnected] = useState(false);
  const [loading, setLoading] = useState(true);
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();
  const username = localStorage.getItem('username');

  // Load message history and users on mount
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/login');
      return;
    }

    // Fetch message history
    api.get('/api/messages')
      .then((response) => {
        const messageHistory = response.data || [];
        // Transform backend messages to match frontend format
        const formattedMessages = messageHistory.map((msg) => ({
          username: msg.User?.username || 'Unknown',
          content: msg.content,
          timestamp: msg.created_at,
        }));
        setMessages(formattedMessages);
        setLoading(false);
      })
      .catch((error) => {
        console.error('Failed to load message history:', error);
        setLoading(false);
      });

    // Fetch users list
    api.get('/api/users')
      .then((response) => {
        setUsers(response.data || []);
      })
      .catch((error) => {
        console.error('Failed to load users:', error);
      });
  }, [navigate]);

  // Setup WebSocket connection
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      return;
    }

    // Connect to WebSocket
    websocketService.connect(token)
      .then(() => {
        setConnected(true);
      })
      .catch((error) => {
        console.error('Failed to connect to WebSocket:', error);
        // Don't block if WebSocket fails - REST endpoints still work
      });

    // Handle incoming messages
    const handleMessage = (data) => {
      setMessages((prev) => [...prev, data]);
    };

    websocketService.onMessage(handleMessage);

    return () => {
      websocketService.removeMessageHandler(handleMessage);
      websocketService.disconnect();
    };
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    websocketService.send({
      type: 'message',
      content: newMessage,
      username: username,
      timestamp: new Date().toISOString(),
    });

    setNewMessage('');
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    websocketService.disconnect();
    navigate('/login');
  };

  return (
    <div className="h-screen flex flex-col bg-gray-100">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-4 shadow-lg">
        <div className="container mx-auto flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold">GoChatApp</h1>
            <p className="text-sm opacity-80">
              {loading ? 'Loading...' : connected ? '🟢 Connected' : '🔴 Offline (using REST API)'}
            </p>
          </div>
          <div className="flex items-center space-x-4">
            <span className="text-sm">Welcome, {username}</span>
            <button
              onClick={handleLogout}
              className="bg-white text-purple-600 px-4 py-2 rounded-lg hover:bg-gray-100 transition duration-200"
            >
              Logout
            </button>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 flex overflow-hidden">
        {/* User List Sidebar */}
        <div className="w-64 bg-white border-r border-gray-200 overflow-y-auto">
          <div className="p-4">
            <h2 className="text-lg font-semibold text-gray-800 mb-4">
              Users ({users.length})
            </h2>
            <div className="space-y-2">
              {users.map((user) => (
                <div
                  key={user.id}
                  className="flex items-center space-x-3 p-2 rounded-lg hover:bg-gray-50"
                >
                  <div className="w-10 h-10 rounded-full bg-gradient-to-r from-blue-400 to-purple-400 flex items-center justify-center text-white font-semibold">
                    {user.username.charAt(0).toUpperCase()}
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-800">
                      {user.username}
                      {user.username === username && (
                        <span className="text-xs text-gray-500"> (you)</span>
                      )}
                    </p>
                    <p className="text-xs text-gray-500">{user.email}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Messages Area */}
        <div className="flex-1 flex flex-col">
          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {loading ? (
              <div className="flex items-center justify-center h-full">
                <p className="text-gray-500">Loading messages...</p>
              </div>
            ) : messages.length === 0 ? (
              <div className="flex items-center justify-center h-full">
                <p className="text-gray-500">No messages yet. Start the conversation!</p>
              </div>
            ) : (
              messages.map((msg, index) => (
                <div
                  key={index}
                  className={`flex ${msg.username === username ? 'justify-end' : 'justify-start'}`}
                >
                  <div
                    className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
                      msg.username === username
                        ? 'bg-blue-600 text-white'
                        : 'bg-white text-gray-800'
                    } shadow`}
                  >
                    <div className="flex items-baseline space-x-2">
                      <span className="font-semibold text-sm">{msg.username}</span>
                      <span className="text-xs opacity-70">
                        {new Date(msg.timestamp).toLocaleTimeString()}
                      </span>
                    </div>
                    <p className="mt-1">{msg.content}</p>
                  </div>
                </div>
              ))
            )}
            <div ref={messagesEndRef} />
          </div>

          {/* Message Input */}
          <div className="bg-white border-t border-gray-200 p-4">
            <form onSubmit={handleSendMessage} className="flex space-x-2">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder={connected ? "Type a message..." : "WebSocket disconnected - messages won't send"}
                className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={!connected}
              />
              <button
                type="submit"
                disabled={!connected}
                className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Send
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Chat;
