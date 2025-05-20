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
                <div class="result-text">
                    <div :class="showFormatted ? 'formatted-text' : 'plain-text'">
                        <pre v-if="!showFormatted" class="plain-content">{{ result }}</pre>
                        <div v-else v-html="formattedResult" class="content"></div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { ProcessStyleSage } from '../../wailsjs/go/main/App';
import { marked } from 'marked';

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

        const formattedResult = computed(() => {
            return marked(result.value);
        });

        const toggleFormatting = () => {
            showFormatted.value = !showFormatted.value;
            console.log('Toggle formatting:', showFormatted.value);
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
            formattedResult,
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
    text-align: left;
    background-color: #1e1e1e;
    color: #e5e7eb;
}

.input-container {
    background: #1e1e1e;
    padding: 16px;
    width: 100%;
    box-sizing: border-box;
    margin-bottom: 16px;
    text-align: left;
}

.input-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.header-left {
    text-align: left;
    flex: 1;
}

.header-left h2 {
    margin: 0;
    color: #60a5fa;
    font-size: 20px;
    font-weight: 600;
    text-align: left;
}

.subtitle {
    margin: 2px 0 0;
    color: #9ca3af;
    font-size: 13px;
}

.header-controls {
    display: flex;
    gap: 6px;
    align-items: center;
}

.control-group {
    min-width: 130px;
}

.styled-select {
    width: 100%;
    padding: 6px 10px;
    border: 1px solid #2c2d32;
    border-radius: 6px;
    font-size: 13px;
    background-color: #1a1b1e;
    cursor: pointer;
    color: #e5e7eb;
    appearance: none;
    background-image: url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23e5e7eb%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E");
    background-repeat: no-repeat;
    background-position: right 8px center;
    background-size: 10px;
    padding-right: 24px;
    height: 32px;
}

.styled-select:hover {
    border-color: #3b82f6;
}

.styled-select:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 1px #3b82f6;
}

.text-area-container {
    width: 100%;
    box-sizing: border-box;
}

.text-input {
    width: 100%;
    min-height: 100px;
    padding: 10px;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    line-height: 1.4;
    resize: none;
    margin: 12px 0 8px;
    background-color: #1a1a1a;
    color: #e5e7eb;
    box-sizing: border-box;
    display: block;
}

.text-input::placeholder {
    color: #6b7280;
}

.process-button {
    padding: 6px 14px;
    background-color: #3b82f6;
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
    height: 32px;
    min-width: 110px;
    transition: all 0.2s ease;
    box-shadow: none;
    cursor: pointer;
}

.process-button:hover:not(:disabled) {
    background-color: #2563eb;
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
    padding: 0 16px;
}

.result-item {
    margin-top: 16px;
}

.result-text {
    font-size: 14px;
    line-height: 1.5;
    color: #e5e7eb;
    text-align: left;
}

.result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
}

.result-header h3 {
    margin: 0;
    color: #e5e7eb;
    font-size: 24px;
    font-weight: 400;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.result-actions {
    display: flex;
    gap: 8px;
}

.format-button, .copy-button {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
    background-color: #2563eb;
    border: none;
    border-radius: 6px;
    color: white;
    font-size: 13px;
    transition: all 0.2s ease;
    height: 32px;
    cursor: pointer;
}

.format-button:hover, .copy-button:hover {
    background-color: #1d4ed8;
}

.format-icon, .copy-icon {
    font-size: 14px;
    line-height: 1;
}

.formatted-text {
    white-space: pre-line;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    line-height: 1.4;
    color: #e5e7eb;
    background-color: #1a1a1a;
    padding: 16px;
    border-radius: 8px;
    margin-top: 12px;
}

.plain-text {
    background-color: #1a1a1a;
    padding: 16px;
    border-radius: 8px;
    margin-top: 12px;
}

.plain-content {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    line-height: 1.5;
    color: #e5e7eb;
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
    font-size: 14px;
}

.content {
    display: flex;
    flex-direction: column;
    gap: 0.75em;
}

.content :deep(p) {
    margin: 0;
    color: #e5e7eb;
}

.content :deep(ul), .content :deep(ol) {
    margin: 0;
    padding-left: 1.25em;
    display: flex;
    flex-direction: column;
    gap: 0.5em;
    color: #e5e7eb;
}

.content :deep(li) {
    margin: 0;
    color: #e5e7eb;
    display: flex;
    gap: 0.5em;
    align-items: baseline;
    line-height: 1.4;
    text-align: match-parent;
}

.content :deep(li)::marker {
    color: #808080;
    unicode-bidi: isolate;
    font-variant-numeric: tabular-nums;
    text-align: match-parent;
    text-transform: none;
    text-indent: 0px !important;
    text-align-last: left !important;
    white-space: pre;
}

.content :deep(strong) {
    font-weight: 600;
    color: #e5e7eb;
}

.content :deep(em) {
    font-style: italic;
    color: #808080;
}

.content :deep(h1), .content :deep(h2), .content :deep(h3) {
    margin: 0;
    font-weight: 600;
    color: #e5e7eb;
}

.content :deep(code) {
    background-color: #1a1a1a;
    padding: 0.2em 0.4em;
    border-radius: 4px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    font-size: 0.9em;
    color: #e5e7eb;
}

.content :deep(pre) {
    background-color: #1a1a1a;
    padding: 1em;
    border-radius: 8px;
    overflow-x: auto;
}

.content :deep(pre code) {
    background-color: transparent;
    padding: 0;
    border-radius: 0;
}

.content :deep(blockquote) {
    margin: 0;
    padding-left: 1em;
    border-left: 3px solid #3b82f6;
    color: #9ca3af;
}

.content :deep(a) {
    color: #60a5fa;
    text-decoration: none;
}

.content :deep(a:hover) {
    text-decoration: underline;
}

.format-button.active {
    background-color: #1a1a1a;
    color: white;
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

.menu-backdrop {
    background-color: rgba(0, 0, 0, 0.4);
}

.side-menu {
    background-color: #1e1e1e;
    border-right: 1px solid #333333;
    color: #e5e7eb;
}

.side-menu a {
    color: #808080;
    text-decoration: none;
}

.side-menu a:hover {
    color: #e5e7eb;
}

.side-menu a.active {
    color: #2563eb;
    background-color: #1a1a1a;
}

.tab-content {
    background-color: #1e1e1e;
}

body {
    background-color: #1e1e1e;
    color: #e5e7eb;
}

#app {
    background-color: #1e1e1e;
    color: #e5e7eb;
}

.menu-item {
    color: #808080;
    padding: 8px 16px;
    text-decoration: none;
    display: block;
    transition: all 0.2s ease;
}

.menu-item:hover {
    color: #e5e7eb;
    background-color: #1a1a1a;
}

.menu-item.active {
    color: #2563eb;
    background-color: #1a1a1a;
}

.menu-header {
    background-color: #2563eb;
    color: white;
    padding: 12px 16px;
    font-weight: 500;
    border-radius: 8px;
    margin: 8px;
}

.menu-close {
    background-color: transparent;
    border: none;
    color: #808080;
    padding: 8px;
    cursor: pointer;
    transition: color 0.2s ease;
}

.menu-close:hover {
    color: #e5e7eb;
}
</style>