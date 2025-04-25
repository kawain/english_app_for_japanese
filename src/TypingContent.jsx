import { useState, useEffect, useRef, useCallback } from 'react'
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
  // --- タイピング速度計測用の状態 ---
  const [startTime, setStartTime] = useState(null) // 問題ごとのタイピング開始時間
  const [currentCPM, setCurrentCPM] = useState(0) // 現在の問題のCPM
  const [averageCPM, setAverageCPM] = useState(0) // 全体の平均CPM
  const [allProblemStats, setAllProblemStats] = useState([]) // 各問題の統計 [{ cpm: number, duration: number, charCount: number }]
  const isTypingStartedForProblem = useRef(false) // 現在の問題でタイピングが開始されたかフラグ

  // WASMの関数で問題のセットアップと問題数を返す
  useEffect(() => {
    const initializeTyping = async () => {
      try {
        const result = await window.CreateTyping()
        console.log('タイピング問題セットアップ完了:', result)
        setMaxIndex(result)
      } catch (error) {
        console.error('タイピング問題のセットアップに失敗しました:', error)
      }
    }
    initializeTyping()
  }, [])

  // フォーカスをする
  useEffect(() => {
    if (progress === 0) return
    typingAreaRef.current?.focus()
  }, [progress])

  // 問題選択ロジック
  const selectQuestion = useCallback(
    async (index, startFlag = false) => {
      setCurrentIndex(index)
      const question = await window.GetTypingQuestion(index)
      setQuestionText(question)
      const array1 = await window.GetTypingQuestionSlice(1)
      setQuestionTextArray1(array1)
      const array2 = await window.GetTypingQuestionSlice(2)
      setQuestionTextArray2(array2)

      setQuestionIndex1(0)
      setQuestionIndex2(0)
      setInputCharacters('')
      setPressedKey('')

      // --- 速度関連の状態リセット ---
      setStartTime(null) // 開始時間をリセット
      // setCurrentCPM(0) // 現在の問題のCPMをリセット
      isTypingStartedForProblem.current = false // タイピング開始フラグをリセット
      whichRef.current = 1 // 英語から開始
      // --- ここまで ---

      if (startFlag) {
        setProgress(1)
      }

      await speak(question.en2, 'en-US')
    },
    [speak]
  )

  // タイピング開始
  const handleStart = useCallback(async () => {
    // --- 全体の統計情報をリセット ---
    setAllProblemStats([])
    setAverageCPM(0)
    // --- ここまで ---
    await selectQuestion(0, true)
    // 開始時にフォーカス
    typingAreaRef.current?.focus()
  }, [selectQuestion])

  // 前の問題へ
  const handlePrevious = useCallback(async () => {
    if (timerIdRef.current) {
      clearTimeout(timerIdRef.current)
      timerIdRef.current = null
    }
    let index = currentIndex
    if (index > 0) {
      index--
    } else {
      index = maxIndex - 1
    }
    await selectQuestion(index)
    typingAreaRef.current?.focus()
  }, [currentIndex, maxIndex, selectQuestion])

  // 次の問題へ
  const handleNext = useCallback(async () => {
    if (timerIdRef.current) {
      clearTimeout(timerIdRef.current)
      timerIdRef.current = null
    }
    let index = currentIndex
    if (index < maxIndex - 1) {
      index++
    } else {
      index = 0
    }
    await selectQuestion(index)
    typingAreaRef.current?.focus()
  }, [currentIndex, maxIndex, selectQuestion])

  // キー入力処理
  const handleKeyDown = useCallback(
    async e => {
      e.preventDefault()

      // 次の問題への遷移中は入力を無視
      if (timerIdRef.current) return

      const moji = e.key

      // --- 最初の有効なキー入力でタイマー開始 ---
      if (!isTypingStartedForProblem.current && moji.length === 1) {
        // 1文字のキー入力で開始
        setStartTime(Date.now())
        isTypingStartedForProblem.current = true
      }

      let input = inputCharacters + moji
      setInputCharacters(input)
      setPressedKey(moji)

      let result = 0
      let currentQuestionIndex = 0
      let currentQuestionArray = []
      let nextWhich = whichRef.current
      let questionCompleted = false // 問題完了フラグ

      if (whichRef.current === 1) {
        // 英語
        currentQuestionIndex = questionIndex1
        currentQuestionArray = questionTextArray1
        result = await window.TypingKeyDown(input, currentQuestionIndex, 1)
        if (result > currentQuestionIndex) {
          setQuestionIndex1(result)
          setInputCharacters('')
          if (result >= currentQuestionArray.length) {
            nextWhich = 2 // 日本語へ
          }
        }
      } else if (whichRef.current === 2) {
        // 日本語
        currentQuestionIndex = questionIndex2
        currentQuestionArray = questionTextArray2
        result = await window.TypingKeyDown(input, currentQuestionIndex, 2)
        if (result > currentQuestionIndex) {
          setQuestionIndex2(result)
          setInputCharacters('')
          if (result >= currentQuestionArray.length) {
            // --- 問題完了時の速度計算 ---
            questionCompleted = true
            const endTime = Date.now()
            if (startTime) {
              // startTimeが記録されている場合のみ計算
              const duration = (endTime - startTime) / 1000 // 秒
              const totalChars =
                questionTextArray1.length + questionTextArray2.length
              const cpm =
                totalChars > 0 && duration > 0
                  ? Math.round((totalChars / duration) * 60)
                  : 0

              setCurrentCPM(cpm) // 現在の問題のCPMをセット

              const newStat = {
                cpm,
                duration,
                charCount: totalChars,
                index: currentIndex
              }
              const updatedStats = [...allProblemStats, newStat]
              setAllProblemStats(updatedStats)

              // 平均CPMを計算
              const totalCPM = updatedStats.reduce(
                (sum, stat) => sum + stat.cpm,
                0
              )
              const avgCPM =
                updatedStats.length > 0
                  ? Math.round(totalCPM / updatedStats.length)
                  : 0
              setAverageCPM(avgCPM)
            }
            // --- 計算ここまで ---

            // 次の問題へ遷移
            timerIdRef.current = setTimeout(async () => {
              const nextIndex =
                currentIndex + 1 >= maxIndex ? 0 : currentIndex + 1
              await selectQuestion(nextIndex)
              timerIdRef.current = null
            }, 500) // 500ms待機

            nextWhich = 1 // 次は英語から
          }
        }
      }

      // whichRef の更新は状態遷移が発生した場合のみ
      if (nextWhich !== whichRef.current) {
        whichRef.current = nextWhich
      }
    },
    [
      inputCharacters,
      questionIndex1,
      questionIndex2,
      questionTextArray1,
      questionTextArray2,
      startTime,
      currentIndex,
      maxIndex,
      selectQuestion,
      allProblemStats
    ]
  )

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
          onClick={async () => await speak(questionText.en2, 'en-US')}
        >
          {Array.isArray(questionTextArray1) &&
            questionTextArray1.map((character, index) => (
              <span
                key={`en-${index}`}
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
                key={`jp-${index}`}
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

        <div className='stats-area'>
          <span>前回のCPM: {currentCPM > 0 ? currentCPM : '-'}</span>
          <span>平均CPM: {averageCPM > 0 ? averageCPM : '-'}</span>
        </div>

        <p>CPM: 1分あたりに入力できる文字数</p>

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
