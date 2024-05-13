# go-micro-tts
利用微软提供的语音服务，可以通过使用该SDK将文本转化为合成语音并获取受支持的声音列表。

## 快速接入

**安装**
```bash
go get -u github.com/xuemingjings/go-micro-tts@latest
```

**使用**

短文本转语音、生成mp3文件、直接写入到本地磁盘
```go
package main

import (
	"context"
	go_micro_tts "github.com/xuemingjings/go-micro-tts"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	speechKey    = "xxx"
	speechRegion = "xxx"
)

func main() {
	ctx := context.TODO()
	tts, err := go_micro_tts.NewGoTTS(
		ctx,
		go_micro_tts.WithSpeechRegion(speechRegion),
		go_micro_tts.WithSpeechKey(speechKey),
	)
	if err != nil {
		panic("初始化报错: " + fmt.Sprintf("%v", err))
	}

	speakXml := go_micro_tts.NewSpeakXml(&go_micro_tts.SpeakXmlReq{
		Lang:   "zh-CN",
		Gender: "Male",
		Name:   "zh-CN-YunxiNeural",
		Text:   "中华兴盛，幸有斌哥。How are you",
	})

	// 创建输出文件
	fileName := time.Now().Format("2006-01-02-15-04-05") + "_file.mp3"
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	err = tts.TextToVoiceDisk(go_micro_tts.Audio48kHz96KbitrateMonoMp3, speakXml, outFile)

	if err != nil {
		log.Printf("TextToVoice调用发生了错误 err:%v \n", err)
		return
	}

	fmt.Println("文本转语音文件已生成", fileName)
}
```

长文本转语音任务创建
```go
package main

import (
	"context"
	go_micro_tts "github.com/xuemingjings/go-micro-tts"
	"fmt"
)

var (
	speechKey    = "xxx"
	speechRegion = "xxx"
)

func main() {
	ctx := context.TODO()
	tts, err := go_micro_tts.NewGoTTS(
		ctx,
		go_micro_tts.WithSpeechRegion(speechRegion),
		go_micro_tts.WithSpeechKey(speechKey),
	)
	if err != nil {
		panic("初始化报错: " + fmt.Sprintf("%v", err))
	}

	inputs := []*go_micro_tts.LongSpeakInputs{
		{
			Text: "中华兴盛，幸有斌哥。How Are you",
		},
	}

	req := go_micro_tts.NewLongSpeak(&go_micro_tts.LongSpeakXmlReq{
		DisplayName:             "长文本输出单个音频文件",
		Inputs:                  inputs,
		OutputFormat:            go_micro_tts.Audio48kHz96KbitrateMonoMp3,
		WordBoundaryEnabled:     false,
		SentenceBoundaryEnabled: false,
		ConcatenateResult:       false,
		DecompressOutputFiles:   false,
		SynthesisConfigVoice:    "zh-CN-YunxiNeural",
	})

	res, err := tts.LongTextToVoiceCreate(req)
	if err != nil {
		panic("方法调用报错: " + fmt.Sprintf("%v", err))
	}

	fmt.Println("任务生成成功，id：", res.Id)
}
```

获取长文本转语音的结果
```go
package main

import (
	"context"
	go_micro_tts "github.com/xuemingjings/go-micro-tts"
	"fmt"
)

var (
	speechKey    = "xxx"
	speechRegion = "xxx"
)

func main() {
	ctx := context.TODO()
	tts, err := go_micro_tts.NewGoTTS(
		ctx,
		go_micro_tts.WithSpeechRegion(speechRegion),
		go_micro_tts.WithSpeechKey(speechKey),
	)
	if err != nil {
		panic("初始化报错: " + fmt.Sprintf("%v", err))
	}

	taskId := " ffef80b4-d8c3-4e15-8d13-95d8c0b0f28d"
	res, err := tts.LongTextToVoiceId(taskId)
	if err != nil {
		panic("方法调用报错: " + fmt.Sprintf("%v", err))
	}
	
	fmt.Println("生成的语音文件zip：", res.Outputs.Result)
}
```

*更新使用方法，请查阅下方的接口*

## 接口

```go
// GetVoiceList 获取语音列表
func (g *GoTTS) GetVoiceList() (*[]VoiceList, error) 

// TextToVoiceDisk 文本转语音
func (g *GoTTS) TextToVoiceDisk(outFormat SsmlOut, ssml *SpeakXml, outFile *os.File) error

// TextToVoice 文本转语音
func (g *GoTTS) TextToVoice(outFormat SsmlOut, ssml *SpeakXml) 

// LongTextToVoiceCreate 创建批处理合成（长语音）
func (g *GoTTS) LongTextToVoiceCreate(longSpeak *LongSpeak) (*LongTextToVoiceCreateRep, error)

// LongTextToVoiceId 获取批处理合成（长语音）
func (g *GoTTS) LongTextToVoiceId(id string) (*LongTextToVoiceGetIdRep, error)

// LongTextToVoice 列出批处理合成（长语音）
func (g *GoTTS) LongTextToVoice(skip, top string) (*LongTextToVoiceGetRep, error) 

// LongTextToVoiceDel 删除批处理合成（长语音）
func (g *GoTTS) LongTextToVoiceDel(id string) (bool, error) 
```

### 语音列表
该JSON列表详细的包括了所有受支持的区域设置、声音、性别、风格和其他详细信息的 JSON 正文的响应。每个语音的 WordsPerMinute 属性可用于估计输出语音的长度。

详细参见 [voiceList](voiceList.md)


## 参考文档
- 文本转语音：https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/rest-text-to-speech?tabs=streaming
- 长文本转语音：https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/batch-synthesis