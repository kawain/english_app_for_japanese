function getEnvironmentFlags () {
  const userAgent = navigator.userAgent.toLowerCase()
  return {
    isUbuntu: userAgent.includes('ubuntu'),
    isWindows: userAgent.includes('windows'),
    isWindows10: userAgent.includes('windows nt 10.0'),
    isMac: userAgent.includes('macintosh'),
    isChrome: userAgent.includes('chrome') && !userAgent.includes('edg'),
    isSafari: userAgent.includes('safari') && !userAgent.includes('chrome'),
    isEdge: userAgent.includes('edg'),
    isFirefox: userAgent.includes('firefox')
  }
}

function validateSpeechParams (text, lang, volumeLevel, isSoundEnabled) {
  if (!isSoundEnabled) {
    console.log('Sound is disabled')
    return false
  }
  if (!window.speechSynthesis) {
    console.log('Speech synthesis is not supported')
    return false
  }
  if (!text) {
    console.log('No text to speak')
    return false
  }
  if (lang !== 'ja-JP' && lang !== 'en-US') {
    console.log('Unsupported language:', lang)
    return false
  }
  if (volumeLevel === 0) {
    console.log('Volume is muted')
    return false
  }

  const { isUbuntu, isFirefox } = getEnvironmentFlags()
  if (isUbuntu && isFirefox && lang === 'ja-JP') {
    console.log(
      'Skipping Japanese on Ubuntu Firefox due to potential compatibility issues.'
    )
    return false
  }
  return true
}

// export function speakText (text, lang, volumeLevel, isSoundEnabled) {
//   if (!validateSpeechParams(text, lang, volumeLevel, isSoundEnabled)) {
//     return
//   }

//   const utterance = new SpeechSynthesisUtterance(text)
//   utterance.lang = lang
//   utterance.pitch = 1
//   utterance.rate = 1
//   utterance.volume = Math.max(0, Math.min(1, volumeLevel / 100))
//   window.speechSynthesis.cancel()
//   window.speechSynthesis.speak(utterance)
// }

export async function speakTextAsync (text, lang, volumeLevel, isSoundEnabled) {
  if (!validateSpeechParams(text, lang, volumeLevel, isSoundEnabled)) {
    return
  }

  const env = getEnvironmentFlags()
  const isJpChromeEdge =
    lang === 'ja-JP' && env.isWindows10 && (env.isChrome || env.isEdge)

  const utterance = new SpeechSynthesisUtterance(text)
  utterance.lang = lang
  utterance.pitch = 1
  utterance.rate = 1
  utterance.volume = Math.max(0, Math.min(1, volumeLevel / 100))

  if (isJpChromeEdge) {
    utterance.volume = Math.max(0, Math.min(1, (volumeLevel * 2) / 100))
  }

  return new Promise(resolve => {
    utterance.onend = () => resolve()
    utterance.onerror = () => resolve()
    window.speechSynthesis.cancel()
    window.speechSynthesis.speak(utterance)
  })
}

// export function checkSyncOrAsync () {
//   const env = getEnvironmentFlags()

//   const syncEnvironments = [
//     { os: env.isUbuntu, browser: env.isChrome }
//     // { os: env.isWindows, browser: env.isFirefox },
//     // 必要に応じて追加
//   ]

//   const isSyncMode = syncEnvironments.some(env => env.os && env.browser)

//   return isSyncMode ? 'sync' : 'async'
// }
