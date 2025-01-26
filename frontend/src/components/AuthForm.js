import React, { useState } from 'react';

const AuthForm = ({ onAuthSuccess }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLogin, setIsLogin] = useState(true);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    const url = isLogin
      ? `${process.env.REACT_APP_URL_DOMAIN}/api/login`
      : `${process.env.REACT_APP_URL_DOMAIN}/api/register`;

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        if (isLogin) {
          const data = await response.json();
          localStorage.setItem('username', username);
          onAuthSuccess(data.token, username);
        } else {
          alert('ユーザー登録成功！ログインしてください。');
          setIsLogin(true);
        }
      } else {
        const errorMessage = await response.text();
        setError(`エラー: ${errorMessage}`);
      }
    } catch (err) {
      setError('通信エラーが発生しました');
    }
  };

  return (
    <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">
        {isLogin ? 'ログイン' : '新規登録'}
      </h2>
      {error && <p className="text-red-500">{error}</p>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          type="text"
          placeholder="ユーザー名"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-400 focus:outline-none"
        />
        <input
          type="password"
          placeholder="パスワード"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-400 focus:outline-none"
        />
        <button
          type="submit"
          className="w-full bg-blue-500 text-white py-2 rounded-lg shadow-md hover:bg-blue-600 transition"
        >
          {isLogin ? 'ログイン' : '登録'}
        </button>
      </form>
      <button
        onClick={() => {
          setIsLogin(!isLogin);
          setError('');
        }}
        className="text-blue-500 hover:underline mt-4 block"
      >
        {isLogin ? '新規登録はこちら' : 'ログインはこちら'}
      </button>
    </div>
  );
};

export default AuthForm;
