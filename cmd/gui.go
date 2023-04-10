package main

import (
	"context"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	// Create a new Fyne app and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Dev reload")

	// Set the size of the window to 650x720 pixels
	myWindow.Resize(fyne.NewSize(650, 720))

	// Create a label widget
	myLabel := widget.NewLabel("Reload your docker containers!")
	contianerText := widget.NewLabel("Running containers")

	// Create a list of Docker container cards
	cards := []*widget.Card{}
	selectedFolders := make(map[string]*widget.Label) // define a map to hold the selected folders for each container
	for _, name := range getDockerContainers() {
		// Create a button widget
		nameLabel := name // move the declaration of nameLabel here
		newButton := widget.NewButton("Select Folder", func() {
			// Show a dialog box to select a folder
			dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
				if err != nil {
					println("Error selecting folder:", err)
					return
				}
				if uri == nil {
					println("No folder selected")
					return
				}
				// Update the label with the selected folder path
				selectedFolders[nameLabel].SetText(uri.Path()) // update the selected folder for this container
			}, myWindow)
		})

		// Create a label widget to display the selected folder for this container
		selectedFolder := widget.NewLabel("No folder selected")
		selectedFolders[name] = selectedFolder // add the selected folder label to the map

		// Create a center layout to center the button and label
		// Add the selected folder entry and newButton to a horizontal box container
		horizontalBox := fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			selectedFolder,
			newButton,
		)

		containerDetails := "Container : " + name
		card := widget.NewCard(
			"",
			containerDetails,
			horizontalBox,
		)
		cards = append(cards, card)
	}

	// Add the cards to a grid layout container
	gridLayout := layout.NewGridLayoutWithColumns(1)

	cardContainer := container.New(gridLayout)
	for _, card := range cards {
		cardContainer.Add(card)
	}

	// Add the label, button, and card container to a vertical box container
	myContainer := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		myLabel,
		contianerText,
		cardContainer,
	)
	scrollContainer := container.NewScroll(myContainer)

	// Set the container as the window's content
	myWindow.SetContent(scrollContainer)

	// Show the window and run the app
	myWindow.ShowAndRun()

	// Quit the app when the window is closed
	os.Exit(0)
}
func getDockerContainers() []string {

	// // Create a new Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println(err)
	}

	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Println(err)
	}

	// Extract the container names
	names := []string{}
	for _, container := range containers {
		trimmedName := strings.TrimLeft(container.Names[0], "/")
		names = append(names, trimmedName)
	}

	return names
}
