import React, { useEffect, useState } from 'react';

const SpaceList = ({ onSpaceSelected }) => {
  const [spaces, setSpaces] = useState([]);

  const fetchSpaces = async () => {
    try {
      const response = await fetch('http://chat-elb-2056070132.ap-northeast-1.elb.amazonaws.com/api/spaces/list');
      const data = await response.json();
      setSpaces(data || []);
    } catch (error) {
      console.error('スペース一覧の取得エラー:', error);
    }
  };

  useEffect(() => {
    fetchSpaces();
  }, []);

  return (
    <div className="space-y-3"> {/* 縦方向に間隔を追加 */}
      <h2 className="text-xl font-semibold">スペース一覧</h2>
      <ul className="space-y-2"> {/* 各スペースの間隔を設定 */}
        {spaces.map((space) => (
          <li key={space.id}>
            <button
              onClick={() => onSpaceSelected(space.id)}
              className="bg-gray-200 px-4 py-2 rounded-lg hover:bg-gray-300 transition w-full text-left"
            >
              {space.name}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default SpaceList;
