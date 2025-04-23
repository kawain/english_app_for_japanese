const localStorageKey = 'excludedWords'

export async function getExcludedWordIds () {
  try {
    const storedIds = await localStorage.getItem(localStorageKey)
    return storedIds ? JSON.parse(storedIds) : []
  } catch (error) {
    console.error('ローカルストレージからのID取得に失敗しました:', error)
    return []
  }
}

export async function addExcludedWordId (wordId) {
  try {
    const excludedIds = await getExcludedWordIds()

    if (!excludedIds.includes(wordId)) {
      excludedIds.push(wordId)
      await localStorage.setItem(localStorageKey, JSON.stringify(excludedIds))
    }
    return true
  } catch (error) {
    console.error('ローカルストレージへのID追加に失敗しました:', error)
    throw error
  }
}

export async function removeExcludedWordId (wordId) {
  try {
    let excludedIds = await getExcludedWordIds()
    excludedIds = excludedIds.filter(id => id !== wordId)
    await localStorage.setItem(localStorageKey, JSON.stringify(excludedIds))
    return true
  } catch (error) {
    console.error('ローカルストレージからのID削除に失敗しました:', error)
    throw error
  }
}

export async function clearExcludedWordIds () {
  try {
    await localStorage.removeItem(localStorageKey)
    return true
  } catch (error) {
    console.error('ローカルストレージのクリアに失敗しました:', error)
    throw error
  }
}

function Storage () {
  // エクスポート機能
  const handleExport = () => {
    const data = localStorage.getItem(localStorageKey)
    if (!data) {
      alert('エクスポートするデータがありません。')
      return
    }

    try {
      // JSON形式かどうかのチェックのためだけに使用
      JSON.parse(data)
      // Blobを作成
      const blob = new Blob([data], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      // ダウンロードリンクを作成
      const a = document.createElement('a')
      a.href = url
      // ファイル名を .json に変更（データ形式に合わせて）
      a.download = 'excludedWords.json'
      document.body.appendChild(a) // Firefoxで必要になることがある
      a.click()
      document.body.removeChild(a) // 後片付け
      // URLを解放
      URL.revokeObjectURL(url)
    } catch (error) {
      // もしJSON形式でない単純な文字列データの場合
      console.warn(
        'データはJSON形式ではありませんでした。テキストとしてエクスポートします。'
      )
      const blob = new Blob([data], { type: 'text/plain' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'excludedWords.txt'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    }
  }

  // インポート機能
  const handleImport = event => {
    const file = event.target.files[0]
    if (!file) {
      // ファイルが選択されなかった場合（キャンセルなど）は何もしない
      return
    }
    // ファイル形式のチェック（任意）
    if (file.type !== 'application/json' && file.type !== 'text/plain') {
      alert('JSONファイルまたはテキストファイルを選択してください。')
      event.target.value = '' // input の値をリセット
      return
    }
    const reader = new FileReader()
    reader.onload = async e => {
      try {
        const content = e.target.result
        let parsedData

        // まずJSONとしてパースを試みる
        try {
          parsedData = JSON.parse(content)
        } catch (jsonError) {
          console.error('JSON パースエラー:', jsonError)
          throw new Error('インポートデータが正しいJSON形式ではありません。')
        }
        // データが配列であることを確認
        if (!Array.isArray(parsedData)) {
          throw new Error('インポートデータが配列形式ではありません。')
        }
        // localStorageに保存
        localStorage.setItem(localStorageKey, JSON.stringify(parsedData))
        // wasmの関数を呼び出して除外単語IDを設定
        await window.SetStorage(parsedData)
        alert('データをインポートしました。')
      } catch (error) {
        alert(`インポートに失敗しました: ${error.message}`)
      } finally {
        event.target.value = ''
      }
    }

    reader.onerror = () => {
      alert('ファイルの読み込みに失敗しました。')
      event.target.value = ''
    }

    // ファイルを読み込む
    reader.readAsText(file)
  }

  return (
    <>
      <div className='storage-container'>
        <h1>LocalStorage エクスポート/インポート</h1>
        <p>他のブラウザにデータを移動できます</p>
        <div className='storage-button-container'>
          <div>
            <button onClick={handleExport}>エクスポート</button>
          </div>
          <div>
            <input
              type='file'
              id='importFile'
              accept='.json,.txt'
              onChange={handleImport}
              className='hidden-input'
            />
            <label htmlFor='importFile' className='custom-file-button'>
              インポート
            </label>
          </div>
        </div>
      </div>
    </>
  )
}

export default Storage
