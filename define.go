package go_micro_tts

import (
	"encoding/xml"
	"time"
)

type VoiceList struct {
	Name            string `json:"Name"`
	DisplayName     string `json:"DisplayName"`
	LocalName       string `json:"LocalName"`
	ShortName       string `json:"ShortName"`
	Gender          string `json:"Gender"`
	Locale          string `json:"Locale"`
	LocaleName      string `json:"LocaleName"`
	SampleRateHertz string `json:"SampleRateHertz"`
	VoiceType       string `json:"VoiceType"`
	Status          string `json:"Status"`
	WordsPerMinute  string `json:"WordsPerMinute"`
}

type SpeakXml struct {
	XMLName xml.Name `xml:"speak"`
	Version string   `xml:"version,attr"`
	Lang    string   `xml:"xml:lang,attr"`
	Voice   VoiceXml `xml:"voice"`
}

type VoiceXml struct {
	XMLName xml.Name `xml:"voice"`
	Lang    string   `xml:"xml:lang,attr"`
	Gender  string   `xml:"xml:gender,attr"`
	Name    string   `xml:"name,attr"`
	Text    string   `xml:",chardata"`
}

// LongSpeak 长语音结构体定义
// https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/batch-synthesis-properties
type LongSpeak struct {
	DisplayName     string               `json:"displayName"` // [必选] 批处理合成的名称
	Description     string               `json:"description"` // [可选] 批处理合成的说明
	TextType        string               `json:"textType"`
	Inputs          []*LongSpeakInputs   `json:"inputs"`
	Properties      *LongSpeakProperties `json:"properties"`
	SynthesisConfig *SynthesisConfig     `json:"synthesisConfig"`
}

type SynthesisConfig struct {
	Voice string `json:"voice"`
}

type LongSpeakProperties struct {
	OutputFormat            string `json:"outputFormat"`
	WordBoundaryEnabled     bool   `json:"wordBoundaryEnabled"`
	SentenceBoundaryEnabled bool   `json:"sentenceBoundaryEnabled"`
	ConcatenateResult       bool   `json:"concatenateResult"`
	DecompressOutputFiles   bool   `json:"decompressOutputFiles"`
}

type LongSpeakInputs struct {
	Text string `json:"text"`
}

type LongSpeakInputTextXml struct {
	XMLName xml.Name `xml:"speak"`
	Version string   `xml:"version,attr"`
	Lang    string   `xml:"xml:lang,attr"`
	Voice   string   `xml:"voice,attr"`
	Text    string   `xml:",innerxml"`
}

