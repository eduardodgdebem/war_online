import { useEffect } from "react";

// Type definitions
export interface Node {
  x: number;
  y: number;
  strengh: number;
  ownerId: number;
}

export interface GameState {
  nodes: Node[][];
}

export interface WarGameMapProps {
  mapData: GameState;
}

const GameMap: React.FC<WarGameMapProps> = ({ mapData }) => {
  // Define colors for different owners
  const ownerColors: Record<number, string> = {
    0: '#cccccc', // Neutral
    1: '#ff9999', // Player 1
    2: '#9999ff', // Player 2
    3: '#99ff99', // Player 3
    4: '#ffcc99', // Player 4
  };

  useEffect(() => console.log(mapData), [mapData])

  if (!mapData) return <></>

  return (
    <div style={{
      display: 'grid',
      gridTemplateColumns: `repeat(${mapData.nodes.length}, 60px)`,
      gap: '4px',
      margin: '20px'
    }}>
      {mapData.nodes.flatMap(row => row).map((node) => (
        <div
          key={`${node.x}-${node.y}`}
          style={{
            width: '60px',
            height: '60px',
            backgroundColor: ownerColors[node.ownerId] || '#cccccc',
            border: '1px solid #999',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            borderRadius: '4px',
            fontSize: '12px'
          }}
        >
          <div>{`(${node.x},${node.y})`}</div>
          <div style={{ fontWeight: 'bold' }}>{node.strengh}</div>
        </div>
      ))}
    </div>
  );
};

export default GameMap
