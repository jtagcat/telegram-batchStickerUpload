This is slapped and bodged together, don't expect anything.

## howto
1. Prepare files:
   - Anything that [@Stickers](https://t.me/stickers) usually takes. Largest side 512, webp is good.
   - Files will be sent added in alphabetical order
   - Add emoji(s) to be associated with the sticker as a subextension to the name of the file: `foobar.webp` → `foobar.✅.webp`
   - Max stickers per pack: 50, max forwarded messages at once: 100 (50x image+emoji)
1. `TGAPIKEY=5667087765:DYTSC_49ScNQHyUJijg5yXc9uc_5A2rto go run . <path to directory with files to send>`
   - You can get a Telegram bot API key from [@botfather](https://t.me/botfather)
1. Send any message to the bot (and get spammed)
1. Start creating a new pack
   1. Send `/newpack` (or `/addsticker`) to [@Stickers](https://t.me/stickers)
   1. Follow the steps until you reach `Alright! Now send me the sticker`
1. Forward messages from the bot to Stickerbot.
   - First selection shall be a (sticker) image, last should be emoji(s).
1. It will take a minute, it seems like nothing was forwarded.
1. Send `/publish` (or `/done`) to [@Stickers](https://t.me/stickers) and follow the steps to complete pack creation.
