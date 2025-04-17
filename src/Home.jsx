import { useState, useMemo } from 'react'

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

  const handleButtonClick = level => {
    setLoading(true)
    setError(null)
    setCurrentLevel(level)

    try {
      const result = window.SearchAndReturnData(level)
      console.log(`レベル${level}の単語データ:`, result)
      // 結果が配列であることを確認（エラーハンドリング）
      if (Array.isArray(result)) {
        setAllData(result)
        setCurrentPage(1) // 新しいデータを取得したら1ページ目に戻る
        console.log(`レベル${level}の単語データ (${result.length}件):`, result)
      } else {
        // 想定外のデータ形式の場合
        console.error('取得したデータが配列ではありません:', result)
        setError('データの取得に失敗しました。予期しない形式です。')
        setAllData([]) // データを空にする
      }
    } catch (err) {
      console.error('SearchAndReturnDataの実行中にエラーが発生しました:', err)
      setError(`データの取得中にエラーが発生しました: ${err.message}`)
      setAllData([]) // エラー時はデータを空にする
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
  }, [allData, currentPage]) // allDataかcurrentPageが変わった時だけ再計算

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
            currentLevel !== null && <p>データがありません。</p>}

          {/* データがある場合のみテーブルとページネーションを表示 */}
          {!loading && !error && allData.length > 0 && (
            <>
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>English</th>
                    <th>Japanese</th>
                    <th>Example (EN/JP)</th>
                    <th>Level</th>
                    <th className='nowrap'>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {displayedData.map(item => (
                    <tr key={item.id}>
                      <td>{item.id}</td>
                      <td>{item.en}</td>
                      <td>{item.jp}</td>
                      <td>
                        {item.en2}
                        <br />
                        {item.jp2}
                      </td>
                      <td className='center-text'>{item.level}</td>
                      <td>
                        <button className='nowrap'>操作</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {renderPaginationButtons()}
            </>
          )}
        </div>
      </div>
    </>
  )
}

export default Home
