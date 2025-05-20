package main

import (
	"context"

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
	styleSage := stylesage.NewStyleSage()
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
	return a.styleSage.ProcessText(a.ctx, req)
}

// ReadClipboard reads the clipboard content
func (a *App) ReadClipboard() (string, error) {
	return clipboard.ReadAll()
}
