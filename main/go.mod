module ss-coding

go 1.26.4

require (
	github.com/PuerkitoBio/goquery v1.12.0
	github.com/bwmarrin/discordgo v0.29.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/jonas747/dca v0.0.0-20210930103944-155f5e5f0cc7
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/cloudflare/circl v1.6.3 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jonas747/ogg v0.0.0-20161220051205-b4f6f4cf3757 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
)

replace github.com/jonas747/dca => ./discord/music_bot/dca

replace github.com/bwmarrin/discordgo => github.com/yeongaori/discordgo-fork v0.0.0-20260324114955-7a1c64e5eb96
