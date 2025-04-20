import { useState, useEffect, useRef } from 'react'
import { useAppContext } from './App.jsx'
import { addExcludedWordId } from './Storage.jsx'
import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'
import QuizChoices from './components/QuizChoices.jsx'
import { tts } from './utils/tts'

// 問題の選択肢の数
const numberOfChoices = 10

function WordQuizContent () {
  const { selectedLevel, volume, isSoundEnabled } = useAppContext()
  // 0: 初期状態スタートボタン表示
  // 1: 問題と選択肢と回答するボタン表示
  // 2: 回答と次の問題ボタン表示
  const [progress, setProgress] = useState(0)
  // WASMからもらうデータ
  const [currentQuiz, setCurrentQuiz] = useState(null)
  // WASMからもらうデータ
  const [quizChoices, setQuizChoices] = useState([])
  // ユーザーの選択
  const [selectedChoiceId, setSelectedChoiceId] = useState(null)
  // ユーザーの答えの正誤
  const [answerResult, setAnswerResult] = useState(null)
  // 正解率統計
  const [correctCount, setCorrectCount] = useState(0)
  const [totalQuestions, setTotalQuestions] = useState(0)
  // CSSエフェクトのハイライト
  const [isHighlighted, setIsHighlighted] = useState(false)
  // 現在のレベルを保存する ref
  const prevLevelRef = useRef(selectedLevel)

  const fetchQuizData = () => {
    try {
      // WASMの関数
      const quizData = window.CreateQuiz(
        parseInt(selectedLevel, 10),
        numberOfChoices
      )
      // WASMの関数
      const choicesData = window.CreateQuizChoices()
      if (!quizData || !choicesData) {
        console.error('クイズデータの取得に失敗しました。')
        return false
      }
      setCurrentQuiz(quizData)
      setQuizChoices(choicesData)
      return true
    } catch (error) {
      console.error('Error fetching quiz data:', error)
      return false
    }
  }

  const handleStart = () => {
    const result = fetchQuizData()
    if (result) {
      setProgress(1)
      setSelectedChoiceId(null)
      setAnswerResult(null)
      setCorrectCount(0)
      setTotalQuestions(0)
    } else {
      alert('クイズデータの読み込みに失敗しました')
      setProgress(0)
    }
  }

  const handleNext = () => {
    const result = fetchQuizData()
    if (result) {
      setProgress(1)
      setSelectedChoiceId(null)
      setAnswerResult(null)
    } else {
      alert('次のクイズデータの読み込みに失敗しました')
      setProgress(0)
    }
  }

  const handleChoiceChange = event => {
    // 回答後は選択を変更できないようにする
    if (answerResult !== null) return
    setSelectedChoiceId(event.target.value)
  }

  const handleAnswer = () => {
    let result
    let isCorrect = false

    if (selectedChoiceId === null) {
      result = 'passed' // 未選択の場合はパス扱い
    } else if (String(selectedChoiceId) === String(currentQuiz.id)) {
      result = 'correct'
      isCorrect = true
    } else {
      result = 'incorrect'
    }

    setAnswerResult(result)
    setTotalQuestions(prev => prev + 1)

    if (isCorrect) {
      setCorrectCount(prev => prev + 1)
      try {
        // 正解した単語をストレージに追加
        addExcludedWordId(String(currentQuiz.id))
        window.AddStorage(currentQuiz.id)
      } catch (error) {
        console.error('Error saving to storage:', error)
      }
    }

    setProgress(2)
  }

  const calculateAccuracy = () => {
    if (totalQuestions === 0) return 0
    return Math.round((correctCount / totalQuestions) * 100)
  }

  // クイズ中のハイライト効果
  useEffect(() => {
    if (progress === 0) return

    let timerId = null

    setIsHighlighted(true)
    timerId = setTimeout(() => {
      setIsHighlighted(false)
    }, 500)

    return () => {
      if (timerId) {
        clearTimeout(timerId)
      }
    }
  }, [progress])

  // レベル変更時にクイズをリセットする useEffect
  useEffect(() => {
    if (progress > 0 && prevLevelRef.current !== selectedLevel) {
      setProgress(0)
      setCurrentQuiz(null)
      setQuizChoices([])
      setSelectedChoiceId(null)
      setAnswerResult(null)
      setCorrectCount(0)
      setTotalQuestions(0)
    }
    // 現在のレベルを ref に保存して次回の比較に使用
    prevLevelRef.current = selectedLevel
  }, [selectedLevel])

  // クイズ内容の読み上げのための TTS
  useEffect(() => {
    if (!currentQuiz || progress === 0) return

    const textToSpeak = progress === 1 ? currentQuiz.en : currentQuiz.en2
    if (textToSpeak) {
      tts(textToSpeak, 'en-US', volume, isSoundEnabled).catch(error => {
        console.error('TTS エラー:', error)
      })
    }
  }, [progress, currentQuiz, volume, isSoundEnabled])

  let content = null

  if (progress === 0) {
    content = (
      <button onClick={handleStart}>
        単語クイズ(レベル {selectedLevel}) スタート
      </button>
    )
  } else if (progress === 1) {
    content = (
      <>
        <h2
          className={isHighlighted ? 'highlight' : ''}
          style={{ cursor: 'pointer' }}
          title='読み上げ'
          onClick={() =>
            tts(currentQuiz.en, 'en-US', volume, isSoundEnabled).catch(
              error => {
                console.error('TTS エラー:', error)
              }
            )
          }
        >
          {currentQuiz.en}
        </h2>
        <div className='quiz-content'>
          <QuizChoices
            quizChoices={quizChoices}
            selectedChoiceId={selectedChoiceId}
            onChoiceChange={handleChoiceChange}
            disabled={answerResult !== null}
          />
          <div>
            <button onClick={handleAnswer}>答える / パス</button>
          </div>
        </div>
      </>
    )
  } else if (progress === 2) {
    content = (
      <>
        <h2
          style={{ cursor: 'pointer' }}
          title='読み上げ'
          onClick={() =>
            tts(currentQuiz.en, 'en-US', volume, isSoundEnabled).catch(
              error => {
                console.error('TTS エラー:', error)
              }
            )
          }
        >
          {currentQuiz.en}
        </h2>
        <div className='quiz-content'>
          <QuizChoices
            quizChoices={quizChoices}
            selectedChoiceId={selectedChoiceId}
            onChoiceChange={handleChoiceChange}
            disabled={answerResult !== null}
          />
          <div>
            <div className='quiz-result'>
              {answerResult === 'correct' && (
                <p
                  className={isHighlighted ? 'highlight' : ''}
                  style={{ color: 'green' }}
                >
                  【正解です】
                </p>
              )}
              {answerResult === 'incorrect' && (
                <p
                  className={isHighlighted ? 'highlight' : ''}
                  style={{ color: 'red' }}
                >
                  【間違いです】
                </p>
              )}
              {answerResult === 'passed' && (
                <p className={isHighlighted ? 'highlight' : ''}>
                  【パスしました】
                </p>
              )}

              <div className='quiz-details'>
                <p>{currentQuiz.en}</p>
                <p>{currentQuiz.jp}</p>
                <p
                  className={isHighlighted ? 'highlight' : ''}
                  style={{ cursor: 'pointer' }}
                  title='読み上げ'
                  onClick={() =>
                    tts(currentQuiz.en2, 'en-US', volume, isSoundEnabled).catch(
                      error => {
                        console.error('TTS エラー:', error)
                      }
                    )
                  }
                >
                  {currentQuiz.en2}
                </p>
                <p>{currentQuiz.jp2}</p>
              </div>

              <div className='quiz-stats'>
                <p>
                  正解数: {correctCount} / {totalQuestions}
                </p>
                <p>正解率: {calculateAccuracy()}%</p>
              </div>
            </div>
            <button onClick={handleNext}>次の問題を行う</button>
          </div>
        </div>
      </>
    )
  }

  return (
    <>
      <div className='quiz-container'>{content}</div>
      <VolumeControl />
      <LevelControl />
    </>
  )
}
export default WordQuizContent
