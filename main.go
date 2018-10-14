package main

import (
	"fmt"
	//"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	Token             = "Bot NDk4NzQzNDkxODYyNTkzNTQ2.DpyMUw.E_QptUvTp0bg7OFFxcykRT5gOtM" //"Bot"という接頭辞がないと401 unauthorizedエラーが起きます
	BotName           = "498743491862593546"
	stopBot           = make(chan bool)
	vcsession         *discordgo.VoiceConnection
	HelloWorld        = "!helloworld"
	ChannelVoiceJoin  = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"
	//EggplantRequest   = "!eg"
	//Doutei            = "!eg dou" // WIP
)

func main() {
	//Discordのセッションを作成
	discord, err := discordgo.New()
	discord.Token = Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate) //全てのWSAPIイベントが発生した時のイベントハンドラを追加
	// websocketを開いてlistening開始
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	<-stopBot //プログラムが終了しないようロック
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID) //チャンネル取得
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	//case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, HelloWorld)): //Bot宛に!helloworld コマンドが実行された時
	case strings.HasPrefix(m.Content, "!helloworld"):
		sendMessage(s, c, "Shut up.")

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, ChannelVoiceJoin)):

		//case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, EggplantRequest)):
	case strings.HasPrefix(m.Content, "!eg"):
		var egTimes = strings.TrimLeft(m.Content, "!eg ")
		var egString = ""
		var times int
		times, _ = strconv.Atoi(egTimes)
		for i := 0; i < times; i++ {
			//slow zone... use []byte append
			egString += ":eggplant:"
		}
		sendMessage(s, c, egString)
	case strings.HasPrefix(m.Content, "!eg dou"):
		// WIP

		//今いるサーバーのチャンネル情報の一覧を喋らせる処理を書いておきますね
		//guildChannels, _ := s.GuildChannels(c.GuildID)
		//var sendText string
		//for _, a := range guildChannels{
		//sendText += fmt.Sprintf("%vチャンネルの%v(IDは%v)\n", a.Type, a.Name, a.ID)
		//}
		//sendMessage(s, c, sendText) チャンネルの名前、ID、タイプ(通話orテキスト)をBOTが話す

		//VOICE CHANNEL IDには、botを参加させたい通話チャンネルのIDを代入してください
		//コメントアウトされた上記の処理を使うことでチャンネルIDを確認できます
		//vcsession, _ = s.ChannelVoiceJoin(c.GuildID, "493359894091530246", false, false)
		//vcsession.AddHandler(onVoiceReceived) //音声受信時のイベントハンドラ

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, ChannelVoiceLeave)):
		vcsession.Disconnect() //今いる通話チャンネルから抜ける
	}
}

//メッセージを受信した時の、声の初めと終わりにPrintされるようだ
//func onVoiceReceived(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
//    log.Print("しゃべったあああああ")
//}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}
