import { useState, useEffect, useMemo } from 'react'
import {
  addExcludedWordId,
  removeExcludedWordId,
  clearExcludedWordIds
} from './Storage.jsx'

// 1ページあたりの表示件数
const ITEMS_PER_PAGE = 100

function Home () {
  // APIから取得した全データ
  const [allData, setAllData] = useState([])
  // 現在表示中のページ番号
  const [currentPage, setCurrentPage] = useState(1)
  // 読み込み中かどうかの状態
  const [loading, setLoading] = useState(false)
  // エラーメッセージ
  const [error, setError] = useState(null)
  // 最後に選択されたレベル（ページネーションでレベルを維持するため）
  const [currentLevel, setCurrentLevel] = useState(null)
  // Setを使うことで、IDの追加・削除・存在チェックを効率的に行えます
  const [processingWordIds, setProcessingWordIds] = useState(new Set())
  // 日本語等を隠す状態
  const [showOffCol, setShowOffCol] = useState(false)
  // セルごとに日本語等を隠す状態を管理
  const [showOffCell, setShowOffCell] = useState({})

  const handleButtonClick = async level => {
    setLoading(true)
    setError(null)
    setCurrentLevel(level)
    setProcessingWordIds(new Set())
    setShowOffCol(false)
    setShowOffCell({})

    try {
      const result = await window.SearchAndReturnDataPromise(level)
      console.log(`レベル${level}の単語データ:`, result)
      setAllData(result)
      setCurrentPage(1)
    } catch (err) {
      console.error(
        'SearchAndReturnDataPromiseの実行中にエラーが発生しました:',
        err
      )
      setError(`データの取得中にエラーが発生しました: ${err.message}`)
      setAllData([])
    } finally {
      setLoading(false)
    }
  }

  // ページ変更時の処理
  const handlePageChange = pageNumber => {
    setCurrentPage(pageNumber)
    window.scrollTo(0, 0)
  }

  // ページネーションの計算
  const totalPages = Math.ceil(allData.length / ITEMS_PER_PAGE)

  // 現在のページに表示するデータを計算 (useMemoで計算結果をキャッシュ)
  const displayedData = useMemo(() => {
    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE
    const endIndex = startIndex + ITEMS_PER_PAGE
    return allData.slice(startIndex, endIndex)
  }, [allData, currentPage])

  // ページネーションボタンを生成する関数
  const renderPaginationButtons = () => {
    if (totalPages <= 1) {
      return null // 1ページ以下の場合はボタンを表示しない
    }

    const buttons = []
    for (let i = 1; i <= totalPages; i++) {
      buttons.push(
        <button
          key={i}
          onClick={() => handlePageChange(i)}
          disabled={currentPage === i}
        >
          {i}
        </button>
      )
    }
    return <div className='pagination-buttons'>{buttons}</div>
  }

  // 復元、除外ボタンの操作
  const handleActionClick = wordId => {
    if (currentLevel === 1 || currentLevel === 2) {
      addExcludedWordId(wordId)
      window.AddStorage(wordId)
    } else if (currentLevel === 0) {
      removeExcludedWordId(wordId)
      window.RemoveStorage(wordId)
    }
    setProcessingWordIds(prev => new Set(prev).add(wordId))
  }

  // 日本語等を隠す状態のチェックボックスのハンドラ
  const handleShowOffCol = event => {
    const isChecked = event.target.checked
    setShowOffCol(isChecked)
    // 現在表示されている displayedData の各要素に対して状態を更新
    setShowOffCell(prev => {
      const newState = { ...prev }
      displayedData.forEach(item => {
        newState[item.id] = isChecked
      })
      return newState
    })
  }

  // セルごとに日本語等を隠す状態のハンドラ
  const handleShowOffCell = wordId => {
    // 以前のステート (prev) を受け取り、新しいオブジェクトを返す
    setShowOffCell(prev => ({
      ...prev, // 既存のキーと値をコピー
      [wordId]: !prev[wordId] // 対象のキーの値だけを反転させる
    }))
  }

  // displayedData または showOffCol (全体表示フラグ) が変更されたときに実行
  useEffect(() => {
    // 新しいページのデータに基づいて showOffCell の状態を更新
    setShowOffCell(prev => {
      const newState = { ...prev } // 既存の状態を引き継ぐ
      const currentIds = new Set(displayedData.map(item => item.id))

      // 古いページのIDを削除 (任意ですが、メモリリークを防ぐために推奨)
      Object.keys(newState).forEach(key => {
        // IDが数値の場合はparseIntが必要
        if (!currentIds.has(parseInt(key, 10))) {
          delete newState[key]
        }
      })

      // 新しいページのIDを追加/更新
      displayedData.forEach(item => {
        // 既に状態が存在しない場合のみ、showOffCol (全体表示フラグ) に基づいて初期化
        if (newState[item.id] === undefined) {
          newState[item.id] = showOffCol
        }
      })
      return newState
    })
  }, [displayedData, showOffCol])

  return (
    <>
      <div className='home-container'>
        <h1>英語アプリ</h1>
        <div className='home-buttons'>
          <button
            onClick={() => handleButtonClick(1)}
            disabled={loading || currentLevel === 1}
          >
            レベル1の単語
          </button>
          <button
            onClick={() => handleButtonClick(2)}
            disabled={loading || currentLevel === 2}
          >
            レベル2の単語
          </button>
          <button
            onClick={() => handleButtonClick(0)}
            disabled={loading || currentLevel === 0}
          >
            除外された単語
          </button>
        </div>

        {/* エラー表示 */}
        {error && <div style={{ color: 'red' }}>エラー: {error}</div>}

        {/* 結果表示エリア */}
        <div className='home-results'>
          {loading && currentLevel !== null && <p>読み込み中...</p>}
          {!loading &&
            !error &&
            allData.length === 0 &&
            currentLevel !== null && (
              <p style={{ textAlign: 'center' }}>データがありません。</p>
            )}

          {/* データがある場合のみテーブルとページネーションを表示 */}
          {!loading && !error && allData.length > 0 && (
            <>
              <p>
                {currentPage}ページ目 / {allData.length}件見つかりました
              </p>
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>English</th>
                    <th>
                      <label className='show-or-hide'>
                        Japanese{' '}
                        <input
                          className='show-or-hide'
                          type='checkbox'
                          checked={showOffCol}
                          onChange={handleShowOffCol}
                        />
                      </label>
                    </th>
                    <th>Example (EN/JA)</th>
                    <th>Level</th>
                    <th className='nowrap'>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {displayedData.map(item => (
                    <tr key={item.id}>
                      <td>{item.id}</td>
                      <td>{item.en}</td>
                      <td
                        className='show-or-hide'
                        style={{ opacity: showOffCell[item.id] ? 1 : 0 }}
                        onClick={() => {
                          handleShowOffCell(item.id)
                        }}
                      >
                        {item.jp}
                      </td>
                      <td style={{ opacity: showOffCell[item.id] ? 1 : 0 }}>
                        {item.en2}
                        <br />
                        {item.jp2}
                      </td>
                      <td className='center-text'>{item.level}</td>
                      <td>
                        <button
                          className='nowrap'
                          onClick={() => handleActionClick(item.id)}
                          disabled={processingWordIds.has(item.id)}
                        >
                          {currentLevel === 0 ? '復元' : '除外'}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {renderPaginationButtons()}
              {currentLevel === 0 ? (
                <div style={{ textAlign: 'center', marginTop: '20px' }}>
                  <button
                    onClick={() => {
                      if (
                        window.confirm('本当に除外リストを全部クリアしますか？')
                      ) {
                        clearExcludedWordIds()
                        window.ClearStorage()
                        setProcessingWordIds(new Set())
                        setAllData([])
                        setCurrentPage(1)
                        window.scrollTo(0, 0)
                      }
                    }}
                  >
                    除外リスト全クリア
                  </button>
                </div>
              ) : (
                ''
              )}
              <div style={{ textAlign: 'right' }}>
                <button onClick={() => window.scrollTo(0, 0)}>
                  ページ先頭
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </>
  )
}

export default Home
