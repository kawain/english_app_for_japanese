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
  // 最後に押されたキー
  const [pressedKey, setPressedKey] = useState('')
  //
  const [times, setTimes] = useState(0)
  //
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
    setTimes(prev => prev + 1)
    setIsTypingStarted(true)
    // 読み上げ
    tts(question.en2, 'en-US', volume, isSoundEnabled)
  }

  const handleKeyDown = e => {
    // デフォルトのキー動作（例: Tabキーでのフォーカス移動など）を防ぐ場合
    e.preventDefault()
    let input = inputCharacters
    const moji = e.key
    input += moji
    setInputCharacters(input)
    setPressedKey(moji)
    const result = window.KeyDown(input, questionIndex)
    if (result > questionIndex) {
      setQuestionIndex(result)
      setInputCharacters('')

      if (result >= questionTextArray.length) {
        setTimeout(() => {
          const question = window.GetTypingQuestion()
          setQuestionText(question)
          const questionArray = window.GetTypingQuestionSlice()
          setQuestionTextArray(questionArray)
          setQuestionIndex(0)
          setInputCharacters('')
          setTimes(prev => prev + 1)
          // 読み上げ
          tts(question.en2, 'en-US', volume, isSoundEnabled)
        }, 500)
      }
    }
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
            <div className='number-area'>{times}回目</div>
            <div className='kanji-area'>
              {questionText.en2} {questionText.jp2}
            </div>
            <div className='hiragana-area'>
              {Array.isArray(questionTextArray) &&
                questionTextArray.map((character, index) => (
                  <span
                    key={index}
                    className={index < questionIndex ? 'correct-char' : ''}
                  >
                    {character}
                  </span>
                ))}
            </div>
            <div className='key-area'>
              最後に押されたキー: <span>{pressedKey}</span>
            </div>
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
