// src/utils/tts.js

// 音声リストを事前にロードする関数
function initializeVoices () {
  return new Promise((resolve, reject) => {
    let voices = window.speechSynthesis.getVoices()
    if (voices.length) {
      resolve(voices)
      return
    }
    window.speechSynthesis.onvoiceschanged = () => {
      voices = window.speechSynthesis.getVoices()
      if (voices.length) {
        resolve(voices)
      } else {
        reject(new Error('No voices available'))
      }
    }
    // タイムアウトを設定（例：5秒）
    setTimeout(() => {
      if (!voices.length) {
        reject(new Error('Voice loading timed out'))
      }
    }, 5000)
  })
}

export function tts (text, lang, volumeLevel) {
  const uttr = new SpeechSynthesisUtterance(text)
  window.speechSynthesis.cancel() // 前の音声をクリア

  initializeVoices()
    .then(voices => {
      console.log('Available voices:', voices)
      // 日本語音声（lang: 'ja-JP'）、英語（en-US）、なければ最初の音声
      const preferredVoice =
        voices.find(voice => voice.lang === lang) || voices[0]
      if (!preferredVoice) {
        throw new Error('No suitable voice found')
      }
      uttr.voice = preferredVoice
      uttr.rate = 1 // 速度（0.1～10）
      uttr.pitch = 1 // ピッチ（0～2）
      uttr.volume = Math.max(0, Math.min(1, volumeLevel / 100)) // 0-1の範囲に正規化
      console.log('Speaking with voice:', preferredVoice.name)
      window.speechSynthesis.speak(uttr)
    })
    .catch(err => {
      console.error('Speech synthesis error:', err)
    })
}
