<template>
  <div id="prompts">
    <div class="input-container">
      <input v-model="prompt" placeholder="Enter prompt" @keyup.enter="submitPrompt" />
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
  justify-content: center;
}

.input-container input {
  margin-right: 10px;
  padding: 10px;
  flex: 1;
}

.results-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    width: 100%;
    justify-content: space-evenly;
    height: calc(100vh - 80px);
    /* Adjust for input/button height */
    overflow-y: auto;
    /* Add scrolling if content exceeds height */
}
</style>
