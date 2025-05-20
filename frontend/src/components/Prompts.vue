<template>
    <div id="prompts">
        <div class="input-container">
            <div class="input-header">
                <div class="header-left">
                    <h2>Text Processor</h2>
                    <p class="subtitle">Process text using customizable patterns</p>
                </div>
                <div class="header-controls">
                    <textarea v-model="echoInput" placeholder="Enter text"></textarea>
                    <div class="controls-row">
                        <v-select v-model="selectedPattern" :options="patterns" placeholder="Select a pattern" label="pattern"
                            class="pattern-select" />
                        <input v-model="language" placeholder="en" class="language" maxlength="2" />
                        <label class="dry-run-toggle">
                            <input type="checkbox" v-model="dryRun" />
                            <span>dry-run</span>
                        </label>
                        <button @click="submitPrompt" :disabled="loading">
                            <span v-if="loading" class="loading-spinner"></span>
                            {{ loading ? 'Processing...' : 'Submit' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div class="results-container">
            <div v-for="(result, index) in results" :key="index" class="result-item">
                <div class="result-header">
                    <div class="input-display">
                        <strong>Input:</strong><code>{{ result.input }}</code>
                    </div>
                </div>
                <div class="result-display" v-html="result.formattedResponse"></div>
            </div>
        </div>

        <button class="clear-button" @click="clearResults" title="Clear All Results">
            <span class="clear-icon">×</span>
        </button>
    </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { FetchFabricResult } from '../../wailsjs/go/main/App';
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
        const language = ref('en');
        const dryRun = ref(false);
        const patterns = ref([]);
        const results = ref([]);
        const loading = ref(false);
        const fabricPath = ref('');

        const detectOSAndSetFabricPath = () => {
            const userAgent = navigator.userAgent;
            if (userAgent.includes('Win')) {
                fabricPath.value = '%USERPROFILE%\\fabric.exe';
            } else {
                fabricPath.value = 'fabric';
            }
        };

        const fetchPatterns = async () => {
            try {
                const command = `${fabricPath.value} --listpatterns`;
                const result = await FetchFabricResult(command);
                patterns.value = result.split('\n').map(pattern => pattern.trim()).filter(pattern => pattern);
            } catch (error) {
                console.error('Error fetching patterns:', error);
            }
        };

        const submitPrompt = async () => {
            loading.value = true;
            let command = `echo "${echoInput.value}" | ${fabricPath.value} -s -p ${selectedPattern.value} -g=${language.value}`;
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

        onMounted(() => {
            detectOSAndSetFabricPath();
            fetchPatterns();
        });

        return {
            echoInput,
            selectedPattern,
            language,
            dryRun,
            patterns,
            results,
            loading,
            fabricPath,
            submitPrompt,
            clearResults,
        };
    }
};
</script>

<style scoped>
#prompts {
    padding: 0;
    max-width: 1200px;
    margin: 0 auto;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    width: 100%;
    box-sizing: border-box;
    text-align: left;
    overflow-x: hidden;
}

.input-container {
    background: white;
    border-radius: 12px;
    padding: 20px;
    border: 1px solid #e5e7eb;
    width: 100%;
    box-sizing: border-box;
    margin-bottom: 20px;
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
    color: #2563eb;
    font-size: 24px;
    font-weight: 600;
}

.subtitle {
    margin: 4px 0 0;
    color: #64748b;
    font-size: 14px;
}

.header-controls {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.header-controls textarea {
    width: 100%;
    min-height: 100px;
    padding: 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.6;
    resize: vertical;
    background-color: #f9fafb;
    color: #1f2937;
    box-sizing: border-box;
}

.header-controls textarea:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 1px #2563eb;
}

.header-controls textarea::placeholder {
    color: #94a3b8;
}

.controls-row {
    display: flex;
    gap: 8px;
    align-items: center;
}

.pattern-select {
    flex: 1;
    min-width: 200px;
}

:deep(.v-select) {
    background: white;
    border-radius: 6px;
    font-size: 14px;
}

