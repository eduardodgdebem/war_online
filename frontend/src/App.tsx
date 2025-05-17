import { useEffect, useRef } from 'react';
import './App.css'
import { usePlayer } from './hooks/usePlayer';
import { useWebSocket } from './hooks/useWebSocket.tsx'

type MessageType = "welcome" | "attack"
const WSURL = "ws://localhost:8080/ws"

function App() {
  const [player, setPlayer] = usePlayer()
  const [connect, _, send] = useWebSocket({
    url: WSURL,
    onMessage: (data) => {
      if (!data) return
      const msgType = data.type as MessageType
      switch (msgType) {
        case "welcome":
          setPlayer(data.player)
          break;
        case 'attack':
          console.log(data)
          break;
      }
    },
    reconnectAttempts: 1
  });


  return (
    <main>
      <p>Name: {player?.name}</p>
      <div className='flex flex-row p-2 gap-2'>
        <button
          className='p-1 bg-red-600 rounded-sm'
          onClick={() => {
            send({ type: "attack", payload: "type" });
          }}>
          click
        </button>
        <button
          className='p-1 bg-green-600 rounded-sm'
          onClick={() => {
            connect()
          }}>
          connect
        </button>
      </div>
    </main>
  )
}

export default App
