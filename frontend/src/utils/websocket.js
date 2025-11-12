// IMPORTANT: The WebSocket endpoint is at /api/ws NOT /ws
// Cache busting: 2025-11-12-v2
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/ws';

class WebSocketService {
  constructor() {
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
    this.messageHandlers = [];
    console.log('[WebSocketService] Initialized with URL:', WS_URL);
  }

  connect(userId, username) {
    return new Promise((resolve, reject) => {
      try {
        const url = `${WS_URL}?user_id=${userId}&username=${encodeURIComponent(username)}`;
        console.log('[WebSocketService] Connecting to:', url);
        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            this.messageHandlers.forEach((handler) => handler(data));
          } catch (error) {
            console.error('Error parsing WebSocket message:', error);
          }
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          reject(error);
        };

        this.ws.onclose = () => {
          console.log('WebSocket disconnected');
          this.attemptReconnect(userId, username);
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  attemptReconnect(userId, username) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
      setTimeout(() => {
        this.connect(userId, username);
      }, this.reconnectDelay * this.reconnectAttempts);
    }
  }

  send(data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  onMessage(handler) {
    this.messageHandlers.push(handler);
  }

  removeMessageHandler(handler) {
    this.messageHandlers = this.messageHandlers.filter((h) => h !== handler);
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

export default new WebSocketService();
