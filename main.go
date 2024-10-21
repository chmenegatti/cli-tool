package main

import (
	"encoding/json"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/go-resty/resty/v2"
	"github.com/rivo/tview"
)

type GitHubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

type App struct {
	client     *resty.Client
	app        *tview.Application
	form       *tview.Form
	textView   *tview.TextView
	inputField *tview.InputField
}

func NewApp() *App {
	return &App{
		client: resty.New(),
		app:    tview.NewApplication(),
	}
}

func (a *App) fetchGitHubUser(username string) (*GitHubUser, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	var user GitHubUser
	resp, err := a.client.R().Get(url)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *App) updateUserInfo(user *GitHubUser) {
	userInfo := fmt.Sprintf(`
Github User Information:

Login: %s
Name: %s
Bio: %s
Location: %s
Followers: %d
Following: %d
`, user.Login, user.Name, user.Bio, user.Location, user.Followers, user.Following)

	a.textView.SetText(userInfo)
}

func (a *App) setupUI() {
	// Setup Form
	a.form = tview.NewForm()
	a.inputField = tview.NewInputField().
		SetLabel("GitHub Username: ").
		SetFieldWidth(30)

	a.form.AddFormItem(a.inputField)
	a.form.AddButton("Search", func() {
		username := a.inputField.GetText()
		user, err := a.fetchGitHubUser(username)
		if err != nil {
			a.textView.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		a.updateUserInfo(user)
	})
	a.form.AddButton("Quit", func() {
		a.app.Stop()
	})

	a.form.SetBorder(true).SetTitle(" Search Form ")

	// Setup TextView
	a.textView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	a.textView.SetBorder(true).SetTitle(" User Information ")

	// Handle Enter key in input field
	a.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			username := a.inputField.GetText()
			user, err := a.fetchGitHubUser(username)
			if err != nil {
				a.textView.SetText(fmt.Sprintf("Error: %v", err))
				return
			}
			a.updateUserInfo(user)
		}
	})

	// Create flex layout
	flex := tview.NewFlex().
		AddItem(a.form, 0, 1, true).
		AddItem(a.textView, 0, 2, false)

	a.app.SetRoot(flex, true)
}

func (a *App) Run() error {
	a.setupUI()
	return a.app.Run()
}

func main() {
	app := NewApp()
	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
	}
}
