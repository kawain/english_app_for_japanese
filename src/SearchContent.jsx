import { useState } from 'react'
import { useAppContext } from './App.jsx'
import VolumeControl from './components/VolumeControl.jsx'

function SearchContent () {
  const { speak } = useAppContext()
  const [keyword, setKeyword] = useState('')
  const [searchResults, setSearchResults] = useState([])

  const handleInputChange = event => {
    setKeyword(event.target.value)
  }

  const handleSearch = async () => {
    let word = keyword.trim()
    if (!word) return
    try {
      let result = await window.SearchWord(word)
      console.log('検索結果:', result)
      setSearchResults(result)
    } catch (error) {
      console.error('検索に失敗しました:', error)
      setSearchResults([])
    }
  }

  return (
    <>
      <div className='search-container'>
        <h2>キーワード検索</h2>
        <div className='search-input-container'>
          <input
            type='text'
            value={keyword}
            onChange={handleInputChange}
            required
            placeholder='検索キーワードを入力してください'
          />
          <button onClick={handleSearch}>検索</button>
        </div>
        {searchResults.length > 0 && (
          <>
            <h3>検索結果</h3>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>English</th>
                  <th>Japanese</th>
                  <th>Example (EN/JA)</th>
                </tr>
              </thead>
              <tbody>
                {searchResults.map(item => (
                  <tr key={item.id}>
                    <td>{item.id}</td>
                    <td
                      style={{ cursor: 'pointer' }}
                      role='button'
                      title='読み上げ'
                      onClick={async () => await speak(item.en, 'en-US')}
                    >
                      {item.en}
                    </td>
                    <td>{item.jp}</td>
                    <td
                      style={{ cursor: 'pointer' }}
                      role='button'
                      title='読み上げ'
                      onClick={async () => await speak(item.en2, 'en-US')}
                    >
                      {item.en2}
                      <br />
                      {item.jp2}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        )}
      </div>
      <VolumeControl />
    </>
  )
}

export default SearchContent
