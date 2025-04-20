import { createContext, useContext, useState, useEffect, useRef } from 'react'
import Home from './Home.jsx'
import ListeningContent from './ListeningContent.jsx'
import WordQuizContent from './WordQuizContent.jsx'
import TypingContent from './TypingContent.jsx'
import Storage, { getExcludedWordIds } from './Storage.jsx'
import { FaHome } from 'react-icons/fa'
import { MdHearing } from 'react-icons/md'
import { MdOutlineQuiz } from 'react-icons/md'
import { TiMessageTyping } from 'react-icons/ti'
import { GrStorage } from 'react-icons/gr'

// コンテキストを作成
const AppContext = createContext()

function App () {
  // WASMとデータの初期化が完了したか
  const [wasmInitialized, setWasmInitialized] = useState(false)
  const [currentContent, setCurrentContent] = useState('home')
  // 初期化処理中/済みフラグ
  const isInitializing = useRef(false)
  // エラーメッセージ
  const [errorMessage, setErrorMessage] = useState('')

  // コンテキストで管理する状態
  const [selectedLevel, setSelectedLevel] = useState('1')
  const [volume, setVolume] = useState(50)
  const [isSoundEnabled, setIsSoundEnabled] = useState(false)

  // レベル変更ハンドラ
  const handleLevelChange = newLevel => {
    setSelectedLevel(newLevel)
    console.log('Level 変更:', newLevel)
  }

  // 音量変更ハンドラ
  const handleVolumeChange = newVolume => {
    setVolume(newVolume)
    console.log('Volume 変更:', newVolume)
  }

  // サウンドのオン/オフを切り替える関数
  const toggleSound = () => {
    setIsSoundEnabled(prev => !prev)
  }

  useEffect(() => {
    if (isInitializing.current) return
    isInitializing.current = true

    const initializeWasmAndData = async () => {
      try {
        console.log('WASM ロード開始...')
        const go = new window.Go()
        const result = await WebAssembly.instantiateStreaming(
          fetch('./main.wasm'),
          go.importObject
        )
        go.run(result.instance)
        console.log('WASM インスタンス実行開始')

        console.log('CSVデータの読み込み開始...')
        const response = await fetch('./word.csv')
        const text = await response.text()
        const rows = text.split('\n')
        const data = []
        for (let i = 1; i < rows.length; i++) {
          const parts = rows[i].split('\t')
          if (parts.length === 8) {
            data.push(parts.map(p => p.trim()))
          }
        }
        // WASMの関数を呼び出して単語オブジェクトを作成
        window.CreateObjects(data)

        console.log('除外単語IDの読み込み開始...')
        const loadedExcludedIds = getExcludedWordIds()
        // WASMの関数を呼び出して除外単語IDを設定
        window.CreateStorage(loadedExcludedIds)

        setWasmInitialized(true)
        console.log('WASM およびデータ初期化完了')
      } catch (error) {
        console.error('WASMのロードまたはデータ初期化に失敗しました:', error)
        setErrorMessage(
          'データ初期化に失敗しました。ブラウザをリロードしてください。'
        )
        isInitializing.current = false
      }
    }

    initializeWasmAndData()
  }, [])

  // 表示するコンテンツを決定する関数
  const renderContent = () => {
    if (!wasmInitialized) {
      return <div className='loading'>Loading data...</div>
    }
    if (errorMessage) {
      return <div className='error'>{errorMessage}</div>
    }
    switch (currentContent) {
      case 'listening':
        return <ListeningContent />
      case 'quiz':
        return <WordQuizContent />
      case 'typing':
        return <TypingContent />
      case 'storage':
        return <Storage />
      case 'home':
        return <Home />
      default:
        return <Home />
    }
  }

  // コンテキストプロバイダーで状態と関数を共有
  return (
    <AppContext.Provider
      value={{
        selectedLevel,
        handleLevelChange,
        volume,
        handleVolumeChange,
        isSoundEnabled,
        toggleSound
      }}
    >
      <div className='container'>
        <nav>
          <button
            onClick={() => setCurrentContent('home')}
            disabled={!wasmInitialized || currentContent === 'home'}
          >
            <FaHome /> ホーム
          </button>
          <button
            onClick={() => setCurrentContent('listening')}
            disabled={!wasmInitialized || currentContent === 'listening'}
          >
            <MdHearing /> リスニング
          </button>
          <button
            onClick={() => setCurrentContent('quiz')}
            disabled={!wasmInitialized || currentContent === 'quiz'}
          >
            <MdOutlineQuiz /> 単語クイズ
          </button>
          <button
            onClick={() => setCurrentContent('typing')}
            disabled={!wasmInitialized || currentContent === 'typing'}
          >
            <TiMessageTyping /> タイピング
          </button>
          <button
            onClick={() => setCurrentContent('storage')}
            disabled={!wasmInitialized || currentContent === 'storage'}
          >
            <GrStorage /> ストレージ
          </button>
        </nav>
        <main>{renderContent()}</main>
      </div>
    </AppContext.Provider>
  )
}

// 子コンポーネントで使うためのカスタムフック
export const useAppContext = () => {
  const context = useContext(AppContext)
  if (!context) {
    throw new Error('useAppContext must be used within an AppContext.Provider')
  }
  return context
}

export function speakText (text, lang, volumeLevel, isSoundEnabled) {
  // 音声がオフならすぐに終了
  if (!isSoundEnabled) {
    console.log('Sound is disabled')
    return
  }
  // ブラウザが読み上げをサポートしていない場合はすぐに終了
  if (!window.speechSynthesis) {
    console.log('Speech synthesis is not supported')
    return
  }
  // textが空文字ならすぐに終了
  if (!text) {
    console.log('No text to speak')
    return
  }
  // langが"ja-JP"か"en-US"以外はすぐに終了
  if (lang !== 'ja-JP' && lang !== 'en-US') {
    console.log('Unsupported language:', lang)
    return
  }
  // volumeLevelが0ならすぐに終了
  if (volumeLevel === 0) {
    console.log('Volume is muted')
    return
  }
  // OSがUbuntuでブラウザがfirefoxでlangが"ja-JP"ならすぐに終了
  const userAgent = navigator.userAgent.toLowerCase()
  const isUbuntu = userAgent.includes('ubuntu')
  const isFirefox = userAgent.includes('firefox')

  if (isUbuntu && isFirefox && lang === 'ja-JP') {
    console.log(
      'Skipping Japanese on Ubuntu Firefox due to potential compatibility issues.'
    )
    return
  }

  const utterance = new SpeechSynthesisUtterance(text)
  utterance.lang = lang
  utterance.rate = 1 // 速度（0.1～10）
  utterance.pitch = 1 // ピッチ（0～2）
  utterance.volume = Math.max(0, Math.min(1, volumeLevel / 100)) // 0-1の範囲に正規化

  window.speechSynthesis.cancel()
  window.speechSynthesis.speak(utterance)
}

export default App
