<template>
    <div id="prompts">
      <div class="input-container">
        <textarea v-model="echoInput" placeholder="Enter text"></textarea>
        <v-select
          v-model="selectedPattern"
          :options="patterns"
          placeholder="Select a pattern"
          label="pattern"
          class="pattern-select"
        />
        <input v-model="language" placeholder="Enter language code" class="language" />
        <label>
          <input type="checkbox" v-model="dryRun" />
          dry-run
        </label>
        <button @click="submitPrompt" :disabled="loading">Submit</button>
      </div>
      <div class="results-container">
        <div v-for="(result, index) in results" :key="index" class="result-item">
          <div class="input-display">
            <p><strong>Input:</strong> {{ result.input }}</p>
          </div>
          <div class="result-display" v-html="result.formattedResponse"></div>
          <hr />
        </div>
      </div>
      <button class="clear-button" @click="clearResults">Clear All</button>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { ListPatterns, FetchFabricResult } from '../../wailsjs/go/main/App';
  import vSelect from 'vue-select';
  import 'vue-select/dist/vue-select.css';
  import { marked } from 'marked';
  
  export default {
    components: {
      vSelect,
    },
    setup() {
      const echoInput = ref('');
      const selectedPattern = ref('');
      const language = ref('english');
      const dryRun = ref(false);
      const patterns = ref([]);
      const results = ref([]);
      const loading = ref(false);
  
      const fetchPatterns = async () => {
        try {
          const result = await ListPatterns();
          patterns.value = result;
        } catch (error) {
          console.error('Error fetching patterns:', error);
        }
      };
  
      const submitPrompt = async () => {
        loading.value = true;
        let command = `echo "${echoInput.value}" | fabric -s -p ${selectedPattern.value} -g=${language.value}`;
        if (dryRun.value) {
          command += ' --dry-run';
        }
        try {
          const result = await FetchFabricResult(command);
          results.value.push({
            input: command,
            formattedResponse: marked(result),
          });
        } catch (error) {
          results.value.push({
            input: command,
            formattedResponse: `Error: ${error.message}`,
          });
        } finally {
          loading.value = false;
        }
      };
  
      const clearResults = () => {
        results.value = [];
      };
  
      onMounted(fetchPatterns);
  
      return {
        echoInput,
        selectedPattern,
        language,
        dryRun,
        patterns,
        results,
        loading,
        submitPrompt,
        clearResults,
      };
    }
  };
  </script>
  
  <style scoped>
  #prompts {
    padding: 10px;
  }
  
  .input-container {
    display: flex;
    justify-content: center;
    align-items: center;
    position: fixed;
    background-color: white;
    z-index: 1000;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    width: 96%;
    padding: 10px;
    margin-top: -10px;
  }
  
  .input-container textarea,
  .input-container input,
  .input-container .pattern-select {
    padding: 10px; 
  }
  
  .input-container textarea {
    width: 50vw;
  }
  
  .input-container .pattern-select {
    width: 300px;
  }
  
  .input-container .language {
    width: 50px;
    text-align: center;
  }
  
  .results-container {
    margin-top: 100px;
    text-align: left;
    overflow-y: auto;
  }
  
  .result-item {
    margin-bottom: 20px;
  }
  
  .input-display {
    background-color: #f0f0f0;
    padding: 10px;
    border-radius: 5px;
  }
  
  .result-display {
    margin-top: 10px;
  }
  
  hr {
    border: 0;
    border-top: 1px solid #ccc;
    margin: 20px 0;
  }
  
  .clear-button {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 50%;
    width: 50px;
    height: 50px;
    font-size: 24px;
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
  
  .clear-button:hover {
    background-color: #0056b3;
  }
  </style>