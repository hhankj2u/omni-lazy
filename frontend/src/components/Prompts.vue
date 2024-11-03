<template>
  <div id="prompts">
    <div class="input-container">
      <input v-model="prompt" placeholder="Enter prompt" />
      <button @click="submitPrompt">Submit</button>
    </div>
    <div class="results-container">
      <p>{{ response }}</p>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  setup() {
    const prompt = ref('');
    const response = ref('');

    const submitPrompt = () => {
      response.value = `You entered: ${prompt.value}`;
    };

    return {
      prompt,
      response,
      submitPrompt
    };
  }
};
</script>

<style scoped>
#prompts {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

.input-container {
  display: flex;
  margin-bottom: 20px;
  justify-content: center;
}

.input-container input {
  margin-right: 10px;
  padding: 10px;
  flex: 1;
}

.results-container {
  display: flex;
  justify-content: center;
  width: 100%;
}
</style>
