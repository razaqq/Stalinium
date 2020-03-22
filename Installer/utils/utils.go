package utils

import (
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type gameInfo struct {
	Game game `xml:"game"`
}

type game struct {
	VersionName string `xml:"version_name"`
}

func getGameVersion(warshipsDirectory string) (string, error) {
	gameInfoBytes, err := ioutil.ReadFile(path.Join(warshipsDirectory, "game_info.xml"))
	if err != nil {
		return "Cannot find xml file to read config, make sure you selected correct WoWs path", err
	}

	gameInfo := gameInfo{}
	err = xml.Unmarshal(gameInfoBytes, &gameInfo)
	if err != nil {
		return "Error while trying to parse game info xml.", err
	}
	i := strings.LastIndex(gameInfo.Game.VersionName, ".")
	return gameInfo.Game.VersionName[:i], nil
}

func GetWarshipsInstalls() (map[string]string, error) {
	cmd := "Get-ItemProperty -Path " +
		"Registry::HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*, " +
		"Registry::HKEY_CURRENT_USER\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*, " +
		"Registry::HKEY_LOCAL_MACHINE\\SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*, " +
		"Registry::HKEY_CURRENT_USER\\SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*" +
		" | Select-Object -Property DisplayName, InstallLocation" +
		" | Where-Object DisplayName -ne $null" +
		" | Where-Object DisplayName -like *warships*" +
		" | Sort-Object -Property DisplayName" +
		" | ConvertTo-CSV -NoType"
	out, err := exec.Command("powershell", cmd).Output()
	if err != nil {
		return nil, err
	}

	entries := map[string]string{}

	r := csv.NewReader(strings.NewReader(string(out)))
	records, err := r.ReadAll()
	for i, record := range records {
		if i == 0 {
			continue
		} // exclude header
		entries[record[0]] = record[1]
	}
	return entries, nil
}

func CreateModDirectory(warshipsDirectory string) (string, error) {
	gameVersion, err := getGameVersion(warshipsDirectory)
	if err != nil {
		return "", err
	}
	modDirectory := path.Join(warshipsDirectory, "res_mods", gameVersion, "banks", "mods", "Stalinium")

	if _, err = os.Stat(modDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(modDirectory, 0755)
		if err != nil {
			return "", err
		}
	}

	return modDirectory, nil
}
