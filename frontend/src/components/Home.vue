<template>
  <div class="home-container">
    <div class="game-header">
      <h1>五子棋游戏</h1>
      <p class="subtitle">Gomoku Game</p>
    </div>

    <div class="game-modes">
      <div class="mode-card" @click="goToAIGame">
        <div class="mode-icon">🤖</div>
        <h3>人机对战</h3>
        <p>与传统AI进行五子棋对战，支持4个难度级别</p>
        <div class="ai-features">
          <span class="feature-tag">初级</span>
          <span class="feature-tag">中级</span>
          <span class="feature-tag">高级</span>
          <span class="feature-tag">专家</span>
        </div>
        <button class="mode-button">开始游戏</button>
      </div>

      <div class="mode-card llm-card" @click="goToLLMBattle">
        <div class="mode-icon">🧠</div>
        <h3>与AI对战</h3>
        <p>与先进的大语言模型AI对战</p>
        <div class="llm-features">
          <span class="feature-tag">DeepSeek</span>
          <span class="feature-tag">ChatGPT</span>
          <span class="feature-tag">Ollama</span>
        </div>
        <button class="mode-button llm-button">挑战AI</button>
      </div>

      <div class="mode-card" @click="goToPVP">
        <div class="mode-icon">👥</div>
        <h3>双人对战</h3>
        <p>与其他玩家在线对战</p>
        <button class="mode-button">进入房间</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

// 方法
function goToAIGame() {
  router.push('/ai-game')
}

function goToPVP() {
  router.push('/pvp')
}

function goToLLMBattle() {
  router.push('/llm-battle')
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.game-header {
  text-align: center;
  margin-bottom: 40px;
  color: white;
}

.game-header h1 {
  font-size: 3rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
}

.game-modes {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.mode-card {
  background: white;
  border-radius: 20px;
  padding: 40px 30px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  min-width: 250px;
}

.mode-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 40px rgba(0,0,0,0.3);
}

.mode-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.mode-card h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
  color: #333;
}

.mode-card p {
  color: #666;
  margin-bottom: 25px;
  line-height: 1.5;
}

.mode-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 30px;
  border-radius: 25px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.mode-button:hover {
  transform: scale(1.05);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.llm-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: 2px solid #ffd700;
  position: relative;
  overflow: hidden;
}

.llm-card::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(45deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transform: rotate(45deg);
  animation: shimmer 3s infinite;
}

.llm-card h3 {
  color: white;
}

.llm-card p {
  color: rgba(255, 255, 255, 0.9);
}

.llm-features, .ai-features {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.feature-tag {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 0.8rem;
  font-weight: 500;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.llm-button {
  background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
  color: #333;
  font-weight: 600;
  border: none;
  box-shadow: 0 4px 15px rgba(255, 215, 0, 0.3);
}

.llm-button:hover {
  background: linear-gradient(135deg, #ffed4e 0%, #ffd700 100%);
  transform: scale(1.05);
  box-shadow: 0 6px 20px rgba(255, 215, 0, 0.4);
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%) translateY(-100%) rotate(45deg);
  }
  100% {
    transform: translateX(100%) translateY(100%) rotate(45deg);
  }
}

@media (max-width: 768px) {
  .game-modes {
    flex-direction: column;
    align-items: center;
  }
}
</style>