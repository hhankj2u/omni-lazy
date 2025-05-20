package stylesage

import (
	"context"
	"fmt"

	"github.com/xyproto/ollamaclient/v2"
)

type StyleSage struct {
	client *ollamaclient.Config
}

type Request struct {
	Text    string `json:"text"`
	Mode    string `json:"mode"`
	Style   string `json:"style"`
	Context string `json:"context"`
}

func NewStyleSage() *StyleSage {
	// Initialize the Ollama client with the model
	client := ollamaclient.New("gemma3:4b")

	// Set reproducible output for consistent results
	client.SetReproducible()

	return &StyleSage{
		client: client,
	}
}

func (s *StyleSage) ProcessText(ctx context.Context, req Request) (string, error) {
	if req.Text == "" {
		return "", fmt.Errorf("input text cannot be empty")
	}

	// Common English language instruction prefix
	englishPrefix := `Process and respond in fluent, professional English only. If the input contains non-English text, translate it to English first, then process the request.

`

	var prompt string

	switch req.Mode {
	case "quick_fix":
		prompt = fmt.Sprintf(`%sCorrect any spelling, grammar, or language issues in this text while keeping its original style and meaning:

Input text:
%s

Format your response using markdown with clear sections:

*Original Text*
[Original text in English]

*Corrected Version*
[Corrected English text with minimal necessary changes]

*Changes Made*
- [List each significant change]
- [Explain why each change was made]

*Style Notes*
- **Tone:** [How the original style was preserved]
- **Flow:** [How natural flow was maintained]`, englishPrefix, req.Text)

	case "slack_message":
		prompt = fmt.Sprintf(`%sHelp me write a professional Slack message in English with the following content, making it %s in tone:

Content to convey:
%s

Format your response exactly like this, using markdown:

*Option 1: Most Diplomatic*

1. **Message Text:** [The actual message in English with emoji]
2. **Tone Level:** [How direct/assertive the message is]
3. **Why It Works:** [Brief explanation of why this approach is effective]
4. **Best Used When:** [Specific scenarios where this option works best]

*Option 2: Moderately Direct*

1. **Message Text:** [The actual message in English with emoji]
2. **Tone Level:** [How direct/assertive the message is]
3. **Why It Works:** [Brief explanation of why this approach is effective]
4. **Best Used When:** [Specific scenarios where this option works best]

[Continue for 3-4 options total]

Guidelines:
- Write all responses in clear, professional English
- Start with most diplomatic option, gradually increase directness
- Use markdown formatting exactly as shown above
- Include appropriate emoji naturally within messages
- Keep explanations concise but informative
- Each option must follow the exact format above`, englishPrefix, req.Style, req.Text)

	case "email_reply":
		prompt = fmt.Sprintf(`%sHelp me write a professional email reply in English with a %s tone:

Original message/context:
%s

Format your response using markdown with clear sections:

*Suggested Reply*

1. **Greeting:**
[Appropriate English greeting]

2. **Main Response:**
[Professional English email body]

3. **Closing:**
[Professional English sign-off]

*Style Analysis*
- **Tone:** [Description of the tone used]
- **Key Points Addressed:** [List of main points covered]
- **Professional Elements:** [Notable professional touches]`, englishPrefix, req.Style, req.Text)

	case "daily_chat":
		prompt = fmt.Sprintf(`%sTranslate (if needed) and respond to the following message/conversation in a %s tone:

Input message:
"%s"

Provide your response in this exact format:

*Translation* (if the input is not in English)
[If the input contains non-English text, provide the English translation here. If input is already in English, write "Input is in English"]

*Suggested Response*
[Write a natural, conversational response that directly addresses the input message]

*Response Details*
**Message Analysis:**
- Message Type: [What kind of message this is - question/statement/request etc.]
- Key Points: [Main points or topics from the message]
- Context Clues: [Any contextual information from the message]

**Response Approach:**
- Tone: [How the %s tone was implemented]
- Style: [Conversational elements used]
- Engagement: [How the response maintains dialogue]

*Alternative Phrasings*
**Casual Version:**
[Same response in a more relaxed style]

**Formal Version:**
[Same response in a more professional style]`, englishPrefix, req.Style, req.Text, req.Style)

	case "meeting_summary":
		prompt = fmt.Sprintf(`%sTransform these meeting notes into a clear English summary with a %s tone:

Meeting content:
%s

Format your response using markdown with clear sections:

*Meeting Summary*

**Key Decisions**
- [Major decision 1 in clear English]
- [Major decision 2 in clear English]
[...]

**Action Items**
1. [Action item with owner in English]
2. [Action item with deadline in English]
[...]

**Discussion Points**
- **Topic 1:** [Summary of discussion in English]
- **Topic 2:** [Summary of discussion in English]
[...]

**Next Steps**
1. [Immediate action needed in English]
2. [Follow-up required in English]
[...]`, englishPrefix, req.Style, req.Text)

	case "status_update":
		prompt = fmt.Sprintf(`%sCreate a clear status update in English with a %s tone:

Update details:
%s

Format your response using markdown with clear sections:

*Status Update Summary*

**Progress Highlights**
- [Key achievement 1 in English]
- [Key achievement 2 in English]
[...]

**Current Status**
- Project: [Overall status in English]
- Timeline: [Schedule status in English]
- Resources: [Resource status in English]

**Challenges & Solutions**
- **Challenge:** [Issue description in English]
  **Solution/Plan:** [How it's being addressed in English]

**Next Steps**
1. [Immediate priority in English]
2. [Following actions in English]
[...]

**Support Needed**
- [Any blockers or required assistance in English]
- [Resource requests if any in English]`, englishPrefix, req.Style, req.Text)

	default:
		return "", fmt.Errorf("invalid mode: %s", req.Mode)
	}

	// Pull the model if needed
	if err := s.client.PullIfNeeded(); err != nil {
		return "", fmt.Errorf("failed to pull model: %w", err)
	}

	// Generate the response using Ollama
	output, err := s.client.GetOutput(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return output, nil
}
