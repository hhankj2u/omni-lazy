<script setup>
import { ref } from 'vue';
import Translators from './components/Translators.vue';
import Prompts from './components/Prompts.vue';
import StyleSage from './components/StyleSage.vue';

const activeTab = ref('StyleSage');
const isMenuOpen = ref(false);

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value;
};
</script>

<template>
  <div class="app-container">
    <button class="menu-toggle" @click="toggleMenu">
      <div class="hamburger" :class="{ 'is-open': isMenuOpen }">
        <span></span>
        <span></span>
        <span></span>
      </div>
    </button>

    <div class="menu-backdrop" :class="{ 'is-visible': isMenuOpen }" @click="toggleMenu"></div>

    <nav class="side-menu" :class="{ 'menu-open': isMenuOpen }">
      <div class="menu-items">
        <button
          v-for="tab in ['StyleSage', 'Translators', 'Prompts']"
          :key="tab"
          :class="{ active: activeTab === tab }"
          @click="() => { activeTab = tab; toggleMenu(); }"
        >
          {{ tab }}
        </button>
      </div>
    </nav>

    <main class="tab-content" :class="{ 'menu-open': isMenuOpen }">
      <Translators v-if="activeTab === 'Translators'" />
      <Prompts v-if="activeTab === 'Prompts'" />
      <StyleSage v-if="activeTab === 'StyleSage'" />
    </main>
  </div>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  overflow-x: hidden;
  background-color: #1e1e1e;
}

.app-container {
  width: 100%;
  height: 100vh;
  background-color: #1e1e1e;
  color: #e5e7eb;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  position: relative;
  overflow-x: hidden;
}

.menu-toggle {
  position: fixed;
  left: 12px;
  top: 12px;
  width: 32px;
  height: 32px;
  background: #1a1a1a;
  border: 1px solid #333333;
  border-radius: 6px;
  cursor: pointer;
  color: #808080;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  z-index: 50;
  padding: 0;
}

.menu-toggle:hover {
  color: #e5e7eb;
  background: #333333;
}

.hamburger {
  width: 16px;
  height: 14px;
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.hamburger span {
  display: block;
  width: 100%;
  height: 2px;
  background-color: currentColor;
  border-radius: 2px;
  transition: all 0.2s ease;
}

.hamburger.is-open span:first-child {
  transform: translateY(6px) rotate(45deg);
}

.hamburger.is-open span:nth-child(2) {
  opacity: 0;
}

.hamburger.is-open span:last-child {
  transform: translateY(-6px) rotate(-45deg);
}

.side-menu {
  position: fixed;
  left: -200px;
  top: 0;
  bottom: 0;
  width: 200px;
  background: #1e1e1e;
  border-right: 1px solid #333333;
  transition: all 0.3s ease;
  z-index: 40;
}

.side-menu.menu-open {
  left: 0;
}

.menu-items {
  padding: 60px 12px 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.menu-items button {
  width: 100%;
  padding: 8px 12px;
  text-align: left;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #808080;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.menu-items button:hover {
  background-color: #333333;
  color: #e5e7eb;
}

.menu-items button.active {
  background-color: #2563eb;
  color: white;
}

.tab-content {
  flex: 1;
  padding: 16px;
  transition: all 0.3s ease;
  width: 100%;
  box-sizing: border-box;
  margin-left: 32px;
  background-color: #1e1e1e;
  max-width: calc(100vw - 48px);
}

.tab-content.menu-open {
  margin-left: 200px;
  max-width: calc(100vw - 216px);
}

.menu-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
  z-index: 30;
  backdrop-filter: blur(2px);
}

.menu-backdrop.is-visible {
  opacity: 1;
  visibility: visible;
}

@media (max-width: 768px) {
  .app-container {
    padding-left: 0;
  }

  .side-menu {
    left: -240px;
    width: 240px;
  }

  .side-menu.menu-open {
    left: 0;
  }

  .tab-content {
    margin-left: 0;
    padding: 12px;
    max-width: 100vw;
  }

  .tab-content.menu-open {
    margin-left: 0;
    max-width: 100vw;
  }
}
</style>