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
  const [selectedSpace, setSelectedSpace] = useState(null); // 選択中のスペース
  const [showSpaceForm, setShowSpaceForm] = useState(false); // スペースフォームの表示状態
  const { sendMessage, lastMessage, readyState } = useWebSocket(
    token && selectedSpace ? `ws://localhost:8080/ws?spaceId=${selectedSpace}` : null,
    { shouldReconnect: () => true }
  );

  // メッセージを初期取得
  useEffect(() => {
    const fetchMessages = async () => {
      if (!selectedSpace) return;

      try {
        const response = await fetch(`http://localhost:8080/messages?spaceId=${selectedSpace}`);
        const data = await response.json();
        setMessages(data || []);
      } catch (error) {
        console.error('メッセージの取得エラー:', error);
      }
    };

    fetchMessages();
  }, [selectedSpace]);

  // WebSocketで受信したメッセージを処理
  useEffect(() => {
    if (lastMessage !== null) {
      const newMessage = JSON.parse(lastMessage.data);
      setMessages((prev) => (Array.isArray(prev) ? [...prev, newMessage] : [newMessage]));
    }
  }, [lastMessage]);

  // 認証成功時の処理
  const handleAuthSuccess = (authToken, authUsername) => {
    setToken(authToken);
    setUsername(authUsername);
    localStorage.setItem('token', authToken);
    localStorage.setItem('username', authUsername);
  };

  // メッセージ送信
  const handleSendMessage = async (text) => {
    if (!selectedSpace) {
      console.error('スペースが選択されていません');
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/messages/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          space_id: selectedSpace,
          username: username || '匿名ユーザー', // ユーザー名
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

  // メッセージ削除
  const handleDeleteMessage = async (id) => {
    if (!id || !selectedSpace) {
      console.error('削除リクエストエラー: 無効なIDまたはスペース未選択');
      alert('削除に失敗しました: 無効なメッセージIDまたはスペースIDです');
      return;
    }

    console.log(`削除リクエスト - メッセージID: ${id}, スペースID: ${selectedSpace}`);

    try {
      setMessages((prev) => prev.filter((message) => message.id !== id));

      const response = await fetch(`http://localhost:8080/delete?id=${id}&spaceId=${selectedSpace}`, {
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

  // ログアウト処理
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    setToken('');
    setUsername('');
  };

  // スペース作成完了時の処理
  const handleSpaceCreated = () => {
    setShowSpaceForm(false); // スペースフォームを閉じる
  };

  if (!token) {
    return <AuthForm onAuthSuccess={handleAuthSuccess} />;
  }

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="max-w-lg w-full bg-white shadow-lg rounded-lg p-6">
        {!selectedSpace ? (
          <>
            <div className="flex space-x-4 mb-4"> {/* ボタン間隔を調整 */}
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
