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

export function tts (text, lang, volumeLevel, isSoundEnabled) {
  // 音声がオフならreturn
  if (!isSoundEnabled) {
    console.log('Sound is disabled')
    return
  }
  // ブラウザが読み上げをサポートしていない場合はreturn
  if (!window.speechSynthesis) {
    console.log('Speech synthesis is not supported')
    return
  }
  // textが空文字ならreturn
  if (!text) {
    console.log('No text to speak')
    return
  }
  // langが"ja-JP"か"en-US"以外はreturn
  if (lang !== 'ja-JP' && lang !== 'en-US') {
    console.log('Unsupported language:', lang)
    return
  }
  // volumeLevelが0ならreturn
  if (volumeLevel === 0) {
    console.log('Volume is muted')
    return
  }
  // OSがUbuntuでブラウザがfirefoxでlangが"ja-JP"ならreturn
  const userAgent = navigator.userAgent.toLowerCase()
  const isUbuntu = userAgent.includes('ubuntu')
  const isFirefox = userAgent.includes('firefox')

  if (isUbuntu && isFirefox && lang === 'ja-JP') {
    console.log(
      'Skipping Japanese TTS on Ubuntu Firefox due to potential compatibility issues.'
    )
    return
  }

  const uttr = new SpeechSynthesisUtterance(text)
  window.speechSynthesis.cancel() // 前の音声をクリア

  initializeVoices()
    .then(voices => {
      // console.log('Available voices:', voices)
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
      console.log('Speech synthesis completed: ', text)
    })
    .catch(err => {
      console.error('Speech synthesis error:', err)
    })
}
