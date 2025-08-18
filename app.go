package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hhankj2u/omni-lazy/internal/prompts"
	"github.com/hhankj2u/omni-lazy/internal/stylesage"
	"github.com/hhankj2u/omni-lazy/internal/translators"

	"github.com/atotto/clipboard"
)

// App struct
type App struct {
	ctx         context.Context
	translators *translators.Translators
	prompt      *prompts.Prompt
	styleSage   *stylesage.StyleSage
}

// NewApp creates a new App application struct
func NewApp() *App {
	translators := translators.NewTranslators()
	styleSage, err := stylesage.NewStyleSage()
	if err != nil {
		// Log error but continue with nil styleSage
		// The app will handle this gracefully in ProcessStyleSage
		log.Printf("WARNING: Failed to initialize StyleSage: %s", err.Error())
	}
	return &App{
		translators: translators,
		styleSage:   styleSage,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SearchDictionary(term string) (map[string]string, error) {
	return a.translators.SearchDictionary(term)
}

func (a *App) FetchFabricResult(command string) (string, error) {
	return a.prompt.FetchFabricResult(command)
}

func (a *App) ProcessStyleSage(req stylesage.Request) (string, error) {
	if a.styleSage == nil {
		return "", fmt.Errorf("StyleSage not initialized - check configuration")
	}
	return a.styleSage.ProcessText(a.ctx, req)
}

// ReadClipboard reads the clipboard content
func (a *App) ReadClipboard() (string, error) {
	return clipboard.ReadAll()
}

// GetStyleSageConfig returns the current StyleSage configuration
func (a *App) GetStyleSageConfig() (map[string]interface{}, error) {
	if a.styleSage == nil {
		return nil, fmt.Errorf("StyleSage not initialized")
	}

	cfg := a.styleSage.GetConfig()
	return map[string]interface{}{
		"ollama_url":                        cfg.Ollama.URL,
		"ollama_model":                       cfg.Ollama.Model,
		"ollama_seed_or_negative":           cfg.Ollama.SeedOrNegative,
		"ollama_temperature_if_negative_seed": cfg.Ollama.TemperatureIfNegativeSeed,
		"ollama_pull_timeout":               cfg.Ollama.PullTimeout,
		"ollama_http_timeout":               cfg.Ollama.HTTPTimeout,
		"ollama_trim_space":                 cfg.Ollama.TrimSpace,
		"ollama_verbose":                    cfg.Ollama.Verbose,
	}, nil
}
