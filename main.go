package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io"
	"io/ioutil"
	"log"
	//"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	//Token             = "Bot NDk4NzQzNDkxODYyNTkzNTQ2.DpyMUw.E_QptUvTp0bg7OFFxcykRT5gOtM" //"Bot"という接頭辞がないと401 unauthorizedエラーが起きます
	Token             = "Bot NDk4NzQzNDkxODYyNTkzNTQ2.Dqdh1w.etTYfhPevzKxzMqjD4eD86S3gM0"
	BotName           = "498743491862593546"
	stopBot           = make(chan bool)
	vcsession         *discordgo.VoiceConnection
	HelloWorld        = "!helloworld"
	ChannelVoiceJoin  = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"
	Folder            = flag.String("f", "sounds", "Folder of files to play.")
	//EggplantRequest   = "!eg"
	//Doutei            = "!eg dou" // WIP
)

func main() {

	var (
		//Token = flag.String("t", "Bot NDk4NzQzNDkxODYyNTkzNTQ2.DpyMUw.E_QptUvTp0bg7OFFxcykRT5gOtM", "Discord token.")
		//GuildID = flag.String("g", "493359894091530242", "Guild ID")
		//GuildID = flag.String("g", "502056641227784204", "Guild ID")
		//ChannelID = flag.String("c", "493359894091530246", "Channel ID")
		//ChannelID = flag.String("c", "502056641680637955", "Channel ID")
		//Folder    = flag.String("f", "sounds", "Folder of files to play.")
		err error
	)

	flag.Parse()
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
	case strings.HasPrefix(m.Content, "!summon"):
		//vcsession, _ = s.ChannelVoiceJoin(c.GuildID, "502056641680637955", false, false)//邪悪な闇の実験場
		vcsession, _ = s.ChannelVoiceJoin(c.GuildID, "493359894091530246", false, false) //ほんちゃん
	case strings.HasPrefix(m.Content, "!s"):
		var fileName = strings.TrimLeft(m.Content, "!s ")
		fileName += ".mp3"
		//var fileName += ".mp3"
		fmt.Println("Reading Folder: ", *Folder)
		files, _ := ioutil.ReadDir("sounds")
		for _, f := range files {
			fmt.Println("PlayAudioFile:", f.Name())
			//
			//fmt.Println("PlayAudioFile:", fileName)
			//s.UpdateStatus(0, f.Name())
			//s.UpdateStatus(0, fileName)
			//fmt.Sprintf("sounds/%s", fileName)
			if fileName == f.Name() {
				PlayAudioFile(vcsession, fmt.Sprintf("%s/%s", *Folder, f.Name()))
			}
			//}
		}
		//vcsession.Close()
		//s.Close()

		//return //とりま

	case strings.HasPrefix(m.Content, "!leave"):
		vcsession.Close()
		//s.Close()

	case strings.HasPrefix(m.Content, "!help"):
		sendMessage(s, c, "!eg <数字> : 指定したナス送信\r\n!summon : VCチャンネルに召喚\r\n!s <ファイル名> : 再生\r\n!leave : VCチャンネルから退出")
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

func PlayAudioFile(v *discordgo.VoiceConnection, filename string) {
	// Send "speaking" packet over the voice websocket
	err := v.Speaking(true)
	if err != nil {
		log.Fatal("Failed setting speaking", err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer v.Speaking(false)

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120
	opts.Volume = 30 //volume

	encodeSession, err := dca.EncodeFile(filename, opts)
	if err != nil {
		log.Fatal("Failed creating an encoding session: ", err)
	}

	done := make(chan error)
	stream := dca.NewStream(encodeSession, v, done)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occured", err)
			}

			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Truncate()
			return
		case <-ticker.C:
			stats := encodeSession.Stats()
			playbackPosition := stream.PlaybackPosition()

			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}
}
