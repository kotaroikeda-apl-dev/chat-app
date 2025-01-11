import React, { useState } from 'react';

const SpaceForm = ({ onSpaceCreated }) => {
  const [spaceName, setSpaceName] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!spaceName.trim()) return;

    try {
      const response = await fetch('http://chat-elb-2056070132.ap-northeast-1.elb.amazonaws.com/spaces', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: spaceName }),
      });
      if (response.ok) {
        onSpaceCreated();
        setSpaceName('');
        alert('スペースが作成されました');
      } else {
        console.error('スペース作成エラー:', await response.text());
      }
    } catch (error) {
      console.error('スペース作成リクエストエラー:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex items-center space-x-3">
      <input
        type="text"
        value={spaceName}
        onChange={(e) => setSpaceName(e.target.value)}
        placeholder="スペース名を入力..."
        className="flex-grow px-4 py-2 border rounded-lg shadow-sm focus:ring-2 focus:ring-blue-400 focus:outline-none"
      />
      <button
        type="submit"
        className="bg-blue-500 text-white px-4 py-2 rounded-lg shadow-md hover:bg-blue-600 transition"
      >
        作成
      </button>
    </form>
  );
};

export default SpaceForm;
