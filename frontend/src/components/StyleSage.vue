<template>
    <div id="stylesage">
        <div class="input-container">
            <div class="input-header">
                <h2>Text Assistant</h2>
                <p class="subtitle">Enhance your writing with AI-powered suggestions</p>
            </div>

            <div class="controls">
                <div class="mode-select control-group">
                    <label>Mode:</label>
                    <select v-model="selectedMode" class="styled-select">
                        <option value="quick_fix">✍️ Quick Fix (Grammar & Spelling)</option>
                        <option value="slack_message">💬 Slack Message</option>
                        <option value="email_reply">📧 Email Reply</option>
                        <option value="daily_chat">🗣️ Daily Chat</option>
                        <option value="meeting_summary">📝 Meeting Summary</option>
                        <option value="status_update">📊 Status Update</option>
                    </select>
                </div>

                <div v-if="selectedMode !== 'quick_fix'" class="style-select control-group">
                    <label>Tone:</label>
                    <select v-model="selectedStyle" class="styled-select">
                        <option value="friendly">😊 Friendly & Warm</option>
                        <option value="neutral">😐 Neutral & Clear</option>
                        <option value="formal">👔 Formal & Professional</option>
                        <option value="enthusiastic">🌟 Enthusiastic & Energetic</option>
                    </select>
                </div>
            </div>

            <div class="text-area-container">
                <textarea
                    v-model="inputText"
                    :placeholder="getPlaceholder()"
                    class="text-input"
                    @input="autoGrow"
                    ref="textarea"
                ></textarea>
                <button
                    @click="processText"
                    :disabled="loading || !inputText.trim()"
                    class="process-button"
                >
                    <span v-if="loading" class="loading-spinner"></span>
                    {{ loading ? 'Processing...' : 'Process Text' }}
                </button>
            </div>
        </div>

        <div v-if="result" class="results-container">
            <div class="result-item">
                <div class="result-header">
                    <h3>Enhanced Text:</h3>
                    <div class="result-actions">
                        <button @click="toggleFormatting" class="format-button" :class="{ active: showFormatted }">
                            <span class="format-icon">🔤</span>
                            {{ showFormatted ? 'Show Plain' : 'Show Formatted' }}
                        </button>
                        <button @click="copyToClipboard" class="copy-button">
                            <span class="copy-icon">📋</span>
                            Copy
                        </button>
                    </div>
                </div>
                <div class="result-text" :class="{ 'formatted': showFormatted }">
                    <template v-if="showFormatted">
                        <div v-for="(block, index) in formattedBlocks" :key="index" class="text-block">
                            <template v-if="block.type === 'main_section'">
                                <div class="main-section">
                                    <h2 class="main-section-title">{{ block.title }}</h2>
                                    <div class="section-content">
                                        <div v-for="(item, iIndex) in block.items" :key="iIndex" class="section-item">
                                            <div v-if="item.type === 'block'" class="content-block">
                                                <div v-if="item.title" class="block-title">{{ item.title }}</div>
                                                <div class="block-content">{{ item.content }}</div>
                                            </div>
                                            <div v-else-if="item.type === 'list'" class="list-block">
                                                <div v-if="item.title" class="list-title">{{ item.title }}</div>
                                                <ul class="list-items">
                                                    <li v-for="(listItem, lIndex) in item.items" :key="lIndex">{{ listItem }}</li>
                                                </ul>
                                            </div>
                                            <div v-else-if="item.type === 'numbered_list'" class="numbered-list-block">
                                                <div v-if="item.title" class="list-title">{{ item.title }}</div>
                                                <ol class="list-items">
                                                    <li v-for="(listItem, lIndex) in item.items" :key="lIndex">{{ listItem }}</li>
                                                </ol>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <template v-else-if="block.type === 'sub_section'">
                                <div class="sub-section">
                                    <h3 class="sub-section-title">{{ block.title }}</h3>
                                    <div class="section-content">
                                        <div v-for="(item, iIndex) in block.items" :key="iIndex" class="section-item">
                                            <div v-if="item.type === 'block'" class="content-block">
                                                <div v-if="item.title" class="block-title">{{ item.title }}</div>
                                                <div class="block-content">{{ item.content }}</div>
                                            </div>
                                            <div v-else-if="item.type === 'list'" class="list-block">
                                                <div v-if="item.title" class="list-title">{{ item.title }}</div>
                                                <ul class="list-items">
                                                    <li v-for="(listItem, lIndex) in item.items" :key="lIndex">{{ listItem }}</li>
                                                </ul>
                                            </div>
                                            <div v-else-if="item.type === 'numbered_list'" class="numbered-list-block">
                                                <div v-if="item.title" class="list-title">{{ item.title }}</div>
                                                <ol class="list-items">
                                                    <li v-for="(listItem, lIndex) in item.items" :key="lIndex">{{ listItem }}</li>
                                                </ol>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </div>
                    </template>
                    <template v-else>
                        {{ result }}
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed } from 'vue';
import { ProcessStyleSage } from '../../wailsjs/go/main/App';

