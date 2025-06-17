import React, { useEffect, useState } from "react"

type Player = {
  name: string,
  id: string
}

export const usePlayer = (): [Player | undefined, React.Dispatch<React.SetStateAction<Player | undefined>>] => {
  const [player, setPlayer] = useState<Player | undefined>()

  useEffect(() => {
    localStorage.setItem("player", JSON.stringify(player))
  }, [player])

  return [player, setPlayer];
}
