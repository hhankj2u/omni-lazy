<template>
    <div id="prompts">
        <div class="input-container">
            <textarea v-model="echoInput" placeholder="Enter text"></textarea>
            <v-select v-model="selectedPattern" :options="patterns" placeholder="Select a pattern" label="pattern"
                class="pattern-select" />
            <input v-model="language" placeholder="Enter language code" class="language" maxlength="2" />
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
    flex-wrap: wrap;
    gap: 8px;
}

.input-container textarea,
.input-container input,
.input-container .pattern-select {
    padding: 10px;
}

.input-container textarea {
    width: 50vw;
    min-height: 80px;
    resize: vertical;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    font-size: 14px;
}

.input-container .pattern-select {
    width: 300px;
}

.input-container .language {
    width: 50px;
    text-align: center;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
}

@media (max-width: 768px) {
    .input-container {
        flex-direction: column;
        align-items: stretch;
        padding: 12px;
        gap: 12px;
        position: static;
        width: auto;
        margin: 0 0 16px;
    }

    .input-container textarea {
        width: 100%;
        font-size: 16px;
    }

    .input-container .pattern-select {
        width: 100%;
    }

    .input-container .language {
        width: 100%;
    }

    .results-container {
        margin-top: 0;
    }

    .clear-button {
        width: 40px;
        height: 40px;
        font-size: 20px;
    }
}

.results-container {
    margin-top: 120px;
    text-align: left;
    overflow-y: auto;
    padding: 0 12px;
}

.result-item {
    margin-bottom: 20px;
    background: white;
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.input-display {
    background-color: #f0f0f0;
    padding: 10px;
    border-radius: 5px;
}

.result-display {
    margin-top: 10px;
    line-height: 1.5;
}

hr {
    border: 0;
    border-top: 1px solid #e5e7eb;
    margin: 20px 0;
}

.clear-button {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background-color: #2563eb;
    color: white;
    border: none;
    border-radius: 50%;
    width: 50px;
    height: 50px;
    font-size: 24px;
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
}

.clear-button:hover {
    background-color: #1d4ed8;
    transform: scale(1.05);
}
</style>