export default {
    setup() {
        const inputText = ref('');
        const selectedMode = ref('quick_fix');
        const selectedStyle = ref('friendly');
        const context = ref('');
        const result = ref('');
        const loading = ref(false);
        const textarea = ref(null);
        const showFormatted = ref(true);

        const getPlaceholder = () => {
            const placeholders = {
                'quick_fix': 'Enter your text here for grammar and spelling check...',
                'slack_message': 'Type your Slack message content...',
                'email_reply': 'Paste the email you want to reply to...',
                'daily_chat': 'Enter the conversation or message you want to respond to...',
                'meeting_summary': 'Paste your meeting notes here...',
                'status_update': 'Enter the details for your status update...'
            };
            return placeholders[selectedMode.value];
        };

        const autoGrow = () => {
            if (textarea.value) {
                textarea.value.style.height = 'auto';
                textarea.value.style.height = textarea.value.scrollHeight + 'px';
            }
        };

        const copyToClipboard = async () => {
            try {
                await navigator.clipboard.writeText(result.value);
                // You could add a toast notification here
            } catch (err) {
                console.error('Failed to copy text:', err);
            }
        };

        const parseOutput = (text) => {
            const blocks = [];
            const lines = text.split('\n');
            let currentBlock = null;
            let currentSection = null;

            const startNewSection = (title, type = 'section') => {
                if (currentSection) {
                    if (currentBlock) {
                        currentSection.items.push(currentBlock);
                    }
                    blocks.push(currentSection);
                } else if (currentBlock) {
                    blocks.push(currentBlock);
                }
                currentSection = {
                    type: type,
                    title: title.replace(/^\*|\*$/g, '').trim(),  // Remove markdown asterisks
                    items: []
                };
                currentBlock = null;
            };

            const startNewBlock = (title, content = '') => {
                if (currentBlock) {
                    if (currentSection) {
                        currentSection.items.push(currentBlock);
                    } else {
                        blocks.push(currentBlock);
                    }
                }
                currentBlock = {
                    type: 'block',
                    title: title.replace(/^\*\*|\*\*$/g, '').trim(),  // Remove markdown bold
                    content: content.trim()
                };
            };

            lines.forEach(line => {
                const trimmedLine = line.trim();

                if (!trimmedLine) {
                    // Empty line handling
                    if (currentBlock && currentBlock.content) {
                        if (currentSection) {
                            currentSection.items.push(currentBlock);
                        } else {
                            blocks.push(currentBlock);
                        }
                        currentBlock = null;
                    }
                    return;
                }

                // Handle different section patterns
                if (trimmedLine.match(/^\*[^*].*[^*]\*$/)) {
                    // Main section headers (e.g., *Meeting Summary*)
                    startNewSection(trimmedLine, 'main_section');
                } else if (trimmedLine.match(/^\*\*[^*].*[^*]\*\*$/)) {
                    // Subsection headers (e.g., **Key Decisions**)
                    if (trimmedLine.endsWith(':**')) {
                        // This is actually a block title
                        const title = trimmedLine.replace(/^\*\*|\*\*$/g, '');
                        startNewBlock(title);
                    } else {
                        startNewSection(trimmedLine, 'sub_section');
                    }
                } else if (trimmedLine.match(/^\d+\.\s+\*\*.*:\*\*/)) {
                    // Numbered items with bold headers (e.g., 1. **Greeting:**)
                    const [, title, ...content] = trimmedLine.match(/^\d+\.\s+\*\*(.*?):\*\*(.*)$/);
                    startNewBlock(title, content.join(''));
                } else if (trimmedLine.match(/^[-•]\s+\*\*.*:\*\*/)) {
                    // Bullet points with bold headers (e.g., - **Challenge:** description)
                    const [, title, ...content] = trimmedLine.match(/^[-•]\s+\*\*(.*?):\*\*(.*)$/);
                    startNewBlock(title, content.join(''));
                } else if (trimmedLine.match(/^[-•]\s+/)) {
                    // Regular bullet points
                    if (!currentBlock || currentBlock.type !== 'list') {
                        startNewBlock('List Items');
                        currentBlock.type = 'list';
                        currentBlock.items = [];
                    }
                    const item = trimmedLine.replace(/^[-•]\s+/, '').trim();
                    if (item) {
                        currentBlock.items.push(item);
                    }
                } else if (trimmedLine.match(/^\d+\.\s+/)) {
                    // Regular numbered items
                    if (!currentBlock || currentBlock.type !== 'numbered_list') {
                        startNewBlock('Numbered Items');
                        currentBlock.type = 'numbered_list';
                        currentBlock.items = [];
                    }
                    const item = trimmedLine.replace(/^\d+\.\s+/, '').trim();
                    if (item) {
                        currentBlock.items.push(item);
                    }
                } else if (currentBlock) {
                    // Content lines
                    if (currentBlock.type === 'list' || currentBlock.type === 'numbered_list') {
                        const item = trimmedLine.trim();
                        if (item) {
                            currentBlock.items.push(item);
                        }
                    } else {
                        const content = trimmedLine.trim();
                        if (content) {
                            currentBlock.content += (currentBlock.content ? ' ' : '') + content;
                        }
                    }
                } else {
                    // Plain text content
                    const content = trimmedLine.trim();
                    if (content) {
                        startNewBlock('Content', content);
                    }
                }
            });

            // Handle any remaining blocks
            if (currentBlock) {
                if (currentSection) {
                    currentSection.items.push(currentBlock);
                } else {
                    blocks.push(currentBlock);
                }
            }
            if (currentSection) {
                blocks.push(currentSection);
            }

            return blocks;
        };

        const formattedBlocks = computed(() => {
            return parseOutput(result.value);
        });

        const toggleFormatting = () => {
            showFormatted.value = !showFormatted.value;
        };

        const processText = async () => {
            if (!inputText.value.trim()) return;

            loading.value = true;
            try {
                const request = {
                    text: inputText.value,
                    mode: selectedMode.value,
                    style: selectedStyle.value,
                    context: context.value
                };

                const response = await ProcessStyleSage(request);
                result.value = response;
            } catch (error) {
                console.error('Error processing text:', error);
                result.value = 'Error: ' + error.message;
            } finally {
                loading.value = false;
            }
        };

        return {
            inputText,
            selectedMode,
            selectedStyle,
            context,
            result,
            loading,
            textarea,
            showFormatted,
            formattedBlocks,
            getPlaceholder,
            autoGrow,
            copyToClipboard,
            toggleFormatting,
            processText
        };
    }
};
</script>

