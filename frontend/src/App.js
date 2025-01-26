import React, { useState, useEffect } from 'react';
import useWebSocket from 'react-use-websocket';
import ChatRoom from './components/ChatRoom';
import MessageForm from './components/MessageForm';
import AuthForm from './components/AuthForm';
import SpaceList from './components/SpaceList';
import SpaceForm from './components/SpaceForm';

const App = () => {
  const [token, setToken] = useState(localStorage.getItem('token') || '');
  const [username, setUsername] = useState(localStorage.getItem('username') || '');
  const [messages, setMessages] = useState([]);
  const [selectedSpace, setSelectedSpace] = useState(
    localStorage.getItem('selectedSpace') || null
  );
  const [showSpaceForm, setShowSpaceForm] = useState(false);
  const { sendMessage, lastMessage, readyState } = useWebSocket(
    token && selectedSpace ? `ws://chat-elb-2056070132.ap-northeast-1.elb.amazonaws.com/ws?spaceId=${selectedSpace}` : null,
    { shouldReconnect: () => true }
  );

  // スペースが変更されたらローカルストレージに保存
  useEffect(() => {
    if (selectedSpace) {
      localStorage.setItem('selectedSpace', selectedSpace);
      fetchMessages(selectedSpace); // 選択したスペースに応じたメッセージを取得
    } else {
      localStorage.removeItem('selectedSpace');
      setMessages([]); // 選択解除時にメッセージをクリア
    }
  }, [selectedSpace]);

  // 選択されたスペースのメッセージを取得
  const fetchMessages = async (spaceId) => {
    try {
      const response = await fetch(`${process.env.REACT_APP_URL_DOMAIN}/api/messages?spaceId=${spaceId}`);
      const data = await response.json();
      setMessages(data || []);
    } catch (error) {
      console.error('メッセージの取得エラー:', error);
    }
  };

  // WebSocketで受信したメッセージを処理
  useEffect(() => {
    if (lastMessage !== null) {
      const newMessage = JSON.parse(lastMessage.data);
      setMessages((prev) => (Array.isArray(prev) ? [...prev, newMessage] : [newMessage]));
    }
  }, [lastMessage]);

  const handleAuthSuccess = (authToken, authUsername) => {
    setToken(authToken);
    setUsername(authUsername);
    localStorage.setItem('token', authToken);
    localStorage.setItem('username', authUsername);
  };

  const handleSendMessage = async (text) => {
    if (!selectedSpace) {
      console.error('スペースが選択されていません');
      return;
    }

    try {
      const response = await fetch(`${process.env.REACT_APP_URL_DOMAIN}/api/messages/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          space_id: selectedSpace,
          username: username || '匿名ユーザー',
          text,
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('メッセージ送信エラー:', errorText);
        return;
      }

      const newMessage = await response.json();
      setMessages((prevMessages) => [...prevMessages, newMessage]);
    } catch (error) {
      console.error('通信エラー:', error);
    }
  };

  const handleDeleteMessage = async (id) => {
    if (!id || !selectedSpace) {
      console.error('削除リクエストエラー: 無効なIDまたはスペース未選択');
      alert('削除に失敗しました: 無効なメッセージIDまたはスペースIDです');
      return;
    }

    console.log(`削除リクエスト - メッセージID: ${id}, スペースID: ${selectedSpace}`);

    try {
      setMessages((prev) => prev.filter((message) => message.id !== id));

      const response = await fetch(`${process.env.REACT_APP_URL_DOMAIN}/api/delete?id=${id}&spaceId=${selectedSpace}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('メッセージ削除エラー:', errorText);
        alert(`削除に失敗しました: ${errorText}`);
      } else {
        console.log(`メッセージ削除成功: ID = ${id}`);
      }
    } catch (error) {
      console.error('削除リクエストエラー:', error);
      alert('削除中にエラーが発生しました');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    localStorage.removeItem('selectedSpace');
    setToken('');
    setUsername('');
    setSelectedSpace(null);
  };

  const handleSpaceCreated = () => {
    setShowSpaceForm(false);
  };

  if (!token) {
    return <AuthForm onAuthSuccess={handleAuthSuccess} />;
  }

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="max-w-lg w-full bg-white shadow-lg rounded-lg p-6">
        {!selectedSpace ? (
          <>
            <div className="flex space-x-4 mb-4">
              <button
                onClick={handleLogout}
                className="bg-red-500 text-white px-4 py-2 rounded-lg shadow-md hover:bg-red-600 transition"
              >
                ログアウト
              </button>
              <button
                onClick={() => setShowSpaceForm(true)}
                className="bg-blue-500 text-white px-4 py-2 rounded-lg shadow-md hover:bg-blue-600 transition"
              >
                新規作成
              </button>
            </div>
            {showSpaceForm ? (
              <SpaceForm onSpaceCreated={handleSpaceCreated} />
            ) : (
              <SpaceList onSpaceSelected={setSelectedSpace} />
            )}
          </>
        ) : (
          <>
            <button
              onClick={() => setSelectedSpace(null)}
              className="bg-gray-300 px-4 py-2 rounded-lg hover:bg-gray-400 transition mb-4"
            >
              スペース選択画面に戻る
            </button>
            <ChatRoom messages={messages} onDeleteMessage={handleDeleteMessage} />
            <MessageForm onSendMessage={handleSendMessage} />
          </>
        )}
        <p className="text-sm text-gray-500">
          WebSocketステータス: {readyState === 1 ? '接続中' : '切断'}
        </p>
      </div>
    </div>
  );
};

export default App;
