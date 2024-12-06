package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	token := os.Getenv("TOKEN")
	guildId := os.Getenv("GUILD_ID")
	channelId := os.Getenv("CHANNEL_ID")
	roleId := os.Getenv("ROLE_ID")
	logId := os.Getenv("LOG_ID")
	joinMsg := os.Getenv("JOIN_MSG")

	log.Println("starting up")

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err.Error())
	}

	discord.Identify.Intents = discordgo.IntentsAll

	joinHandler := func(_ *discordgo.Session, join *discordgo.GuildMemberAdd) {
		log.Printf("%s joined %s\n", join.DisplayName(), guildId)

		if join.GuildID != guildId {
			return
		}

		log.Printf("%s join event is relevant", join.DisplayName())

		uc, err := discord.UserChannelCreate(join.User.ID)
		if err != nil {
			log.Println(err.Error())
			discord.ChannelMessageSend(logId, fmt.Sprintf("error while messaging <@%s>: %s", join.User.ID, err.Error()))
			return
		}

		_, err = discord.ChannelMessageSend(uc.ID, joinMsg)
		if err != nil {
			log.Println(err.Error())
			discord.ChannelMessageSend(logId, fmt.Sprintf("error while messaging <@%s>: %s", join.User.ID, err.Error()))
			return
		}
	}

	messageHandler := func(_ *discordgo.Session, message *discordgo.MessageCreate) {
		if message == nil || message.Member == nil || message.Author == nil {
			return
		}

		log.Printf("%s - %s: %s\n", message.ID, message.Author.Username, message.Content)

		if message.ChannelID != channelId {
			return
		}

		if hasRole(message.Member, roleId) {
			return
		}

		log.Printf("%s is relevant, granting role\n", message.ID)

		nRoles := append(message.Member.Roles, roleId)
		_, err := discord.GuildMemberEdit(guildId, message.Author.ID, &discordgo.GuildMemberParams{
			Roles: &nRoles,
		})
		if err != nil {
			log.Println(err.Error())
			discord.ChannelMessageSend(logId, fmt.Sprintf("error while granting role to <@%s>: %s", message.Author.ID, err.Error()))
			return
		}

		discord.ChannelMessageSend(logId, fmt.Sprintf("granted role to <@%s>", message.Author.ID))
	}

	discord.AddHandler(joinHandler)
	discord.AddHandler(messageHandler)

	healthCheckServer()

	if err := discord.Open(); err != nil {
		log.Fatalln(err.Error())
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	log.Println("shutting down")

	if err := discord.Close(); err != nil {
		log.Fatalln(err.Error())
	}
}

func hasRole(member *discordgo.Member, roleId string) bool {
	for _, role := range member.Roles {
		if role == roleId {
			return true
		}
	}
	return false
}

func healthCheckServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Service is healthy.")
	})

	go func() {
		addr := "0.0.0.0:80"
		log.Printf("starting health check service on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("service failed: %s\n", err)
		}
	}()
}
