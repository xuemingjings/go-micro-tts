package go_micro_tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	speechKey    = ""
	speechRegion = ""

	// 长文本
	longTest = `
	先帝创业未半而中道崩殂，今天下三分，益州疲弊，此诚危急存亡之秋也。然侍卫之臣不懈于内，忠志之士忘身于外者，盖追先帝之殊遇，欲报之于陛下也。诚宜开张圣听，以光先帝遗德，恢弘志士之气，不宜妄自菲薄，引喻失义，以塞忠谏之路也。

　　宫中府中，俱为一体；陟罚臧否，不宜异同。若有作奸犯科及为忠善者，宜付有司论其刑赏，以昭陛下平明之理，不宜偏私，使内外异法也。

　　侍中、侍郎郭攸之、费祎、董允等，此皆良实，志虑忠纯，是以先帝简拔以遗陛下。愚以为宫中之事，事无大小，悉以咨之，然后施行，必能裨补阙漏，有所广益。

　　将军向宠，性行淑均，晓畅军事，试用于昔日，先帝称之曰能，是以众议举宠为督。愚以为营中之事，悉以咨之，必能使行阵和睦，优劣得所。

　　亲贤臣，远小人，此先汉所以兴隆也；亲小人，远贤臣，此后汉所以倾颓也。先帝在时，每与臣论此事，未尝不叹息痛恨于桓、灵也。侍中、尚书、长史、参军，此悉贞良死节之臣，愿陛下亲之信之，则汉室之隆，可计日而待也。

　　臣本布衣，躬耕于南阳，苟全性命于乱世，不求闻达于诸侯。先帝不以臣卑鄙，猥自枉屈，三顾臣于草庐之中，咨臣以当世之事，由是感激，遂许先帝以驱驰。后值倾覆，受任于败军之际，奉命于危难之间，尔来二十有一年矣。

　　先帝知臣谨慎，故临崩寄臣以大事也。受命以来，夙夜忧叹，恐托付不效，以伤先帝之明；故五月渡泸，深入不毛。今南方已定，兵甲已足，当奖率三军，北定中原，庶竭驽钝，攘除奸凶，兴复汉室，还于旧都。此臣所以报先帝而忠陛下之职分也。至于斟酌损益，进尽忠言，则攸之、祎、允之任也。

　　愿陛下托臣以讨贼兴复之效，不效，则治臣之罪，以告先帝之灵。若无兴德之言，则责攸之、祎、允等之慢，以彰其咎；陛下亦宜自谋，以咨诹善道，察纳雅言，深追先帝遗诏。臣不胜受恩感激。今当远离，临表涕零，不知所言。
	`
)

func TestGetToken(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	token := tts.token
	log.Printf("得到的token: %v", token)
}

func TestGetVoiceList(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	list, err := tts.GetVoiceList()
	if err != nil {
		log.Fatalf("获取VoiceList报错 err:%v", err)
	}

	log.Printf("获取语音列表 %+v", list)
}

func TestTextToVoiceDisk(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	outFormat := Audio48kHz96KbitrateMonoMp3

	ssml := NewSpeakXml(&SpeakXmlReq{
		Lang:   "zh-CN",
		Gender: "Male",
		Name:   "zh-CN-YunxiNeural",
		Text:   "中华兴盛 幸有斌哥 how are you",
	})

	// 创建输出文件
	fileName := time.Now().Format("2006-01-02-15-04-05") + "_file.mp3"
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	err = tts.TextToVoiceDisk(outFormat, ssml, outFile)

	if err != nil {
		log.Printf("TextToVoice调用发生了错误 err:%v \n", err)
		return
	}

	log.Println("文本转语音文件已生成： " + fileName)
}

func TestLongTextToVoiceCreate(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	inputs := []*LongSpeakInputs{
		{
			Text: "中华兴盛，幸有斌哥",
		},
	}

	req := NewLongSpeak(&LongSpeakXmlReq{
		DisplayName:             "长文本输出单个音频文件01",
		Inputs:                  inputs,
		OutputFormat:            Audio48kHz96KbitrateMonoMp3,
		WordBoundaryEnabled:     false,
		SentenceBoundaryEnabled: false,
		ConcatenateResult:       false,
		DecompressOutputFiles:   false,
		SynthesisConfigVoice:    "zh-CN-YunxiNeural",
	})

	res, err := tts.LongTextToVoiceCreate(req)
	if err != nil {
		log.Fatalf("方法返回错误 err:%v", err)
	}

	log.Printf("id:%s, res:%+v", res.Id, res)
}

