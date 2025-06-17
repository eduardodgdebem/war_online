import { useState } from 'react'
import './App.css'
import { useWebSocket } from './hooks/useWebSocket'
import GameMap from './components/GameMap'

type MessageType = "welcome" | "getGameState" | "setIsReady"
const WSURL = "ws://localhost:8080/ws"

function App() {
  const [p, setP] = useState<any>(undefined)
  const [c, setC] = useState<any>(undefined)
  const { connect, disconnect, send, state: wsState } = useWebSocket({
    url: WSURL,
    enabled: false,
    onMessage: (data) => {
      if (!data) return
      const msgType = data.type as MessageType
      switch (msgType) {
        case "welcome":
          setP(data.payload.player)
          break;
        case "getGameState":
          console.log(data.payload)
          setC(data.payload)
          break;
      }
    },
    onOpen() {
      console.log("WebSocket connected")
    },
    onError(error) {
      console.error("WebSocket error:", error)
    },
  });

  const handleGetGameState = () => {
    if (wsState.isConnected) {
      send({ type: "getGameState" });
    } else {
      console.warn("Cannot send - WebSocket not connected");
    }
  };

  return (
    <main>
      <p>Name: {JSON.stringify(p?.name)}</p>
      <GameMap mapData={c?.map}></GameMap>
      <div className='flex flex-row p-2 gap-2'>
        <button
          className='p-1 bg-red-600 rounded-sm'
          onClick={handleGetGameState}
          disabled={!wsState.isConnected}
        >
          Get Game State
        </button>
        {wsState.isConnected ? (
          <button
            className='p-1 bg-yellow-600 rounded-sm'
            onClick={disconnect}
          >
            Disconnect
          </button>
        ) : (
          <button
            className='p-1 bg-green-600 rounded-sm'
            onClick={connect}
          >
            Connect
          </button>
        )}
      </div>
      {wsState.error && <p className="text-red-500">Error: {wsState.error.type}</p>}
    </main>
  )
}

export default App
