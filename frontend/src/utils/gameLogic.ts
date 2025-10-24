import { Player, GameStatus, type Position, type BoardState } from '../types/game'

// 棋盘大小
export const BOARD_SIZE = 15

// 创建空棋盘
export function createEmptyBoard(): number[][] {
  return Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(Player.NONE))
}

// 检查位置是否有效
export function isValidPosition(x: number, y: number): boolean {
  return x >= 0 && x < BOARD_SIZE && y >= 0 && y < BOARD_SIZE
}

// 检查位置是否为空
export function isEmpty(board: number[][], x: number, y: number): boolean {
  return isValidPosition(x, y) && board[y][x] === Player.NONE
}

// 检查是否有五子连珠
export function checkWin(board: number[][], x: number, y: number, player: Player): boolean {
  if (player === Player.NONE) return false

  const directions = [
    [1, 0],   // 水平
    [0, 1],   // 垂直
    [1, 1],   // 主对角线
    [1, -1]   // 副对角线
  ]

  for (const [dx, dy] of directions) {
    let count = 1 // 包含当前棋子

    // 向正方向检查
    for (let i = 1; i < 5; i++) {
      const nx = x + dx * i
      const ny = y + dy * i
      if (!isValidPosition(nx, ny) || board[ny][nx] !== player) break
      count++
    }

    // 向负方向检查
    for (let i = 1; i < 5; i++) {
      const nx = x - dx * i
      const ny = y - dy * i
      if (!isValidPosition(nx, ny) || board[ny][nx] !== player) break
      count++
    }

    if (count >= 5) return true
  }

  return false
}

// 检查棋盘是否已满
export function isBoardFull(board: number[][]): boolean {
  for (let y = 0; y < BOARD_SIZE; y++) {
    for (let x = 0; x < BOARD_SIZE; x++) {
      if (board[y][x] === Player.NONE) return false
    }
  }
  return true
}

// 获取游戏状态
export function getGameStatus(board: number[][], lastMove?: Position): GameStatus {
  if (lastMove) {
    const { x, y } = lastMove
    const player = board[y][x]
    
    if (checkWin(board, x, y, player)) {
      return player === Player.HUMAN ? GameStatus.HUMAN_WIN : GameStatus.AI_WIN
    }
  }

  if (isBoardFull(board)) {
    return GameStatus.DRAW
  }

  return GameStatus.PLAYING
}

// 执行移动
export function makeMove(board: number[][], x: number, y: number, player: Player): boolean {
  if (!isEmpty(board, x, y)) return false
  
  board[y][x] = player
  return true
}

// 创建初始游戏状态
export function createInitialGameState(): BoardState {
  return {
    board: createEmptyBoard(),
    currentPlayer: Player.HUMAN,
    gameStatus: GameStatus.PLAYING,
    winner: Player.NONE
  }
}

// 复制棋盘
export function copyBoard(board: number[][]): number[][] {
  return board.map(row => [...row])
}

// 获取所有空位置
export function getEmptyPositions(board: number[][]): Position[] {
  const positions: Position[] = []
  for (let y = 0; y < BOARD_SIZE; y++) {
    for (let x = 0; x < BOARD_SIZE; x++) {
      if (board[y][x] === Player.NONE) {
        positions.push({ x, y })
      }
    }
  }
  return positions
}

// 格式化游戏状态文本
export function formatGameStatus(status: GameStatus): string {
  switch (status) {
    case GameStatus.PLAYING:
      return '游戏进行中'
    case GameStatus.HUMAN_WIN:
      return '玩家获胜！'
    case GameStatus.AI_WIN:
      return 'AI获胜！'
    case GameStatus.DRAW:
      return '平局！'
    default:
      return '未知状态'
  }
}

// 获取当前玩家文本
export function getCurrentPlayerText(player: Player): string {
  switch (player) {
    case Player.HUMAN:
      return '玩家回合'
    case Player.AI:
      return 'AI回合'
    default:
      return ''
  }
}