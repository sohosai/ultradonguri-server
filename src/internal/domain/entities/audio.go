package entities

// Audio は音声のビジネスルールを表す
type Audio struct {
	IsMuted bool
}

// // Mute は音声をミュートするビジネスロジック
// func (a *Audio) Mute() {
// 	a.IsMuted = true
// }

// // Unmute は音声のミュートを解除するビジネスロジック
// func (a *Audio) Unmute() {
// 	a.IsMuted = false
// }
