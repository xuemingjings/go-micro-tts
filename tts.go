package go_micro_tts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jefferyjob/go-easy-utils/jsonUtil"
	"github.com/xuemingjings/go-micro-tts/internal"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	apiToken           = "https://%s.api.cognitive.microsoft.com/sts/v1.0/issueToken"
	apiVoiceList       = "https://%s.tts.speech.microsoft.com/cognitiveservices/voices/list"
	apiTextToVoice     = "https://%s.tts.speech.microsoft.com/cognitiveservices/v1"
	apiLongTextToVoice = "https://%s.customvoice.api.speech.microsoft.com/api/texttospeech/3.1-preview1/batchsynthesis"
)

type GoTTS struct {
	ctx          context.Context
	speechKey    string // SPEECH_KEY 必填
	speechRegion string // SPEECH_REGION 必填
	token        string // 自动生成
}

type Option func(*GoTTS)

func NewGoTTS(ctx context.Context, opts ...Option) (*GoTTS, error) {
	g := &GoTTS{ctx: ctx}

	for _, o := range opts {
		o(g)
	}

	// 参数验证
	if g.speechRegion == "" {
		return nil, errors.New("the parameter speechRegion is defined as")
	}
	if g.speechKey == "" {
		return nil, errors.New("the parameter speechKey is defined as")
	}

	return g, nil
}

func WithSpeechKey(speechKey string) Option {
	return func(g *GoTTS) {
		g.speechKey = speechKey
	}
}

func WithSpeechRegion(speechRegion string) Option {
	return func(g *GoTTS) {
		g.speechRegion = speechRegion
	}
}

func (g *GoTTS) setToken() error {
	if g.token != "" {
		return nil
	}

	uri := fmt.Sprintf(apiToken, g.speechRegion)

	header := map[string]any{
		"Ocp-Apim-Subscription-Key": g.speechKey,
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpFormUrlencoded),
	)
	resp, funcClose, err := client.SendRequest(http.MethodPost, uri, nil)
	defer funcClose()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("set token function request code:" + resp.Status)
	}

	req, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	g.token = string(req)

	return nil
}

// GetVoiceList 获取语音列表
func (g *GoTTS) GetVoiceList() (*[]VoiceList, error) {
	url := fmt.Sprintf(apiVoiceList, g.speechRegion)

	header := map[string]any{
		"Ocp-Apim-Subscription-Key": g.speechKey,
	}
	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpFormUrlencoded),
	)
	resp, funcClose, err := client.SendRequest(http.MethodGet, url, nil)
	defer funcClose()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	req, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &[]VoiceList{}
	err = json.Unmarshal(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// TextToVoiceDisk 文本转语音
func (g *GoTTS) TextToVoiceDisk(outFormat SsmlOut, ssml *SpeakXml, outFile *os.File) error {
	resp, funcClose, err := g.TextToVoice(outFormat, ssml)
	defer funcClose()
	if err != nil {
		return err
	}

	return internal.WriteToDisk(resp.Body, outFile)
}

// TextToVoice 文本转语音
func (g *GoTTS) TextToVoice(outFormat SsmlOut, ssml *SpeakXml) (*http.Response, func(), error) {
	uri := fmt.Sprintf(apiTextToVoice, g.speechRegion)

	err := g.setToken()
	if err != nil {
		return nil, func() {}, err
	}

	header := map[string]any{
		"Ocp-Apim-Subscription-Key": g.speechKey,
		"X-Microsoft-OutputFormat":  outFormat,
		"User-Agent":                "DouShen",
		"Authorization":             "Bearer " + g.token,
	}

	body := map[string]any{
		"xml": ssml,
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithTimeout(time.Second*60),
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpSsml),
	)
	resp, funcClose, err := client.SendRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, funcClose, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, funcClose, errors.New(resp.Status)
	}

	if resp.ContentLength == 0 {
		return nil, funcClose, errors.New("http response ContentLength=0")
	}

	return resp, funcClose, nil
}

// LongTextToVoiceCreate 创建批处理合成（长语音）
func (g *GoTTS) LongTextToVoiceCreate(longSpeak *LongSpeak) (*LongTextToVoiceCreateRep, error) {
	uri := fmt.Sprintf(apiLongTextToVoice, g.speechRegion)
	jsonData, _ := json.Marshal(longSpeak)
	header := map[string]any{
		"Ocp-Apim-Subscription-Key": g.speechKey,
	}
	body := map[string]any{
		"json": string(jsonData),
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpJson),
	)
	resp, funcClose, err := client.SendRequest(http.MethodPost, uri, body)
	defer funcClose()
	if err != nil {
		return nil, err
	}

	// 201 是成功，其他都是失败
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(resp.Status)
	}

	req, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &LongTextToVoiceCreateRep{}
	err = json.Unmarshal(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// LongTextToVoiceId 获取批处理合成（长语音）
func (g *GoTTS) LongTextToVoiceId(id string) (*LongTextToVoiceGetIdRep, error) {
	uri := fmt.Sprintf(apiLongTextToVoice, g.speechRegion) + "/" + id

	header := map[string]any{
		"Ocp-Apim-Subscription-Key": g.speechKey,
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpJson),
	)
	resp, funcClose, err := client.SendRequest(http.MethodGet, uri, nil)
	defer funcClose()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	req, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &LongTextToVoiceGetIdRep{}
	err = jsonUtil.JsonToStruct(string(req), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// LongTextToVoice 列出批处理合成（长语音）
func (g *GoTTS) LongTextToVoice(skip, top string) (*LongTextToVoiceGetRep, error) {
	params := fmt.Sprintf("?skip=%s&top=%s", skip, top)
	uri := fmt.Sprintf(apiLongTextToVoice, g.speechRegion) + params

	header := map[string]any{
		"Ocp-Apim-Subscription-Key": g.speechKey,
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpJson),
	)
	resp, funcClose, err := client.SendRequest(http.MethodGet, uri, nil)
	defer funcClose()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	req, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &LongTextToVoiceGetRep{}
	err = json.Unmarshal(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// LongTextToVoiceDel 删除批处理合成（长语音）
func (g *GoTTS) LongTextToVoiceDel(id string) (bool, error) {
	uri := fmt.Sprintf(apiLongTextToVoice, g.speechRegion) + "/" + id

	header := map[string]any{
		"Ocp-Apim-Subscription-Key": g.speechKey,
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpJson),
	)
	resp, funcClose, err := client.SendRequest(http.MethodDelete, uri, nil)
	defer funcClose()

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	}

	return false, err
}
