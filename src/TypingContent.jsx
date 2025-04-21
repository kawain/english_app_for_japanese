import { useState, useEffect, useRef } from 'react'
import { useAppContext } from './App.jsx'
import VolumeControl from './components/VolumeControl.jsx'
import { GrLinkPrevious } from 'react-icons/gr'
import { GrLinkNext } from 'react-icons/gr'

function TypingContent () {
  const { speak } = useAppContext()
  // 0: 初期状態スタートボタン表示
  // 1: タイピング表示
  const [progress, setProgress] = useState(0)
  // 最大問題数（インデックスの数）
  const [maxIndex, setMaxIndex] = useState(0)
  // 今の問題のインデックス
  const [currentIndex, setCurrentIndex] = useState(0)
  // 表示用の問題の文字列(日本語は漢字)
  const [questionText, setQuestionText] = useState({})
  // 問題の文字配列(英語)
  const [questionTextArray1, setQuestionTextArray1] = useState([])
  // 問題の文字配列(日本語)
  const [questionTextArray2, setQuestionTextArray2] = useState([])
  // questionTextArray1のインデックス
  const [questionIndex1, setQuestionIndex1] = useState(0)
  // questionTextArray2のインデックス
  const [questionIndex2, setQuestionIndex2] = useState(0)
  // 入力途中の文字列
  const [inputCharacters, setInputCharacters] = useState('')
  // 最後に押されたキー
  const [pressedKey, setPressedKey] = useState('')
  // typing-contentのref
  const typingAreaRef = useRef(null)
  // 英語か日本語かどちらを行っているか
  const whichRef = useRef(1)
  const timerIdRef = useRef(null)

  // WASMの関数で問題のセットアップと問題数を返す
  useEffect(() => {
    const result = window.CreateTyping()
    console.log(result)
    setMaxIndex(result)
  }, [])

  // フォーカスをする
  useEffect(() => {
    if (progress === 0) return
    typingAreaRef.current?.focus()
  }, [progress])

  // 問題配列から任意のインデックスで問題を抽出
  const selectQuestion = (index, startFlag = false) => {
    setCurrentIndex(index)
    // WASMの関数
    const question = window.GetTypingQuestion(index)
    setQuestionText(question)
    // WASMの関数(英語の配列)
    const array1 = window.GetTypingQuestionSlice(1)
    setQuestionTextArray1(array1)
    // WASMの関数(日本語の配列)
    const array2 = window.GetTypingQuestionSlice(2)
    setQuestionTextArray2(array2)

    setQuestionIndex1(0)
    setQuestionIndex2(0)
    setInputCharacters('')
    setPressedKey('')

    if (startFlag) {
      setProgress(1)
    }

    speak(question.en2, 'en-US')
  }

  // タイピング開始
  const handleStart = () => {
    selectQuestion(0, true)
  }

  const handlePrevious = () => {
    let index = currentIndex
    if (index > 0) {
      index--
      selectQuestion(index)
    } else {
      selectQuestion(maxIndex - 1)
    }
  }

  const handleNext = () => {
    let index = currentIndex
    if (index < maxIndex - 1) {
      index++
      selectQuestion(index)
    } else {
      selectQuestion(0)
    }
  }

  const handleKeyDown = e => {
    // デフォルトのキー動作（例: Tabキーでのフォーカス移動など）を防ぐ場合
    e.preventDefault()

    let input = inputCharacters
    const moji = e.key
    input += moji
    setInputCharacters(input)
    setPressedKey(moji)

    let result = 0
    if (whichRef.current === 1) {
      result = window.KeyDown(input, questionIndex1, 1)
      if (result > questionIndex1) {
        setQuestionIndex1(result)
        setInputCharacters('')
        if (result >= questionTextArray1.length) {
          whichRef.current = 2
        }
      }
    } else if (whichRef.current === 2) {
      result = window.KeyDown(input, questionIndex2, 2)
      if (result > questionIndex2) {
        setQuestionIndex2(result)
        setInputCharacters('')
        if (result >= questionTextArray2.length) {
          timerIdRef.current = setTimeout(() => {
            whichRef.current = 1
            selectQuestion(currentIndex + 1)
            timerIdRef.current = null
          }, 500)
        }
      }
    }
  }

  useEffect(() => {
    return () => {
      if (timerIdRef.current) {
        clearTimeout(timerIdRef.current)
      }
    }
  }, [])

  let content = null

  if (progress === 0) {
    content = <button onClick={handleStart}>タイピング開始</button>
  } else if (progress === 1) {
    content = (
      <div
        ref={typingAreaRef}
        className='typing-content'
        tabIndex={0}
        onKeyDown={handleKeyDown}
      >
        <div className='number-area'>
          {currentIndex + 1} / {maxIndex} 回目
        </div>

        <div
          className='english-area'
          style={{ cursor: 'pointer' }}
          title='読み上げ'
          onClick={() => speak(questionText.en2, 'en-US')}
        >
          {Array.isArray(questionTextArray1) &&
            questionTextArray1.map((character, index) => (
              <span
                key={index}
                className={index < questionIndex1 ? 'correct-char' : ''}
              >
                {character}
              </span>
            ))}
        </div>

        <div className='hiragana-area'>
          {Array.isArray(questionTextArray2) &&
            questionTextArray2.map((character, index) => (
              <span
                key={index}
                className={index < questionIndex2 ? 'correct-char' : ''}
              >
                {character}
              </span>
            ))}
        </div>

        <div className='kanji-area'>{questionText.jp2}</div>

        <div className='key-area'>
          最後に押されたキー: <span>{pressedKey}</span>
        </div>
        <div className='button-container'>
          <button onClick={handlePrevious}>
            <GrLinkPrevious />
            <span>前の問題</span>
          </button>
          <button onClick={handleNext}>
            <span>次の問題</span>
            <GrLinkNext />
          </button>
        </div>
      </div>
    )
  }

  return (
    <>
      <div className='typing-container'>{content}</div>
      <VolumeControl />
    </>
  )
}
export default TypingContent