<style scoped>
#stylesage {
    padding: 20px;
    max-width: 1000px;
    margin: 0 auto;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

.input-header {
    text-align: center;
    margin-bottom: 24px;
}

.input-header h2 {
    margin: 0;
    color: #2563eb;
    font-size: 28px;
}

.subtitle {
    margin: 8px 0 0;
    color: #64748b;
    font-size: 16px;
}

.input-container {
    background: white;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
}

.controls {
    display: flex;
    gap: 20px;
    flex-wrap: wrap;
    margin-bottom: 20px;
}

.control-group {
    flex: 1;
    min-width: 200px;
}

.control-group label {
    display: block;
    margin-bottom: 8px;
    color: #4b5563;
    font-weight: 500;
}

.styled-select {
    width: 100%;
    padding: 10px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    background-color: white;
    cursor: pointer;
    transition: all 0.2s;
}

.styled-select:hover {
    border-color: #2563eb;
}

.text-area-container {
    position: relative;
}

.text-input {
    width: 100%;
    min-height: 150px;
    padding: 16px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 16px;
    line-height: 1.5;
    resize: none;
    transition: all 0.2s;
    margin-bottom: 16px;
    text-align: left;
}

.text-input:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.process-button {
    width: 100%;
    padding: 12px 24px;
    background-color: #2563eb;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.process-button:hover:not(:disabled) {
    background-color: #1d4ed8;
}

.process-button:disabled {
    background-color: #93c5fd;
    cursor: not-allowed;
}

.loading-spinner {
    width: 20px;
    height: 20px;
    border: 3px solid #ffffff;
    border-radius: 50%;
    border-top-color: transparent;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.results-container {
    margin-top: 24px;
}

.result-item {
    background: white;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
}

.result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    text-align: left;
}

