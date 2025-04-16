// src/utils/tts.js

/**
 * 指定されたテキストを読み上げる関数 (Text-to-Speech)
 * @param {string} text 読み上げるテキスト
 * @param {string} lang 言語コード (例: 'en-US', 'ja-JP')
 * @param {number} volumeLevel 音量 (0から100の範囲)
 * @returns {Promise<void>} 発話が完了したら解決される Promise
 */
export function tts (text, lang, volumeLevel) {
  return new Promise((resolve, reject) => {
    // Speech Synthesis APIがサポートされているか確認
    if (!window.speechSynthesis) {
      const errorMsg = 'Speech Synthesis API is not supported in this browser.'
      console.error(errorMsg)
      // サポートされていない場合はエラーで Promise を reject
      reject(new Error(errorMsg))
      return
    }

    const uttr = new SpeechSynthesisUtterance()
    uttr.text = text
    uttr.lang = lang
    uttr.rate = 1.0 // 読み上げ速度 (デフォルト: 1.0)
    uttr.pitch = 1.0 // 声の高さ (デフォルト: 1.0)

    // volumeLevel (0-100) を SpeechSynthesisUtterance の volume (0-1) に変換
    // Math.max/min で範囲外の値を丸める
    uttr.volume = Math.max(0, Math.min(1, volumeLevel / 100))

    // 発話終了時の処理
    uttr.onend = () => {
      resolve() // Promise を解決
    }

    // エラー発生時の処理
    uttr.onerror = error => {
      console.error('Speech synthesis error:', error)
      reject(error) // Promise を reject
    }

    // 発話を実行
    window.speechSynthesis.speak(uttr)
  })
}
