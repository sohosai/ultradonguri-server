package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"slices"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
)

func FindSceneByName(client *goobs.Client, name string) (string, error) {
	scenes, err := client.Scenes.GetSceneList()

	if err != nil {
		return "", err
	}

	for _, scene := range scenes.Scenes {
		if scene.SceneName == name {
			return scene.SceneUuid, nil
		}
	}

	return "", fmt.Errorf("The scene named %s doesn't exist", name)
}

func FindInputByName(client *goobs.Client, name string) (string, error) {
	inputs_, err := client.Inputs.GetInputList()

	if err != nil {
		return "", err
	}

	for _, input := range inputs_.Inputs {
		if input.InputName == name {
			return input.InputUuid, nil
		}
	}

	return "", fmt.Errorf("The input named %s doesn't exist", name)
}

func CreateDummyScene(client *goobs.Client) (string, error) {
	sceneName, err := GetRandomString(10)
	if err != nil {
		return "", err
	}

	res, err := client.Scenes.CreateScene(scenes.NewCreateSceneParams().WithSceneName(sceneName))
	if err != nil {
		return "", err
	}

	return res.SceneUuid, nil
}

func DetermineTextInputKind(client *goobs.Client) string {
	kind := "text_ft2_source_v2" // for linux/mac
	if kinds, _ := client.Inputs.GetInputKindList(inputs.NewGetInputKindListParams()); kinds != nil {
		is_on_linux_or_mac := slices.Contains(kinds.InputKinds, kind)

		if !is_on_linux_or_mac {
			kind = "text_gdiplus_v2" // for windows
		}
	}

	return kind
}

// websocket apiのinputs.CreateInputはInputを作成した後にSceneに自動でSceneに追加してしまう。
// この関数はダミーのSceneに追加するためのラッパー
func CreateInputToDummyScene(client *goobs.Client, kind string, name string) (string, error) {
	dummySceneUuid, err := CreateDummyScene(client)

	if err != nil {
		return "", err
	}

	res, err := client.Inputs.CreateInput(
		inputs.NewCreateInputParams().
			WithInputKind(kind).
			WithInputName(name).
			WithSceneUuid(dummySceneUuid).
			WithSceneItemEnabled(false))

	if err != nil {
		return "", err
	}

	return res.InputUuid, nil
}

func GetRandomString(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