func TestLongTextToVoiceCreate2(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	inputs := []*LongSpeakInputs{
		{
			Text: longTest + longTest + longTest + longTest + longTest + longTest + longTest + longTest,
		},
	}

	req := NewLongSpeak(&LongSpeakXmlReq{
		DisplayName:             "长文本输出单个音频文件02",
		Inputs:                  inputs,
		OutputFormat:            Audio48kHz96KbitrateMonoMp3,
		WordBoundaryEnabled:     false,
		SentenceBoundaryEnabled: false,
		ConcatenateResult:       false,
		DecompressOutputFiles:   false,
		SynthesisConfigVoice:    "zh-CN-YunxiNeural",
	})

	res, err := tts.LongTextToVoiceCreate(req)
	if err != nil {
		log.Fatalf("方法返回错误 err:%v", err)
	}

	log.Printf("id:%s, res:%+v", res.Id, res)
}

func TestLongTextToVoiceCreate3(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	inputs := []*LongSpeakInputs{
		{
			Text: longTest + longTest + longTest,
		},
		{
			Text: longTest + longTest + longTest,
		},
		{
			Text: longTest + longTest + longTest,
		},
	}

	req := NewLongSpeak(&LongSpeakXmlReq{
		DisplayName:             "长文本输出多个音频文件",
		Inputs:                  inputs,
		OutputFormat:            Audio48kHz96KbitrateMonoMp3,
		WordBoundaryEnabled:     false,
		SentenceBoundaryEnabled: false,
		ConcatenateResult:       false,
		DecompressOutputFiles:   false,
		SynthesisConfigVoice:    "zh-CN-YunxiNeural",
	})

	res, err := tts.LongTextToVoiceCreate(req)
	if err != nil {
		log.Fatalf("方法返回错误 err:%v", err)
	}

	log.Printf("id:%s, res:%+v", res.Id, res)
}

func TestLongTextToVoiceId(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	res, err := tts.LongTextToVoiceId("df052889-caa3-4e70-ba2a-a0d3d9abd8eb")
	if err != nil {
		log.Fatalf("方法返回错误 err:%v", err)
	}

	jsonByte, _ := json.Marshal(res)

	log.Println("下载的文件", res.Outputs.Result)
	log.Printf("方法返回 res: %v", string(jsonByte))
}

func TestLongTextToVoice(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	res, err := tts.LongTextToVoice("0", "100")
	if err != nil {
		log.Fatalf("方法返回错误 err:%v", err)
	}

	jsonByte, _ := json.Marshal(res)

	log.Printf("列表方法返回 res: %v", string(jsonByte))
}

func TestLongTextToVoiceDel(t *testing.T) {
	ctx := context.TODO()
	tts, err := NewGoTTS(
		ctx,
		WithSpeechRegion(speechRegion),
		WithSpeechKey(speechKey),
	)
	if err != nil {
		log.Fatalf("初始化报错 err:%v", err)
	}

	res, err := tts.LongTextToVoiceDel("67d8f902-1daa-4971-9292-76692ad96556")
	if err != nil {
		log.Fatalf("方法返回错误 err:%v", err)
	}

	log.Printf("res:%v", res)
}

// 创建批处理合成
// https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/batch-synthesis#create-batch-synthesis
func TestLongBatchSynthesisPost(t *testing.T) {
	//ctx := context.TODO()
	//tts, err := NewGoTTS(
	//	ctx,
	//	WithSpeechRegion(speechRegion),
	//	WithSpeechKey(speechKey),
	//)
	//if err != nil {
	//	log.Fatalf("初始化报错 err:%v", err)
	//}
	//token := tts.token

	requestBody := []byte(`{
    "displayName": "libin name",
    "description": "libin description",
    "textType": "SSML",
    "inputs": [
        {
            "text": "<speak version='1.0' xml:lang='en-US'><voice name='en-US-JennyNeural'>The rainbow has seven colors.</voice></speak>"
        }
    ],
    "properties": {
        "outputFormat": "riff-24khz-16bit-mono-pcm",
        "wordBoundaryEnabled": false,
        "sentenceBoundaryEnabled": false,
        "concatenateResult": false,
        "decompressOutputFiles": false
    }
}`)

	// 创建批处理合成
	uri := "https://%s.customvoice.api.speech.microsoft.com/api/texttospeech/3.1-preview1/batchsynthesis"
	url := fmt.Sprintf(uri, speechRegion)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", speechKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Body:", string(body))
}

// 列出批处理合成
func TestLongBatchSynthesisGet(t *testing.T) {
	// 列出批处理合成
	//uri := "https://%s.customvoice.api.speech.microsoft.com/api/texttospeech/3.1-preview1/batchsynthesis?skip=0&top=100"

	// 获取批处理合成
	uri := "https://%s.customvoice.api.speech.microsoft.com/api/texttospeech/3.1-preview1/batchsynthesis/67d8f902-1daa-4971-9292-76692ad96556"

	url := fmt.Sprintf(uri, speechRegion)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", speechKey)
	//req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Body:", string(body))
}