// LongTextToVoiceCreateRep 长语音任务创建后返回
type LongTextToVoiceCreateRep struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	InnerError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"innerError"`
	TextType     string `json:"textType"`
	CustomVoices struct {
	} `json:"customVoices"`
	Properties struct {
		TimeToLive              string `json:"timeToLive"`
		OutputFormat            string `json:"outputFormat"`
		ConcatenateResult       bool   `json:"concatenateResult"`
		DecompressOutputFiles   bool   `json:"decompressOutputFiles"`
		WordBoundaryEnabled     bool   `json:"wordBoundaryEnabled"`
		SentenceBoundaryEnabled bool   `json:"sentenceBoundaryEnabled"`
		Customized              bool   `json:"customized"`
	} `json:"properties"`
	LastActionDateTime time.Time `json:"lastActionDateTime"`
	Status             string    `json:"status"`
	Id                 string    `json:"id"`
	CreatedDateTime    time.Time `json:"createdDateTime"`
	DisplayName        string    `json:"displayName"`
	Description        string    `json:"description"`
}

// LongTextToVoiceGetIdRep 根据ID查询任务返回
type LongTextToVoiceGetIdRep struct {
	TextType        string `json:"textType"`
	SynthesisConfig struct {
		Voice string `json:"voice"`
	} `json:"synthesisConfig,omitempty"`
	CustomVoices struct {
	} `json:"customVoices"`
	Properties struct {
		AudioSize           int    `json:"audioSize,omitempty"`
		DurationInTicks     int    `json:"durationInTicks,omitempty"`
		SucceededAudioCount int    `json:"succeededAudioCount,omitempty"`
		Duration            string `json:"duration,omitempty"`
		BillingDetails      struct {
			CustomNeural int `json:"customNeural"`
			Neural       int `json:"neural"`
		} `json:"billingDetails"`
		TimeToLive              string `json:"timeToLive"`
		OutputFormat            string `json:"outputFormat"`
		ConcatenateResult       bool   `json:"concatenateResult"`
		DecompressOutputFiles   bool   `json:"decompressOutputFiles"`
		WordBoundaryEnabled     bool   `json:"wordBoundaryEnabled"`
		SentenceBoundaryEnabled bool   `json:"sentenceBoundaryEnabled"`
		Customized              bool   `json:"customized"`
		Error                   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	} `json:"properties"`
	Outputs struct {
		Result string `json:"result"`
	} `json:"outputs"`
	LastActionDateTime time.Time `json:"lastActionDateTime"`
	Status             string    `json:"status"`
	Id                 string    `json:"id"`
	CreatedDateTime    time.Time `json:"createdDateTime"`
	DisplayName        string    `json:"displayName"`
	Description        string    `json:"description"`
}

type LongTextToVoiceGetRep struct {
	Values []LongTextToVoiceGetIdRep `json:"values"`
}

// SsmlOut 语音输出格式
type SsmlOut string

var (
	AmrWb16000Hz                  SsmlOut = "amr-wb-16000hz"
	Audio16kHz16Bit32kbpsMonoOpus SsmlOut = "audio-16khz-16bit-32kbps-mono-opus"
	Audio16kHz32KbitrateMonoMp3   SsmlOut = "audio-16khz-32kbitrate-mono-mp3"
	Audio16kHz64KbitrateMonoMp3   SsmlOut = "audio-16khz-64kbitrate-mono-mp3"
	Audio16kHz128KbitrateMonoMp3  SsmlOut = "audio-16khz-128kbitrate-mono-mp3"
	Audio24kHz16Bit24kbpsMonoOpus SsmlOut = "audio-24khz-16bit-24kbps-mono-opus"
	Audio24kHz16Bit48kbpsMonoOpus SsmlOut = "audio-24khz-16bit-48kbps-mono-opus"
	Audio24kHz48KbitrateMonoMp3   SsmlOut = "audio-24khz-48kbitrate-mono-mp3"
	Audio24kHz96KbitrateMonoMp3   SsmlOut = "audio-24khz-96kbitrate-mono-mp3"
	Audio24kHz160KbitrateMonoMp3  SsmlOut = "audio-24khz-160kbitrate-mono-mp3"
	Audio48kHz96KbitrateMonoMp3   SsmlOut = "audio-48khz-96kbitrate-mono-mp3"
	Audio48kHz192KbitrateMonoMp3  SsmlOut = "audio-48khz-192kbitrate-mono-mp3"
	Ogg16kHz16BitMonoOpus         SsmlOut = "ogg-16khz-16bit-mono-opus"
	Ogg24kHz16BitMonoOpus         SsmlOut = "ogg-24khz-16bit-mono-opus"
	Ogg48kHz16BitMonoOpus         SsmlOut = "ogg-48khz-16bit-mono-opus"
	Raw8kHz8BitMonoAlaw           SsmlOut = "raw-8khz-8bit-mono-alaw"
	Raw8kHz8BitMonoMulaw          SsmlOut = "raw-8khz-8bit-mono-mulaw"
	Raw8kHz16BitMonoPcm           SsmlOut = "raw-8khz-16bit-mono-pcm"
	Raw16kHz16BitMonoPcm          SsmlOut = "raw-16khz-16bit-mono-pcm"
	Raw16kHz16BitMonoTruesilk     SsmlOut = "raw-16khz-16bit-mono-truesilk"
	Raw22050Hz16BitMonoPcm        SsmlOut = "raw-22050hz-16bit-mono-pcm"
	Raw24kHz16BitMonoPcm          SsmlOut = "raw-24khz-16bit-mono-pcm"
	Raw24kHz16BitMonoTruesilk     SsmlOut = "raw-24khz-16bit-mono-truesilk"
	Raw44100Hz16BitMonoPcm        SsmlOut = "raw-44100hz-16bit-mono-pcm"
	Raw48kHz16BitMonoPcm          SsmlOut = "raw-48khz-16bit-mono-pcm"
	Webm16kHz16BitMonoOpus        SsmlOut = "webm-16khz-16bit-mono-opus"
	Webm24kHz16Bit24kbpsMonoOpus  SsmlOut = "webm-24khz-16bit-24kbps-mono-opus"
	Webm24kHz16BitMonoOpus        SsmlOut = "webm-24khz-16bit-mono-opus"
)
