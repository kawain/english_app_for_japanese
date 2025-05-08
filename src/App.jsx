import {
  createContext,
  useContext,
  useState,
  useEffect,
  useRef,
  useCallback
} from 'react'
import Home from './Home.jsx'
import ListeningContent from './ListeningContent.jsx'
import WordQuizContent from './WordQuizContent.jsx'
import TypingContent from './TypingContent.jsx'
import SearchContent from './SearchContent.jsx'
import Storage from './Storage.jsx'
import { speakTextAsync } from './utils/tts.js'
import { FaHome } from 'react-icons/fa'
import { MdHearing } from 'react-icons/md'
import { MdOutlineQuiz } from 'react-icons/md'
import { TiMessageTyping } from 'react-icons/ti'
import { GrStorage } from 'react-icons/gr'
import { FaSearch } from "react-icons/fa";

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
  const [volume, setVolume] = useState(30)
  const [isSoundEnabled, setIsSoundEnabled] = useState(false)

  // 初期化
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

        const success = await window.InitializeAppData()
        if (success) {
          setWasmInitialized(true)
          console.log('WASM およびデータ初期化完了')
        }
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

  // speak関数
  const speak = useCallback(
    async (text, lang) => {
      try {
        await speakTextAsync(text, lang, volume, isSoundEnabled)
      } catch (error) {
        console.error('speakでエラーが発生しました:', error)
      }
    },
    [volume, isSoundEnabled]
  )

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
      case 'search':
        return <SearchContent />
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
        toggleSound,
        speak
      }}
    >
      <div className='container'>
        <nav>
          <button
            onClick={() => setCurrentContent('home')}
            disabled={!wasmInitialized || currentContent === 'home'}
          >
            <FaHome />
            <span>ホーム</span>
          </button>
          <button
            onClick={() => setCurrentContent('listening')}
            disabled={!wasmInitialized || currentContent === 'listening'}
          >
            <MdHearing />
            <span>リスニング</span>
          </button>
          <button
            onClick={() => setCurrentContent('quiz')}
            disabled={!wasmInitialized || currentContent === 'quiz'}
          >
            <MdOutlineQuiz />
            <span>単語クイズ</span>
          </button>
          <button
            onClick={() => setCurrentContent('typing')}
            disabled={!wasmInitialized || currentContent === 'typing'}
          >
            <TiMessageTyping />
            <span>タイピング</span>
          </button>
          <button
            onClick={() => setCurrentContent('search')}
            disabled={!wasmInitialized || currentContent === 'search'}
          >
            <FaSearch />
            <span>検索</span>
          </button>
          <button
            onClick={() => setCurrentContent('storage')}
            disabled={!wasmInitialized || currentContent === 'storage'}
          >
            <GrStorage />
            <span>ストレージ</span>
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

export default App
