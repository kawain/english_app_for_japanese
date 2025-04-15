import { useState, useEffect, useRef } from 'react'
import ListeningContent from './ListeningContent.jsx'
import WordQuizContent from './WordQuizContent.jsx'
import TypingContent from './TypingContent.jsx'
import Storage, { getExcludedWordIds } from './Storage.jsx'

function App () {
  // WASMとデータの初期化が完了したか
  const [wasmInitialized, setWasmInitialized] = useState(false)
  const [currentContent, setCurrentContent] = useState('home')
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
        if (window.CreateObjects) {
          window.CreateObjects(data)
          console.log('CreateObjects 呼び出し完了')
        } else {
          console.error('window.CreateObjects is not defined yet.')
          return
        }

        console.log('除外単語IDの読み込み開始...')
        const loadedExcludedIds = getExcludedWordIds()
        console.log('読み込んだ除外単語ID:', loadedExcludedIds)
        // WASMの関数を呼び出して除外単語IDを設定
        if (window.CreateStorage) {
          window.CreateStorage(loadedExcludedIds)
          console.log('CreateStorage 呼び出し完了')
        } else {
          console.error('window.CreateStorage is not defined yet.')
          return
        }

        setWasmInitialized(true)
        console.log('WASM およびデータ初期化完了')
      } catch (error) {
        console.error('WASMのロードまたはデータ初期化に失敗しました:', error)
        isInitializing.current = false
      }
    }

    initializeWasmAndData()
  }, [])

  // 表示するコンテンツを決定する関数
  const renderContent = () => {
    // WASMの準備ができるまでローディング表示などを出す
    if (!wasmInitialized) {
      return <div>Loading WASM and data...</div>
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
      default:
        return <h1>英語アプリ</h1>
    }
  }

  const handleTest1Click = () => {
    if (!wasmInitialized) {
      console.log('WASM is not initialized yet.')
      return
    }
    if (window.test1) {
      console.log('Calling test1...')
      window.test1()
    } else {
      console.error('window.test1 is not defined.')
    }
  }

  return (
    <>
      <div className='container'>
        <nav>
          <button
            onClick={() => setCurrentContent('listening')}
            disabled={!wasmInitialized}
          >
            ヒアリング
          </button>
          <button
            onClick={() => setCurrentContent('quiz')}
            disabled={!wasmInitialized}
          >
            単語クイズ
          </button>
          <button
            onClick={() => setCurrentContent('typing')}
            disabled={!wasmInitialized}
          >
            タイピング
          </button>
          <button
            onClick={() => setCurrentContent('storage')}
            disabled={!wasmInitialized}
          >
            ストレージ
          </button>
          {currentContent !== 'home' && (
            <button
              onClick={() => setCurrentContent('home')}
              disabled={!wasmInitialized}
            >
              ホームに戻る
            </button>
          )}
        </nav>
        <main>{renderContent()}</main>
        {wasmInitialized && <button onClick={handleTest1Click}>テスト1</button>}
      </div>
    </>
  )
}

export default App
