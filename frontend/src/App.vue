<script setup>
import { ref } from 'vue';
import Translators from './components/Translators.vue';
import Prompts from './components/Prompts.vue';

const activeTab = ref('Translators');
</script>

<template>
  <div>
    <div class="tabs">
      <button :class="{ active: activeTab === 'Translators' }" @click="activeTab = 'Translators'">Translators</button>
      <button :class="{ active: activeTab === 'Prompts' }" @click="activeTab = 'Prompts'">Prompts</button>
    </div>
    <div class="tab-content">
      <Translators v-if="activeTab === 'Translators'" />
      <Prompts v-if="activeTab === 'Prompts'" />
    </div>
  </div>
</template>

<style>
.tabs {
  display: flex;
  justify-content: center;
  /* margin-bottom: 20px; */
}

.tabs button {
  padding: 10px 20px;
  margin: 0 5px;
  cursor: pointer;
}

.tabs button.active {
  background-color: #007bff;
  color: white;
}

.tab-content {
  /* padding: 20px; */
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
