<template>
    <div id="stylesage">
        <div class="input-container">
            <div class="input-header">
                <div class="header-left">
                    <h2>Text Assistant</h2>
                    <p class="subtitle">Enhance your writing with AI-powered suggestions</p>
                </div>
                <div class="header-controls">
                    <div class="mode-select control-group">
                        <select v-model="selectedMode" class="styled-select">
                            <option value="quick_fix">✍️ Quick Fix</option>
                            <option value="slack_message">💬 Slack</option>
                            <option value="email_reply">📧 Email</option>
                            <option value="daily_chat">🗣️ Chat</option>
                            <option value="meeting_summary">📝 Summary</option>
                            <option value="status_update">📊 Status</option>
                        </select>
                    </div>

                    <div v-if="selectedMode !== 'quick_fix'" class="style-select control-group">
                        <select v-model="selectedStyle" class="styled-select">
                            <option value="friendly">😊 Friendly</option>
                            <option value="neutral">😐 Neutral</option>
                            <option value="formal">👔 Formal</option>
                            <option value="enthusiastic">🌟 Energetic</option>
                        </select>
                    </div>

                    <button
                        @click="processText"
                        :disabled="loading || !inputText.trim()"
                        class="process-button"
                    >
                        <span v-if="loading" class="loading-spinner"></span>
                        {{ loading ? 'Processing...' : 'Process (⌃↵)' }}
                    </button>
                </div>
            </div>

            <div class="text-area-container">
                <textarea
                    v-model="inputText"
                    :placeholder="getPlaceholder()"
                    class="text-input"
                    @input="autoGrow"
                    @keydown.ctrl.enter.prevent="processText"
                    ref="textarea"
                ></textarea>
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
import { ref, computed, onMounted, onUnmounted } from 'vue';
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
        const processingTimeout = ref(null);

        // Debounced process function
        const debouncedProcess = () => {
            if (processingTimeout.value) {
                clearTimeout(processingTimeout.value);
            }

            processingTimeout.value = setTimeout(async () => {
                if (!inputText.value.trim() || loading.value) return;

                loading.value = true;
                try {
                    const request = {
                        text: inputText.value,
                        mode: selectedMode.value,
                        style: selectedMode.value === 'quick_fix' ? 'neutral' : selectedStyle.value,
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
            }, 300); // 300ms debounce delay
        };

        // Add keyboard event handler for Mac Command+Enter
        const handleKeydown = (e) => {
            if ((e.metaKey || e.ctrlKey) && e.key === 'Enter' && inputText.value.trim() && !loading.value) {
                e.preventDefault();
                debouncedProcess();
            }
        };

        onMounted(() => {
            window.addEventListener('keydown', handleKeydown);
        });

        onUnmounted(() => {
            window.removeEventListener('keydown', handleKeydown);
            if (processingTimeout.value) {
                clearTimeout(processingTimeout.value);
            }
        });

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
            processText: debouncedProcess
        };
    }
};
</script>

<style scoped>
#stylesage {
    padding: 0;
    max-width: 1200px;
    margin: 0 auto;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    width: 100%;
    box-sizing: border-box;
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
    align-items: center;
    margin-bottom: 16px;
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
    display: flex;
    gap: 8px;
    align-items: center;
}

.control-group {
    min-width: 140px;
}

.styled-select {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    font-size: 13px;
    background-color: white;
    cursor: pointer;
    color: #4b5563;
    appearance: none;
    background-image: url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23464A54%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E");
    background-repeat: no-repeat;
    background-position: right 8px center;
    background-size: 10px;
    padding-right: 24px;
}

.styled-select:hover {
    border-color: #94a3b8;
}

.styled-select:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 1px #2563eb;
}

.text-area-container {
    width: 100%;
    box-sizing: border-box;
}

.text-input {
    width: 100%;
    min-height: 120px;
    padding: 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.5;
    resize: none;
    margin: 16px 0;
    background-color: #f9fafb;
    box-sizing: border-box;
    display: block;
}

.process-button {
    padding: 8px 16px;
    background-color: #2563eb;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    white-space: nowrap;
    height: 36px;
    min-width: 120px;
    transition: all 0.2s ease;
    box-shadow: none;
}

.process-button:hover:not(:disabled) {
    background-color: #1d4ed8;
}

.process-button:disabled {
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
    margin-top: 16px;
}

.result-item {
    margin-top: 20px;
}

.result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.result-header h3 {
    margin: 0;
    font-size: 16px;
}

.result-actions {
    display: flex;
    gap: 6px;
}

.format-button, .copy-button {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
    background-color: #f3f4f6;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    color: #4b5563;
    font-size: 13px;
    transition: all 0.2s ease;
    height: 32px;
}

.format-button:hover, .copy-button:hover {
    background-color: #e5e7eb;
    color: #1f2937;
}

.format-icon, .copy-icon {
    font-size: 14px;
    line-height: 1;
}

.result-text {
    font-size: 14px;
    line-height: 1.5;
    color: #1f2937;
    background-color: white;
    text-align: left;
}

.main-section {
    background: white;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
    border: 1px solid #e5e7eb;
}

.main-section-title {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #e5e7eb;
}

.section-content {
    display: grid;
    gap: 16px;
}

.section-item {
    background: #f9fafb;
    border-radius: 8px;
    padding: 16px;
    position: relative;
}

.section-item::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: #2563eb;
    border-top-left-radius: 8px;
    border-bottom-left-radius: 8px;
}

.block-title {
    font-weight: 600;
    color: #2563eb;
    font-size: 14px;
    margin-bottom: 8px;
}

.block-content {
    color: #4b5563;
    line-height: 1.5;
    font-size: 14px;
}

.list-items {
    margin: 8px 0 0;
    padding-left: 24px;
    color: #4b5563;
}

.list-items li {
    margin-bottom: 6px;
}

.list-items li:last-child {
    margin-bottom: 0;
}

.sub-section {
    background: white;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
    border: 1px solid #e5e7eb;
}

.sub-section-title {
    font-size: 15px;
    font-weight: 600;
    color: #2563eb;
    margin: 0 0 12px;
}

@media (max-width: 640px) {
    #stylesage {
        padding: 0;
    }

    .input-container {
        padding: 16px;
        border-radius: 8px;
        margin-bottom: 16px;
    }

    .header-left h2 {
        font-size: 20px;
    }

    .input-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 16px;
    }

    .header-controls {
        width: 100%;
        flex-wrap: wrap;
        gap: 8px;
    }

    .control-group {
        flex: 1;
        min-width: 140px;
    }

    .process-button {
        width: 100%;
        margin-top: 8px;
    }

    .text-input {
        margin: 12px 0;
        font-size: 16px;
    }

    .main-section {
        padding: 16px;
        margin-bottom: 16px;
        border-radius: 8px;
    }

    .section-item {
        padding: 12px;
    }

    .main-section-title {
        font-size: 15px;
        margin: 0 0 12px;
        padding-bottom: 8px;
    }
}

@media (max-width: 480px) {
    .control-group {
        min-width: 100%;
    }

    .styled-select {
        width: 100%;
    }
}
</style>