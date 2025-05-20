<template>
  <div id="translators">
    <div class="input-container">
      <div class="input-header">
        <div class="header-left">
          <h2>Dictionary Lookup</h2>
          <p class="subtitle">Search multiple dictionaries at once</p>
        </div>
        <div class="header-controls">
          <div class="search-group">
            <input v-model="term" placeholder="Enter word" @keyup.enter="search" />
            <button @click="search">Search</button>
          </div>
          <label class="clipboard-toggle">
            <input type="checkbox" v-model="disableClipboard" />
            <span>Disable Clipboard</span>
          </label>
        </div>
      </div>
    </div>

    <div class="results-container">
      <div v-for="(html, name) in results" :key="name" class="result-item">
        <h2>{{ name }} Dictionary</h2>
        <iframe :srcdoc="formattedResultHtml(html, name)"></iframe>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { SearchDictionary, ReadClipboard } from '../../wailsjs/go/main/App';

export default {
  data() {
    return {
      term: '',
      results: {},
      disableClipboard: ref(true),
      previousClipboardContent: '',
      baseUrls: {
        WEBSTER: 'https://www.merriam-webster.com',
        CAMBRIDGE: 'https://dictionary.cambridge.org',
        SOHA: 'http://tratu.soha.vn'
      }
    };
  },
  computed: {
    formattedResultHtml() {
      return (html, name) => {
        const baseUrl = this.baseUrls[name];
        // Inject a base URL tag if it's not already present in the HTML
        if (!html.includes('<base')) {
          return `<base href="${baseUrl}">` + html;
        }
        return html;
      };
    }
  },
  methods: {
    async search() {
      try {
        this.results = await SearchDictionary(this.term);
      } catch (error) {
        console.error('Error fetching dictionary result:', error);
      }
    },
    async readClipboard() {
      if (!this.disableClipboard) {
        try {
          const text = await ReadClipboard();
          if (text !== this.previousClipboardContent) {
            this.previousClipboardContent = text;
            const words = text.split(/\s+/).slice(0, 3).join(' ');
            console.log('Clipboard contents:', words);
            this.term = words;
            this.search();
          }
        } catch (err) {
          console.log('Failed to read clipboard contents: ', err);
        }
      }
    },
  },
  mounted() {
    setInterval(this.readClipboard, 2000);
  }
};
</script>

<style scoped>
#translators {
  padding: 0;
  max-width: 1200px;
  margin: 0 auto;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  width: 100%;
  box-sizing: border-box;
  background-color: #1e2227;
  color: #a9b1ba;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.input-container {
  background: #1e2227;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #252931;
  width: 100%;
  box-sizing: border-box;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.input-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
}

.header-left {
  text-align: left;
}

.header-left h2 {
  margin: 0;
  color: #60a5fa;
  font-size: 24px;
  font-weight: 600;
}

.subtitle {
  margin: 4px 0 0;
  color: #636973;
  font-size: 14px;
}

.header-controls {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 300px;
}

.search-group {
  display: flex;
  gap: 8px;
}

.search-group input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #252931;
  border-radius: 6px;
  font-size: 14px;
  background-color: #252931;
  color: #a9b1ba;
}

.search-group input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.search-group button {
  padding: 8px 16px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.search-group button:hover {
  background-color: #2563eb;
}

.clipboard-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #636973;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 0;
}

.clipboard-toggle input[type="checkbox"] {
  width: 16px;
  height: 16px;
  border: 1px solid #252931;
  border-radius: 4px;
  cursor: pointer;
  margin: 0;
  background-color: #252931;
}

.results-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
  box-sizing: border-box;
  flex: 1;
  min-height: 0; /* Important for Firefox */
}

.result-item {
  background: #1e2227;
  border-radius: 12px;
  border: 1px solid #252931;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  height: auto;
  min-height: 0; /* Important for Firefox */
}

.result-item h2 {
  margin: 0;
  padding: 16px 20px;
  font-size: 16px;
  font-weight: 600;
  color: #a9b1ba;
  background-color: #252931;
  border-bottom: 1px solid #252931;
  flex-shrink: 0;
}

iframe {
  flex: 1;
  border: none;
  width: 100%;
  height: 100%;
  min-height: 0; /* Important for Firefox */
  background: white;
}

@media (max-width: 768px) {
  .input-container {
    padding: 16px;
  }

  .input-header {
    flex-direction: column;
    gap: 16px;
  }

  .header-left h2 {
    font-size: 20px;
  }

  .header-controls {
    width: 100%;
    min-width: 0;
  }

  .search-group {
    flex-direction: column;
  }

  .search-group input {
    width: 100%;
    font-size: 16px;
  }

  .search-group button {
    width: 100%;
  }

  .results-container {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .result-item h2 {
    padding: 12px 16px;
  }
}
</style>