.result-header h3 {
    margin: 0;
    color: #1f2937;
    font-size: 18px;
    text-align: left;
}

.result-actions {
    display: flex;
    gap: 8px;
}

.format-button {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    background-color: #f3f4f6;
    border: none;
    border-radius: 6px;
    color: #4b5563;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
}

.format-button:hover {
    background-color: #e5e7eb;
}

.format-button.active {
    background-color: #2563eb;
    color: white;
}

.copy-button {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    background-color: #f3f4f6;
    border: none;
    border-radius: 6px;
    color: #4b5563;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
}

.copy-button:hover {
    background-color: #e5e7eb;
}

.copy-icon {
    font-size: 16px;
}

.result-text {
    white-space: pre-wrap;
    font-size: 16px;
    line-height: 1.6;
    color: #1f2937;
    background-color: #f9fafb;
    padding: 16px;
    border-radius: 6px;
    text-align: left;
}

.result-text.formatted {
    white-space: normal;
}

.text-block {
    margin-bottom: 20px;
    padding: 0 4px;
}

.text-block:last-child {
    margin-bottom: 0;
}

.block-header {
    font-weight: 600;
    color: #1f2937;
    font-size: 1.15em;
    margin: 24px 0 16px;
    padding-bottom: 8px;
    border-bottom: 2px solid #e5e7eb;
    text-align: left;
}

.block-option {
    margin: 12px 0;
    padding: 16px;
    background-color: #ffffff;
    border-radius: 8px;
    border: 1px solid #e5e7eb;
    text-align: left;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.option-title {
    font-weight: 600;
    color: #2563eb;
    margin-bottom: 8px;
    text-align: left;
    font-size: 1.05em;
}

.option-content {
    color: #4b5563;
    line-height: 1.6;
    text-align: left;
    padding-left: 8px;
    border-left: 2px solid #e5e7eb;
}

.block-list {
    margin: 8px 0;
    padding-left: 16px;
    text-align: left;
}

.list-item {
    margin: 4px 0;
    color: #4b5563;
    line-height: 1.5;
    text-align: left;
    position: relative;
    padding-left: 8px;
}

.block-text {
    color: #4b5563;
    line-height: 1.6;
    margin: 8px 0;
    text-align: left;
}

.block-separator {
    height: 1px;
    background-color: #e5e7eb;
    margin: 24px 0;
}

.option-group {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 32px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.option-group-title {
    font-size: 1.25em;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 20px;
    padding-bottom: 12px;
    border-bottom: 2px solid #e5e7eb;
}

.option-sections {
    display: grid;
    gap: 20px;
}

.option-section {
    background: #f9fafb;
    padding: 16px;
    border-radius: 8px;
    border-left: 3px solid #2563eb;
}

.section-title {
    font-weight: 600;
    color: #2563eb;
    margin-bottom: 8px;
    font-size: 1.1em;
}

.section-content {
    color: #4b5563;
    line-height: 1.6;
    white-space: pre-wrap;
}

.standalone-section {
    background: white;
    padding: 20px;
    border-radius: 8px;
    margin-bottom: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.main-section {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 32px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.main-section-title {
    font-size: 1.4em;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 20px;
    padding-bottom: 12px;
    border-bottom: 2px solid #e5e7eb;
}

.sub-section {
    background: #f8fafc;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 24px;
}

.sub-section-title {
    font-size: 1.2em;
    font-weight: 600;
    color: #2563eb;
    margin: 0 0 16px;
}

.section-content {
    display: grid;
    gap: 16px;
}

.section-item {
    background: white;
    border-radius: 8px;
    padding: 16px;
    border-left: 3px solid #2563eb;
}

.content-block {
    display: grid;
    gap: 8px;
}

.block-title {
    font-weight: 600;
    color: #2563eb;
    font-size: 1.1em;
}

.block-content {
    color: #4b5563;
    line-height: 1.6;
}

.list-block, .numbered-list-block {
    display: grid;
    gap: 8px;
}

.list-title {
    font-weight: 500;
    color: #4b5563;
}

.list-items {
    margin: 0;
    padding-left: 24px;
    color: #4b5563;
    line-height: 1.6;
}

.list-items li {
    margin-bottom: 8px;
}

.list-items li:last-child {
    margin-bottom: 0;
}

@media (max-width: 640px) {
    #stylesage {
        padding: 16px;
    }

    .control-group {
        min-width: 100%;
    }

    .input-container,
    .result-item {
        padding: 16px;
    }
}
</style>