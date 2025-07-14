## About
![image](https://github.com/user-attachments/assets/f728f24e-70e8-4b05-8046-e3ce198208cb)

A versatile and flexible application designed to handle a wide variety of tasks or data inputs without the need for specialized tools.
A Swiss Army knife that can perform multiple functions with just one tool.
It pertains to the idea of simplifying task management by centralizing different functionalities in a single, easy-to-use application.

## Features

- Translators: https://github.com/hhankj2u/go-translators
- Prompts: https://github.com/danielmiessler/fabric
- StyleSage: Text Assistant - Enhance your writing with AI-powered suggestions
![image](https://github.com/user-attachments/assets/609c7f37-cbe3-482d-a6fe-c1e30d463c33)

## Configuration

The application uses a configuration file to manage Ollama settings. The config file is automatically created at `~/.omni-lazy/config.json` on first run.

### Configuration Format

```json
{
  "ollama": {
    "url": "http://localhost:11434",
    "model": "gemma3:4b",
    "seed_or_negative": 0,
    "temperature_if_negative_seed": 0.0,
    "pull_timeout": 0,
    "http_timeout": 0,
    "trim_space": true,
    "verbose": false
  }
}
```

### Configuration Options

- **url**: The Ollama server URL (default: `http://localhost:11434`)
- **model**: The Ollama model to use (default: `gemma3:4b`)
- **seed_or_negative**: Seed for reproducible output, or negative for random (default: `0`)
- **temperature_if_negative_seed**: Temperature when seed is negative (default: `0.0`)
- **pull_timeout**: Timeout for model pulling in seconds (default: `0` = no timeout)
- **http_timeout**: HTTP request timeout in seconds (default: `0` = no timeout)
- **trim_space**: Whether to trim whitespace from responses (default: `true`)
- **verbose**: Enable verbose logging (default: `false`)

You can modify the configuration file to use different models, connect to remote Ollama instances, or fine-tune the client behavior.

## Build

To build the application for desktop production, use:

```sh
wails build -tags desktop,production -ldflags "-w -s"
```

## Development

For development with live reload, use the following workflow:

### 1. Start the Wails dev server (backend + desktop app)

This will launch the desktop app with hot reload for Go/backend and frontend changes:

```sh
wails dev
```

### 2. (Optional) Start the frontend dev server only

If you want to work on the frontend (Vue/Vite) in isolation, you can run:

```sh
cd frontend
npm install   # Only needed once
npm run dev
```

This will start the Vite dev server at http://localhost:5173 (or similar), with live reload for Vue components and styles.

### Notes
- When using `wails dev`, both backend and frontend changes are hot-reloaded in the desktop app.
- You can use the browser-based Vite dev server for rapid frontend iteration, but integration with the desktop app is best tested via `wails dev`.

