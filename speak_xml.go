package go_micro_tts

import "encoding/xml"

type SpeakXmlReq struct {
	Lang   string
	Gender string
	Name   string
	Text   string
}

func NewSpeakXml(req *SpeakXmlReq) *SpeakXml {
	return &SpeakXml{
		Version: "1.0",
		Lang:    req.Lang,
		Voice: VoiceXml{
			Lang:   req.Lang,
			Gender: req.Gender,
			Name:   req.Name,
			Text:   req.Text,
		},
	}
}

type LongSpeakXmlReq struct {
	DisplayName             string             // 批处理合成的名称
	Inputs                  []*LongSpeakInputs // 如果需要多个音频输出文件，则最多包含 1,000 个文本对象。
	OutputFormat            SsmlOut            // 音频输出格式
	WordBoundaryEnabled     bool               // 确定是否生成字边界数据
	SentenceBoundaryEnabled bool               // 确定是否生成句子边界数据
	ConcatenateResult       bool               // 确定是否要连接结果（设置为 true，则每个合成结果会写入同一音频输出文件）
	DecompressOutputFiles   bool               // 确定是否解压缩目标容器中的合成结果文件
	SynthesisConfigVoice    string             // 说出音频输出内容的语音
}

func NewLongSpeak(req *LongSpeakXmlReq) *LongSpeak {
	return &LongSpeak{
		DisplayName: req.DisplayName,
		TextType:    "PlainText",
		Inputs:      req.Inputs,
		Properties: &LongSpeakProperties{
			OutputFormat:            string(req.OutputFormat),
			WordBoundaryEnabled:     req.WordBoundaryEnabled,
			SentenceBoundaryEnabled: req.SentenceBoundaryEnabled,
			ConcatenateResult:       req.ConcatenateResult,
			DecompressOutputFiles:   req.DecompressOutputFiles,
		},
		SynthesisConfig: &SynthesisConfig{
			Voice: req.SynthesisConfigVoice,
		},
	}
}

func NewLongSpeakInputTextXml(reqs ...*LongSpeakInputTextXml) []*LongSpeakInputs {
	var res []*LongSpeakInputs
	for _, v := range reqs {
		xmlBytes, _ := xml.Marshal(v)
		xmlStr := &LongSpeakInputs{
			Text: string(xmlBytes),
		}
		res = append(res, xmlStr)
	}
	return res
}
