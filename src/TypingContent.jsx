import { useState, useEffect, useRef } from 'react'
import VolumeControl from './components/VolumeControl.jsx'
import { tts } from './utils/tts'

function TypingContent ({
  volume,
  onVolumeChange,
  isSoundEnabled,
  onToggleSound
}) {
  const [isTypingStarted, setIsTypingStarted] = useState(false)
  // 表示用の問題の文字列(日本語は漢字)
  const [questionText, setQuestionText] = useState({})
  // 問題の文字配列(日本語はひらがな)
  const [questionTextArray, setQuestionTextArray] = useState([])
  // questionTextArrayのインデックス
  const [questionIndex, setQuestionIndex] = useState(0)
  // 入力途中の文字列
  const [inputCharacters, setInputCharacters] = useState('')

  const [pressedKey, setPressedKey] = useState('')
  const typingAreaRef = useRef(null)

  useEffect(() => {
    typingAreaRef.current?.focus()
  }, [isTypingStarted])

  const handleStartTyping = () => {
    const question = window.GetTypingQuestion()
    setQuestionText(question)
    const questionArray = window.GetTypingQuestionSlice()
    setQuestionTextArray(questionArray)
    setQuestionIndex(0)

    console.log(question)
    console.log(questionArray)

    setIsTypingStarted(true)
  }

  const handleKeyDown = e => {
    // デフォルトのキー動作（例: Tabキーでのフォーカス移動など）を防ぐ場合
    // e.preventDefault();

    let input = inputCharacters
    input += e.key
    setInputCharacters(input)

    setPressedKey(e.key)

    // setInputCharacters(prev => [...prev, moji])

    // console.log('押されたキー:', moji)
    console.log('入力途中の文字配列:', inputCharacters)
  }

  return (
    <>
      <div className='typing-container'>
        {!isTypingStarted ? (
          <button onClick={handleStartTyping}>タイピング開始</button>
        ) : (
          <div
            ref={typingAreaRef}
            className='typing-content'
            tabIndex={0}
            onKeyDown={handleKeyDown}
          >
            <p>
              {/* questionTextArray が配列であることを確認してから map を使う */}
              {Array.isArray(questionTextArray) &&
                questionTextArray.map((character, index) => (
                  // 各文字 (character) を <span> タグで囲む
                  // React のリスト表示では、各要素に一意な key を指定する必要があるため、index を key として使用
                  <span key={index}>{character}</span>
                ))}
            </p>
            <p>
              {questionText.en2} {questionText.jp2}
            </p>
            <p>
              最後に押されたキー: <strong>{pressedKey}</strong>
            </p>
          </div>
        )}
      </div>
      <VolumeControl
        volume={volume}
        onVolumeChange={onVolumeChange}
        isSoundEnabled={isSoundEnabled}
        onToggleSound={onToggleSound}
      />
    </>
  )
}
export default TypingContent
