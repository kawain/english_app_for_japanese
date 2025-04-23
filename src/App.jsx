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
import Storage, {
  getExcludedWordIds,
  addExcludedWordId,
  removeExcludedWordId,
  clearExcludedWordIds
} from './Storage.jsx'
import { speakTextAsync } from './utils/tts.js'
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
        await window.CreateObject(data)

        console.log('除外単語IDの読み込み開始...')
        const loadedExcludedIds = await getExcludedWordIds()
        // WASMの関数を呼び出して除外単語IDを設定
        await window.SetStorage(loadedExcludedIds)

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

  // addStorage
  const addStorage = useCallback(async wordId => {
    try {
      await addExcludedWordId(wordId)
      await window.AddStorage(wordId)
    } catch (error) {
      console.error('addStorageでエラーが発生しました:', error)
    }
  }, [])

  // removeStorage
  const removeStorage = useCallback(async wordId => {
    try {
      await removeExcludedWordId(wordId)
      await window.RemoveStorage(wordId)
    } catch (error) {
      console.error('removeStorageでエラーが発生しました:', error)
    }
  }, [])

  // clearStorage
  const clearStorage = useCallback(async () => {
    try {
      await clearExcludedWordIds()
      await window.ClearStorage()
    } catch (error) {
      console.error('clearStorageでエラーが発生しました:', error)
    }
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
        volume,
        isSoundEnabled,
        addStorage,
        removeStorage,
        clearStorage,
        speak,
        handleLevelChange,
        handleVolumeChange,
        toggleSound
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