:deep(.vs__dropdown-toggle) {
    border-color: #e5e7eb;
    border-radius: 6px;
    padding: 4px 8px;
    background-color: #f9fafb;
    min-height: 36px;
}

:deep(.vs__dropdown-toggle:hover) {
    border-color: #94a3b8;
}

:deep(.vs__dropdown-toggle:focus) {
    border-color: #2563eb;
    box-shadow: 0 0 0 1px #2563eb;
}

:deep(.vs__selected-options) {
    padding: 2px 0;
    gap: 4px;
}

:deep(.vs__deselect) {
    margin-left: 4px;
    fill: #94a3b8;
    padding: 0;
    height: 14px;
    width: 14px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s ease;
}

:deep(.vs__deselect:hover) {
    background: #f1f5f9;
    fill: #64748b;
}

:deep(.vs__selected) {
    position: relative;
    display: inline-flex;
    align-items: center;
    background: #f1f5f9;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 13px;
    color: #1f2937;
    line-height: 1.2;
}

.language {
    width: 60px;
    padding: 8px 12px;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    font-size: 14px;
    text-align: center;
    background-color: #f9fafb;
}

.language:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 1px #2563eb;
}

.dry-run-toggle {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #64748b;
    font-size: 13px;
    cursor: pointer;
    padding: 4px 0;
    white-space: nowrap;
}

.dry-run-toggle input[type="checkbox"] {
    width: 16px;
    height: 16px;
    border: 1px solid #e5e7eb;
    border-radius: 4px;
    cursor: pointer;
    margin: 0;
}

button {
    padding: 8px 16px;
    background-color: #2563eb;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    white-space: nowrap;
    min-width: 100px;
    justify-content: center;
}

button:hover:not(:disabled) {
    background-color: #1d4ed8;
}

button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.loading-spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 0.8s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.results-container {
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
}

.result-item {
    background: white;
    border-radius: 12px;
    padding: 20px;
    border: 1px solid #e5e7eb;
    width: 100%;
    box-sizing: border-box;
}

.result-header {
    margin-bottom: 16px;
}

.input-display {
    background-color: #f8fafc;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 13px;
    color: #64748b;
    line-height: 1.5;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    border: 1px solid #e5e7eb;
    display: block;
    width: 100%;
    box-sizing: border-box;
}

.input-display strong {
    color: #475569;
    font-weight: 600;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    margin-right: 8px;
}

.input-display code {
    color: #1e40af;
    font-weight: 500;
    white-space: pre-wrap;
    word-wrap: break-word;
}

.result-display {
    font-size: 14px;
    line-height: 1.6;
    color: #1f2937;
}

/* Style JSON-like output */
.result-display > *:not(ul):not(li) {
    display: block;
    width: 100%;
    box-sizing: border-box;
    white-space: pre-wrap;
    word-wrap: break-word;
    overflow-wrap: break-word;
}

.result-display ul {
    margin: 12px 0;
    padding-left: 24px;
    list-style-type: disc;
}

.result-display li {
    margin: 8px 0;
    line-height: 1.5;
    color: #374151;
}

.result-display li::marker {
    color: #94a3b8;
}

.clear-button {
    position: fixed;
    bottom: 20px;
    right: 20px;
    width: 40px;
    height: 40px;
    min-width: 40px;
    padding: 0;
    border-radius: 50%;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
}

.clear-icon {
    font-size: 24px;
    line-height: 1;
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

    .controls-row {
        flex-direction: column;
        width: 100%;
    }

    .pattern-select {
        width: 100%;
        min-width: 0;
    }

    .language {
        width: 100%;
    }

    button {
        width: 100%;
    }

    .result-item {
        padding: 16px;
    }

    .clear-button {
        width: 36px;
        height: 36px;
        min-width: 36px;
        bottom: 16px;
        right: 16px;
    }

    .clear-icon {
        font-size: 20px;
    }

    .header-controls textarea {
        font-size: 16px;
        padding: 12px;
    }

    .input-display {
        padding: 12px;
    }
}
</style>