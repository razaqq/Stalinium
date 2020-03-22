package gui

import (
	"Stalinium/Installer/bridge"
	"Stalinium/Installer/downloader"
	"Stalinium/Installer/utils"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/multimedia"
	"github.com/therecipe/qt/widgets"
	"os"
)

type StaliniumApp struct {
	app        *widgets.QApplication
	mainWindow *widgets.QMainWindow

	ProgressBar *widgets.QProgressBar
	Installs    map[string]string
}

func (sa *StaliniumApp) Init() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	sa.app = widgets.NewQApplication(len(os.Args), os.Args)

	b := bridge.NewAppBridge(nil)
	b.ConnectProgress(func(percentDone, totalDone, total, speed float64) {
		sa.updateProgress(percentDone, totalDone, total, speed)
	})

	b.ConnectSuccess(func() { sa.showSuccess() })
	b.ConnectError(func(msg string) { sa.showError(msg) })

	sa.mainWindow = widgets.NewQMainWindow(nil, 0)
	sa.mainWindow.SetFixedSize2(333, 200)
	sa.mainWindow.SetWindowTitle("Stalinium Installer")

	centralWidget := widgets.NewQLabel(sa.mainWindow, 0)
	movie := gui.NewQMovie3(":/qml/stalin.gif", core.NewQByteArray2("gif", len("gif")), centralWidget)
	// movie := gui.NewQMovie3("qml/stalin.gif", core.NewQByteArray2("gif", len("gif")), centralWidget)
	movie.SetSpeed(120)
	movie.Start()
	centralWidget.SetMovie(movie)

	playlist := multimedia.NewQMediaPlaylist(nil)
	playlist.AddMedia(multimedia.NewQMediaContent2(core.NewQUrl3("qrc:/qml/hardbass.mp3", core.QUrl__TolerantMode)))
	// playlist.AddMedia(multimedia.NewQMediaContent2(core.NewQUrl3("qml/hardbass.mp3", core.QUrl__TolerantMode)))
	playlist.SetPlaybackMode(multimedia.QMediaPlaylist__Loop)
	player := multimedia.NewQMediaPlayer(nil, multimedia.QMediaPlayer__LowLatency)
	player.SetPlaylist(playlist)
	player.SetVolume(30)
	player.Play()

	centralLayout := widgets.NewQVBoxLayout()
	centralLayout.AddStretch(100)
	centralWidget.SetLayout(centralLayout)
	centralWidget.SetFixedWidth(sa.mainWindow.Width())
	sa.mainWindow.SetCentralWidget(centralWidget)

	selector := widgets.NewQComboBox(centralWidget)
	keys := make([]string, 0, len(sa.Installs))
	for k := range sa.Installs {
		keys = append(keys, k)
	}
	selector.AddItems(keys)
	selector.SetFixedWidth(280)
	centralLayout.AddWidget(selector, 0, core.Qt__AlignCenter)

	button := widgets.NewQPushButton2("Install", centralWidget)
	button.ConnectClicked(func(bool) {
		modDir, err := utils.CreateModDirectory(sa.Installs[selector.CurrentText()])
		if err != nil {
			sa.showError("Failed to create mod directory.")
		}
		go downloader.DownloadMod(modDir, b)
		button.SetDisabled(true)
	})
	centralLayout.AddWidget(button, 0, core.Qt__AlignCenter)

	sa.ProgressBar = widgets.NewQProgressBar(centralWidget)
	sa.ProgressBar.SetTextVisible(true)
	sa.ProgressBar.SetAlignment(core.Qt__AlignCenter)
	sa.ProgressBar.SetFixedWidth(280)
	centralLayout.AddWidget(sa.ProgressBar, 0, core.Qt__AlignCenter)

	sa.mainWindow.Show()
	sa.app.Exec()
}

func (sa *StaliniumApp) updateProgress(percentDone, totalDone, total, speed float64) {
	text := fmt.Sprintf("Downloading %.2f%s (%.1f/%.1f MB) %.2f MB/s", percentDone, "%", totalDone, total, speed)
	sa.ProgressBar.SetValue(int(percentDone))
	sa.ProgressBar.SetFormat(text)
}

func (sa *StaliniumApp) showError(msg string) {
	msgBox := widgets.NewQMessageBox2(widgets.QMessageBox__Critical, "Error", fmt.Sprintf("Error encountered:\n%s", msg), widgets.QMessageBox__Ok, sa.mainWindow, 0)
	msgBox.Button(widgets.QMessageBox__Ok).ConnectClicked(func(bool) { sa.app.Exit(1) })
	msgBox.Exec()
}

func (sa *StaliniumApp) showSuccess() {
	msgBox := widgets.NewQMessageBox2(widgets.QMessageBox__Information, "Ok", "Install successful!", widgets.QMessageBox__Ok, sa.mainWindow, 0)
	msgBox.Button(widgets.QMessageBox__Ok).ConnectClicked(func(bool) { sa.app.Exit(0) })
	msgBox.Exec()
}
