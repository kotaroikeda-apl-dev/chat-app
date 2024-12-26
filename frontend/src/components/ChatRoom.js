import React from 'react';

const ChatRoom = ({ messages, onDeleteMessage }) => {
  // メッセージがない場合の表示
  if (!messages || messages.length === 0) {
    return (
      <div className="bg-white shadow-md rounded-lg p-4 mb-6">
        <h2 className="text-xl font-semibold mb-4 text-gray-800">チャットルーム</h2>
        <p className="text-gray-500">メッセージがありません。</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow-md rounded-lg p-4 mb-6">
      <h2 className="text-xl font-semibold mb-4 text-gray-800">チャットルーム</h2>
      <ul className="space-y-2">
        {messages.map((msg) => (
          <li key={msg.id} className="flex justify-between items-center">
            <span className="text-gray-700">
              <strong className="text-blue-500">{msg.username}</strong>: {msg.text}
            </span>
            <button
              onClick={() => {
                console.log('削除対象メッセージID:', msg.id); // デバッグ用ログ
                if (!msg.id) {
                  console.error('無効なメッセージID: ', msg);
                  alert('削除に失敗しました: 無効なメッセージIDです');
                } else if (window.confirm('このメッセージを削除しますか？')) {
                  onDeleteMessage(msg.id);
                }
              }}
              className="bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600 transition"
            >
              削除
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ChatRoom;
