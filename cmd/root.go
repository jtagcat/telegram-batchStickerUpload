package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	simpleretry "github.com/jtagcat/simpleretry/pkg"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "telegram-batchStickerUpload",
	Short: "Spam user with stuff",
	Long: `
	This is a quick dirty thing, don't expect anything.

	Filenames will be added in alphabetical order.
	Filename format: foo.✨.webp (emojis used)

	Usage: <sourcedir>`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal().Msg("Expected exactly 1 argument")
		}
		if strings.HasPrefix(args[0], "~/") {
			dirname, _ := os.UserHomeDir()
			args[0] = filepath.Join(dirname, args[0][2:])
		}
		inDir := args[0]

		dirList, err := ioutil.ReadDir(inDir)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read directory")
		}
		var files []fileMin
		for _, file := range dirList {
			if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") { // not dir or hidden file
				extension := path.Ext(file.Name())
				if extension == "" {
					log.Error().Str("filename", file.Name()).Msg("File has no extension, skipping")
					continue
				}
				extLess := strings.TrimSuffix(file.Name(), extension)
				emojiext := path.Ext(extLess)
				if emojiext == "" {
					log.Error().Str("filename", file.Name()).Msg("File has no emoji, skipping")
				}
				emoji := emojiext[1:]
				files = append(files, fileMin{path.Join(inDir, file.Name()), emoji})
			}
		}

		bot, err := tgbotapi.NewBotAPI(os.Getenv("TGAPIKEY"))
		if err != nil {
			log.Fatal().Err(err).
				Msg("Failed to initialize bot; is your TGAPIKEY valid?")
		}

		// bot to bot communication is disallowed, but tgbotapi.NewStickerSetConfig and stuff exists.
		// No idea on how to use it, so we do the following horseshit horseshit:

		update := tgbotapi.NewUpdate(0)
		update.Timeout = 60

		log.Info().Str("bot_link", "https://t.me/"+bot.Self.UserName).Msg("Send any message to bot initiating receiving raw stickers.")
		// bugbug: updates are not flushed
		updates := bot.GetUpdatesChan(update)
		var cid int64
		for u := range updates {
			if u.Message != nil {
				cid = u.Message.Chat.ID
				log.Info().Str("sending_to", u.Message.Chat.UserName).Msg("Sending raw stickers…")
				break
			}
		}

		for _, f := range files {
			backoff := wait.Backoff{
				Duration: time.Second,
				Steps:    5,
				Factor:   1.5,
			}
			serr := simpleretry.OnError(backoff, func() (bool, error) {
				file, err := os.ReadFile(f.path)
				if err != nil {
					return false, err
				}

				_, err = bot.Send(tgbotapi.NewDocument(cid, tgbotapi.FileBytes{
					Name:  path.Base(f.path),
					Bytes: file,
				}))
				return true, err
			})
			if serr != nil {
				log.Error().Err(serr).Msg("error sending file")
			}

			serr = simpleretry.OnError(backoff, func() (bool, error) {
				_, err := bot.Send(tgbotapi.NewMessage(cid, f.emoji))
				return true, err
			})
			if serr != nil {
				log.Error().Err(serr).Msg("error sending emoji")
			}

			time.Sleep(200 * time.Millisecond) // else stuff can get out of order
		}
	},
}

type fileMin struct {
	path  string
	emoji string
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
