import { useState, useEffect, useRef } from 'react'
import ListeningContent from './ListeningContent.jsx'
import WordQuizContent from './WordQuizContent.jsx'
import TypingContent from './TypingContent.jsx'
import Storage, { getExcludedWordIds } from './Storage.jsx'

function App () {
  // WASMとデータの初期化が完了したか
  const [wasmInitialized, setWasmInitialized] = useState(false)
  const [currentContent, setCurrentContent] = useState('home')
  // 音量状態を追加 (初期値: 50)
  const [volume, setVolume] = useState(50)
  // 選択されたレベルの状態を追加
  const [selectedLevel, setSelectedLevel] = useState('1')
  // 初期化処理中/済みフラグ
  const isInitializing = useRef(false)

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
        console.log('読み込んだ除外単語ID:', loadedExcludedIds)
        // WASMの関数を呼び出して除外単語IDを設定
        window.CreateStorage(loadedExcludedIds)

        setWasmInitialized(true)
        console.log('WASM およびデータ初期化完了')
      } catch (error) {
        console.error('WASMのロードまたはデータ初期化に失敗しました:', error)
        isInitializing.current = false
      }
    }

    initializeWasmAndData()
  }, [])

  // 音量変更ハンドラを追加
  const handleVolumeChange = newVolume => {
    setVolume(newVolume)
    console.log('Volume 変更:', newVolume)
  }

  // レベル変更ハンドラを追加
  const handleLevelChange = newLevel => {
    setSelectedLevel(newLevel)
    console.log('Level 変更:', newLevel)
  }

  // 表示するコンテンツを決定する関数
  const renderContent = () => {
    // WASMの準備ができるまでローディング表示などを出す
    if (!wasmInitialized) {
      return <div className='loading'>Loading WASM and data...</div>
    }
    switch (currentContent) {
      case 'listening':
        return (
          <ListeningContent
            volume={volume}
            onVolumeChange={handleVolumeChange}
            level={selectedLevel}
            onLevelChange={handleLevelChange}
          />
        )
      case 'quiz':
        return (
          <WordQuizContent
            volume={volume}
            onVolumeChange={handleVolumeChange}
            level={selectedLevel}
            onLevelChange={handleLevelChange}
          />
        )
      case 'typing':
        return (
          <TypingContent volume={volume} onVolumeChange={handleVolumeChange} />
        )
      case 'storage':
        return <Storage />
      case 'home':
      default:
        return <h1>英語アプリ</h1>
    }
  }

  return (
    <>
      <div className='container'>
        <nav>
          <button
            onClick={() => setCurrentContent('home')}
            disabled={!wasmInitialized || currentContent === 'home'}
          >
            ホーム
          </button>
          <button
            onClick={() => setCurrentContent('listening')}
            disabled={!wasmInitialized || currentContent === 'listening'}
          >
            ヒアリング
          </button>
          <button
            onClick={() => setCurrentContent('quiz')}
            disabled={!wasmInitialized || currentContent === 'quiz'}
          >
            単語クイズ
          </button>
          <button
            onClick={() => setCurrentContent('typing')}
            disabled={!wasmInitialized || currentContent === 'typing'}
          >
            タイピング
          </button>
          <button
            onClick={() => setCurrentContent('storage')}
            disabled={!wasmInitialized || currentContent === 'storage'}
          >
            ストレージ
          </button>
        </nav>
        <main>{renderContent()}</main>
      </div>
    </>
  )
}

export default App
