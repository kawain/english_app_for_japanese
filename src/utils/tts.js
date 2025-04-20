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
  // Promiseを返すように変更
  return new Promise((resolve, reject) => {
    // 音声がオフならすぐに解決して終了
    if (!isSoundEnabled) {
      console.log('Sound is disabled')
      resolve()
      return
    }
    // ブラウザが読み上げをサポートしていない場合はすぐに解決して終了
    if (!window.speechSynthesis) {
      console.log('Speech synthesis is not supported')
      resolve()
      return
    }
    // textが空文字ならすぐに解決して終了
    if (!text) {
      console.log('No text to speak')
      resolve()
      return
    }
    // langが"ja-JP"か"en-US"以外はすぐに解決して終了
    if (lang !== 'ja-JP' && lang !== 'en-US') {
      console.log('Unsupported language:', lang)
      resolve()
      return
    }
    // volumeLevelが0ならすぐに解決して終了
    if (volumeLevel === 0) {
      console.log('Volume is muted')
      resolve()
      return
    }
    // OSがUbuntuでブラウザがfirefoxでlangが"ja-JP"ならすぐに解決して終了
    const userAgent = navigator.userAgent.toLowerCase()
    const isUbuntu = userAgent.includes('ubuntu')
    const isFirefox = userAgent.includes('firefox')

    if (isUbuntu && isFirefox && lang === 'ja-JP') {
      console.log(
        'Skipping Japanese TTS on Ubuntu Firefox due to potential compatibility issues.'
      )
      resolve()
      return
    }

    window.speechSynthesis.cancel()

    const uttr = new SpeechSynthesisUtterance(text)

    initializeVoices()
      .then(voices => {
        // console.log('Available voices:', voices);
        const preferredVoice =
          voices.find(voice => voice.lang === lang) || voices[0]
        if (!preferredVoice) {
          // 適切な音声が見つからない場合もエラーとして扱う
          throw new Error(`No suitable voice found for lang: ${lang}`)
        }
        uttr.voice = preferredVoice
        uttr.rate = 1 // 速度（0.1～10）
        uttr.pitch = 1 // ピッチ（0～2）
        uttr.volume = Math.max(0, Math.min(1, volumeLevel / 100)) // 0-1の範囲に正規化

        // 読み上げ完了時のイベントリスナー
        uttr.onend = () => {
          console.log('Speech synthesis completed: ', text)
          resolve() // Promiseを解決
        }

        // 読み上げエラー時のイベントリスナー
        uttr.onerror = event => {
          console.error('Speech synthesis error:', event.error)
          reject(new Error(`Speech synthesis error: ${event.error}`)) // Promiseを拒否
        }

        console.log('Speaking with voice:', preferredVoice.name, 'Text:', text)
        window.speechSynthesis.speak(uttr)
      })
      .catch(err => {
        console.error(
          'Speech synthesis initialization or voice finding error:',
          err
        )
        reject(err) // 初期化や音声検索のエラーも拒否
      })
  })
}
