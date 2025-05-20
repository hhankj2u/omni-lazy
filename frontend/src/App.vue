<script setup>
import { ref } from 'vue';
import Translators from './components/Translators.vue';
import Prompts from './components/Prompts.vue';
import StyleSage from './components/StyleSage.vue';

const activeTab = ref('Translators');
</script>

<template>
  <div>
    <div class="tabs">
      <button :class="{ active: activeTab === 'Translators' }" @click="activeTab = 'Translators'">Translators</button>
      <button :class="{ active: activeTab === 'Prompts' }" @click="activeTab = 'Prompts'">Prompts</button>
      <button :class="{ active: activeTab === 'StyleSage' }" @click="activeTab = 'StyleSage'">StyleSage</button>
    </div>
    <div class="tab-content">
      <Translators v-if="activeTab === 'Translators'" />
      <Prompts v-if="activeTab === 'Prompts'" />
      <StyleSage v-if="activeTab === 'StyleSage'" />
    </div>
  </div>
</template>

<style>
.tabs {
  display: flex;
  justify-content: center;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  background-color: white;
  z-index: 1000;
  padding: 5px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.tabs button {
  padding: 5px 10px;
  margin: 0 5px;
  cursor: pointer;
  font-size: 14px;
}

.tabs button.active {
  background-color: #007bff;
  color: white;
}

.tab-content {
  margin-top: 40px; /* Adjust this value to match the height of the tabs */
}

#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>