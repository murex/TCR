package tcr

import "github.com/daspoet/gowinkey"

func mainMenu() {
	printOptionsMenu()
	keyboardEvents, _ := gowinkey.Listen()

	for event := range keyboardEvents {
		if event.Type == gowinkey.KeyPressed {
			switch event.VirtualKey {
			case gowinkey.VK_D:
				runAsDriver()
			case gowinkey.VK_N:
				runAsNavigator()
			case gowinkey.VK_Q:
				quit()
			}
			printOptionsMenu()
		}
	}
}