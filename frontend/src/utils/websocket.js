/**
 * WebSocket服务
 * 用于处理实时更新，如扫描任务进度、系统通知等
 */

class WebSocketService {
  constructor() {
    this.socket = null
    this.isConnected = false
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectInterval = 3000
    this.eventListeners = {}
    this.subscribedTopics = new Set()
    this.pingInterval = null
    this.pongTimeout = null
  }

  /**
   * 连接到WebSocket服务器
   */
  connect() {
    if (this.socket && (this.socket.readyState === WebSocket.CONNECTING || this.socket.readyState === WebSocket.OPEN)) {
      console.log('WebSocket已连接或正在连接中')
      return
    }

    // 获取API的基本URL，并转换为WebSocket URL
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsHost = process.env.VUE_APP_API_BASE_URL 
      ? process.env.VUE_APP_API_BASE_URL.replace(/^https?:\/\//, '')
      : window.location.host
    const wsUrl = `${wsProtocol}//${wsHost}/ws`

    try {
      this.socket = new WebSocket(wsUrl)
      
      this.socket.onopen = this.onOpen.bind(this)
      this.socket.onclose = this.onClose.bind(this)
      this.socket.onmessage = this.onMessage.bind(this)
      this.socket.onerror = this.onError.bind(this)
      
      console.log('WebSocket连接中...')
    } catch (error) {
      console.error('WebSocket连接失败:', error)
      this.scheduleReconnect()
    }
  }

  /**
   * 断开WebSocket连接
   */
  disconnect() {
    this.clearPingInterval()
    this.clearPongTimeout()
    this.subscribedTopics.clear()
    
    if (this.socket) {
      this.socket.close()
      this.socket = null
    }
    
    this.isConnected = false
    console.log('WebSocket已断开连接')
  }

  /**
   * 当WebSocket连接打开时调用
   * @param {Event} event 
   */
  onOpen(event) {
    this.isConnected = true
    this.reconnectAttempts = 0
    console.log('WebSocket连接成功')
    
    // 重新订阅之前的主题
    this.resubscribeTopics()
    
    // 设置心跳检测
    this.setupPing()
    
    // 分发连接事件
    this.dispatchEvent('connection', { connected: true })
  }

  /**
   * 当WebSocket连接关闭时调用
   * @param {CloseEvent} event 
   */
  onClose(event) {
    this.isConnected = false
    console.log(`WebSocket连接关闭: ${event.code} ${event.reason}`)
    
    this.clearPingInterval()
    this.clearPongTimeout()
    
    // 自动重连
    if (!event.wasClean) {
      this.scheduleReconnect()
    }
    
    // 分发连接事件
    this.dispatchEvent('connection', { connected: false })
  }

  /**
   * 当收到WebSocket消息时调用
   * @param {MessageEvent} event 
   */
  onMessage(event) {
    try {
      const message = JSON.parse(event.data)
      
      // 处理pong响应
      if (message.type === 'pong') {
        this.clearPongTimeout()
        return
      }
      
      // 处理普通事件消息
      if (message.type) {
        this.dispatchEvent(message.type, message.data)
      }
      
      // 日志接收的消息
      console.log('收到WebSocket消息:', message)
    } catch (error) {
      console.error('解析WebSocket消息失败:', error, event.data)
    }
  }

  /**
   * 当WebSocket发生错误时调用
   * @param {Event} event 
   */
  onError(event) {
    console.error('WebSocket错误:', event)
  }

  /**
   * 安排重新连接
   */
  scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('达到最大重连次数，停止重连')
      return
    }

    this.reconnectAttempts++
    
    console.log(`尝试重新连接 (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)
    
    setTimeout(() => {
      this.connect()
    }, this.reconnectInterval)
  }

  /**
   * 发送消息到WebSocket服务器
   * @param {Object} message 消息对象
   */
  sendMessage(message) {
    if (!this.isConnected) {
      console.warn('WebSocket未连接，无法发送消息')
      return false
    }

    try {
      this.socket.send(JSON.stringify(message))
      return true
    } catch (error) {
      console.error('发送WebSocket消息失败:', error)
      return false
    }
  }

  /**
   * 订阅主题
   * @param {string} topic 主题名称
   */
  subscribe(topic) {
    if (!topic) return false

    this.subscribedTopics.add(topic)
    
    if (this.isConnected) {
      return this.sendMessage({
        type: 'subscribe',
        topic: topic
      })
    }
    
    return false
  }

  /**
   * 取消订阅主题
   * @param {string} topic 主题名称
   */
  unsubscribe(topic) {
    if (!topic) return false

    this.subscribedTopics.delete(topic)
    
    if (this.isConnected) {
      return this.sendMessage({
        type: 'unsubscribe',
        topic: topic
      })
    }
    
    return false
  }

  /**
   * 重新订阅之前的主题
   */
  resubscribeTopics() {
    this.subscribedTopics.forEach(topic => {
      this.sendMessage({
        type: 'subscribe',
        topic: topic
      })
    })
  }

  /**
   * 添加事件监听器
   * @param {string} eventType 事件类型
   * @param {Function} listener 监听器函数
   */
  addEventListener(eventType, listener) {
    if (!eventType || typeof listener !== 'function') return

    if (!this.eventListeners[eventType]) {
      this.eventListeners[eventType] = []
    }

    // 避免重复添加相同的监听器
    if (!this.eventListeners[eventType].includes(listener)) {
      this.eventListeners[eventType].push(listener)
    }
  }

  /**
   * 移除事件监听器
   * @param {string} eventType 事件类型
   * @param {Function} listener 监听器函数
   */
  removeEventListener(eventType, listener) {
    if (!eventType || !this.eventListeners[eventType]) return

    const index = this.eventListeners[eventType].indexOf(listener)
    if (index !== -1) {
      this.eventListeners[eventType].splice(index, 1)
    }
    
    // 如果没有监听器了，清理数组
    if (this.eventListeners[eventType].length === 0) {
      delete this.eventListeners[eventType]
    }
  }

  /**
   * 分发事件到监听器
   * @param {string} eventType 事件类型
   * @param {any} data 事件数据
   */
  dispatchEvent(eventType, data) {
    if (!eventType || !this.eventListeners[eventType]) return

    this.eventListeners[eventType].forEach(listener => {
      try {
        listener(data)
      } catch (error) {
        console.error(`执行事件监听器发生错误 (${eventType}):`, error)
      }
    })
  }

  /**
   * 设置心跳检测
   */
  setupPing() {
    this.clearPingInterval()
    this.clearPongTimeout()
    
    // 每15秒发送一次ping
    this.pingInterval = setInterval(() => {
      if (this.isConnected) {
        this.sendMessage({ type: 'ping' })
        
        // 设置30秒的pong超时
        this.pongTimeout = setTimeout(() => {
          console.warn('WebSocket心跳超时，重新连接')
          this.disconnect()
          this.connect()
        }, 30000)
      }
    }, 15000)
  }

  /**
   * 清除ping间隔
   */
  clearPingInterval() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }
  }

  /**
   * 清除pong超时
   */
  clearPongTimeout() {
    if (this.pongTimeout) {
      clearTimeout(this.pongTimeout)
      this.pongTimeout = null
    }
  }
}

// 导出单例实例
const wsService = new WebSocketService()
export default wsService 